// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package cluster handles the resource creation to run Kubernetes clusters and bootstrap of the
// Kubernetes cluster. Additional provides can be introduced by implementing the ClusterManager
// interface.
package cluster

import (
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/config"
)

const (
	NoneClusterManagerProvider = "none"
	KindClusterManagerProvider = "kind"
)

// KubernetesCluster represents a defines k8s cluster.
type KubernetesCluster struct {
	// Name is the name of the cluster.
	Name string
	// KubeConfig contains the Kubeconfig data for the cluster.
	Kubeconfig []byte
}

// Manager provides methods for creating and managing Kubernetes
// clusters.
type Manager interface {
	// Create will create a new cluster or return an error indicating a problem
	// during creation.
	Create(c *config.UnmanagedClusterConfig) (*KubernetesCluster, error)
	// Get retrieves cluster information or return an error indicating a problem.
	Get(clusterName string) (*KubernetesCluster, error)
	// Delete will destroy a cluster or return an error indicating a problem.
	Delete(c *config.UnmanagedClusterConfig) error
	// Prepare will fetch an image or perform any pre-steps that can be done
	// prior to actually creating the cluster.
	Prepare(c *config.UnmanagedClusterConfig) error
	// PreflightCheck performs any pre-checks that can find issues up front that
	// would cause problems for cluster creation. Returns nil if there are no
	// errors found, otherwise a list of the errors that need to be resolved.
	PreflightCheck() []error
}

// NewClusterManager provides a way to dynamically get a cluster manager based on the unmanaged cluster config provider
func NewClusterManager(c *config.UnmanagedClusterConfig) Manager {
	switch c.Provider {
	case KindClusterManagerProvider:
		return NewKindClusterManager()
	case NoneClusterManagerProvider:
		return NewNoopClusterManager()
	}

	// By default, return a noop cluster manager in cases we can't parse what provider the users gave
	return NewNoopClusterManager()
}

// NewNoopClusterManager creates a new noop cluster manager - intended for use with "none" provider
func NewNoopClusterManager() Manager {
	return NoopClusterManager{}
}

// NewKindClusterManager gets a ClusterManager implementation for the kind provider.
func NewKindClusterManager() Manager {
	// For now, just hard coding to return our KindClusterManager.
	return KindClusterManager{}
}
