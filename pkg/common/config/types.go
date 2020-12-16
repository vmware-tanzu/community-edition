// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

// Config struct
type Config struct {
	// Init config file
	configFile string
	// GitHub token
	githubToken string

	// Raw config file
	byRaw []byte
}
