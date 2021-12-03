// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/onsi/ginkgo"
)

var (
	WorkingDir string
)

func init() {
	WorkingDir = getCurrentDir()
	if WorkingDir == "" {
		log.Fatal("current working directory is empty")
	}
}

func Kubectl(input io.Reader, args ...string) (string, error) {
	return cliRunner("kubectl", input, args...)
}

func Tanzu(input io.Reader, args ...string) (string, error) {
	return cliRunner("tanzu", input, args...)
}

func Aws(args ...string) (string, error) {
	return cliRunner("aws", nil, args...)
}

func cliRunner(name string, input io.Reader, args ...string) (string, error) {
	fmt.Fprintf(ginkgo.GinkgoWriter, "+ %s %s\n", name, strings.Join(args, " "))

	var stdOut bytes.Buffer
	mwriter := io.MultiWriter(&stdOut, os.Stderr)
	cmd := exec.Command(name, args...)
	cmd.Stdin = input
	cmd.Stdout = mwriter
	cmd.Stderr = mwriter

	err := cmd.Run()
	if err != nil {
		rc := -1
		if ee, ok := err.(*exec.ExitError); ok {
			rc = ee.ExitCode()
		}

		return "", fmt.Errorf("%s\nexit status: %d", err.Error(), rc)
	}

	return stdOut.String(), err
}

func GetClusterContext(clusterName string) string {
	clusterContext := clusterName + "-admin@" + clusterName
	return clusterContext
}

func getCurrentDir() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Println("error while getting current working directory", err)
	}

	return wd
}
