// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kubeconfig

import (
	"path/filepath"
	"testing"
)

func TestRemoveKIND(t *testing.T) {
	err := RemoveKIND("kind-upstream-k8s", filepath.Join(".", "/testdata/config"))
	if err != nil {
		t.Errorf("Failed to remove kind cluster. Expected error: nil, got: %s", err)
	}
	config, err := Read(filepath.Join(".", "/testdata/config"))
	if err != nil {
		t.Errorf("failed to read kubeconfig file. Expected error: nil, got %s", err)
	}
	if len(config.Clusters) != 1 {
		t.Errorf("RemoteKind did not cleanup the kind cluster from config. Expected: %d, got %d", 1, len(config.Clusters))
	}
}
