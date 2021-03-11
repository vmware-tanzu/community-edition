// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kapp

const (
	// DefaultWorkingDirectory is the working directory
	DefaultWorkingDirectory string = "working"

	// DefaultAppCrdNamespace is the default App Crd namespace
	DefaultAppCrdNamespace string = "tanzu-extensions"
	// DefaultServiceAccountPostfix is the default Service Account postfix name
	DefaultServiceAccountPostfix string = "-extension-sa"
	// DefaultRoleBindingPostfix is the default Role Binding postfix name
	DefaultRoleBindingPostfix string = "-extension"

	// DefaultRepositoryName is tce-main.tanzu.vmware
	DefaultRepositoryName string = "tce-main.tanzu.vmware"
	// DefaultRepositoryImage is projects.registry.vmware.com/tce/main:dev
	DefaultRepositoryImage string = "projects.registry.vmware.com/tce/main:dev"
)
