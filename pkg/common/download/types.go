// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package download

// Config struct
type Config struct {
	// GitHubOrg default is vmware-tanzu
	GitHubOrg string
	// GitHubRepo default is tce
	GitHubRepo string
}

// Manager for downloads
type Manager struct {
	cfg *Config
}
