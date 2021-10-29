// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cluster

import (
	"os"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/config"
)

type NoopClusterManager struct{}

// Create will create a new KubernetesCluster that points to the default
func (ncm NoopClusterManager) Create(c *config.LocalClusterConfig) (*KubernetesCluster, error) {
	// readkubeconfig in bytes
	kcBytes, err := os.ReadFile(c.KubeconfigPath)
	if err != nil {
		return nil, err
	}

	kc := &KubernetesCluster{
		Name:       c.ClusterName,
		Kubeconfig: kcBytes,
	}

	return kc, nil
}

// Get retrieves cluster information or return an error indicating a problem.
func (ncm NoopClusterManager) Get(clusterName string) (*KubernetesCluster, error) {
	return nil, nil
}

// Delete for noop does nothing since these clusters have no provider and are not lifecycled
func (ncm NoopClusterManager) Delete(c *config.LocalClusterConfig) error {
	return nil
}
