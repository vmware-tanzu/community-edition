// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

// PortMap is the mapping between a host port and a container port.
type PortMap struct {
	// HostPort is the port on the host machine.
	HostPort int
	// ContainerPort is the port on the container to map to.
	ContainerPort int
}

// LocalClusterConfig contains all of the configuration settings for creating a
// local Tanzu cluster.
type LocalClusterConfig struct {
	// ClusterName is the name of the cluster.
	ClusterName string
	// Provider is the local infastructure provider to use (e.g. kind).
	Provider string
	// CNI is the networking CNI to use in the cluster. Default is antrea.
	CNI string
	// PodCidr is the Pod CIDR range to assign pod IP addresses.
	PodCidr string
	// ServiceCidr is the Service CIDR range to assign service IP addresses.
	ServiceCidr string
	// TkrLocation is the path to the Tanzu Kubernetes Release (TKR) data.
	TkrLocation string
	// PortsToForward contains a mapping of host to container ports that should
	// be exposed.
	PortsToForward []PortMap
}
