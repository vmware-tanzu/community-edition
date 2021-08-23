// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func ReadFileAndReplaceContents(filename string, findReplaceMap map[string]string) (string, error) {
	byteContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	contents := string(byteContents)
	for k, v := range findReplaceMap {
		contents = strings.Replace(contents, k, v, -1)
	}

	return contents, nil
}

func ReadFileAndReplaceContentsTempFile(filename string, findReplaceMap map[string]string) (string, error) {
	contents, err := ReadFileAndReplaceContents(filename, findReplaceMap)
	if err != nil {
		return "", err
	}

	file, err := ioutil.TempFile("", fmt.Sprintf("%s-*%s", strings.TrimPrefix(filepath.Base(filename), filepath.Ext(filename)), filepath.Ext(filename)))
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write([]byte(contents))
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}
