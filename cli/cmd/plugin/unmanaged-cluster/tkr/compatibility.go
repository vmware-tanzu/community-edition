// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkr

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"gopkg.in/yaml.v3"
)

// TkrImage is the image:path combo of a compatible TKr
type TkrImage struct { //nolint:revive
	Path string `yaml:"image"`
}

// UnmanagedClusterPluginVersion is the version of the unmanaged cluster CLI plugin
// and versions of TKrs it's compatible with
type UnmanagedClusterPluginVersion struct {
	Version              string     `yaml:"version"`
	SupportedTkrVersions []TkrImage `yaml:"supportedTkrVersions"`
}

// TkrCompatibility is the top level compatibility structure
type Compatibility struct {
	Version                        string                          `yaml:"version"`
	UnmanagedClusterPluginVersions []UnmanagedClusterPluginVersion `yaml:"unmanagedClusterPluginVersions"`
}

// GetLatestCompatibilityTag expects a non-tagged registry URL that references a compatibility file image
// For example: projects.registry.vmware.com/tce/compatibility
// And returns the latest major tagged version
// For example, if a given registry is tagged [v1, v2, v3, v4],
// this function will return "v4"
// Otherwise, it returns an error
func GetLatestCompatibilityTag(registryURL string) (string, error) {
	img := Image{
		RegistryURL: registryURL,
	}

	tags, err := img.GetTags()
	if err != nil {
		return "", err
	}

	// Sorts tags based on major version
	// Reference similar TKG compatibility code:
	// https://github.com/vmware-tanzu/tanzu-framework/blob/4e383ca5760a67b41ad27dd52c23ae78378d302c/pkg/v1/tkg/tkgconfigbom/bom.go#L461
	// Expects tags to be formatted as `vN` where N is a number.
	// Example: [v1, v2, v3, v4]
	tagNum := []int{}
	for _, tag := range tags {
		ver, err := strconv.Atoi(tag[1:])
		if err == nil {
			tagNum = append(tagNum, ver)
		}
	}

	sort.Ints(tagNum)
	if len(tagNum) == 0 {
		return "", fmt.Errorf("failed to get valid image tags for compatibility image. Expected tags to be formatted as `vN` where N is a number. Actual tags: %s", tags)
	}

	// Re-format last tag in list (since it is the "latest")
	return fmt.Sprintf("v%d", tagNum[len(tagNum)-1]), nil
}

// ReadCompatibilityFile will process a given file on the filesystem
// and return a TkrCompatibility struct
func ReadCompatibilityFile(filePath string) (*Compatibility, error) {
	c := &Compatibility{}
	rawC, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read compatibility file. Error: %s", err.Error())
	}

	err = yaml.Unmarshal(rawC, c)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal local compatibility file yaml. Error: %s", err.Error())
	}

	return c, nil
}
