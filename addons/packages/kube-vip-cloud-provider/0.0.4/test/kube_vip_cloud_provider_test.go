// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kube_vip_cloud_provider_test

import (
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/repo"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/ytt"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/providers/tests/unit/matchers"
)

var _ = Describe("kube-vip CloudProvider Ytt Templates", func() {
	var (
		filePaths    []string
		values       string
		output       string
		yttRenderErr error

		configDir = filepath.Join(repo.RootDir(), "addons/packages/kube-vip-cloud-provider/0.0.4/bundle/config")

		file01rbac                  = filepath.Join(configDir, "upstream/kube-vip-cloud-provider/rbac.yaml")
		file02config                = filepath.Join(configDir, "upstream/kube-vip-cloud-provider/configmap.yaml")
		file03deployment            = filepath.Join(configDir, "upstream/kube-vip-cloud-provider/deployment.yaml")
		fileOverlayUpdateConfig     = filepath.Join(configDir, "overlays/overlay-configmap.yaml")
		fileOverlayUpdateDeployment = filepath.Join(configDir, "overlays/overlay-deployment.yaml")

		fileValuesYaml = filepath.Join(configDir, "values.yaml")
		fileValuesStar = filepath.Join(configDir, "values.star")
	)

	JustBeforeEach(func() {
		filePaths = []string{
			file01rbac,
			file02config,
			file03deployment,
			fileOverlayUpdateConfig,
			fileOverlayUpdateDeployment,
			fileValuesYaml,
			fileValuesStar,
		}
		output, yttRenderErr = ytt.RenderYTTTemplate(ytt.CommandOptions{}, filePaths, strings.NewReader(values))
	})

	Context("Config L4 LB", func() {
		BeforeEach(func() {
			values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
kubevipCloudProvider:
`
		})

		When("ConfigMap has all setttings", func() {
			BeforeEach(func() {
				values = values + `
  loadbalancerCIDRs: "10.0.0.1/24"
  loadbalancerIPRanges: "10.0.0.1-10.0.0.2"

`
			})

			It("should render configmap for kube-vip-cloud-provider", func() {
				Expect(yttRenderErr).NotTo(HaveOccurred())
				docs, err := matchers.FindDocsMatchingYAMLPath(output, map[string]string{
					"$.kind":               "ConfigMap",
					"$.metadata.name":      "kubevip",
					"$.metadata.namespace": "kube-system",
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(docs).To(HaveLen(1))
				cfg := unmarshalConfigMap(output)
				Expect(getConfigKeyFromData(cfg, "cidr-global")).To(Equal("10.0.0.1/24"))
				Expect(getConfigKeyFromData(cfg, "range-global")).To(Equal("10.0.0.1-10.0.0.2"))
			})
		})

		When("Config doesn't have either loadbalancerCIDRs and loadbalancerIPRanges", func() {
			It("should error out", func() {
				Expect(yttRenderErr).To(HaveOccurred())
			})
		})
	})
})

func getConfigKeyFromData(cfg corev1.ConfigMap, key string) string {
	cfgValue, exists := cfg.Data[key]
	Expect(exists).To(Equal(true))
	return cfgValue
}

func unmarshalConfigMap(output string) corev1.ConfigMap {
	docs := findDocsWithString(output, "kind: ConfigMap")
	var cm corev1.ConfigMap
	err := yaml.Unmarshal([]byte(docs[0]), &cm)
	Expect(err).NotTo(HaveOccurred())
	return cm
}

func findDocsWithString(output, selector string) []string {
	var docs []string
	for _, doc := range strings.Split(output, "---") {
		if strings.Contains(doc, selector) {
			docs = append(docs, doc)
		}
	}
	return docs
}
