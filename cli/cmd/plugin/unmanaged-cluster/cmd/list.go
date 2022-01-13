// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/tanzu"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/component"
)

const listDesc = `
List known unmanaged clusters. This list is produced by locating clusters saved to
$HOME/.config/tanzu/tkg/unmanaged`

var outputFormat = "table"

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
	ListCmd.Flags().BoolVar(&co.tty, "tty", true, "Specify whether terminal is tty; set to false to disable styled output.")
	ListCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format (yaml|json|table); default is table")
}

// list outputs a list of all unmanaged clusters on the system.
func list(cmd *cobra.Command, args []string) error {
	log := logger.NewLogger(TtySetting(cmd.Flags()), 0)
	tClient := tanzu.New(log)
	clusters, err := tClient.List()
	if err != nil {
		return fmt.Errorf("unable to list clusters. Error: %s\n", err.Error()) //nolint:revive,stylecheck
	}

	t := component.NewOutputWriter(cmd.OutOrStdout(), outputFormat, "NAME", "PROVIDER")
	for _, c := range clusters {
		t.AddRow(c.Name, c.Provider)
	}
	t.Render()
	log.Info("")

	return nil
}
