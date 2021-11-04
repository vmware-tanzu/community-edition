// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/tanzu"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/component"
)

const listDesc = `
List known standalone clusters. This list is produced by locating clusters saved to
$HOME/.config/tanzu/tkg/standalone`

// ListCmd returns a list of existing clusters.
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List standalone environments",
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
	ListCmd.Flags().BoolVar(&co.tty, "tty", true, "Specify whether terminal is tty. Set to false to disable styled output.")
}

// list outputs a list of all standalone clusters on the system.
func list(cmd *cobra.Command, args []string) error {
	log := logger.NewLogger(TtySetting(cmd.Flags()), 0)
	tClient := tanzu.New(log)
	clusters, err := tClient.List()
	if err != nil {
		return fmt.Errorf("unable to list clusters. Error: %s", err.Error())
	}

	t := component.NewOutputWriter(cmd.OutOrStdout(), "table", "NAME", "PROVIDER")
	for _, c := range clusters {
		t.AddRow(c.Name, c.Provider)
	}
	t.Render()

	return nil
}
