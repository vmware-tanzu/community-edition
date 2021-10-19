// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"github.com/spf13/cobra"
	kindCluster "sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cluster/nodes"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/cluster"
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
	// TODO(joshrosso)
}

func list(cmd *cobra.Command, args []string) error {
	clusterManager := cluster.NewClusterManager()
	clusters, err := clusterManager.List()
	if err != nil {
		fmt.Printf("Unable to list clusters. Error: %s", err.Error())
	}

	// TODO(stmcginnis): Pull in table output formatting from tanzu-framework.
	for _, c := range clusters {
		fmt.Println(c.Name)
	}

	return nil
}

// ListNodes returns a list of nodes for the current cluster.
// If the cluster doesn't exist, an empty list is returned.
func ListNodes(clusterName string) []nodes.Node {
	// TODO(stmcginnis): Remove kind interaction and abstract node information.
	provider := kindCluster.NewProvider()
	nodes := []nodes.Node{}
	nodes, _ = provider.ListNodes(clusterName)
	return nodes
}
