// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetTCEContainerName() string {
	return DefaultClusterName + "-control-plane"
}

func GetUserHome() string {
	dirname, _ := os.UserHomeDir()
	return dirname
}

func GetClusterConfigFileName() string {
	return filepath.Join(GetUserHome(), ClusterConfigFileName)
}

func GetKubeconfigFileName() string {
	return filepath.Join(GetUserHome(), ".kube", "config")
}

func GetLogsFileName() string {
	return filepath.Join(GetUserHome(), ClusterLogFile)
}

func GetInternalLogsFileName() string {
	return filepath.Join(GetUserHome(), ".config", "tanzu", "tkg", "unmanaged", DefaultClusterName, "bootstrap.log")
}

func WriteConfigFile(configBytes []byte, fileNamePath string) error {
	err := os.WriteFile(fileNamePath, configBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file to (%s). Error: %s", fileNamePath, err.Error())
	}

	return nil
}
