// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package types

// File - yep, it's that
type File struct {
	Name        string `json:"filename"`
	Description string `json:"description,omitempty"`
}

// Extension - yep, it's that
type Extension struct {
	Name                   string `json:"name"`
	Description            string `json:"description,omitempty"`
	Version                string `json:"version"`
	KubernetesMinSupported string `json:"minsupported,omitempty"`
	KubernetesMaxSupported string `json:"maxsupported,omitempty"`
	Files                  []File `json:"files"`
}

// Metadata outer container for metadata
type Metadata struct {
	Extensions      []Extension `json:"extensions"`
	Version         string      `json:"version"`
	GitHubRepo      string      `json:"repo,omitempty"`
	GitHubBranchTag string      `json:"branch,omitempty"`

	ExtensionLookup map[string]*Extension
}

// Version - released version
type Version struct {
	Version string `json:"version"`
	Date    string `json:"date,omitempty"`
}

// Release outer container for metadata
type Release struct {
	Versions []Version `json:"versions"`
	Stable   string    `json:"stable"`
	Date     string    `json:"date"`

	VersionLookup map[string]*Version
}

// IManager interface for Manager
type IManager interface {
	RawMetadata() ([]byte, error)
	InitMetadata() (*Metadata, error)
	InitRelease() (*Release, error)
}
