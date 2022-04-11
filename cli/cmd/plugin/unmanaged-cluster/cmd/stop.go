// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd // nolint:dupl

import (
	"fmt"

	"github.com/spf13/cobra"

	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/tanzu"
)

const stopDesc = `
Stop a Tanzu unmanaged cluster. If supported, you can start it at a later point
in time.`

// StopCmd stops a running cluster.
var StopCmd = &cobra.Command{
	Use:   "stop <cluster name>",
	Short: "(experimental) Stop an unmanaged cluster",
	Long:  stopDesc,
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: stop,
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func init() {
	StopCmd.Flags().Bool("tty-disable", false, "Disable log stylization and emojis")
}

func stop(cmd *cobra.Command, args []string) error { // nolint:dupl
	var clusterName string

	// validate a cluster name was passed
	if len(args) < 1 {
		return fmt.Errorf("must specify cluster name to stop")
	} else if len(args) == 1 {
		clusterName = args[0]
	}
	log := logger.NewLogger(TtySetting(cmd.Flags()), 0)

	log.Eventf(logger.TestTubeEmoji, "Attempting to stop cluster: %s\n", clusterName)
	tClient := tanzu.New(log)
	err := tClient.Stop(clusterName)
	if err != nil {
		log.Errorf("Failed to stop cluster. Error: %s\n", err.Error())
		return nil
	}

	log.Eventf(logger.TestTubeEmoji, "Stopped cluster: %s\n", clusterName)

	return nil
}
