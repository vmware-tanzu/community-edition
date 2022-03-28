// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//nolint:goconst
package tkr

import (
	"io/ioutil"
	"os"
	"testing"
)

var sampleTkrYamls = `version: v1
unmanagedClusterPluginVersions:
- version: dev
  supportedTkrVersions:
  - image: projects.registry.vmware.com/tce/tkr:v0.17.0
  - image: projects.registry.vmware.com/tce/tkr:v0.16.0
- version: v0.11.0
  supportedTkrVersions:
  - image: projects.registry.vmware.com/tce/tkr:v0.16.0
  - image: projects.registry.vmware.com/tce/tkr:v0.15.0
- version: v0.10.0
  supportedTkrVersions:
  - image: projects.registry.vmware.com/tce/tkr:v0.0.1
`

//nolint:gocyclo
func TestOrderOfCompatibleTkrs(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "compatibility-test-")
	if err != nil {
		t.Errorf(err.Error())
	}

	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(sampleTkrYamls))
	if err != nil {
		t.Errorf(err.Error())
	}

	c, err := ReadCompatibilityFile(tmpFile.Name())
	if err != nil {
		t.Errorf(err.Error())
	}

	if c.Version != "v1" {
		t.Errorf("expected version of TKr compatibility to be v1, was: %s\n", c.Version)
	}

	if len(c.UnmanagedClusterPluginVersions) != 3 {
		t.Errorf("expected to find 3 CLI plugin versions, was: %v\n", c.UnmanagedClusterPluginVersions)
	}

	if c.UnmanagedClusterPluginVersions[0].Version != "dev" {
		t.Errorf("expected first CLI plugin version to be dev, was: %s\n", c.UnmanagedClusterPluginVersions[0])
	}

	if len(c.UnmanagedClusterPluginVersions[0].SupportedTkrVersions) != 2 {
		t.Errorf("expected first CLI plugin to have 2 compatibility TKrs, was: %v\n", c.UnmanagedClusterPluginVersions[0].SupportedTkrVersions)
	}

	if c.UnmanagedClusterPluginVersions[0].SupportedTkrVersions[0].Path != "projects.registry.vmware.com/tce/tkr:v0.17.0" {
		t.Errorf("expected first compatible TKR for first CLI plugin to be projects.registry.vmware.com/tce/tkr:v0.17.0, was: %s\n", c.UnmanagedClusterPluginVersions[0].SupportedTkrVersions[0].Path)
	}

	if c.UnmanagedClusterPluginVersions[0].SupportedTkrVersions[1].Path != "projects.registry.vmware.com/tce/tkr:v0.16.0" {
		t.Errorf("expected second compatible TKR for first CLI plugin to be projects.registry.vmware.com/tce/tkr:v0.16.0, was: %s\n", c.UnmanagedClusterPluginVersions[0].SupportedTkrVersions[1].Path)
	}

	if c.UnmanagedClusterPluginVersions[1].Version != "v0.11.0" {
		t.Errorf("expected second CLI plugin version to be v0.11.0, was: %s\n", c.UnmanagedClusterPluginVersions[1])
	}

	if len(c.UnmanagedClusterPluginVersions[1].SupportedTkrVersions) != 2 {
		t.Errorf("expected second CLI plugin to have 2 compatibility TKrs, was: %v\n", c.UnmanagedClusterPluginVersions[1].SupportedTkrVersions)
	}

	if c.UnmanagedClusterPluginVersions[1].SupportedTkrVersions[0].Path != "projects.registry.vmware.com/tce/tkr:v0.16.0" {
		t.Errorf("expected first compatible TKR for second CLI plugin to be projects.registry.vmware.com/tce/tkr:v0.16.0, was: %s\n", c.UnmanagedClusterPluginVersions[1].SupportedTkrVersions[0].Path)
	}

	if c.UnmanagedClusterPluginVersions[1].SupportedTkrVersions[1].Path != "projects.registry.vmware.com/tce/tkr:v0.15.0" {
		t.Errorf("expected second compatible TKR for first CLI plugin to be projects.registry.vmware.com/tce/tkr:v0.15.0, was: %s\n", c.UnmanagedClusterPluginVersions[1].SupportedTkrVersions[1].Path)
	}

	if c.UnmanagedClusterPluginVersions[2].Version != "v0.10.0" {
		t.Errorf("expected first CLI plugin version to be v0.10.0, was: %s\n", c.UnmanagedClusterPluginVersions[2])
	}

	if len(c.UnmanagedClusterPluginVersions[2].SupportedTkrVersions) != 1 {
		t.Errorf("expected second CLI plugin to have 1 compatibility TKrs, was: %v\n", c.UnmanagedClusterPluginVersions[2].SupportedTkrVersions)
	}

	if c.UnmanagedClusterPluginVersions[2].SupportedTkrVersions[0].Path != "projects.registry.vmware.com/tce/tkr:v0.0.1" {
		t.Errorf("expected first compatible TKR for third CLI plugin to be projects.registry.vmware.com/tce/tkr:v0.0.1, was: %s\n", c.UnmanagedClusterPluginVersions[2].SupportedTkrVersions[0].Path)
	}
}
