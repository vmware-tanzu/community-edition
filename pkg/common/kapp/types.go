// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kapp

// Repository object
type Repository struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Config struct
type Config struct {
	// Kubeconfig is the users kubeconfig
	Kubeconfig string

	// WorkingDirectory is the users working directory
	WorkingDirectory string
	// ExtensionNamespace is the extension namespace to install into
	ExtensionNamespace string
	// ExtensionServiceAccountPostfix is the extension postfix for service account
	ExtensionServiceAccountPostfix string
	// ExtensionRoleBindingPostfix is the extension postfix for role binding
	ExtensionRoleBindingPostfix string

	// Repositories list of repos
	Repositories []*Repository `json:"repositories"`
}

// AppCrdInput for creating an app
type AppCrdInput struct {
	//Common between UserDefined and Secret
	Namespace string

	// UserDefined
	Name       string
	URL        string
	Paths      map[string]string
	Version    string
	ConfigPath string
	Config     []byte

	// From Secret
	ClusterName string

	// Force delete
	Force bool
	// Teardown by deleting ServiceAccount and RoleBinding
	Teardown bool
}

// Kapp object for kapp
type Kapp struct {
	config *Config

	// localWorkingDirectory is where all modification based on user overrides is defined
	localWorkingDirectory string
}
