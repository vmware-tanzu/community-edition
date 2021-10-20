// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/spf13/cobra"

	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/tanzu"
)

type initStandaloneOptions struct {
	clusterConfigFile string
	ui                bool
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
	CreateCmd.Flags().StringVarP(&iso.clusterConfigFile, "kind-config", "k", "", "Kind configuration file; fully overwrites Tanzu defaults")
	CreateCmd.Flags().StringVarP(&iso.clusterConfigFile, "port-forward", "p", "", "Port to forward from host to container")
	CreateCmd.Flags().BoolVar(&createOpts.tty, "tty", true, "Specify whether terminal is tty;\\nSet to false to disable styled output; default: true")
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

	tm := tanzu.New(clusterName)
	err := tm.Deploy("")
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	return nil
}
