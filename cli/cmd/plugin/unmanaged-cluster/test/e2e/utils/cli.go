// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/tanzu"
)

func InstallTCE() error {
	topDirPath := os.Getenv("ROOT_DIR")
	err := os.Chdir(topDirPath)
	if err != nil {
		log.Println("error while changing directory to the top:", err)
		return err
	}

	// installing dependencies i.e docker and kubectl
	err = runDeployScript("hack/ensure-deps/ensure-docker.sh")
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = runDeployScript("hack/ensure-deps/ensure-kubectl.sh")
	if err != nil {
		log.Fatal(err)
		return err
	}

	return runDeployScript("test/download-or-build-tce.sh")
}

func runDeployScript(filename string) error {
	mwriter := io.MultiWriter(os.Stdout)
	cmd := exec.Command("/bin/bash", filename)
	cmd.Stderr = mwriter
	cmd.Stdout = mwriter
	err := cmd.Run() // blocks until sub process is complete
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func Kubectl(input io.Reader, args ...string) (string, error) {
	return cliRunner("kubectl", input, args...)
}

func Tanzu(input io.Reader, args ...string) (string, error) {
	return cliRunner("tanzu", input, args...)
}

func cliRunner(name string, input io.Reader, args ...string) (string, error) {
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

func ContainsUC(ucLists []tanzu.Cluster, e string) bool {
	for _, uc := range ucLists {
		if uc.Name == e {
			return true
		}
	}
	return false
}
