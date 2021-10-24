// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"github.com/spf13/cobra"

	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/tanzu"
)

const (
	configDir             = ".config"
	tanzuConfigDir        = "tanzu"
	tkgConfigDir          = "tkg"
	tkgSysNamespace       = "tkg-system"
	tkgSvcAcctName        = "core-pkgs"
	tkgCoreRepoName       = "tkg-core-repository"
	tkgGlobalPkgNamespace = "tanzu-package-repo-global"
	tceRepoName           = "community-repository"
	tceRepoUrl            = "projects.registry.vmware.com/tce/main:0.9.1"
)

type initStandaloneOptions struct {
	clusterConfigFile      string
	infrastructureProvider string
	ui                     bool
	bind                   string
	browser                string
	timeout                time.Duration
}

type createOptions struct {
	tty bool
}

// CreateCmd creates a standalone workload cluster.
var CreateCmd = &cobra.Command{
	Use:   "create <cluster name> -f <configuration location>",
	Short: "create a local tanzu cluster",
	RunE:  create,
}

var iso = initStandaloneOptions{}
var createOpts = createOptions{}

func init() {
	CreateCmd.Flags().StringVarP(&iso.clusterConfigFile, "config", "f", "", "Configuration file for local cluster creation")
	CreateCmd.Flags().StringVarP(&iso.infrastructureProvider, "provider", "p", "", "The infrastructure provider to use for cluster creation. Default is 'kind'")
	CreateCmd.Flags().BoolVar(&createOpts.tty, "tty", true, "Specify whether terminal is tty;\\nSet to false to disable styled ouput; default: true")
}

func create(cmd *cobra.Command, args []string) error {
	var clusterName string

	// validate a cluster name was passed when not using the kickstart UI
	if len(args) < 1 && !iso.ui {
		return fmt.Errorf("Cluster name not specified.")
	} else if len(args) == 1 {
		clusterName = args[0]
	}
	log := logger.NewLogger(createOpts.tty, 0)

	// Determine our configuration to use
	configArgs := map[string]string{
		"clusterconfigfile": iso.clusterConfigFile,
		"clustername":       clusterName,
	}
	clusterConfig, err := tanzu.InitializeConfiguration(configArgs)
	if err != nil {
		log.Errorf("Failed to initialize configuration. Error %s\n", clusterConfig)
		return nil
	}

	tm := tanzu.New(log)
	err = tm.Deploy(clusterConfig)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	return nil
    }
