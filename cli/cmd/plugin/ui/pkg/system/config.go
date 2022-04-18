// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package system

import (
	"os"
	"path/filepath"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/config"
)

// GetConfigDir gets the local filepath where we store our configuration.
func GetConfigDir() string {
	tanzuConfigDir, err := config.LocalDir()
	if err != nil {
		// Fall back to a hard coded default for now
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".config", "tanzu", "tkg")
	}

	return filepath.Join(tanzuConfigDir, "tkg")
}
