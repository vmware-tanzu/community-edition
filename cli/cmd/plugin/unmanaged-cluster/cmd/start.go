// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd // nolint:dupl

import (
	"fmt"

	"github.com/spf13/cobra"

	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/tanzu"
)

const startDesc = `
Start an existing Tanzu unmanaged cluster.`

// StartCmd starts a stopped cluster.
var StartCmd = &cobra.Command{
	Use:   "start <cluster name>",
	Short: "(experimental) Start an unmanaged cluster",
	Long:  startDesc,
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: start,
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func init() {
	StartCmd.Flags().Bool("tty-disable", false, "Disable log stylization and emojis")
}

func start(cmd *cobra.Command, args []string) error {
	var clusterName string

	// validate a cluster name was passed
	if len(args) < 1 {
		return fmt.Errorf("must specify cluster name to start")
	} else if len(args) == 1 {
		clusterName = args[0]
	}
	log := logger.NewLogger(TtySetting(cmd.Flags()), LoggingVerbosity(cmd))

	log.Eventf(logger.TestTubeEmoji, "Attempting to start cluster: %s\n", clusterName)
	tClient := tanzu.New(log)
	err := tClient.Start(clusterName)
	if err != nil {
		log.Errorf("Failed to start cluster. Error: %s\n", err.Error())
		return nil
	}

	log.Eventf(logger.TestTubeEmoji, "Started cluster: %s\n", clusterName)

	return nil
}
