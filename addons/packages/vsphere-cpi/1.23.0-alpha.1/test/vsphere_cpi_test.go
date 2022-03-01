// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package vspherecpi_test

import (
	"fmt"
	"path/filepath"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/cloud-provider-vsphere/pkg/cloudprovider/vsphere/config"
	"sigs.k8s.io/yaml"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/repo"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/ytt"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/providers/tests/unit/matchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("vSphere CPI Ytt Templates", func() {
	var (
		filePaths    []string
		values       string
		output       string
		yttRenderErr error

		configDir = filepath.Join(repo.RootDir(), "addons/packages/vsphere-cpi/1.23.0-alpha.1/bundle/config")

		// vsphere-cpi mode
		file01rbac                 = filepath.Join(configDir, "upstream/vsphere-cpi/01-rbac.yaml")
		file02config               = filepath.Join(configDir, "upstream/vsphere-cpi/02-config.yaml")
		file03secret               = filepath.Join(configDir, "upstream/vsphere-cpi/03-secret.yaml")
		file04daemonset            = filepath.Join(configDir, "upstream/vsphere-cpi/04-daemonset.yaml")
		fileOverlayUpdateConfig    = filepath.Join(configDir, "overlays/vsphere-cpi/update-config.yaml")
		fileOverlayAddSecret       = filepath.Join(configDir, "overlays/vsphere-cpi/add-secret.yaml")
		fileOverlayUpdateSecret    = filepath.Join(configDir, "overlays/vsphere-cpi/update-secret.yaml")
		fileOverlayUpdateDaemonset = filepath.Join(configDir, "overlays/vsphere-cpi/update-daemonset.yaml")

		// vsphere-paravirtual-cpi mode
		file01ParavirtualRbac       = filepath.Join(configDir, "upstream/vsphere-paravirtual-cpi/01-rbac.yaml")
		file02ParavirtualConfig     = filepath.Join(configDir, "upstream/vsphere-paravirtual-cpi/02-config.yaml")
		file03ParavirtualDeployment = filepath.Join(configDir, "upstream/vsphere-paravirtual-cpi/03-deployment.yaml")

		fileOverlayParavirtualUpdateDeployment = filepath.Join(configDir, "overlays/vsphere-paravirtual-cpi/update-deployment.yaml")

		fileValuesYaml        = filepath.Join(configDir, "values.yaml")
		fileValuesStar        = filepath.Join(configDir, "values.star")
		fileVsphereconfLibTxt = filepath.Join(configDir, "vsphereconf.lib.txt")
	)

	const (
		defaultValues = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  server: fake-server.com
  datacenter: dc0
  username: my-user
  password: my-password
  insecureFlag: True`
		defaultParavirtualValues = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  mode: vsphereParavirtualCPI
  clusterName: "tkg-cluster"
  clusterUID: "57341fa8-0983-472f-b744-00cf724dd307"
  supervisorMasterEndpointIP: "192.168.123.2"
  supervisorMasterPort: "6443"
`
	)

	JustBeforeEach(func() {
		filePaths = []string{
			file01rbac,
			file02config,
			file03secret,
			file04daemonset,
			fileOverlayUpdateConfig,
			fileOverlayAddSecret,
			fileOverlayUpdateSecret,
			fileOverlayUpdateDaemonset,
			file01ParavirtualRbac,
			file02ParavirtualConfig,
			file03ParavirtualDeployment,
			fileOverlayParavirtualUpdateDeployment,
			fileValuesYaml,
			fileValuesStar,
			fileVsphereconfLibTxt,
		}
		output, yttRenderErr = ytt.RenderYTTTemplate(ytt.CommandOptions{}, filePaths, strings.NewReader(values))
	})

	Context("vSphere mode", func() {
		BeforeEach(func() {
			values = defaultValues
		})

		When("mode is empty", func() {
			It("should render resources for vsphereCPI", func() {
				Expect(yttRenderErr).NotTo(HaveOccurred())
				vsphereCPIDocuments(output)
			})
		})

		When("mode is vsphereCPI", func() {
			BeforeEach(func() {
				values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  mode: vsphereCPI
  server: fake-server.com
  datacenter: dc0
  username: my-user
  password: my-password
  insecureFlag: True`
			})
			It("should render resources for vsphereCPI", func() {
				Expect(yttRenderErr).NotTo(HaveOccurred())
				vsphereCPIDocuments(output)
			})
		})

		When("mode is vsphereParavirtualCPI", func() {
			BeforeEach(func() {
				values = defaultParavirtualValues
			})
			It("should render resources for vsphereParavirtualCPI", func() {
				Expect(yttRenderErr).NotTo(HaveOccurred())
				vsphereParavirtualDocuments(output)
			})

			Context("pod routing by Antrea NSX", func() {
				When("antreaNSXPodRoutingEnabled is true", func() {
					BeforeEach(func() {
						values = defaultParavirtualValues + "\n  antreaNSXPodRoutingEnabled: true"
					})

					It("should add arguments to the ccm deployment container", func() {
						docs, err := matchers.FindDocsMatchingYAMLPath(output, map[string]string{
							"$.kind":               "Deployment",
							"$.metadata.name":      "guest-cluster-cloud-provider",
							"$.metadata.namespace": "vmware-system-cloud-provider",
						})
						Expect(err).ToNot(HaveOccurred())
						Expect(docs).To(HaveLen(1))
						deployment := docs[0]
						Expect(deployment).To(ContainSubstring("--controllers=route"))
						Expect(deployment).To(ContainSubstring("--configure-cloud-routes=true"))
						Expect(deployment).To(ContainSubstring("--allocate-node-cidrs=true"))
						Expect(deployment).To(ContainSubstring("--cluster-cidr=0.0.0.0/0"))
					})
				})
			})
		})

		When("mode is nonsense", func() {
			BeforeEach(func() {
				values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  mode: vsphereCPINonSense
  server: fake-server.com
  datacenter: dc0
  username: my-user
  password: my-password
  insecureFlag: True`
			})
			It("should throw render error", func() {
				Expect(yttRenderErr).To(HaveOccurred())
				Expect(yttRenderErr.Error()).To(ContainSubstring("vsphereCPI mode should be either vsphereCPI or vsphereParavirtualCPI"))
			})
		})
	})

	Context("ccm-cloud-config owner reference", func() {
		When("mode is vsphereParavirtualCPI and tkg cluster values are provided", func() {
			BeforeEach(func() {
				values = defaultParavirtualValues
			})

			It("should render correct owner reference for the config map", func() {
				configMaps := unmarshalConfigMaps(output)
				Expect(configMaps).NotTo(BeEmpty())
				ownerRefConfigMap := findConfigMapByName(configMaps, "ccm-owner-reference")
				ownerRef, exists := ownerRefConfigMap.Data["owner-reference"]
				Expect(exists).To(BeTrue())
				Expect(ownerRef).To(Equal("{\"apiVersion\": \"cluster.x-k8s.io/v1beta1\",\n\"kind\": \"Cluster\",\n\"name\": \"tkg-cluster\",\n\"uid\": \"57341fa8-0983-472f-b744-00cf724dd307\"}\n"))
			})
		})
	})

	Context("guest-cluster-cloud-provider deployment", func() {
		BeforeEach(func() {
			values = defaultParavirtualValues
		})

		When("mode is vsphereParavirtualCPI and supervisor api-server is provided", func() {
			It("renders correct supervisor api-server info", func() {
				deployments := unmarshalDeployment(output)
				Expect(deployments).NotTo(BeEmpty())
				deployment := findDeploymentByName(deployments, "guest-cluster-cloud-provider")
				Expect(deployment.Spec.Template.Spec.Containers[0].Env[0].Value).To(Equal("192.168.123.2"))
				Expect(deployment.Spec.Template.Spec.Containers[0].Env[1].Value).To(Equal("6443"))
				Expect(deployment.Spec.Template.Spec.Containers[0].Args[3]).To(Equal("--cluster-name=tkg-cluster"))
			})
		})
	})

	Context("DaemonSet env vars", func() {
		BeforeEach(func() {
			values = defaultValues
		})

		It("renders a DaemonSet with ENABLE_ALPHA_DUAL_STACK env var feature flag", func() {
			Expect(yttRenderErr).NotTo(HaveOccurred())
			daemonSet := parseDaemonSet(output)
			containerEnvVars := daemonSet.Spec.Template.Spec.Containers[0].Env
			Expect(transformEnvVarsToMap(containerEnvVars)).To(HaveKeyWithValue("ENABLE_ALPHA_DUAL_STACK", "true"))
			Expect(transformEnvVarsToMap(containerEnvVars)).NotTo(HaveKey("HTTP_PROXY"))
			Expect(transformEnvVarsToMap(containerEnvVars)).NotTo(HaveKey("HTTPS_PROXY"))
			Expect(transformEnvVarsToMap(containerEnvVars)).NotTo(HaveKey("NO_PROXY"))
		})
		When("http proxy is configured", func() {
			BeforeEach(func() {
				values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  server: fake-server.com
  datacenter: dc0
  username: my-user
  password: my-password
  insecureFlag: True
  http_proxy: http://10.10.10.1:8080
  https_proxy: https://10.10.10.1:8080
  no_proxy: 10.10.10.2,example.com`
			})
			It("includes http proxy env vars", func() {

				Expect(yttRenderErr).NotTo(HaveOccurred())
				daemonSet := parseDaemonSet(output)
				containerEnvVars := daemonSet.Spec.Template.Spec.Containers[0].Env
				Expect(transformEnvVarsToMap(containerEnvVars)).To(HaveKeyWithValue("ENABLE_ALPHA_DUAL_STACK", "true"))
				Expect(transformEnvVarsToMap(containerEnvVars)).To(HaveKeyWithValue("HTTP_PROXY", "http://10.10.10.1:8080"))
				Expect(transformEnvVarsToMap(containerEnvVars)).To(HaveKeyWithValue("HTTPS_PROXY", "https://10.10.10.1:8080"))
				Expect(transformEnvVarsToMap(containerEnvVars)).To(HaveKeyWithValue("NO_PROXY", "10.10.10.2,example.com"))
			})
		})

		It("renders a data value hash starts with prefix and not parsable as an integer", func() {
			Expect(yttRenderErr).NotTo(HaveOccurred())
			daemonSet := parseDaemonSet(output)
			hash, exist := daemonSet.Spec.Template.Annotations["vsphere-cpi/data-values-hash"]
			Expect(exist).To(BeTrue())
			Expect(hash).Should(HavePrefix("h-"))
		})
	})

	Context("vsphere-cloud-config ConfigMap INI file", func() {
		Context("IPFamilies", func() {
			Context("unset", func() {
				BeforeEach(func() {
					values = defaultValues
				})

				It("defaults to ipv4", func() {
					Expect(yttRenderErr).NotTo(HaveOccurred())
					Expect(configuredIPFamily(output)).To(ConsistOf("ipv4"))
				})
			})
			Context("IPv4", func() {
				BeforeEach(func() {
					values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  server: fake-server.com
  datacenter: dc0
  username: my-user
  password: my-password
  insecureFlag: True
  ipFamily: ipv4`
				})

				It("configures ipv4", func() {
					Expect(yttRenderErr).NotTo(HaveOccurred())
					Expect(configuredIPFamily(output)).To(ConsistOf("ipv4"))
				})
			})
			Context("IPv6", func() {
				BeforeEach(func() {
					values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  server: fake-server.com
  datacenter: dc0
  username: my-user
  password: my-password
  insecureFlag: True
  ipFamily: ipv6`
				})

				It("configures ipv6", func() {
					Expect(yttRenderErr).NotTo(HaveOccurred())
					Expect(configuredIPFamily(output)).To(ConsistOf("ipv6"))
				})
			})
			Context("IPv4,IPv6", func() {
				BeforeEach(func() {
					values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  server: fake-server.com
  datacenter: dc0
  username: my-user
  password: my-password
  insecureFlag: True
  ipFamily: ipv4,ipv6`
				})

				It("configures ipv4,ipv6", func() {
					Expect(yttRenderErr).NotTo(HaveOccurred())
					Expect(configuredIPFamily(output)).To(Equal([]string{"ipv4", "ipv6"}))
				})
			})
			Context("IPv6,IPv4", func() {
				BeforeEach(func() {
					values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  server: fake-server.com
  datacenter: dc0
  username: my-user
  password: my-password
  insecureFlag: True
  ipFamily: ipv6,ipv4`
				})

				It("configures ipv6,ipv4", func() {
					Expect(yttRenderErr).NotTo(HaveOccurred())
					Expect(configuredIPFamily(output)).To(Equal([]string{"ipv6", "ipv4"}))
				})
			})
			Context("Invalid IP Family", func() {
				BeforeEach(func() {
					values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  server: fake-server.com
  datacenter: dc0
  username: my-user
  password: my-password
  insecureFlag: True
  ipFamily: ipv5`
				})

				It("errors when rendering", func() {
					Expect(yttRenderErr).To(MatchError(ContainSubstring("vsphereCPI ipFamily should be one of \"ipv4\", \"ipv6\", \"ipv4,ipv6\", or \"ipv6,ipv4\" if provided")))
				})
			})
		})

		Context("Node IP selectors", func() {
			When("vmExcludeInternalNetworkSubnetCidr is set", func() {
				BeforeEach(func() {
					values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  server: fake-server.com
  datacenter: dc0
  username: my-user
  password: my-password
  vmExcludeInternalNetworkSubnetCidr: 192.0.2.0/24,fe80::1/128
  insecureFlag: True`
				})

				It("correctly sets the value in the INI", func() {
					Expect(yttRenderErr).NotTo(HaveOccurred())
					Expect(nodesConfiguration(output).ExcludeInternalNetworkSubnetCIDR).To(Equal("192.0.2.0/24,fe80::1/128"))
				})
			})

			When("vmExcludeInternalNetworkSubnetCidr is not set", func() {
				BeforeEach(func() {
					values = defaultValues
				})

				It("does not set the value in the INI", func() {
					Expect(yttRenderErr).NotTo(HaveOccurred())
					Expect(nodesConfiguration(output).ExcludeInternalNetworkSubnetCIDR).To(BeEmpty())
				})
			})

			When("vmExcludeExternalNetworkSubnetCidr is set", func() {
				BeforeEach(func() {
					values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  server: fake-server.com
  datacenter: dc0
  username: my-user
  password: my-password
  vmExcludeExternalNetworkSubnetCidr: 192.0.3.0/24,fe80::2/128
  insecureFlag: True`
				})

				It("correctly sets the value in the INI", func() {
					Expect(yttRenderErr).NotTo(HaveOccurred())
					Expect(nodesConfiguration(output).ExcludeExternalNetworkSubnetCIDR).To(Equal("192.0.3.0/24,fe80::2/128"))
				})
			})

			When("vmExcludeExternalNetworkSubnetCidr is not set", func() {
				BeforeEach(func() {
					values = defaultValues
				})

				It("does not set the value in the INI", func() {
					Expect(yttRenderErr).NotTo(HaveOccurred())
					Expect(nodesConfiguration(output).ExcludeExternalNetworkSubnetCIDR).To(BeEmpty())
				})
			})

			When("vmInternalNetwork is set", func() {
				BeforeEach(func() {
					values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  server: fake-server.com
  datacenter: dc0
  username: my-user
  password: my-password
  vmInternalNetwork: meow
  insecureFlag: True`
				})

				It("correctly sets the value in the INI", func() {
					Expect(yttRenderErr).NotTo(HaveOccurred())
					Expect(nodesConfiguration(output).InternalVMNetworkName).To(Equal("meow"))
				})
			})

			When("vmInternalNetwork is not set", func() {
				BeforeEach(func() {
					values = defaultValues
				})

				It("does not set the value in the INI", func() {
					Expect(yttRenderErr).NotTo(HaveOccurred())
					Expect(nodesConfiguration(output).InternalVMNetworkName).To(BeEmpty())
				})
			})

			When("vmExternalNetwork is set", func() {
				BeforeEach(func() {
					values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  server: fake-server.com
  datacenter: dc0
  username: my-user
  password: my-password
  vmExternalNetwork: meow
  insecureFlag: True`
				})

				It("correctly sets the value in the INI", func() {
					Expect(yttRenderErr).NotTo(HaveOccurred())
					Expect(nodesConfiguration(output).ExternalVMNetworkName).To(Equal("meow"))
				})
			})

			When("vmExternalNetwork is not set", func() {
				BeforeEach(func() {
					values = defaultValues
				})

				It("does not set the value in the INI", func() {
					Expect(yttRenderErr).NotTo(HaveOccurred())
					Expect(nodesConfiguration(output).ExternalVMNetworkName).To(BeEmpty())
				})
			})
		})

		Context("Secure connection", func() {
			When("When TLSthumbprint is set and insecure is false", func() {
				BeforeEach(func() {
					values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  server: fake-server.com
  datacenter: dc0
  username: my-user
  password: my-password
  vmExternalNetwork: meow
  tlsThumbprint: fake-thumbprint
  insecureFlag: False`
				})

				It("The output should contain the thumbprint", func() {
					Expect(yttRenderErr).NotTo(HaveOccurred())
					Expect(tlsThumbprint(output)).NotTo(BeEmpty())
				})
			})

			When("When TLSthumbprint is set and insecure is true", func() {
				BeforeEach(func() {
					values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  server: fake-server.com
  datacenter: dc0
  username: my-user
  password: my-password
  vmExternalNetwork: meow
  tlsThumbprint: fake-thumbprint
  insecureFlag: True`
				})

				It("The output should not contain the thumbprint", func() {
					Expect(yttRenderErr).NotTo(HaveOccurred())
					Expect(tlsThumbprint(output)).To(BeEmpty())
				})
			})

			When("When TLSthumbprint is not set and insecure is false", func() {
				BeforeEach(func() {
					values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  server: fake-server.com
  datacenter: dc0
  username: my-user
  password: my-password
  vmExternalNetwork: meow
  tlsThumbprint: ""
  insecureFlag: False`
				})

				It("should error out", func() {
					Expect(yttRenderErr).To(HaveOccurred())
				})
			})

			When("When TLSthumbprint is not set and insecure is True", func() {
				BeforeEach(func() {
					values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  server: fake-server.com
  datacenter: dc0
  username: my-user
  password: my-password
  vmExternalNetwork: meow
  tlsThumbprint: ""
  insecureFlag: True`
				})

				It("The output should not contain the thumbprint, and succeed", func() {
					Expect(yttRenderErr).NotTo(HaveOccurred())
					Expect(tlsThumbprint(output)).To(BeEmpty())
				})
			})
		})
	})
})

func transformEnvVarsToMap(envVars []corev1.EnvVar) map[string]string {
	envVarMap := map[string]string{}
	for _, envVar := range envVars {
		if _, exists := envVarMap[envVar.Name]; !exists {
			envVarMap[envVar.Name] = envVar.Value
		} else {
			Fail(fmt.Sprintf("Unexpected duplicate EnvVar Name found: %s.", envVar.Name))
		}
	}
	return envVarMap
}

func parseDaemonSet(output string) appsv1.DaemonSet {
	daemonSetDocIndex := 6
	daemonSetDoc := strings.Split(output, "---")[daemonSetDocIndex]
	var daemonSet appsv1.DaemonSet
	err := yaml.Unmarshal([]byte(daemonSetDoc), &daemonSet)
	Expect(err).NotTo(HaveOccurred())
	return daemonSet
}

func findConfigMapByName(cms []corev1.ConfigMap, name string) *corev1.ConfigMap {
	for _, cm := range cms {
		if cm.Name == name {
			return &cm
		}
	}
	return nil
}

func findDeploymentByName(deploys []appsv1.Deployment, name string) *appsv1.Deployment {
	for _, deploy := range deploys {
		if deploy.Name == name {
			return &deploy
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

func unmarshalDeployment(output string) []appsv1.Deployment {
	docs := findDocsWithString(output, "kind: Deployment")
	deploys := make([]appsv1.Deployment, len(docs))
	for i, doc := range docs {
		var deploy appsv1.Deployment
		err := yaml.Unmarshal([]byte(doc), &deploy)
		Expect(err).NotTo(HaveOccurred())
		deploys[i] = deploy
	}
	return deploys
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

func cpiConfig(output string) *config.CPIConfig {
	configMaps := unmarshalConfigMaps(output)
	Expect(configMaps).NotTo(BeEmpty())
	vsphereConf := findConfigMapByName(configMaps, "vsphere-cloud-config")
	Expect(vsphereConf).NotTo(BeNil())
	rawConfigINI := []byte(vsphereConf.Data["vsphere.conf"])
	Expect(rawConfigINI).NotTo(BeNil())
	cpiConfig, err := config.ReadCPIConfigINI(rawConfigINI)
	Expect(err).NotTo(HaveOccurred())
	return cpiConfig
}

func configuredIPFamily(output string) []string {
	return cpiConfig(output).VirtualCenter["fake-server.com"].IPFamilyPriority
}

func nodesConfiguration(output string) config.Nodes {
	return cpiConfig(output).Nodes
}

func assertFoundOne(docs []string, err error) {
	Expect(err).NotTo(HaveOccurred())
	Expect(docs).To(HaveLen(1))
}

func vsphereCPIDocuments(output string) {
	// asserting the resources by upstream/vsphere-cpi/01-rbac.yaml
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "ServiceAccount",
		"$.metadata.name": "cloud-controller-manager",
	}))
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "ClusterRole",
		"$.metadata.name": "system:cloud-controller-manager",
	}))
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "ClusterRoleBinding",
		"$.metadata.name": "system:cloud-controller-manager",
	}))

	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "RoleBinding",
		"$.metadata.name": "servicecatalog.k8s.io:apiserver-authentication-reader",
	}))

	// asserting the resources by upstream/vsphere-cpi/02-config.yaml
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "ConfigMap",
		"$.metadata.name": "vsphere-cloud-config",
	}))

	// asserting the resources by upstream/vsphere-cpi/03-secret.yaml
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "Secret",
		"$.metadata.name": "cloud-provider-vsphere-credentials",
	}))

	// asserting the resources by upstream/vsphere-cpi/04-daemonset.yaml
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "DaemonSet",
		"$.metadata.name": "vsphere-cloud-controller-manager",
	}))
}

func vsphereParavirtualDocuments(output string) {
	// asserting the resources by upstream/vsphere-paravirtual-cpi/01-rbac.yaml
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "Namespace",
		"$.metadata.name": "vmware-system-cloud-provider",
	}))
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":               "ServiceAccount",
		"$.metadata.name":      "cloud-provider-svc-account",
		"$.metadata.namespace": "vmware-system-cloud-provider",
	}))
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "ClusterRole",
		"$.metadata.name": "cloud-provider-cluster-role",
	}))
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "ClusterRole",
		"$.metadata.name": "cloud-provider-patch-cluster-role",
	}))
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "ClusterRoleBinding",
		"$.metadata.name": "cloud-provider-cluster-role-binding",
	}))
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "ClusterRoleBinding",
		"$.metadata.name": "cloud-provider-patch-cluster-role-binding",
	}))

	// asserting the resources by upstream/vsphere-paravirtual-cpi/02-rbac.yaml
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":               "ConfigMap",
		"$.metadata.name":      "ccm-cloud-config",
		"$.metadata.namespace": "vmware-system-cloud-provider",
	}))
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":               "ConfigMap",
		"$.metadata.name":      "ccm-owner-reference",
		"$.metadata.namespace": "vmware-system-cloud-provider",
	}))

	// asserting the resources by upstream/vsphere-paravirtual-cpi/03-deployment.yaml
	assertFoundOne(matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":               "Deployment",
		"$.metadata.name":      "guest-cluster-cloud-provider",
		"$.metadata.namespace": "vmware-system-cloud-provider",
	}))
}

func tlsThumbprint(output string) string {
	return cpiConfig(output).Global.Thumbprint
}
