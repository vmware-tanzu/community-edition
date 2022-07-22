// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package oraclecpi_test

import (
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/repo"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/ytt"
)

var _ = Describe("vSphere CPI Ytt Templates", func() {
	var (
		filePaths    []string
		values       string
		output       string
		yttRenderErr error

		configDir = filepath.Join(repo.RootDir(), "addons/packages/oracle-cpi/1.22.1/bundle/config")
	)

	BeforeEach(func() {
		filePaths = []string{
			// base templates
			filepath.Join(configDir, "upstream/oci-cloud-controller-manager.yaml"),
			filepath.Join(configDir, "upstream/oci-cloud-controller-manager-rbac.yaml"),

			filepath.Join(configDir, "upstream/oci-csi-controller-driver.yaml"),
			filepath.Join(configDir, "upstream/oci-csi-node-driver.yaml"),
			filepath.Join(configDir, "upstream/oci-csi-node-rbac.yaml"),
			filepath.Join(configDir, "upstream/oci-csi-storage-class.yaml"),

			filepath.Join(configDir, "upstream/oci-volume-provisioner.yaml"),
			filepath.Join(configDir, "upstream/oci-volume-provisioner-fss.yaml"),
			filepath.Join(configDir, "upstream/oci-volume-provisioner-rbac.yaml"),
			filepath.Join(configDir, "upstream/oci-volume-storage-class.yaml"),

			filepath.Join(configDir, "upstream/provider-config.yaml"),

			// overlays
			filepath.Join(configDir, "overlay/provider-config.yaml"),

			// default values
			filepath.Join(configDir, "values.star"),
			filepath.Join(configDir, "values.yaml"),

			// data lib
			filepath.Join(configDir, "provider-config.lib.txt"),
		}
	})

	JustBeforeEach(func() {
		output, yttRenderErr = ytt.RenderYTTTemplate(ytt.CommandOptions{}, filePaths, strings.NewReader(values))
	})

	When("values are provided", func() {

		BeforeEach(func() {
			// choose to use the sample values
			values = `
#@data/values
#@overlay/match-child-defaults missing_ok=True

---
compartment: sample-compartment
vcn: sample-vcn
loadBalancer:
  subnet1: sample-subnet1
  subnet2: sample-subnet2
`
		})

		It("can render cloud config secret successfully", func() {
			Expect(yttRenderErr).ToNot(HaveOccurred())
			Expect(output).ToNot(BeEmpty())
		})
	})
})
