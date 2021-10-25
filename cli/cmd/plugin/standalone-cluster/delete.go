// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/tanzu"

	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/log"
)

// DeleteCmd deletes a standalone workload cluster.
var DeleteCmd = &cobra.Command{
	Use:   "delete <cluster name>",
	Short: "delete a local tanzu cluster",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE:    destroy,
	Aliases: []string{"del", "rm"},
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func destroy(cmd *cobra.Command, args []string) error {
	var clusterName string

	// validate a cluster name was passed
	if len(args) < 1 {
		return fmt.Errorf("must specify cluster name to delete")
	} else if len(args) == 1 {
		clusterName = args[0]
	}
	log := logger.NewLogger(true, 0)

	log.Eventf("\\U+1F5D1", " Deleting cluster: %s\n", clusterName)
	tClient := tanzu.New(log)
	err := tClient.Delete(clusterName)
	if err != nil {
		log.Errorf("Failed to delete cluster. Error: %s\n", err.Error())
		return nil
	}

	log.Eventf("\\U+1F5D1", " Deleted cluster: %s\n", clusterName)

	return nil
}
