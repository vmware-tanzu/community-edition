// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package gcp

// Config struct
type Config struct {
	// TceBucket location of the TCE GCP bucket for Tanzu CLI
	TceBucket string
	// MetadataDirectory is the metadata directory
	MetadataDirectory string
	// MetadataFileName is the file name for the metadata
	MetadataFileName string
	// VersionTag to use
	VersionTag string `json:"version"`
	// ReleasesFileName is the file name defining all releases
	ReleasesFileName string
}

// Bucket object for GCP
type Bucket struct {
	config *Config

	// metadata for current version
	remoteMetadataFile string
	localMetadataFile  string

	// release management file
	remoteReleaseFile string
}
