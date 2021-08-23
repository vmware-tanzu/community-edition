// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"bytes"
	"os/exec"
	"strings"
)

var rootDir string

// RootDir returns the root directory for this git repository
func RootDir() string {
	if rootDir != "" {
		return rootDir
	}

	var stdout bytes.Buffer
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	rootDir = strings.TrimSpace(stdout.String())
	return rootDir
}
