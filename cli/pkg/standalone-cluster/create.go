// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package standalone

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu-private/tkg-cli/pkg/tkgctl"
	"github.com/vmware-tanzu-private/tkg-cli/pkg/types"

	"github.com/vmware-tanzu/tce/cli/utils"
)

type initStandaloneOptions struct {
	clusterConfigFile      string
	infrastructureProvider string
}

// CreateCmd creates a standalone workload cluster.
var CreateCmd = &cobra.Command{
	Use:   "create <cluster name> -f <configuration location>",
	Short: "create a standalone workload cluster",
	RunE:  create,
}

var iso = initStandaloneOptions{}

func init() {
	CreateCmd.Flags().StringVarP(&iso.clusterConfigFile, "file", "f", "", "Configuration file from which to create a standalone cluster")

	CreateCmd.Flags().StringVarP(&iso.infrastructureProvider, "infrastructure", "i", "", "Infrastructure to deploy the standalone cluster on ['aws', 'vsphere', 'docker']")
	CreateCmd.Flags().MarkHidden("infrastructure") //nolint
}

func create(cmd *cobra.Command, args []string) error {
	// validate a package name was passed
	if len(args) < 1 {
		return fmt.Errorf("no cluster name specified")
	}

	fmt.Println(tkgctl.CreateClusterOptions{})

	homedir, err := os.UserHomeDir()
	if err != nil {
		return utils.NonUsageError(cmd, err, "unable to determine user home directory.")
	}

	// setup client options
	opt := tkgctl.Options{
		KubeConfig:        "",
		KubeContext:       "",
		ConfigDir:         fmt.Sprintf("%s/.tanzu", homedir),
		LogOptions:        tkgctl.LoggingOptions{Verbosity: 10},
		ProviderGetter:    nil,
		CustomizerOptions: types.CustomizerOptions{},
		SettingsFile:      "",
	}

	// create new client
	c, err := tkgctl.New(opt)
	if err != nil {
		return utils.NonUsageError(cmd, err, "unable to create Tanzu management client.")
	}

	// create a new standlone cluster
	initRegionOpts := tkgctl.InitRegionOptions{
		ClusterConfigFile: iso.clusterConfigFile,
	}

	if iso.infrastructureProvider != "" {
		initRegionOpts.InfrastructureProvider = iso.infrastructureProvider
	}

	err = c.InitStandalone(initRegionOpts)
	if err != nil {
		return utils.NonUsageError(cmd, err, "failed to initialize standalone cluster.")
	}

	return nil
}
