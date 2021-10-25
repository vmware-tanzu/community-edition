// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/config"
	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/tanzu"
)

type createLocalOpts struct {
	clusterConfigFile      string
	infrastructureProvider string
	ui                     bool
	tty                    bool
}

// CreateCmd creates a standalone workload cluster.
var CreateCmd = &cobra.Command{
	Use:   "create <cluster name> -f <configuration location>",
	Short: "create a local tanzu cluster",
	RunE:  create,
}

var co = createLocalOpts{}

func init() {
	CreateCmd.Flags().StringVarP(&co.clusterConfigFile, "config", "f", "", "Configuration file for local cluster creation")
	CreateCmd.Flags().StringVarP(&co.infrastructureProvider, "provider", "p", "", "The infrastructure provider to use for cluster creation. Default is 'kind'")
	CreateCmd.Flags().BoolVar(&co.tty, "tty", true, "Specify whether terminal is tty. Set to false to disable styled output; default: true")
}

func create(cmd *cobra.Command, args []string) error {
	var clusterName string

	// validate a cluster name was passed when not using the kickstart UI
	if len(args) < 1 && !co.ui {
		return fmt.Errorf("cluster name not specified")
	} else if len(args) == 1 {
		clusterName = args[0]
	}
	// initial logger, needed for logging if something goes wrong
	log := logger.NewLogger(false, 0)

	// Determine our configuration to use
	configArgs := map[string]string{
		"clusterconfigfile": co.clusterConfigFile,
		"clustername":       clusterName,
	}
	clusterConfig, err := config.InitializeConfiguration(configArgs)
	if err != nil {
		log.Errorf("Failed to initialize configuration. Error %s\n", clusterConfig)
		return nil
	}
	ttySetting, err := strconv.ParseBool(clusterConfig.Tty)
	if err != nil {
		log.Errorf("TTY setting was invalid. Error: %s", err.Error())
		return nil
	}
	// reset logger here based on parsed configuration
	log = logger.NewLogger(ttySetting, 0)

	tm := tanzu.New(log)
	err = tm.Deploy(clusterConfig)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	return nil
}
