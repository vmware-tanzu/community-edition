// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package standalone

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu-private/tkg-cli/pkg/tkgctl"
)

// CreateCmd creates a standalone workload cluster.
var CreateCmd = &cobra.Command{
	Use:   "create <cluster name> -f <configuration location>",
	Short: "create a standalone workload cluster",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: create,
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func init() {
	fmt.Println(tkgctl.CreateClusterOptions{})
}

func create(cmd *cobra.Command, args []string) error {
	// validate a package name was passed
	if len(args) < 1 {
		return fmt.Errorf("no cluster name specified")
	}

	return nil
}
