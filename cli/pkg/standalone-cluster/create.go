// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package standalone

import (
	"fmt"
	"github.com/vmware-tanzu-private/tkg-cli/pkg/types"
	"os/user"

	"github.com/vmware-tanzu-private/tkg-cli/pkg/tkgctl"

	"github.com/spf13/cobra"
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
}

func create(cmd *cobra.Command, args []string) error {
	// validate a package name was passed
	if len(args) < 1 {
		return fmt.Errorf("no cluster name specified")
	}

	fmt.Println(tkgctl.CreateClusterOptions{})

	usr, err := user.Current()
	if err != nil {
		return err
	}

	// setup client options
	opt := tkgctl.Options{
		KubeConfig:        "",
		KubeContext:       "",
		ConfigDir:         usr.HomeDir+"/.tanzu",
		LogOptions:        tkgctl.LoggingOptions{},
		ProviderGetter:    nil,
		CustomizerOptions: types.CustomizerOptions{},
		SettingsFile:      "",
	}

	// create new client
	c, err := tkgctl.New(opt)
	if err != nil {
		fmt.Println(err.Error())
	}

	// create a new management cluster
	initRegionOpts := tkgctl.InitRegionOptions{
		ClusterConfigFile:           "/home/josh/d/confdir/test.yaml",
	}
	err = c.InitStandalone(initRegionOpts)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}
