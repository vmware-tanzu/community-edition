// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kubeconfig

import (
	"fmt"
	"os"
	"path/filepath"
)

// write writes cfg to configPath
// it will ensure the directories in the path if necessary
func write(cfg *Config, configPath string) error {
	encoded, err := Encode(cfg)
	if err != nil {
		return err
	}
	// NOTE: 0755 / 0600 are to match client-go
	dir := filepath.Dir(configPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s for kubeconfig: %v", dir, err)
		}
	}
	if err := os.WriteFile(configPath, encoded, 0600); err != nil {
		return fmt.Errorf("failed to write kubeconfig to %s: %v", configPath, err)
	}
	return nil
}
