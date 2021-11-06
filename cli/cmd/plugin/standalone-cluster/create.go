// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/config"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/constants"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/region"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgctl"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/types"
)

type initStandaloneOptions struct {
	clusterConfigFile      string
	infrastructureProvider string
	ui                     bool
	bind                   string
	browser                string
	timeout                time.Duration
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
	CreateCmd.Flags().DurationVarP(&iso.timeout, "timeout", "t", constants.DefaultLongRunningOperationTimeout, "Time duration to wait for an operation before timeout. Timeout duration in hours(h)/minutes(m)/seconds(s) units or as some combination of them (e.g. 2h, 30m, 2h30m10s)")
}

func create(cmd *cobra.Command, args []string) error {
	var clusterName string

	// validate a cluster name was passed when not using the kickstart UI
	if len(args) < 1 && !iso.ui {
		return fmt.Errorf("no cluster name specified")
	} else if len(args) == 1 {
		clusterName = args[0]
	}

	// create new client
	c, err := newTKGCtlClient(false)
	if err != nil {
		return NonUsageError(cmd, err, "unable to create Tanzu Standalone Cluster client")
	}

	// create a new standlone cluster
	initRegionOpts := tkgctl.InitRegionOptions{
		ClusterConfigFile: iso.clusterConfigFile,
		ClusterName:       clusterName,
		Plan:              defaultPlan,
		UI:                iso.ui,
		Bind:              iso.bind,
		Browser:           iso.browser,
		Edition:           BuildEdition,
		Timeout:           iso.timeout,
		// all tce-based clusters should opt out of CEIP
		// since standalone-clusters are specific to TCE, we'll
		// always set this to "false"
		CeipOptIn: "false",
	}

	if iso.infrastructureProvider != "" {
		initRegionOpts.InfrastructureProvider = iso.infrastructureProvider
	}

	err = c.InitStandalone(initRegionOpts)
	if err != nil {
		return Error(err, "failed to initialize standalone cluster.")
	}

	err = saveStandaloneClusterConfig(clusterName, iso.clusterConfigFile)
	if err != nil {
		return Error(err, "failed to store standalone bootstrap cluster config")
	}

	return nil
}

func newTKGCtlClient(forceUpdateTKGCompatibilityImage bool) (tkgctl.TKGClient, error) {
	configDir, err := getTKGConfigDir()
	if err != nil {
		return nil, Error(err, "unable to determine Tanzu configuration directory.")
	}

	return tkgctl.New(tkgctl.Options{
		ConfigDir: configDir,
		CustomizerOptions: types.CustomizerOptions{
			RegionManagerFactory: region.NewFactory(),
		},

		LogOptions:                       tkgctl.LoggingOptions{Verbosity: logLevel, File: logFile},
		ForceUpdateTKGCompatibilityImage: forceUpdateTKGCompatibilityImage,
	})
}

func getTKGConfigDir() (string, error) {
	tanzuConfigDir, err := config.LocalDir()
	if err != nil {
		return "", Error(err, "unable to get home directory")
	}
	return filepath.Join(tanzuConfigDir, "tkg"), nil
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

	// Save the cluster configuration for future restore cycle
	configDir, err := getTKGConfigDir()
	if err != nil {
		return err
	}

	clusterConfigDir := filepath.Join(configDir, "clusterconfigs")
	err = os.MkdirAll(clusterConfigDir, 0755)
	if err != nil {
		return err
	}

	clusterConfigFile := clusterName + ".yaml"
	writeConfigPath := filepath.Join(clusterConfigDir, clusterConfigFile)

	log.Infof("Saving bootstrap cluster config for standalone cluster at '%v'", writeConfigPath)
	err = os.WriteFile(writeConfigPath, clusterConfigBytes, constants.ConfigFilePermissions)
	if err != nil {
		return fmt.Errorf("cannot write cluster config file for standalone bootstrap cluster: %v", err)
	}

	return nil
}
