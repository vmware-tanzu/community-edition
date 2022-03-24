// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/internal/hack"
	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/tanzu"
)

const listDesc = `
List known unmanaged clusters. This list is produced by locating clusters saved to
$HOME/.config/tanzu/tkg/unmanaged`

type listUnmanagedOptions struct {
	outputFormat string
	quiet        bool
}

var lo = listUnmanagedOptions{}

// ListCmd returns a list of existing clusters.
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List unmanaged environments",
	Long:  listDesc,
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
	ListCmd.Flags().Bool("tty-disable", false, "Disable log stylization and emojis")
	ListCmd.Flags().StringVarP(&lo.outputFormat, "output", "o", "table", "Output format (yaml|json|table)")
	ListCmd.Flags().BoolVarP(&lo.quiet, "quiet", "q", false, "Only display cluster names")
}

// list outputs a list of all unmanaged clusters on the system.
func list(cmd *cobra.Command, args []string) error {
	log := logger.NewLogger(TtySetting(cmd.Flags()), 0)
	tClient := tanzu.New(log)
	clusters, err := tClient.List()
	if err != nil {
		return fmt.Errorf("unable to list clusters. Error: %s\n", err.Error()) //nolint:revive,stylecheck
	}

	if lo.quiet {
		printQuiet(clusters)
		return nil
	}

	t := hack.NewOutputWriter(cmd.OutOrStdout(), lo.outputFormat, "NAME", "PROVIDER")
	for _, c := range clusters {
		t.AddRow(c.Name, c.Provider)
	}
	t.Render()

	return nil
}

func printQuiet(clusters []tanzu.Cluster) {
	for _, c := range clusters {
		fmt.Printf("%s\n", c.Name)
	}
}
