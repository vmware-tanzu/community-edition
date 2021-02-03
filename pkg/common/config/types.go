// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

// Config struct
type Config struct {
	// Init config file
	configFile string
	// Release version
	ReleaseVersion string `json:"version"`
	// GitHub token
	GithubToken string `json:"token,omitempty"`
}
