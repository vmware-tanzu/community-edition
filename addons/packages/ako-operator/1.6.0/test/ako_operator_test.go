// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package akooperator_test

import (
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/repo"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/ytt"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/providers/tests/unit/matchers"
)

func assertFoundOne(docs []string, err error) {
	Expect(err).NotTo(HaveOccurred())
	Expect(docs).To(HaveLen(1))
}

// AKODeploymentConfigExists asserts the existence of all AKODeploymentConfig in the ako-operator package
func AKODeploymentConfigExists(output string) {
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "AKODeploymentConfig",
		"$.metadata.name": "install-ako-for-all",
	}))
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "AKODeploymentConfig",
		"$.metadata.name": "install-ako-for-management-cluster",
	}))
}

// AviInfraSettingCRDExists asserts the existence of the AviInfraSetting in the ako-operator package
func AviInfraSettingCRDExists(output string) {
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "CustomResourceDefinition",
		"$.metadata.name": "aviinfrasettings.ako.vmware.com",
	}))
}

// AKOODeploymentExists asserts the existence of the AKOO deployment in the ako-operator package
func AKOODeploymentExists(output string) {
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "Deployment",
		"$.metadata.name": "ako-operator-controller-manager",
	}))
}

// AVISecretsExist asserts the existence of the credentials and CA in the ako-operator package
func AVISecretsExist(output string) {
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "Secret",
		"$.metadata.name": "avi-controller-credentials",
	}))
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "Secret",
		"$.metadata.name": "avi-controller-ca",
	}))
}

var _ = Describe("AKO Operator Ytt Templates", func() {
	var (
		filePaths    []string
		values       string
		output       string
		yttRenderErr error

		configDir = filepath.Join(repo.RootDir(), "addons/packages/ako-operator/1.6.0/bundle/config")

		file01AKODeploymentConfig = filepath.Join(configDir, "upstream/akooperator/akodeploymentconfig.yaml")
		file02AVIInfrasettings    = filepath.Join(configDir, "upstream/akooperator/aviinfrasettings.yaml")
		file03Deployment          = filepath.Join(configDir, "upstream/akooperator/deployment.yaml")
		file04Secret              = filepath.Join(configDir, "upstream/akooperator/secret.yaml")
		file05Static              = filepath.Join(configDir, "upstream/akooperator/static.yaml")

		fileOverlayAKODeploymentConfig = filepath.Join(configDir, "overlays/overlay-akodeploymentconfig.yaml")
		fileOverlayDeployment          = filepath.Join(configDir, "overlays/overlay-deployment.yaml")
		fileOverlaySecret              = filepath.Join(configDir, "overlays/overlay-secret.yaml")

		fileValuesYaml  = filepath.Join(configDir, "values.yaml")
		fileValuesStar  = filepath.Join(configDir, "values.star")
		filesKappConfig = filepath.Join(configDir, "kapp-config.yaml")
	)

	const (
		defaultValues = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
akoOperator:
  avi_enable: true
  namespace: tkg-system-networking
  cluster_name: ""
`
	)

	JustBeforeEach(func() {
		filePaths = []string{
			file01AKODeploymentConfig,
			file02AVIInfrasettings,
			file03Deployment,
			file04Secret,
			file05Static,
			fileOverlayAKODeploymentConfig,
			fileOverlayDeployment,
			fileOverlaySecret,
			fileValuesYaml,
			fileValuesStar,
			filesKappConfig,
		}
		output, yttRenderErr = ytt.RenderYTTTemplate(ytt.CommandOptions{}, filePaths, strings.NewReader(values))
	})

	Context("render ytt templates for ako-operator", func() {
		When("avi is enabled", func() {
			BeforeEach(func() {
				values = defaultValues
			})

			It("should render ako-operator related objects", func() {
				Expect(yttRenderErr).NotTo(HaveOccurred())
				AKODeploymentConfigExists(output)
				AviInfraSettingCRDExists(output)
				AKOODeploymentExists(output)
				AVISecretsExist(output)
			})
		})
	})

	Context("kapp rebase rule for avi credentials and CA", func() {
		When("avi_admin_credential_name and avi_ca_name is specified", func() {
			BeforeEach(func() {
				values = defaultValues + "\n" + "  config:\n" + "    avi_admin_credential_name: avi-controller-credentials\n" + "    avi_ca_name: avi-controller-ca"
			})
			It("should have rebase rule to skip their reconciliation in kapp-config", func() {
				cms, err := matchers.FindDocsMatchingYAMLPath(output, map[string]string{
					"$.kind":          "ConfigMap",
					"$.metadata.name": "ako-operator-kapp-config",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(cms).To(HaveLen(1))

				Expect(cms[0]).To(ContainSubstring("[data, username]"))
				Expect(cms[0]).To(ContainSubstring("[data, password]"))
				Expect(cms[0]).To(ContainSubstring("avi-controller-credentials"))
				Expect(cms[0]).To(ContainSubstring("[data, certificateAuthorityData]"))
				Expect(cms[0]).To(ContainSubstring("avi-controller-ca"))
			})
		})
	})

	Context("avi_nsxt_t1_lr field", func() {
		When("avi_nsxt_t1_lr field is specified", func() {
			BeforeEach(func() {
				values = defaultValues + "\n" + "  config:\n" + "    avi_nsxt_t1_lr: test_lr"
			})
			It("should have it set in AKODeploymentConfig spec", func() {
				adcs, err := matchers.FindDocsMatchingYAMLPath(output, map[string]string{
					"$.kind": "AKODeploymentConfig",
				})
				Expect(err).NotTo(HaveOccurred())
				for _, adc := range adcs {
					Expect(adc).To(ContainSubstring("nsxtT1LR: test_lr"))
				}
			})
		})

		When("avi_nsxt_t1_lr field is not specified", func() {
			BeforeEach(func() {
				values = defaultValues
			})
			It("should have it set in AKODeploymentConfig spec", func() {
				adcs, err := matchers.FindDocsMatchingYAMLPath(output, map[string]string{
					"$.kind": "AKODeploymentConfig",
				})
				Expect(err).NotTo(HaveOccurred())
				for _, adc := range adcs {
					Expect(adc).ToNot(ContainSubstring("nsxtT1LR"))
				}
			})
		})
	})
})
