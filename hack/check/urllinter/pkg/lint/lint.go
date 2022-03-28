// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package lint

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"

	xurls "mvdan.cc/xurls/v2"
)

type LinkLintConfig struct {
	IncludeExts       []string              `yaml:"includeExts"`
	ExcludeLinks      []string              `yaml:"excludeLinks"`
	ExcludePaths      []string              `yaml:"excludePaths"`
	AcceptStatusCodes []int                 `yaml:"acceptStatusCodes"`
	LinkMap           map[string][]LinkLint // consists map as the key and file details as values
}

type LinkLint struct {
	Path     string
	Line     string
	Position Position
	Message  string
	Status   string
}

type Position struct {
	Row, Col int
}

func New(configFile string) (*LinkLintConfig, error) {
	if configFile == "" {
		return nil, errors.New("configuration file cannot be empty")
	}
	file, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	llc := &LinkLintConfig{}
	err = yaml.Unmarshal(file, llc)
	if err != nil {
		return nil, err
	}
	llc.LinkMap = make(map[string][]LinkLint)
	return llc, nil
}

func NewFromContent(content []byte) (*LinkLintConfig, error) {
	llc := &LinkLintConfig{}
	err := yaml.Unmarshal(content, llc)
	if err != nil {
		return nil, err
	}
	llc.LinkMap = make(map[string][]LinkLint)
	return llc, nil
}

func (llc *LinkLintConfig) Init(dir string) error {
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip directories
			if info.IsDir() {
				return nil
			}
			// tmp variable to manipulate the filename based on a simple match
			// if using a full qualified path:
			// 		/home/bob/go/src/github.com/vmware-tanzu/community-edition/addons/packages/harbor/2.2.3/bundle/.imgpkg/images.yml
			// wouldnt match because it was leading with:
			// 		/home/bob/go/src/github.com/vmware-tanzu/community-edition/
			// the solution was to remove the initial path and the leading "/" for the
			// filepath.Match(match, tmp) to return true
			tmp := path

			// fix for fully qualified path
			if strings.Index(tmp, dir) == 0 {
				tmp = strings.Replace(tmp, dir, "", 1)
			}
			// remove leading /
			if strings.Index(tmp, "/") == 0 {
				tmp = strings.Replace(tmp, "/", "", 1)
			}

			for _, exclude := range llc.ExcludePaths {
				if strings.HasPrefix(tmp, exclude) {
					return nil
				} else if strings.HasPrefix(exclude, "*.") {
					if filepath.Ext(tmp) == filepath.Ext(exclude) {
						return nil
					}
				} else if string(exclude[len(exclude)-1]) != "/" { // its a file
					if tmp == exclude {
						return nil
					}
				}
			}
			ext := filepath.Ext(path)
			for _, ex := range llc.IncludeExts {
				if ext == ex {
					err = llc.ReadFile(path)
					if err != nil {
						return err
					}
				}
			}
			return nil
		})
	if err != nil {
		return err
	}

	// TODO:  Put this behind a debug option in the future
	// Keep around for debugging. It's a dump of world.
	// fmt.Printf("\n\nDump All Links:\n")
	// for link, objs := range llc.LinkMap {
	// 	fmt.Printf("Link: %s\n", link)
	// 	for _, obj := range objs {
	// 		fmt.Printf("\tPath: %s\n", obj.Path)
	// 	}
	// }
	// fmt.Printf("\n")

	return err
}

func (llc *LinkLintConfig) ReadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	count := 1
	rxStrict := xurls.Strict()
	for s.Scan() {
		skip := false
		line := strings.Trim(s.Text(), " ")
		link := rxStrict.FindString(line)
		col := strings.Index(s.Text(), link)

		// do not consider it as url if it does not start with http or https
		if len(link) >= 8 && strings.ToLower(strings.Trim(link, " ")[0:4]) != "http" {
			continue
		}

		for _, l := range llc.ExcludeLinks {
			if strings.Contains(link, l) {
				skip = true
				break
			}
		}
		if link != "" && !skip {
			llints := llc.LinkMap[link]
			llc.LinkMap[link] = append(llints, LinkLint{Path: path, Line: link, Position: Position{Row: count, Col: col}, Status: "", Message: ""})
		}
		count++
	}
	return s.Err()
}

func checkURL(link string) (int, error) {
	if !IsURL(link) {
		return 0, errors.New("invalid URL")
	}

	/* #nosec G402 */
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cli := &http.Client{Transport: tr}
	ctx := context.Background()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, link, bytes.NewBuffer([]byte("")))
	resp, err := cli.Do(req)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	return resp.StatusCode, nil
}

type checkResult struct {
	StatusCode int
	URL        string
	Error      error
}

func (llc *LinkLintConfig) isValidCode(statusCode int) bool {
	valid := false
	for _, code := range llc.AcceptStatusCodes {
		if code == statusCode {
			valid = true
			break
		}
	}
	return valid
}

func (llc *LinkLintConfig) LintAll() bool {
	linkCount := len(llc.LinkMap)

	// Limit the number of GET requests so we don't get rate limited
	results := make(chan checkResult, 2)
	wg := sync.WaitGroup{}

	isFatal := false
	count := 0
	go func() {
		// Loop through our results as they come in and print out the results
		for res := range results {
			wg.Done()
			count++
			fmt.Printf("Result %d of %d url(s)\n", count, linkCount)
			if res.Error != nil {
				isFatal = true
				llc.OnFail(res.Error.Error(), res.URL)
				continue
			}

			// No error, but check that the status code was acceptable
			if llc.isValidCode(res.StatusCode) {
				llc.OnPass(fmt.Sprintf("HTTP Status Code: %d", res.StatusCode), res.URL)
			} else {
				isFatal = true
				llc.OnFail(fmt.Sprintf("HTTP Status Code: %d", res.StatusCode), res.URL)
			}
		}
	}()

	// Go through each found URL and validate it exists
	for key := range llc.LinkMap {
		wg.Add(1)

		go func(key string) {
			statusCode, err := checkURL(key)
			results <- checkResult{
				StatusCode: statusCode,
				Error:      err,
				URL:        key,
			}
		}(key)
	}
	wg.Wait()

	return isFatal
}

func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
