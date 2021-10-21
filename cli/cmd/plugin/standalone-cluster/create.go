// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"strings"

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
	CreateCmd.Flags().StringVarP(&iso.infrastructureProvider, "infraprovider", "i", "", "The infrastructure provider to use for cluster creation. Default is 'kind'")
	CreateCmd.Flags().StringVarP(&iso.clusterConfigFile, "kind-config", "k", "", "Kind configuration file; fully overwrites Tanzu defaults")
	CreateCmd.Flags().StringVarP(&iso.clusterConfigFile, "port-forward", "p", "", "Port to forward from host to container")
	CreateCmd.Flags().BoolVar(&createOpts.tty, "tty", true, "Specify whether terminal is tty;\\nSet to false to disable styled ouput; default: true")
}

func create(cmd *cobra.Command, args []string) error {
	var clusterName string

	// validate a cluster name was passed when not using the kickstart UI
	if len(args) < 1 && !iso.ui {
		return Error(nil, "no cluster name specified")
	} else if len(args) == 1 {
		clusterName = args[0]
	}
	log := logger.NewLogger(createOpts.tty, 0)

	// Read in the configuration we should use
	clusterConfig, err := initializeConfiguration(clusterName, &iso)
	if err != nil {
		return Error(err, "error processing configuration")
	}

	tm := tanzu.New(log)
	err = tm.Deploy(clusterConfig)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	return nil
}

// initializeConfiguration determines the configuration to use for cluster creation.
//
// There are three places where configuration comes from:
// - configuration file
// - environment variables
// - command line arguments
//
// The effective configuration is determined by combining these sources, in ascending
// order of preference listed. So env variables override values in the config file,
// and explicit CLI arguments override config file and env variable values.
func initializeConfiguration(clusterName string, iso *initStandaloneOptions) (*tanzu.LocalClusterConfig, error) {
	// TODO(stmcginnis): handle loading values from iso.clusterConfigFile
	config := &tanzu.LocalClusterConfig{ClusterName: clusterName}

	// Check what provider to use for creating cluster
	config.Provider = strings.ToLower(os.Getenv("LOCAL_INFRA_PROVIDER"))
	if iso.infrastructureProvider != "" {
		config.Provider = strings.ToLower(iso.infrastructureProvider)
	}
	if config.Provider == "" {
		config.Provider = "kind"
	}
	config.Provider = iso.infrastructureProvider

	return config, Error(nil, "placeholder")
}
