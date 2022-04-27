// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package antrea_test

import (
	"path/filepath"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
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

		configDir                     = filepath.Join(repo.RootDir(), "addons/packages/antrea/1.2.3/bundle/config")
		fileAntreaYaml                = filepath.Join(configDir, "upstream/antrea.yaml")
		fileAntreaOverlayYaml         = filepath.Join(configDir, "overlay/antrea-overlay.yaml")
		fileUpdateStrategyOverlayYaml = filepath.Join(configDir, "overlay/update-strategy-overlay.yaml")
		fileValuesYaml                = filepath.Join(configDir, "values.yaml")
		fileValuesStar                = filepath.Join(configDir, "values.star")
	)
	BeforeEach(func() {
		values = ""
	})

	JustBeforeEach(func() {
		filePaths = []string{fileAntreaYaml, fileAntreaOverlayYaml, fileUpdateStrategyOverlayYaml, fileValuesYaml, fileValuesStar}
		output, err = ytt.RenderYTTTemplate(ytt.CommandOptions{}, filePaths, strings.NewReader(values))
		Expect(err).NotTo(HaveOccurred())
	})

	Context("default configuration", func() {
		It("renders a ConfigMap with a default ipam configuration", func() {
			configMap := findConfigMapByName(unmarshalConfigMaps(output), "antrea-config-822fk25299")
			Expect(configMap).NotTo(BeNil())
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

	Context("antrea-agent with serviceCIDRv6 configuration", func() {
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
			configMap := findConfigMapByName(unmarshalConfigMaps(output), "antrea-config-822fk25299")
			Expect(configMap).NotTo(BeNil())
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

	Context("antrea-agent-tweaker with default configuration", func() {
		It("render disabled UDP tunnel offload feature", func() {
			configMap := findConfigMapByName(unmarshalConfigMaps(output), "antrea-agent-tweaker-g56hc6fh8t")
			Expect(configMap).NotTo(BeNil())
			Expect(configMap.Data["antrea-agent-tweaker.conf"]).To(MatchYAML(`---
disableUdpTunnelOffload: false
`))
		})
	})

	Context("antrea-agent-tweaker with enabled UDP tunnel configuration", func() {
		BeforeEach(func() {
			values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
antrea:
  config:
    disableUdpTunnelOffload: true
`
		})

		It("render enabled UDP tunnel offload feature", func() {
			configMap := findConfigMapByName(unmarshalConfigMaps(output), "antrea-agent-tweaker-g56hc6fh8t")
			Expect(configMap).NotTo(BeNil())
			Expect(configMap.Data["antrea-agent-tweaker.conf"]).To(MatchYAML(`---
disableUdpTunnelOffload: true
`))
		})
	})

	Context("configures nodeSelector and updateStrategy", func() {
		BeforeEach(func() {
			values = `#@data/values
---
nodeSelector:
  tanzuKubernetesRelease: 1.22.3
deployment:
  updateStrategy: RollingUpdate
  rollingUpdate:
    maxUnavailable: 0
    maxSurge: 1
daemonset:
  updateStrategy: OnDelete
`
		})
		It("renders the DaemonSet and Deployment with desired nodeSelector and updateStrategy", func() {
			Expect(err).NotTo(HaveOccurred())
			daemonSet := parseDaemonSet(output)
			deployment := parseDeployment(output)
			Expect(deployment.Spec.Template.Spec.NodeSelector).ToNot(BeNil())
			Expect(deployment.Spec.Template.Spec.NodeSelector["tanzuKubernetesRelease"]).To(Equal("1.22.3"))
			Expect(deployment.Spec.Strategy.Type).To(Equal(appsv1.RollingUpdateDeploymentStrategyType))
			Expect(deployment.Spec.Strategy.RollingUpdate).ToNot(BeNil())
			Expect(deployment.Spec.Strategy.RollingUpdate.MaxUnavailable.IntVal).To(Equal(int32(0)))
			Expect(deployment.Spec.Strategy.RollingUpdate.MaxSurge.IntVal).To(Equal(int32(1)))
			Expect(daemonSet.Spec.UpdateStrategy.Type).To(Equal(appsv1.OnDeleteDaemonSetStrategyType))
		})
	})

})

func findConfigMapByName(cms []corev1.ConfigMap, name string) *corev1.ConfigMap {
	for _, cm := range cms {
		if cm.Name == name {
			return &cm
		}
	}
	return nil
}

func unmarshalConfigMaps(output string) []corev1.ConfigMap {
	docs := findDocsWithString(output, "kind: ConfigMap")
	cms := make([]corev1.ConfigMap, len(docs))
	for i, doc := range docs {
		var cm corev1.ConfigMap
		err := yaml.Unmarshal([]byte(doc), &cm)
		Expect(err).NotTo(HaveOccurred())
		cms[i] = cm
	}
	return cms
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

func parseDaemonSet(output string) appsv1.DaemonSet {
	daemonSetDocIndex := 47
	daemonSetDoc := strings.Split(output, "---")[daemonSetDocIndex]
	var daemonSet appsv1.DaemonSet
	err := yaml.Unmarshal([]byte(daemonSetDoc), &daemonSet)
	Expect(err).NotTo(HaveOccurred())
	return daemonSet
}
func parseDeployment(output string) appsv1.Deployment {
	deploymentDocIndex := 39
	deploymentDoc := strings.Split(output, "---")[deploymentDocIndex]
	var deployment appsv1.Deployment
	err := yaml.Unmarshal([]byte(deploymentDoc), &deployment)
	Expect(err).NotTo(HaveOccurred())
	return deployment
}
