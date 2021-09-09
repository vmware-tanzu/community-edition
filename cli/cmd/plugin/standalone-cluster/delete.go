// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgctl"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/types"
)

type teardownStandaloneOptions struct {
	force      bool
	skip       bool
	configFile string
}

// DeleteCmd deletes a standalone workload cluster.
var DeleteCmd = &cobra.Command{
	Use:   "delete <cluster name>",
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
	DeleteCmd.Flags().StringVarP(&tso.configFile, "config", "f", "", "Optional cluster configuration file. Defaults to config used during standalone-cluster create")
	DeleteCmd.Flags().BoolVar(&tso.force, "force", false, "Force delete")
	DeleteCmd.Flags().BoolVarP(&tso.skip, "yes", "y", false, "Delete workload cluster without asking for confirmation")
}

func teardown(cmd *cobra.Command, args []string) error {
	// validate a cluster name was passed
	if len(args) < 1 {
		return fmt.Errorf("no cluster name specified")
	}
	clusterName := args[0]

	// Check for TKG_HTTPS_PROXY, TKG_HTTP_PROXY, TKG_NO_PROXY env variables
	// Set to empty if not present
	_, ok := os.LookupEnv("TKG_NO_PROXY")
	if !ok {
		os.Setenv("TKG_NO_PROXY", "")
	}

	_, ok = os.LookupEnv("TKG_HTTP_PROXY")
	if !ok {
		os.Setenv("TKG_HTTP_PROXY", "")
	}

	_, ok = os.LookupEnv("TKG_HTTPS_PROXY")
	if !ok {
		os.Setenv("TKG_HTTPS_PROXY", "")
	}

	configDir, err := getTKGConfigDir()
	if err != nil {
		return NonUsageError(cmd, err, "unable to determine Tanzu configuration directory.")
	}

	if tso.configFile == "" {
		clusterConfigPath, err := getStandaloneClusterConfig(clusterName)
		if err != nil {
			return Error(err, "unable to load standalone cluster configuration")
		}
		tso.configFile = clusterConfigPath
	}

	// setup client options
	opt := tkgctl.Options{
		KubeConfig:        "",
		KubeContext:       "",
		ConfigDir:         configDir,
		LogOptions:        tkgctl.LoggingOptions{Verbosity: logLevel, File: logFile},
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
		ClusterName:   clusterName,
		Force:         tso.force,
		SkipPrompt:    tso.skip,
		ClusterConfig: tso.configFile,
	}

	err = c.DeleteStandalone(teardownRegionOpts)
	if err != nil {
		return Error(err, "standalone cluster deletion failed")
	}

	err = removeStandaloneClusterConfig(clusterName)
	if err != nil {
		return Error(err, "could not remove temorary standalone cluster config")
	}

	return nil
}

func getStandaloneClusterConfig(clusterName string) (string, error) {
	configDir, err := getTKGConfigDir()
	if err != nil {
		return "", err
	}

	// fetch the expected cluster configuration for the restore cycle
	clusterConfigDir := filepath.Join(configDir, "clusterconfigs")
	clusterConfigFile := clusterName + ".yaml"
	readConfigPath := filepath.Join(clusterConfigDir, clusterConfigFile)

	log.Infof("Loading bootstrap cluster config for standalone cluster at '%v'", readConfigPath)

	_, err = os.Stat(readConfigPath)
	if os.IsNotExist(err) {
		log.Infof("no bootstrap cluster config found - using default config")
		return "", nil
	}

	return readConfigPath, nil
}

func removeStandaloneClusterConfig(clusterName string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".config", "tanzu", "tkg", "configs")
	clusterConfigFile := clusterName + "_ClusterConfig"
	deleteConfigPath := filepath.Join(configDir, clusterConfigFile)

	log.Infof("Removing temporary bootstrap cluster config for standalone cluster at '%v'", deleteConfigPath)

	// If file doesn't exist, try UI clusterconfig and skip
	_, err = os.Stat(deleteConfigPath)
	if os.IsNotExist(err) {
		log.Infof("no bootstrap cluster config found - skipping")
		err = removeUIStandaloneClusterConfig(clusterName)
		if err != nil {
			return fmt.Errorf("error removing UI bootstrap clusterconfig: %v", err)
		}
		return nil
	}

	err = os.Remove(deleteConfigPath)
	if err != nil {
		return fmt.Errorf("could not delete file: %v", deleteConfigPath)
	}

	return nil
}

func removeUIStandaloneClusterConfig(clusterName string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".config", "tanzu", "clusterconfigs")
	clusterConfigFile := clusterName + ".yaml"
	deleteConfigPath := filepath.Join(configDir, clusterConfigFile)

	log.Infof("Removing temporary UI bootstrap cluster config for standalone cluster at '%v'", deleteConfigPath)

	// If file doesn't exist, assume CAPD and skip
	_, err = os.Stat(deleteConfigPath)
	if os.IsNotExist(err) {
		log.Infof("no UI bootstrap cluster config found - skipping")
		return nil
	}

	err = os.Remove(deleteConfigPath)
	if err != nil {
		return fmt.Errorf("could not delete file: %v", deleteConfigPath)
	}

	return nil
}
