//go:build e2e

// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"os/exec"
	"testing"
)

var myDir string
var tceRepoPath string
var tanzuDiagnosticPluginDir string
var tanzuDiagnosticBin string
var clusterName string
var clusterKubeContext string
var home string

func initializeVariable() {
	// variables used across functions
	myDir, _ = os.Getwd()
	tceRepoPath = os.Getenv("ROOT_DIR")
	tanzuDiagnosticPluginDir = myDir + "/.."
	tanzuDiagnosticBin = myDir + "/tanzu-diagnostics-e2e-bin"
	home = os.Getenv("HOME")
	rand, _ := rand.Int(rand.Reader, big.NewInt(1000))
	clusterName = "uc" + rand.String() + "test"
	clusterKubeContext = "kind-" + clusterName
}

func TestSetup(t *testing.T) {
	initializeVariable()
	// Installing TCE
	_, err := exec.LookPath("tanzu")
	if err != nil {
		err := os.Chdir(tceRepoPath)
		if err != nil {
			t.Errorf("error while changing directory %v", err)
		}
		err = runDeployScript("test/download-or-build-tce.sh")
		if err != nil {
			t.Errorf("Error while installing TCE: %v", err)
		}
	}

	err = os.Chdir(tceRepoPath + "/cli/cmd/plugin/unmanaged-cluster")
	if err != nil {
		t.Errorf("error while changing directory to unmanged-cluster: %v", err)
	}
	// creating Unmanaged cluster
	// Future Changes https://github.com/vmware-tanzu/community-edition/pull/4106#discussion_r860183492
	err = cliRunner("go", nil, "run", ".", "create", clusterName)
	if err != nil {
		t.Errorf("Error while Unmanaged Cluster creation: %v", err)
	}
}

func TestDiagnostic(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "diagnosticE2Etest")
	if err != nil {
		t.Errorf("error while creating temp directory: %v", err)
	}

	err = os.Chdir(tanzuDiagnosticPluginDir)
	if err != nil {
		t.Errorf("error while changing directory to Diagnostic plugin: %v", err)
	}

	diagnostic := []string{"run", ".", "collect",
		"--bootstrap-cluster-name", clusterName,
		"--management-cluster-kubeconfig", home + "/.kube/config",
		"--management-cluster-context", clusterKubeContext,
		"--management-cluster-name", clusterName,
		"--workload-cluster-infra", "docker",
		"--workload-cluster-kubeconfig", home + "/.kube/config",
		"--workload-cluster-context", clusterKubeContext,
		"--workload-cluster-name", clusterName,
		"--unmanaged-cluster-kubeconfig", home + "/.kube/config",
		"--unmanaged-cluster-context", clusterKubeContext,
		"--unmanaged-cluster-name", clusterName,
		"--output-dir", tempDir}

	// Collecting logs for all the cluster
	err = cliRunner("go", nil, diagnostic[:]...)
	if err != nil {
		t.Errorf("error while collectiong diagnostic %v", err)
	}

	// Checking all the log files
	diagnosticFiles := [4]string{"/bootstrap." + clusterName + ".diagnostics.tar.gz",
		"/management-cluster." + clusterName + ".diagnostics.tar.gz",
		"/workload-cluster." + clusterName + ".diagnostics.tar.gz",
		"/unmanaged-cluster." + clusterName + ".diagnostics.tar.gz"}
	for _, element := range diagnosticFiles {
		_, err = os.Stat(tempDir + element)
		if err != nil {
			t.Errorf("error while looking for %v %v", element, err)
		}
	}

	// cleaning th temp directory
	err = os.RemoveAll(tempDir)
	if err != nil {
		t.Error(err)
	}

	// deleting Unmanaged cluster
	// Future Changes https://github.com/vmware-tanzu/community-edition/pull/4106#discussion_r860183492
	err = os.Chdir(tceRepoPath + "/cli/cmd/plugin/unmanaged-cluster")
	if err != nil {
		t.Errorf("error while changing directory to unmanged-cluster: %v", err)
	}

	err = cliRunner("go", nil, "run", ".", "delete", clusterName)
	if err != nil {
		t.Errorf("Error while Unmanaged Cluster deletion: %v", err)
	}
}

func runDeployScript(filename string, args ...string) error {
	mwriter := io.MultiWriter(os.Stdout)
	cmd := exec.Command(filename, args...)
	cmd.Stderr = mwriter
	cmd.Stdout = mwriter
	err := cmd.Run() // blocks until sub process is complete
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func cliRunner(name string, input io.Reader, args ...string) error {
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

		return fmt.Errorf("%s\nexit status: %d", err.Error(), rc)
	}
	return err
}
