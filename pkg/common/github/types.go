// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package github

import (
	download "github.com/vmware-tanzu/tce/pkg/common/download"
	types "github.com/vmware-tanzu/tce/pkg/common/types"
)

// Config struct
type Config struct {
	// GitHubRepo to access extension data
	GitHubRepo string
	// GitHubBranchTag to use
	GitHubBranchTag string `json:"version"`
	// GitHubURI formatted URI for root location of extensions
	GitHubURI string
	// ExtensionDirectory is the extension directory
	ExtensionDirectory string

	// github token
	Token string `json:"token"`

	// originalBranchTag provided
	originalBranchTag string
}

// Manager for downloads
type Manager struct {
	cfg *Config
	dl  *download.Manager

	extensionDirectoryRoot string
	extensionDirectory     string
	extMgr                 types.IManager
}
