// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package tanzu is responsible for orchestrating the various packages that satisfy unmanaged
// operations such as create, configure, list, and delete. This package is meant to be the API
// entrypoint for those calling unmanaged programmatically.
package tanzu

const (
	// 0  - Success!
	Success int = iota

	// 1  - Provided configurations could not be validated
	InvalidConfig

	// 2  - Could not create local cluster directories
	ErrCreatingClusterDirs

	// 3  - Unable to get TKR BOM
	ErrTkrBom

	// 4  - Could not render config
	ErrRenderingConfig

	// 5  - TKR BOM not parseable
	ErrTkrBomParsing

	// 6  - Could not resolve kapp controller bundle
	ErrKappBundleResolving

	// 7  - Unable to create new cluster
	ErrCreateCluster

	// 8  - Unable to use existing cluster (if provided)
	ErrExistingCluster

	// 9  - Could not install kapp controller to cluster
	ErrKappInstall

	// 10 - Could not install core package repo to cluster
	ErrCorePackageRepoInstall

	// 11 - Could not install additional package repo
	ErrOtherPackageRepoInstall

	// 12 - Could not install CNI package
	ErrCniInstall

	// 13 - Failed to merge kubeconfig and set context
	ErrKubeconfigContextSet

	// 14 - Cound not install profile
	ErrProfileInstall
)
