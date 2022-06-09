// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetTCEContainerName gets the name of the cluster's container. Currently we
// only create single node clusters, so this is the control plane container.
func GetTCEContainerName() string {
	return DefaultClusterName + "-control-plane"
}

// GetUserHome is a helper function to get a user's home directory path.
func GetUserHome() string {
	dirname, _ := os.UserHomeDir()
	return dirname
}

// GetClusterConfigFileName gets the full path to the cluster config file.
func GetClusterConfigFileName() string {
	return filepath.Join(GetUserHome(), ClusterConfigFileName)
}

// GetKubeconfigFileName gets the full path to the kubeconfig file.
func GetKubeconfigFileName() string {
	return filepath.Join(GetUserHome(), ".kube", "config")
}

// GetLogsFileName gets the full path to the cluster log file.
func GetLogsFileName() string {
	return filepath.Join(GetUserHome(), ClusterLogFile)
}

// GetInternalLogsFileName gets the full path to the cluster bootstrap log file.
func GetInternalLogsFileName() string {
	return filepath.Join(GetUserHome(), ".config", "tanzu", "tkg", "unmanaged", DefaultClusterName, "bootstrap.log")
}

// WriteConfigFile writes the config content to a file.
func WriteConfigFile(configBytes []byte, fileNamePath string) error {
	err := os.WriteFile(fileNamePath, configBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file to (%s). Error: %s", fileNamePath, err.Error())
	}

	return nil
}
