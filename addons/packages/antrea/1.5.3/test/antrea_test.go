// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package antrea_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	goyaml "gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/repo"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/ytt"
)

const portRange = "60000-61000"

var (
	configName    = "antrea-config-f5d8g47b88"
	configTweaker = "antrea-agent-tweaker-g56hc6fh8t"

	// data header overwritten
	dataHeader = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---`
)

var _ = Describe("Antrea YTT Templates", func() {
	var (
		filePaths []string
		values    string
		output    string
		err       error

		configDir                     = filepath.Join(repo.RootDir(), "addons/packages/antrea/1.5.3/bundle/config")
		fileAntreaYaml                = filepath.Join(configDir, "upstream/antrea.yaml")
		fileAntreaOverlayYaml         = filepath.Join(configDir, "overlay/antrea-overlay.yaml")
		fileAntreaStrategyOverlayYaml = filepath.Join(configDir, "overlay/update-strategy-overlay.yaml")
		fileValuesSchema              = filepath.Join(configDir, "schema.yaml")
		fileValuesYaml                = filepath.Join(configDir, "values.yaml")
		fileValuesStar                = filepath.Join(configDir, "values.star")
	)

	BeforeEach(func() {
		values = ""
	})

	JustBeforeEach(func() {
		filePaths = []string{fileValuesSchema, fileAntreaYaml, fileAntreaOverlayYaml, fileAntreaStrategyOverlayYaml, fileValuesYaml, fileValuesStar}
		output, err = ytt.RenderYTTTemplate(ytt.CommandOptions{}, filePaths, strings.NewReader(values))
	})

	Context("antrea components with default configuration", func() {
		It("renders multiple uonfigMap with a default IPAM configuration", func() {
			Expect(err).NotTo(HaveOccurred())
			configMap := findConfigMapByName(unmarshalConfigMaps(output), configName)
			Expect(configMap).NotTo(BeNil())
			Expect(configMap.Data["antrea-agent.conf"]).To(MatchYAML(`---
antreaProxy:
  nodePortAddresses: []
  proxyAll: false
  proxyLoadBalancerIPs: false
  skipServices: []
egress:
  exceptCIDRs: []
featureGates:
  AntreaIPAM: false
  AntreaPolicy: true
  AntreaProxy: true
  Egress: true
  EndpointSlice: false
  FlowExporter: false
  Multicast: false
  NetworkPolicyStats: false
  NodePortLocal: true
  ServiceExternalIP: false
  Traceflow: true
noSNAT: false
nodePortLocal:
  enable: true
  portRange: 61000-62000
serviceCIDR: 10.96.0.0/12
tlsCipherSuites: TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384
trafficEncapMode: encap
trafficEncryptionMode: none
tunnelType: geneve
wireGuard:
  port: 51820
`))

			Expect(configMap.Data["antrea-controller.conf"]).To(MatchYAML(`---
featureGates:
  AntreaIPAM: false
  AntreaPolicy: true
  Egress: true
  NetworkPolicyStats: false
  ServiceExternalIP: false
  Traceflow: true
nodeIPAM: null
tlsCipherSuites: TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384
`))
		})
	})

	Context("antrea-agent with serviceCIDRv6 configuration", func() {
		BeforeEach(func() {
			values, err = marshalAntreaConfig(fileValuesYaml, func(config *AntreaConfig) {
				config.Antrea.Config.ServiceCIDRv6 = "[fe80::1]/64"
			})
			Expect(err).NotTo(HaveOccurred())
		})
		It("renders a ConfigMap with IPv6 IPAM configuration", func() {
			Expect(err).NotTo(HaveOccurred())
			configMap := findConfigMapByName(unmarshalConfigMaps(output), configName)
			Expect(configMap).NotTo(BeNil())
			Expect(configMap.Data["antrea-agent.conf"]).To(ContainSubstring(`serviceCIDRv6: '[fe80::1]/64'`))
		})
	})

	Context("antrea-agent-tweaker with default configuration", func() {
		It("render disabled UDP tunnel offload feature", func() {
			Expect(err).NotTo(HaveOccurred())
			configMap := findConfigMapByName(unmarshalConfigMaps(output), configTweaker)
			Expect(configMap).NotTo(BeNil())
			Expect(configMap.Data["antrea-agent-tweaker.conf"]).To(ContainSubstring("disableUdpTunnelOffload: false"))
		})
	})

	Context("antrea-agent-tweaker with enabled UDP tunnel configuration", func() {
		BeforeEach(func() {
			values, err = marshalAntreaConfig(fileValuesYaml, func(config *AntreaConfig) {
				config.Antrea.Config.DisableUDPTunnelOffload = true
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("render enabled UDP tunnel offload feature", func() {
			Expect(err).NotTo(HaveOccurred())
			configMap := findConfigMapByName(unmarshalConfigMaps(output), configTweaker)
			Expect(configMap).NotTo(BeNil())
			Expect(configMap.Data["antrea-agent-tweaker.conf"]).To(ContainSubstring("disableUdpTunnelOffload: true"))
		})
	})

	Context("antrea configuration has wrong fields", func() {
		BeforeEach(func() {
			values = `#@data/values
	---
	antrea:
	  config:
	invalid: "option"`
		})

		It("fails to generate manifests", func() {
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Configuring Egress", func() {
		Context("without feature gate", func() {
			BeforeEach(func() {
				values, err = marshalAntreaConfig(fileValuesYaml, func(config *AntreaConfig) {
					config.Antrea.Config.FeatureGates.Egress = false
					config.Antrea.Config.Egress.ExceptCIDRs = []string{"10.0.0.0/16"}
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("ignores the configuration", func() {
				Expect(err).NotTo(HaveOccurred())
				configMap := findConfigMapByName(unmarshalConfigMaps(output), configName)
				Expect(configMap).NotTo(BeNil())
				Expect(configMap.Data["antrea-agent.conf"]).To(ContainSubstring("Egress: false"))
				Expect(configMap.Data["antrea-agent.conf"]).ToNot(ContainSubstring("exceptCIDR"))
			})
		})

		Context("with the feature gate enabled", func() {
			BeforeEach(func() {
				values, err = marshalAntreaConfig(fileValuesYaml, func(config *AntreaConfig) {
					config.Antrea.Config.FeatureGates.Egress = true
					config.Antrea.Config.Egress.ExceptCIDRs = []string{"10.0.0.0/16"}
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("settings are configured", func() {
				Expect(err).NotTo(HaveOccurred())
				configMap := findConfigMapByName(unmarshalConfigMaps(output), configName)
				Expect(configMap).NotTo(BeNil())
				for _, s := range []string{"Egress: true", "10.0.0.0/16"} {
					Expect(configMap.Data["antrea-agent.conf"]).To(ContainSubstring(s))
				}
			})
		})
	})

	Describe("Configuring NodePortLocal", func() {
		Context("without feature gate", func() {
			BeforeEach(func() {
				values, err = marshalAntreaConfig(fileValuesYaml, func(config *AntreaConfig) {
					config.Antrea.Config.FeatureGates.NodePortLocal = false
					config.Antrea.Config.NodePortLocal.Enabled = true
					config.Antrea.Config.NodePortLocal.PortRange = portRange
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("ignores the configuration", func() {
				Expect(err).NotTo(HaveOccurred())
				configMap := findConfigMapByName(unmarshalConfigMaps(output), configName)
				Expect(configMap).NotTo(BeNil())
				Expect(configMap.Data["antrea-agent.conf"]).To(ContainSubstring("NodePortLocal: false"))
				Expect(configMap.Data["antrea-agent.conf"]).ToNot(ContainSubstring("portRange"))
			})
		})

		Context("with the feature gate enabled", func() {
			portRange := "60000-61000"
			BeforeEach(func() {
				values, err = marshalAntreaConfig(fileValuesYaml, func(config *AntreaConfig) {
					config.Antrea.Config.FeatureGates.NodePortLocal = true
					config.Antrea.Config.NodePortLocal.PortRange = portRange
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("settings are configured", func() {
				Expect(err).NotTo(HaveOccurred())
				configMap := findConfigMapByName(unmarshalConfigMaps(output), configName)
				Expect(configMap).NotTo(BeNil())

				for _, s := range []string{"NodePortLocal: true", portRange} {
					Expect(configMap.Data["antrea-agent.conf"]).To(ContainSubstring(s))
				}
			})
		})
	})

	Describe("Configuring AntreaProxy", func() {
		Context("without feature gate", func() {
			BeforeEach(func() {
				values, err = marshalAntreaConfig(fileValuesYaml, func(config *AntreaConfig) {
					config.Antrea.Config.FeatureGates.AntreaProxy = false
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("ignores the configuration", func() {
				Expect(err).NotTo(HaveOccurred())
				configMap := findConfigMapByName(unmarshalConfigMaps(output), configName)
				Expect(configMap).NotTo(BeNil())
				Expect(configMap.Data["antrea-agent.conf"]).To(ContainSubstring("AntreaProxy: false"))
				for _, s := range []string{"proxyAll", "nodePortAddresses", "skipServices", "proxyLoadBalancersIPs"} {
					Expect(configMap.Data["antrea-agent.conf"]).ToNot(ContainSubstring(s))
				}
			})
		})

		Context("with the feature gate enabled", func() {
			BeforeEach(func() {
				values, err = marshalAntreaConfig(fileValuesYaml, func(config *AntreaConfig) {
					config.Antrea.Config.FeatureGates.AntreaProxy = true
					config.Antrea.Config.AntreaProxy.ProxyAll = true
					config.Antrea.Config.AntreaProxy.NodePortAddresses = []string{"10.0.0.0/24"}
					config.Antrea.Config.AntreaProxy.ProxyLoadBalancerIPS = true
					config.Antrea.Config.AntreaProxy.SkipServices = []string{"kube-system/kube-dns"}
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("settings are configured", func() {
				Expect(err).NotTo(HaveOccurred())
				configMap := findConfigMapByName(unmarshalConfigMaps(output), configName)
				Expect(configMap).NotTo(BeNil())

				for _, s := range []string{"AntreaProxy: true", "nodePortAddresses", "skipServices", "proxyLoadBalancerIPs: true"} {
					Expect(configMap.Data["antrea-agent.conf"]).To(ContainSubstring(s))
				}
			})
		})
	})

	Describe("Configuring FlowExporter", func() {
		Context("with the feature gate enabled", func() {
			seconds := "10s"
			BeforeEach(func() {
				values, err = marshalAntreaConfig(fileValuesYaml, func(config *AntreaConfig) {
					config.Antrea.Config.FeatureGates.FlowExporter = true
					config.Antrea.Config.FlowExporter.Address = "0.0.0.0"
					config.Antrea.Config.FlowExporter.PollInterval = seconds
					config.Antrea.Config.FlowExporter.ExportTimeout = seconds
					config.Antrea.Config.FlowExporter.IdleExportTimeout = seconds
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("settings are configured", func() {
				Expect(err).NotTo(HaveOccurred())
				configMap := findConfigMapByName(unmarshalConfigMaps(output), configName)
				Expect(configMap).NotTo(BeNil())
				for _, s := range []string{"flowPollInterval: " + seconds, "flowCollectorAddr: 0.0.0.0", "idleFlowExportTimeout: " + seconds, "activeFlowExportTimeout: " + seconds} {
					Expect(configMap.Data["antrea-agent.conf"]).To(ContainSubstring(s))
				}
			})
		})
	})

	Describe("Changing root settings", func() {
		Context("should be allowed", func() {
			BeforeEach(func() {
				values, err = marshalAntreaConfig(fileValuesYaml, func(config *AntreaConfig) {
					config.Antrea.Config.TransportInterface = "eth0"
					config.Antrea.Config.TransportInterfaceCIDRs = []string{"10.0.0.0/24"}
					config.Antrea.Config.MulticastInterfaces = []string{"eth0"}
					config.Antrea.Config.WireGuard.Port = 51821
					config.Antrea.Config.KubeAPIServerOverride = "10.0.0.1"
					config.Antrea.Config.EnableUsageReporting = true
					config.Antrea.Config.TunnelType = "vxlan"
					config.Antrea.Config.TrafficEncryptionMode = "ipsec"
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("and args must be rendered", func() {
				Expect(err).NotTo(HaveOccurred())
				configMap := findConfigMapByName(unmarshalConfigMaps(output), configName)
				Expect(configMap).NotTo(BeNil())
				for _, s := range []string{
					"transportInterface: eth0",
					"multicastInterfaces",
					"transportInterfaceCIDRs",
					"port: 51821",
					"kubeAPIServerOverride: 10.0.0.1",
					"tunnelType: vxlan",
					"trafficEncryptionMode: ipsec",
				} {
					Expect(configMap.Data["antrea-agent.conf"]).To(ContainSubstring(s))
				}
				Expect(configMap.Data["antrea-controller.conf"]).To(ContainSubstring("enableUsageReporting"))
			})
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

func marshalAntreaConfig(valuesFile string, settingFunc func(config *AntreaConfig)) (string, error) {
	var (
		err    error
		config *AntreaConfig
	)
	if config, err = loadAntreaConfig(valuesFile); err != nil {
		return "", err
	}

	// Overwrite values in the config pointer
	settingFunc(config)

	content, err := goyaml.Marshal(config)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s\n%s", dataHeader, content), nil
}

// loadAntreaConfig unmarshal the configuration file into AntreaConfig
func loadAntreaConfig(configFile string) (*AntreaConfig, error) {
	var (
		err     error
		content []byte
		config  = AntreaConfig{}
	)

	if content, err = os.ReadFile(configFile); err != nil {
		return nil, err
	}
	if err := goyaml.Unmarshal(content, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

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
	deploymentDocIndex := 40
	deploymentDoc := strings.Split(output, "---") //[deploymentDocIndex]
	var deployment appsv1.Deployment
	err := yaml.Unmarshal([]byte(deploymentDoc[deploymentDocIndex]), &deployment)
	Expect(err).NotTo(HaveOccurred())
	return deployment
}
