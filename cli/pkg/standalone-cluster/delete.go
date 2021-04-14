// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package standalone

import (
	"fmt"
	"os/user"

	"github.com/vmware-tanzu-private/tkg-cli/pkg/tkgctl"
	"github.com/vmware-tanzu-private/tkg-cli/pkg/types"

	"github.com/spf13/cobra"
)

type teardownStandaloneOptions struct {
	force bool
	skip  bool
}

// DeleteCmd deletes a standalone workload cluster.
var DeleteCmd = &cobra.Command{
	Use:   "delete <cluster name> -f <configuration location>",
	Short: "delete a standalone workload cluster",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: teardown,
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

var tso = teardownStandaloneOptions{}

func init() {
	DeleteCmd.Flags().BoolVarP(&tso.force, "force", "f", false, "Force delete")
	DeleteCmd.Flags().BoolVarP(&tso.skip, "skip", "s", false, "Skip user deletion prompt")
}

func teardown(cmd *cobra.Command, args []string) error {
	// validate a cluster name was passed
	if len(args) < 1 {
		return fmt.Errorf("no cluster name specified")
	}
	clusterName := args[0]

	fmt.Println(tkgctl.CreateClusterOptions{})

	usr, err := user.Current()
	if err != nil {
		return err
	}

	// setup client options
	opt := tkgctl.Options{
		KubeConfig:        "",
		KubeContext:       "",
		ConfigDir:         usr.HomeDir + "/.tanzu",
		LogOptions:        tkgctl.LoggingOptions{Verbosity: 10},
		ProviderGetter:    nil,
		CustomizerOptions: types.CustomizerOptions{},
		SettingsFile:      "",
	}

	// create new client
	c, err := tkgctl.New(opt)
	if err != nil {
		fmt.Println(err.Error())
	}

	// delete a new standlone cluster
	teardownRegionOpts := tkgctl.DeleteRegionOptions{
		ClusterName: clusterName,
		Force:       tso.force,
		SkipPrompt:  tso.skip,
	}

	err = c.DeleteStandalone(teardownRegionOpts)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}
