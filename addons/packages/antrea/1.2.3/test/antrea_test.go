// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package antrea_test

import (
	"fmt"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/repo"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/ytt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Antrea Ytt Templates", func() {
	var (
		filePaths []string
		values    string
		output    string
		err       error

		configDir             = filepath.Join(repo.RootDir(), "addons/packages/antrea/1.2.3/bundle/config")
		fileAntreaYaml        = filepath.Join(configDir, "upstream/antrea.yaml")
		fileAntreaOverlayYaml = filepath.Join(configDir, "overlay/antrea_overlay.yaml")
		fileValuesYaml        = filepath.Join(configDir, "values.yaml")
		fileValuesStar        = filepath.Join(configDir, "values.star")
	)
	BeforeEach(func() {
		values = ""
	})

	JustBeforeEach(func() {
		filePaths = []string{fileAntreaYaml, fileAntreaOverlayYaml, fileValuesYaml, fileValuesStar}
		output, err = ytt.RenderYTTTemplate(ytt.CommandOptions{}, filePaths, strings.NewReader(values))
	})

	Context("default configuration", func() {
		It("renders a ConfigMap with a default ipam configuration", func() {
			Expect(err).NotTo(HaveOccurred())
			configMap := parseConfigMap(output)
			Expect(configMap.Data["antrea-agent.conf"]).To(MatchYAML(`---
featureGates:
  AntreaProxy: true
  EndpointSlice: false
  Traceflow: true
  NodePortLocal: false
  AntreaPolicy: true
  FlowExporter: false
  NetworkPolicyStats: false
  Egress: false
trafficEncapMode: encap
noSNAT: false
serviceCIDR: 10.96.0.0/12
tlsCipherSuites: TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384
`))
		})
	})

	Context("serviceCIDRv6 configuration", func() {
		BeforeEach(func() {
			values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
antrea:
  config:
    serviceCIDRv6: "[fe80::1]/64"
`
		})

		It("renders a ConfigMap with IPv6 ipam configuration", func() {
			Expect(err).NotTo(HaveOccurred())
			configMap := parseConfigMap(output)
			Expect(configMap.Data["antrea-agent.conf"]).To(MatchYAML(`---
featureGates:
  AntreaProxy: true
  EndpointSlice: false
  Traceflow: true
  NodePortLocal: false
  AntreaPolicy: true
  FlowExporter: false
  NetworkPolicyStats: false
  Egress: false
trafficEncapMode: encap
noSNAT: false
serviceCIDR: 10.96.0.0/12
serviceCIDRv6: "[fe80::1]/64"
tlsCipherSuites: TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384
`))
		})
	})
})

func parseConfigMap(output string) corev1.ConfigMap {
	configMapDocStr := findDocWithString(output, "kind: ConfigMap")
	var configMap corev1.ConfigMap
	err := yaml.Unmarshal([]byte(configMapDocStr), &configMap)
	Expect(err).NotTo(HaveOccurred())
	return configMap
}

func findDocWithString(output, selector string) string {
	docStrs := strings.Split(output, "---")
	for _, docStr := range docStrs {
		if strings.Contains(docStr, selector) {
			return docStr
		}
	}
	Fail(fmt.Sprintf("Expected to find doc containing substring %q, but none was found", selector))
	return "this won't be returned because it would have failed, but compilers gotta be compilers"
}
