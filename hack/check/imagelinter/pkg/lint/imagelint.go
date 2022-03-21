// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package lint

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type ImageLintConfig struct {
	IncludeExts       []string               `yaml:"includeExts"`
	MatchPattern      []string               `yaml:"matchPattern"`
	IncludeLines      []string               `yaml:"includeLines"`
	IgnoreImages      []string               `yaml:"ignoreImages"`
	SuccessValidators []string               `yaml:"succesValidators"`
	FailureValidators []string               `yaml:"failureValidators"`
	ImageMap          map[string][]ImageLint // consists map as the key and file details as values
}

func New(configFile string) (*ImageLintConfig, error) {
	file, err := os.ReadFile(configFile)

	if err != nil {
		return nil, err
	}

	ilc := &ImageLintConfig{}

	err = yaml.Unmarshal(file, ilc)
	if err != nil {
		return nil, err
	}
	ilc.ImageMap = make(map[string][]ImageLint)
	return ilc, nil
}

func NewFromContent(content []byte) (*ImageLintConfig, error) {
	ilc := &ImageLintConfig{}
	err := yaml.Unmarshal(content, ilc)
	if err != nil {
		return nil, err
	}
	ilc.ImageMap = make(map[string][]ImageLint)
	return ilc, nil
}

type ImageLint struct {
	Path     string
	Status   string
	Position Position
}
type Position struct {
	Row, Col int
}

func (imc *ImageLintConfig) Init(dir string) error {
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
			// skip .git directory
			if strings.Index(tmp, ".git") == 0 {
				return nil
			}
			// end

			// if the pattern is not match move next
			matched := false
			for _, match := range imc.MatchPattern {
				m, _ := filepath.Match(match, tmp)
				if m {
					for _, ext := range imc.IncludeExts {
						if ext == filepath.Ext(path) {
							err = imc.ReadFile(path)
							if err != nil {
								return err
							}
						}
					}
					matched = true
					break
				}
			}
			if !matched {
				goto jumpOut
			}
		jumpOut:
			return nil
		})
	if err != nil {
		return err
	}
	return err
}

func (imc *ImageLintConfig) ReadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	count := 1
	skip := false
	for s.Scan() {
		// ignore lines
		// if the line is commented then skip it
		line := strings.Trim(s.Text(), " ")

		// if it is a comment line just ignore it
		if IsComment(line) {
			continue
		}
		// comments for yaml or yml files. If there is comment in the line take only uncommented part
		index := strings.Index(line, "#")
		if index > 0 {
			line = line[0:index]
		}
		// This is for go or programming code only as comments in yaml files start with #
		// start
		if len(line) >= 2 && line[:2] == "/*" {
			skip = true
		}
		if strings.Contains(line, "*/") {
			skip = false
		}
		if skip {
			continue
		}

		for _, searchterm := range imc.IncludeLines {
			if strings.Contains(line, searchterm) {
				index := strings.Index(line, searchterm) + len(searchterm)
				if strings.Trim(line[index:], " ") != "" {
					ln := removeChars(line[index:])
					if CanIgnore(ln) {
						continue
					}
					if imc.CanIgnoreImage(ln) {
						continue
					}
					ilints := imc.ImageMap[ln]
					imc.ImageMap[ln] = append(ilints, ImageLint{Path: path, Position: Position{Row: count, Col: index}, Status: "YetToLint"})
				}
			}
		}
		count++
	}
	err = s.Err()
	if err != nil {
		return err
	}
	return nil
}

func CanIgnore(line string) bool {
	minLineLen := 5
	ignores := []string{"%", "$", "%", "{", "}", "...", ",", " "}
	if len(line) < minLineLen {
		return true
	}
	for _, s := range ignores {
		if strings.Contains(strings.Trim(line, " "), s) {
			return true
		}
	}
	return false
}

func IsComment(line string) bool {
	if len(line) >= 2 && line[:2] == "//" {
		return true
	}
	if len(line) > 1 && line[:1] == "#" {
		return true
	}
	return false
}

func (imc *ImageLintConfig) CanIgnoreImage(line string) bool {
	for _, ln := range imc.IgnoreImages {
		if ln == line {
			return true
		}
	}
	return false
}
func removeChars(line string) string {
	line = strings.Trim(line, " ")
	line = strings.Trim(line, `"`)
	line = strings.Trim(line, `'`)
	return line
}
