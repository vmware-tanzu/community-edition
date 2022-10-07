// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package oraclecpi_test

import (
	"encoding/base64"
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
			filepath.Join(configDir, "overlay/proxy.yaml"),

			// default values
			filepath.Join(configDir, "values.star"),
			filepath.Join(configDir, "values.yaml"),
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

	When("static authentication creds are provided", func() {
		BeforeEach(func() {
			values = `
#@data/values
#@overlay/match-child-defaults missing_ok=True

---
auth:
  region: us-sanjose-1
  tenancy: ocid1.tenancy.oc1..aaaaaaaaxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
  user: ocid1.user.oc1..aaaaaaaaveptnubesjspvkqzohqkjsdblv5dmnkyvbolbna6rf76io3uox2a
  key: |
      -----BEGIN PRIVATE KEY-----

      -----END PRIVATE KEY-----
  fingerprint: eb:02:ee:4b:4c:xx:xx:xx:xx:55:df:54:00:db:be:0f
  passphrase: ""
` //#nosec
		})

		It("can render credentials successfully", func() {
			Expect(yttRenderErr).ToNot(HaveOccurred())
			Expect(output).ToNot(BeEmpty())
			Expect(output).To(ContainSubstring("cloud-provider.yaml:"))

			/* #nosec */
			cred := "auth:\n  region: us-sanjose-1\n  tenancy: ocid1.tenancy.oc1..aaaaaaaaxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\n  user: ocid1.user.oc1..aaaaaaaaveptnubesjspvkqzohqkjsdblv5dmnkyvbolbna6rf76io3uox2a\n  fingerprint: eb:02:ee:4b:4c:xx:xx:xx:xx:55:df:54:00:db:be:0f\n  passphrase: \"\"\n  useInstancePrincipals: false\n  key: |\n    -----BEGIN PRIVATE KEY-----\n\n    -----END PRIVATE KEY-----\n"
			encodedCred := base64.StdEncoding.EncodeToString([]byte(cred))
			Expect(output).To(ContainSubstring(encodedCred))
		})
	})

	When("proxy configurations are provided", func() {
		BeforeEach(func() {
			values = `
#@data/values
#@overlay/match-child-defaults missing_ok=True

---
compartment: sample-compartment
vcn: sample-vcn
loadBalancer:
  subnet1: sample-subnet1
  subnet2: sample-subnet2
http_proxy: 10.0.0.1
https_proxy: 10.0.0.2
no_proxy: 10.0.0.1
`
		})

		It("can render cloud config secret successfully", func() {
			Expect(yttRenderErr).ToNot(HaveOccurred())
			Expect(output).ToNot(BeEmpty())
			Expect(output).To(ContainSubstring(` env:
        - name: HTTP_PROXY
          value: 10.0.0.1
        - name: HTTPS_PROXY
          value: 10.0.0.2
        - name: NO_PROXY
          value: 10.0.0.1`))
		})
	})

})
