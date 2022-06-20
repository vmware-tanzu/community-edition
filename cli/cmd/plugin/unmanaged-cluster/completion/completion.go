// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package completion

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/cluster"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/cmd"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/config"
	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/tanzu"
)

const clusterAnyStatus = ""

func Setup(rootCmd *cobra.Command) {
	// The below commands accept the names of existing clusters
	cmd.DeleteCmd.ValidArgsFunction = completeExistingClusterName(clusterAnyStatus)
	cmd.StartCmd.ValidArgsFunction = completeExistingClusterName(cluster.StatusStopped)
	cmd.StopCmd.ValidArgsFunction = completeExistingClusterName(cluster.StatusRunning)

	// The below commands should not provide any shell completion choices nor file names
	cmd.ConfigureCmd.ValidArgsFunction = cobra.NoFileCompletions
	cmd.CreateCmd.ValidArgsFunction = cobra.NoFileCompletions
	cmd.ListCmd.ValidArgsFunction = cobra.NoFileCompletions
}

// completeExistingClusterName returns a completion function that lists all clusters that match
// the specified clusterStatus.  If clusterStatus is clusterAnyStatus, all clusters will match.
func completeExistingClusterName(clusterStatus string) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			// No more args accepted
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		log := logger.NewLogger(false, 0) // No need for logs
		tClient := tanzu.New(log)

		// Find clusters and their status
		clusters, err := tClient.List()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var clusterNames []string
		for _, c := range clusters {
			// retrieve status from provider
			conf := &config.UnmanagedClusterConfig{
				Provider: c.Provider,
			}
			clusterManager := cluster.NewClusterManager(conf)
			kc, err := clusterManager.Get(c.Name)

			if err != nil {
				continue
			}

			// Only complete clusters of the right status, or if the status is unknown.
			// A value of clusterAnyStatus indicates a request for all clusters.
			if kc.Status == clusterStatus || kc.Status == cluster.StatusUnknown || clusterStatus == clusterAnyStatus {
				clusterNames = append(clusterNames, fmt.Sprintf("%s\tstatus: %s - provider: %s", c.Name, kc.Status, c.Provider))
			}
		}
		return clusterNames, cobra.ShellCompDirectiveNoFileComp
	}
}
