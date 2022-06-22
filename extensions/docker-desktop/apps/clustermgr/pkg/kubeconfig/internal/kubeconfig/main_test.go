// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kubeconfig

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func copyFile(src, dst string) {
	bytesRead, err := os.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(dst, bytesRead, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	copyFile(filepath.Join(".", "testdata", "config.orig"), filepath.Join(".", "testdata", "config"))
	exitCode := m.Run()
	os.Remove(filepath.Join(".", "testdata", "config"))
	os.Exit(exitCode)
}
