// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package cluster implements functionality for interacting with clusters to
// perform CRUD and other operations.
package cluster

import (
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/config"
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

// ClusterManager provides methods for creating and managing Kubernetes
// clusters.
//nolint:golint
type ClusterManager interface {
	// Create will create a new cluster or return an error indicating a problem
	// during creation.
	Create(c *config.StandaloneClusterConfig) (*KubernetesCluster, error)
	// Get retrieves cluster information or return an error indicating a problem.
	Get(clusterName string) (*KubernetesCluster, error)
	// Delete will destroy a cluster or return an error indicating a problem.
	Delete(c *config.StandaloneClusterConfig) error
}

// NewClusterManager provides a way to dynamically get a cluster manager based on the standalone cluster config provider
func NewClusterManager(c *config.StandaloneClusterConfig) ClusterManager {
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
func NewNoopClusterManager() ClusterManager {
	return NoopClusterManager{}
}

// NewKindClusterManager gets a ClusterManager implementation for the kind provider.
func NewKindClusterManager() ClusterManager {
	// For now, just hard coding to return our KindClusterManager.
	return KindClusterManager{}
}
