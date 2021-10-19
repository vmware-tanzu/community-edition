// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package cluster implements functionality for interacting with clusters to
// perform CRUD and other operations.
package cluster

// KubernetesCluster represents a defines k8s cluster.
type KubernetesCluster struct {
	// Name is the name of the cluster.
	Name string
	// KubeConfig contains the Kubeconfig data for the cluster.
	Kubeconfig string
}

// CreateOpts contains data to be used when creating a new cluster.
type CreateOpts struct {
	// Name is the name for the new cluster.
	Name string
	// KubeconfigPath is the path to the kubeconfig to use.
	KubeconfigPath string
	// Config contains the raw configuration data to use when creating the cluster.
	Config []byte
}

// ClusterManager provides methods for creating and managing Kubernetes
// clusters.
type ClusterManager interface {
	// Create will create a new cluster or return an error indicating a problem
	// during creation.
	Create(opts *CreateOpts) (*KubernetesCluster, error)
	// Get retrieves cluster information or return an error indicating a problem.
	Get(name string) (*KubernetesCluster, error)
	// List gets a list of all local clusters.
	List() ([]*KubernetesCluster, error)
	// Delete will destroy a cluster or return an error indicating a problem.
	Delete(name string) error
}

// NewClusterManager gets a ClusterManager implementation.
func NewClusterManager() ClusterManager {
	// For now, just hard coding to return our KindClusterManager.
	return KindClusterManager{}
}
