// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package cluster implements functionality for interacting with clusters to
// perform CRUD and other operations.
package cluster

import (
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/config"
)

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
	// Config contains the full cluster creation details passed in from the user when calling
	// create.
	Config *config.LocalClusterConfig
}

// ClusterManager provides methods for creating and managing Kubernetes
// clusters.
//nolint:golint
type ClusterManager interface {
	// Create will create a new cluster or return an error indicating a problem
	// during creation.
	Create(opts *CreateOpts) (*KubernetesCluster, error)
	// Get retrieves cluster information or return an error indicating a problem.
	Get(clusterName string) (*KubernetesCluster, error)
	// List gets a list of all local clusters.
	List() ([]*KubernetesCluster, error)
	// Delete will destroy a cluster or return an error indicating a problem.
	Delete(clusterName string) error
}

// NewClusterManager gets a ClusterManager implementation.
func NewClusterManager() ClusterManager {
	// For now, just hard coding to return our KindClusterManager.
	return KindClusterManager{}
}
