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
	"strconv"
	"strings"

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
			for _, exclude := range llc.ExcludePaths {
				if strings.HasPrefix(path, exclude) {
					return nil
				} else if strings.HasPrefix(exclude, "*.") {
					if filepath.Ext(path) == filepath.Ext(exclude) {
						return nil
					}
				} else if string(exclude[len(exclude)-1]) != "/" { // its a file
					if path == exclude {
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
		// ignore lines
		// if the line is commented then skip it
		// This is for go or programming comments only

		// if len(line) >= 2 && line[:2] == "//" {
		// 	continue
		// }
		// // comments for yaml or yml files. Do not consider that line if line start with a comment
		// if len(line) > 1 && line[:1] == "#" {
		// 	continue
		// }
		// // comments for yaml or yml files. If there is comment in the line take only uncommented part
		// index := strings.Index(line, "#")
		// if index > 0 {
		// 	line = line[0:index]
		// }
		// // This is for go or programming code only as comments in yaml files start with #
		// // start
		// if len(line) >= 2 && line[:2] == "/*" {
		// 	//TODO here
		// 	skip = true
		// }
		// if strings.Contains(line, "*/") {
		// 	skip = false
		// }
		// if skip {
		// 	continue
		// }

		if len(link) >= 8 && strings.ToLower(strings.Trim(link, " ")[0:4]) != "http" { // do not consider it as url if it dies not start with http or https
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
	err = s.Err()
	if err != nil {
		return err
	}
	return nil
}

func checkURL(link string) (int, error) {
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

func (llc *LinkLintConfig) LintAll() bool {
	isFatal := false
	count := 0
	for key := range llc.LinkMap {
		count++
		fmt.Println("Currently checking ", count, " url(s) out of ", len(llc.LinkMap))
		if !IsURL(key) {
			isFatal = true
			llc.OnFail("Invalid URL", key)
			continue
		}

		statusCode, err := checkURL(key)
		if err != nil {
			isFatal = true
			llc.OnFail(err.Error(), key)
			continue
		}

		accepted := false
		for _, code := range llc.AcceptStatusCodes {
			if code == statusCode {
				llc.OnPass("http Status-code "+strconv.Itoa(statusCode), key)
				accepted = true
				break
			}
		}

		if accepted {
			continue
		} else {
			isFatal = true
			llc.OnFail("http Status-code "+strconv.Itoa(statusCode), key)
		}
	}

	return isFatal
}

func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
