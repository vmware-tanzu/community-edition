// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/cluster"
	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/log"
)

// ListCmd returns a list of existing clusters.
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "list tanzu clusters",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE:    list,
	Aliases: []string{"ls"},
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func init() {
	ListCmd.Flags().BoolVar(&co.tty, "tty", true, "Specify whether terminal is tty. Set to false to disable styled output.")
}

// list outputs a list of all local clusters on the system.
func list(cmd *cobra.Command, args []string) error {
	log := logger.NewLogger(TtySetting(cmd.Flags()), 0)

	clusterManager := cluster.NewKindClusterManager()
	clusters, err := clusterManager.List()
	if err != nil {
		log.Errorf("Unable to list clusters. Error: %s", err.Error())
	}

	// TODO(stmcginnis): Pull in table output formatting from tanzu-framework
	// and determine what else should be shown in addition to the name.
	log.Info("NAME\n")
	for _, c := range clusters {
		log.Infof("%s\n", c.Name)
	}

	return nil
}
