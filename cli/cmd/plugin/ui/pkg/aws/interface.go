// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package aws has functions for interacting with AWS services. This is a temporary
// solution until necessary functionality can be merged into tanzu-framework's
// aws client.
package aws

import "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/bootstrap"

// Client defines methods to access AWS inventory
type Client interface {
	VerifyAccount() error
	ListVPCs() ([]*VPC, error)
	EncodeCredentials() (string, error)
	ListAvailabilityZones() ([]*AvailabilityZone, error)
	ListRegionsByUser() ([]string, error)
	GetSubnetGatewayAssociations(vpcID string) (map[string]bool, error)
	ListSubnets(vpcID string) ([]*Subnet, error)
	CreateCloudFormationStack() error
	CreateCloudFormationStackWithTemplate(template *bootstrap.Template) error
	GenerateBootstrapTemplate(GenerateBootstrapTemplateInput) (*bootstrap.Template, error)
	ListInstanceTypes(optionalAZName string) ([]string, error)
	ListCloudFormationStacks() ([]string, error)
	ListEC2KeyPairs() ([]*KeyPair, error)
}

// VPC contains information about an AWS VPC.
type VPC struct {
	// CIDR is the IP range for the VPC.
	CIDR string `json:"cidr,omitempty"`
	// ID is the VPC ID.
	ID string `json:"id,omitempty"`
}

// AvailabilityZone contains information about an AWS az.
type AvailabilityZone struct {
	// ID is the availability zone ID.
	ID string `json:"id,omitempty"`
	// Name is the availability zone name.
	Name string `json:"name,omitempty"`
}

// Subnet contains information abot an AWS subnet.
type Subnet struct {
	// AvailabilityZoneID is the availability zone ID.
	AvailabilityZoneID string `json:"availabilityZoneId,omitempty"`
	// AvailabilityZoneName is the availability zone name.
	AvailabilityZoneName string `json:"availabilityZoneName,omitempty"`
	// CIDR is the subnet's IP range.
	CIDR string `json:"cidr,omitempty"`
	// ID is the subnet ID.
	ID string `json:"id,omitempty"`
	// IsPublic indicates if the subnet is public.
	IsPublic *bool `json:"isPublic"`
	// State is the current state of the subnet.
	State string `json:"state,omitempty"`
	// VPCID is the ID of the subnet's VPC.
	VPCID string `json:"vpcId,omitempty"`
}

// GenerateBootstrapTemplateInput is the input to the GenerateBootstrapTemplate func
type GenerateBootstrapTemplateInput struct {
	// BootstrapConfigFile is the path to a CAPA bootstrapv1 configuration file that can be used
	// to customize IAM policies
	BootstrapConfigFile string
	// DisableTanzuMissionControlPermissions if true will remove IAM permissions for use by Tanzu Mission Control
	// from all nodes
	DisableTanzuMissionControlPermissions bool
}

type KeyPair struct {
	// ID is the key pair's ID.
	ID string
	// Name is the name of the key pair.
	Name string
	// Thumbprint is the key pairs thumbprint.
	Thumbprint string
}
