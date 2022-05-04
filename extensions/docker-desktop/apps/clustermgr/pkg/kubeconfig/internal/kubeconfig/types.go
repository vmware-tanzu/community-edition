// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kubeconfig

/*
NOTE: all of these types are based on the upstream v1 types from client-go
https://github.com/kubernetes/client-go/blob/0bdba2f9188006fc64057c2f6d82a0f9ee0ee422/tools/clientcmd/api/v1/types.go
We've forked them to:
- remove types and fields kind does not need to inspect / modify
- generically support fields kind doesn't inspect / modify using yaml.v3
- have clearer names (AuthInfo -> User)
*/

// Config represents a KUBECONFIG, with the fields kind is likely to use
// Other fields are handled as unstructured data purely read for writing back
// to disk via the OtherFields field
type Config struct {
	// Clusters is a map of referenceable names to cluster configs
	Clusters []NamedCluster `yaml:"clusters,omitempty"`
	// Users is a map of referenceable names to user configs
	Users []NamedUser `yaml:"users,omitempty"`
	// Contexts is a map of referenceable names to context configs
	Contexts []NamedContext `yaml:"contexts,omitempty"`
	// CurrentContext is the name of the context that you would like to use by default
	CurrentContext string `yaml:"current-context,omitempty"`
	// OtherFields contains fields kind does not inspect or modify, these are
	// read purely for writing back
	OtherFields map[string]interface{} `yaml:",inline,omitempty"`
}

// NamedCluster relates nicknames to cluster information
type NamedCluster struct {
	// Name is the nickname for this Cluster
	Name string `yaml:"name"`
	// Cluster holds the cluster information
	Cluster Cluster `yaml:"cluster"`
}

// Cluster contains information about how to communicate with a kubernetes cluster
type Cluster struct {
	// Server is the address of the kubernetes cluster (https://hostname:port).
	Server string `yaml:"server,omitempty"`
	// OtherFields contains fields kind does not inspect or modify, these are
	// read purely for writing back
	OtherFields map[string]interface{} `yaml:",inline,omitempty"`
}

// NamedUser relates nicknames to user information
type NamedUser struct {
	// Name is the nickname for this User
	Name string `yaml:"name"`
	// User holds the user information
	// We do not touch this and merely write it back
	User map[string]interface{} `yaml:"user"`
}

// NamedContext relates nicknames to context information
type NamedContext struct {
	// Name is the nickname for this Context
	Name string `yaml:"name"`
	// Context holds the context information
	Context Context `yaml:"context"`
}

// Context is a tuple of references to a cluster (how do I communicate with a kubernetes cluster), a user (how do I identify myself), and a namespace (what subset of resources do I want to work with)
type Context struct {
	// Cluster is the name of the cluster for this context
	Cluster string `yaml:"cluster"`
	// User is the name of the User for this context
	User string `yaml:"user"`
	// OtherFields contains fields kind does not inspect or modify, these are
	// read purely for writing back
	OtherFields map[string]interface{} `yaml:",inline,omitempty"`
}
