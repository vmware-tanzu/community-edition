//go:build e2e

// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"testing"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/config"
	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/tanzu"
	e2e "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/test/e2e/utils"
)

var clusterName string

const (
	colorReset = "\033[0m" // Reset
	colorBlue  = "\033[34m"
	colorRed   = "\033[31m" // Fail
	colorGreen = "\033[32m" // Pass
	passConst  = "GREEN"
	failConst  = "RED"
)

// installing TCE
func TestTCEInstallation(t *testing.T) {
	signal := passConst
	fmt.Println("-------------------------------", string(colorBlue), "Unmanged cluster e2e Test", string(colorReset), "---------------------------------------------")
	err := e2e.InstallTCE()
	if err != nil {
		t.Errorf("Error while installing TCE: %v", err)
		signal = failConst
	}
	if signal == failConst {
		fmt.Println("-------------------------------", string(colorRed), "TCE installation Failed", string(colorReset), "---------------------------------------------")
	} else {
		fmt.Println("-------------------------------", string(colorGreen), "TCE installation Passed", string(colorReset), "---------------------------------------------")
	}
}

// Installing Unmanahged cluster
func TestUCInstallation(t *testing.T) {
	signal := passConst

	// Random Unmanged Cluster Name
	rand, _ := rand.Int(rand.Reader, big.NewInt(1000))
	clusterName = "uc" + rand.String() + "test"

	ucConfig := map[string]interface{}{config.ClusterName: clusterName}
	clusterConfig, _ := config.InitializeConfiguration(ucConfig)
	log := logger.NewLogger(false, 0)
	tm := tanzu.New(log)
	_, err := tm.Deploy(clusterConfig)
	if err != nil {
		signal = failConst
		t.Errorf("error while Unmanaged Cluster creation: %v", err)
	}

	if signal == failConst {
		fmt.Println("-------------------------------", string(colorRed), "Unmanged Cluster creation Failed see above logs", string(colorReset), "---------------------------------------------")
	} else {
		fmt.Println("-------------------------------", string(colorGreen), "Unmanged Cluster created Successfully", string(colorReset), "---------------------------------------------")
	}
}

// Checking Unmanaged cluster working
func TestUCWorking(t *testing.T) {
	signal := passConst

	repoList, err := e2e.Tanzu(nil, "package", "repository", "list", "-A")
	if err != nil || repoList == "" {
		t.Errorf("error while checking for package repositories: %v", err)
		signal = failConst
	}

	registryExist, err := regexp.MatchString("\\btanzu-package-repo-global\\b", repoList)
	if registryExist == false || err != nil {
		t.Errorf("Package registry not present or %v", err)
		signal = failConst
	}

	coreRepo, err := regexp.MatchString("\\btkg-system\\b", repoList)
	if coreRepo == false || err != nil {
		t.Errorf("Core repository not present or %v", err)
		signal = failConst
	}

	_, err = e2e.Kubectl(nil, "get", "pods", "-A")
	if err != nil {
		t.Errorf("error while checking for pods: %v", err)
		signal = failConst
	}

	if signal == failConst {
		fmt.Println("-------------------------------", string(colorRed), "Unmanaged Cluster is not healthy see above logs", string(colorReset), "---------------------------------------------")
	} else {
		fmt.Println("-------------------------------", string(colorGreen), "Unmanaged Cluster is healthy", string(colorReset), "---------------------------------------------")
	}
}

// Deleting unmanage cluster
func TestUCDeletion(t *testing.T) {
	signal := passConst

	log := logger.NewLogger(false, 0)
	tm := tanzu.New(log)
	err := tm.Delete(clusterName)
	if err != nil {
		t.Errorf("error while unmanaged cluster deletion")
		signal = failConst
	}

	ucLists, err := tm.List()
	if err != nil {
		t.Errorf("error while fetching unmanaged cluster list: %v", err)
		signal = failConst
	}

	if e2e.ContainsUC(ucLists, clusterName) {
		t.Errorf("error while unmanaged cluster deletion: %v", err)
		signal = failConst
	}

	if signal == failConst {
		fmt.Println("-------------------------------", string(colorRed), "Unmanaged Cluster deletion Failed see above logs", string(colorReset), "---------------------------------------------")
	} else {
		fmt.Println("-------------------------------", string(colorGreen), "Unmanaged Cluster deleted Successfully", string(colorReset), "---------------------------------------------")
	}
}
