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
	Kubeconfig string
}

// ClusterManager provides methods for creating and managing Kubernetes
// clusters.
//nolint:golint
type ClusterManager interface {
	// Create will create a new cluster or return an error indicating a problem
	// during creation.
	Create(c *config.LocalClusterConfig) (*KubernetesCluster, error)
	// Get retrieves cluster information or return an error indicating a problem.
	Get(clusterName string) (*KubernetesCluster, error)
	// List gets a list of all local clusters.
	List() ([]*KubernetesCluster, error)
	// Delete will destroy a cluster or return an error indicating a problem.
	Delete(c *config.LocalClusterConfig) error
}

// NewKindClusterManager gets a ClusterManager implementation.
func NewKindClusterManager() ClusterManager {
	// For now, just hard coding to return our KindClusterManager.
	return KindClusterManager{}
}
