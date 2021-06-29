// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package standalone

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu-private/core/pkg/v1/client"
	"github.com/vmware-tanzu-private/tkg-cli/pkg/constants"
	"github.com/vmware-tanzu-private/tkg-cli/pkg/log"
	"github.com/vmware-tanzu-private/tkg-cli/pkg/tkgctl"
	"github.com/vmware-tanzu-private/tkg-cli/pkg/types"

	"github.com/vmware-tanzu/tce/cli/pkg/utils"
)

type initStandaloneOptions struct {
	clusterConfigFile      string
	infrastructureProvider string
	ui                     bool
	bind                   string
	browser                string
}

// CreateCmd creates a standalone workload cluster.
var CreateCmd = &cobra.Command{
	Use:   "create <cluster name> -f <configuration location>",
	Short: "create a standalone workload cluster",
	RunE:  create,
}

var iso = initStandaloneOptions{}

const defaultPlan = "dev"

func init() {
	CreateCmd.Flags().StringVarP(&iso.clusterConfigFile, "file", "f", "", "Configuration file from which to create a standalone cluster")
	CreateCmd.Flags().BoolVarP(&iso.ui, "ui", "u", false, "Launch interactive standalone cluster provisioning UI")
	CreateCmd.Flags().StringVarP(&iso.infrastructureProvider, "infrastructure", "i", "", "Infrastructure to deploy the standalone cluster on. Only needed when using -i docker.")
	CreateCmd.Flags().StringVarP(&iso.bind, "bind", "b", "127.0.0.1:8080", "Specify the IP and port to bind the Kickstart UI against (e.g. 127.0.0.1:8080).")
	CreateCmd.Flags().StringVarP(&iso.browser, "browser", "", "", "Specify the browser to open the Kickstart UI on. Use 'none' for no browser. Defaults to OS default browser. Supported: ['chrome', 'firefox', 'safari', 'ie', 'edge', 'none']")
}

func create(cmd *cobra.Command, args []string) error {
	var clusterName string

	// validate a cluster name was passed when not using the kickstart UI
	if len(args) < 1 && !iso.ui {
		return fmt.Errorf("no cluster name specified")
	} else if len(args) == 1 {
		clusterName = args[0]
	}

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

	// create a new standlone cluster
	initRegionOpts := tkgctl.InitRegionOptions{
		ClusterConfigFile: iso.clusterConfigFile,
		ClusterName:       clusterName,
		Plan:              defaultPlan,
		UI:                iso.ui,
		Bind:              iso.bind,
		Browser:           iso.browser,
	}

	if iso.infrastructureProvider != "" {
		initRegionOpts.InfrastructureProvider = iso.infrastructureProvider
	}

	err = c.InitStandalone(initRegionOpts)
	if err != nil {
		return utils.Error(err, "failed to initialize standalone cluster.")
	}

	err = saveStandaloneClusterConfig(clusterName, iso.clusterConfigFile)
	if err != nil {
		return utils.Error(err, "failed to store standalone bootstrap cluster config")
	}

	return nil
}

func saveStandaloneClusterConfig(clusterName, clusterConfigPath string) error {
	// If there is no cluster config provided
	// assume CAPD and don't try to save anything
	if clusterConfigPath == "" {
		return nil
	}

	// Get config contents
	clusterConfigBytes, err := os.ReadFile(clusterConfigPath)
	if err != nil {
		return fmt.Errorf("cannot read cluster config file: %v", err)
	}

	// get the user homedir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Save the cluster configuration for future restore cycle
	configDir := filepath.Join(homeDir, ".config", "tanzu", "tce", "configs")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return err
	}

	clusterConfigFile := clusterName + "_ClusterConfig"
	writeConfigPath := filepath.Join(configDir, clusterConfigFile)

	log.Infof("Saving bootstrap cluster config for standalone cluster at '%v'", writeConfigPath)
	err = os.WriteFile(writeConfigPath, clusterConfigBytes, constants.ConfigFilePermissions)
	if err != nil {
		return fmt.Errorf("cannot write cluster config file for standalone bootstrap cluster: %v", err)
	}

	return nil
}
