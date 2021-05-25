// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package standalone

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu-private/core/pkg/v1/client"
	"github.com/vmware-tanzu-private/tkg-cli/pkg/tkgctl"
	"github.com/vmware-tanzu-private/tkg-cli/pkg/types"

	"github.com/vmware-tanzu/tce/cli/pkg/utils"
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
	Args:  cobra.ExactArgs(1),
}

var iso = initStandaloneOptions{}

const defaultPlan = "dev"

func init() {
	CreateCmd.Flags().StringVarP(&iso.clusterConfigFile, "file", "f", "", "Configuration file from which to create a standalone cluster")
	CreateCmd.Flags().StringVarP(&iso.infrastructureProvider, "infrastructure", "i", "", "Infrastructure to deploy the standalone cluster on. Only needed when using -i docker.")
}

func create(cmd *cobra.Command, args []string) error {
	// validate a cluster name was passed
	if len(args) < 1 {
		return fmt.Errorf("no cluster name specified")
	}
	clusterName := args[0]
	cmd.Println(tkgctl.CreateClusterOptions{})

	configDir, err := client.LocalDir()
	if err != nil {
		return utils.NonUsageError(cmd, err, "unable to determine Tanzu configuration directory.")
	}

	// setup client options
	opt := tkgctl.Options{
		KubeConfig:        "",
		KubeContext:       "",
		ConfigDir:         configDir,
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

	// when using CAPD, set default plan to "dev"
	var plan string
	if iso.infrastructureProvider == "docker" {
		plan = defaultPlan
	}

	// create a new standlone cluster
	initRegionOpts := tkgctl.InitRegionOptions{
		ClusterConfigFile: iso.clusterConfigFile,
		ClusterName:       clusterName,
		Plan:              plan,
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
