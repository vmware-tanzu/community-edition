//go:build e2e

// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/exec"
	"testing"
)

var tceRepoPath string
var clusterName string
var conformanceDirPath string
var unmanagedClusterDirPath string
var conformanceResultFilename string

func initializeVariable() {
	// Variables used across functions
	tceRepoPath = os.Getenv("ROOT_DIR")
	conformanceDirPath = tceRepoPath + "/cli/cmd/plugin/conformance"
	unmanagedClusterDirPath = tceRepoPath + "/cli/cmd/plugin/unmanaged-cluster"
	rand, _ := rand.Int(rand.Reader, big.NewInt(1000))
	clusterName = "uc" + rand.String() + "conformance-test"
}

func TestSetup(t *testing.T) {
	initializeVariable()

	err := os.Chdir(unmanagedClusterDirPath)
	if err != nil {
		t.Errorf("error while changing directory to unmanged-cluster: %v", err)
	}
	// Creating Unmanaged cluster
	// Future task check the way cliRunner() is printing the log so that the space
	// between the logs get removed and logs look nice
	_, err = cliRunner("go", nil, "run", ".", "create", "--tty-disable", clusterName)
	if err != nil {
		t.Errorf("Error while Unmanaged Cluster creation: %v", err)
	}
}

func TestConformance(t *testing.T) {
	err := os.Chdir(conformanceDirPath)
	if err != nil {
		t.Errorf("error while changing directory to conformance: %v", err)
	}

	// Running Conformance Plugin
	_, err = cliRunner("go", nil, "run", ".", "run", "--wait")
	if err != nil {
		t.Errorf("Error while running Conformance plugin: %v", err)
	}

	// Retrieving Conformance result file
	conformanceResultFilename, err = cliRunner("go", nil, "run", ".", "retrieve")
	if err != nil {
		t.Errorf("Error while retrieving the conformance result filename: %v", err)
	}
	conformanceResultFilename = "/" + conformanceResultFilename[:len(conformanceResultFilename)-1]

	// Extracting the result
	conformanceResult, err := cliRunner("go", nil, "run", ".", "results", conformanceDirPath+conformanceResultFilename)
	if err != nil {
		t.Errorf("Error while fetching the result from tarball produced by conformance: %v", err)
	}
	fmt.Printf("%v", conformanceResult)
}

func TestCleanUp(t *testing.T) {
	// Deleting Conformance result file
	err := os.Remove(conformanceDirPath + conformanceResultFilename)
	if err != nil {
		t.Errorf("Error while deleting conformance result file: %v", err)
	}

	// Deleting Conformance
	_, err = cliRunner("go", nil, "run", ".", "delete", "--wait")
	if err != nil {
		t.Errorf("Error while conformance deletion: %v", err)
	}

	// Deleting Unmanaged cluster
	err = os.Chdir(unmanagedClusterDirPath)
	if err != nil {
		t.Errorf("error while changing directory to unmanged-cluster: %v", err)
	}

	_, err = cliRunner("go", nil, "run", ".", "delete", clusterName)
	if err != nil {
		t.Errorf("Error while Unmanaged Cluster deletion: %v", err)
	}
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
	return stdOut.String(), nil
}
