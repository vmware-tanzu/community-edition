// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cluster

import (
	"os"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/config"
)

type NoopClusterManager struct{}

// Create will create a new KubernetesCluster that points to the default
func (ncm NoopClusterManager) Create(c *config.UnmanagedClusterConfig) (*KubernetesCluster, error) {
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
func (ncm NoopClusterManager) Delete(c *config.UnmanagedClusterConfig) error {
	return nil
}

// Prepare doesn't perform any preparation steps before cluster creation.
func (ncm NoopClusterManager) Prepare(c *config.UnmanagedClusterConfig) error {
	return nil
}

// PreflightCheck performs any pre-checks that can find issues up front that
// would cause problems for cluster creation.
func (ncm NoopClusterManager) PreflightCheck() ([]string, []error) {
	return nil, nil
}

// ProviderNotify is a noop. Nothing to notify about for the noop provider
func (ncm NoopClusterManager) ProviderNotify() []string {
	return []string{}
}
