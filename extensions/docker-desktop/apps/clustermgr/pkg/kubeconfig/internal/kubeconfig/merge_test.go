// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kubeconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteMerged(t *testing.T) {
	tests := []struct {
		name          string
		sourceFile    string
		secondaryFile string
		expectedCount int
	}{
		{
			name:          "Add New Cluster to config",
			sourceFile:    "config.orig",
			secondaryFile: "config.secondary",
			expectedCount: 3,
		},
		{
			name:          "Append existing cluster config",
			sourceFile:    "config.orig",
			secondaryFile: "config.secondary-amend",
			expectedCount: 2,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			copyFile(filepath.Join(".", "testdata", tt.sourceFile), filepath.Join(".", "testdata", "config.merge"))
			defer func() {
				os.Remove(filepath.Join(".", "testdata", "config.merge"))
			}()
			config, err := Read(filepath.Join(".", "testdata", tt.secondaryFile))
			if err != nil {
				t.Error("Failed to read secondary config file")
			}
			err = WriteMerged(config, filepath.Join(".", "testdata", "config.merge"))
			if err != nil {
				t.Errorf("WriteMerge failed. Expected Error %v, Got %v", nil, err)
			}
			config, _ = Read(filepath.Join(".", "testdata", "config.merge"))
			if len(config.Clusters) != tt.expectedCount {
				t.Errorf("WriteMerge failed. Expected Cluster Count %d, Got %d", 3, len(config.Clusters))
			}
		})
	}
}
