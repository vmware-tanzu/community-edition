// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tanzu

import (
	"os"
	"path/filepath"
)

// PortMap is the mapping between a host port and a container port.
type PortMap struct {
	// HostPort is the port on the host machine.
	HostPort int `yaml:"HostPort"`
	// ContainerPort is the port on the container to map to.
	ContainerPort int `yaml:"ContainerPort"`
}

// LocalClusterConfig contains all the configuration settings for creating a
// local Tanzu cluster.
type LocalClusterConfig struct {
	// ClusterName is the name of the cluster.
	ClusterName string `yaml:"ClusterName"`
	// Provider is the local infrastructure provider to use (e.g. kind).
	Provider string `yaml:"Provider"`
	// ProviderConfiguration offers optional provider-specific configuration.
	// The exact keys and values accepted are determined by the provider.
	ProviderConfiguration map[string]interface{}
	// CNI is the networking CNI to use in the cluster. Default is antrea.
	CNI string `yaml:"CNI"`
	// CNIConfiguration offers optional cni-plugin specific configuration.
	// The exact keys and values accepted are determined by the CNI choice.
	CNIConfiguration map[string]interface{}
	// PodCidr is the Pod CIDR range to assign pod IP addresses.
	PodCidr string `yaml:"PodCidr"`
	// ServiceCidr is the Service CIDR range to assign service IP addresses.
	ServiceCidr string `yaml:"ServiceCidr"`
	// TkrLocation is the path to the Tanzu Kubernetes Release (TKR) data.
	TkrLocation string `yaml:"TkrLocation"`
	// PortsToForward contains a mapping of host to container ports that should
	// be exposed.
	PortsToForward []PortMap `yaml:"PortsToForward"`
}

// KubeConfigPath gets the full path to the KubeConfig for this local cluster.
func (lcc *LocalClusterConfig) KubeConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), configDir, tanzuConfigDir, lcc.ClusterName+".yaml")
}
