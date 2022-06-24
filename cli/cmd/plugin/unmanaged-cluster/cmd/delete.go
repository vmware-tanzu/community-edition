// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"os"

	"github.com/spf13/cobra"

	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/tanzu"
)

const deleteDesc = `
Delete a Tanzu unmanaged cluster. This will attempt to destroy the running cluster
and remove the configuration stored in $HOME/.config/tanzu/tkg/unmanaged/${CLUSTER_NAME}.`

// DeleteCmd deletes an unmanaged workload cluster.
var DeleteCmd = &cobra.Command{
	Use:   "delete <cluster name>",
	Short: "Delete an unmanaged cluster",
	Long:  deleteDesc,
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE:    destroy,
	Aliases: []string{"rm"},
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	Args: checkSingleClusterArg,
}

func init() {
	DeleteCmd.Flags().Bool("tty-disable", false, "Disable log stylization and emojis")
}

func destroy(cmd *cobra.Command, args []string) error {
	// args have already been checked by DeleteCmd.Args()
	clusterName := args[0]

	log := logger.NewLogger(TtySetting(cmd.Flags()), LoggingVerbosity(cmd))

	log.Eventf(logger.TestTubeEmoji, "Deleting cluster: %s\n", clusterName)
	tClient := tanzu.New(log)
	err := tClient.Delete(clusterName)
	if err != nil {
		log.Errorf("Failed delete operation. Error: %s\n", err.Error())
		os.Exit(1)
	}

	log.Eventf(logger.TestTubeEmoji, "Deleted cluster: %s\n", clusterName)

	return nil
}
