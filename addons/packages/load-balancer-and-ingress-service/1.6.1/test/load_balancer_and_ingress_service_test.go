// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package loadbalancer_and_ingress_service_test

import (
	"path/filepath"
	"strings"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/providers/tests/unit/matchers"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/repo"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/ytt"
)

var _ = Describe("LoadBalancer And Ingress Service Ytt Templates", func() {
	var (
		filePaths    []string
		values       string
		output       string
		yttRenderErr error

		configDir = filepath.Join(repo.RootDir(), "addons/packages/load-balancer-and-ingress-service/1.6.1/bundle/config")

		file01rbac                   = filepath.Join(configDir, "upstream/loadbalancerandingressservice/rbac.yaml")
		file02config                 = filepath.Join(configDir, "upstream/loadbalancerandingressservice/configmap.yaml")
		file03secret                 = filepath.Join(configDir, "upstream/loadbalancerandingressservice/secret.yaml")
		file04statefulset            = filepath.Join(configDir, "upstream/loadbalancerandingressservice/statefulset.yaml")
		file05AviInfraSettings       = filepath.Join(configDir, "upstream/loadbalancerandingressservice/crds/aviinfrasettings.yaml")
		fileOverlayUpdateConfig      = filepath.Join(configDir, "overlays/overlay-configmap.yaml")
		fileOverlayUpdateRbac        = filepath.Join(configDir, "overlays/overlay-rbac.yaml")
		fileOverlayUpdateSecret      = filepath.Join(configDir, "overlays/overlay-secret.yaml")
		fileOverlayUpdateStatefulSet = filepath.Join(configDir, "overlays/overlay-statefulset.yaml")

		fileValuesYaml     = filepath.Join(configDir, "values.yaml")
		fileValuesStar     = filepath.Join(configDir, "values.star")
		fileKappConfigYaml = filepath.Join(configDir, "kapp-config.yaml")
	)

	JustBeforeEach(func() {
		filePaths = []string{
			file01rbac,
			file02config,
			file03secret,
			file04statefulset,
			file05AviInfraSettings,
			fileOverlayUpdateConfig,
			fileOverlayUpdateRbac,
			fileOverlayUpdateSecret,
			fileOverlayUpdateStatefulSet,
			fileValuesYaml,
			fileValuesStar,
			fileKappConfigYaml,
		}
		output, yttRenderErr = ytt.RenderYTTTemplate(ytt.CommandOptions{}, filePaths, strings.NewReader(values))
	})

	Context("L7 ingress config", func() {
		BeforeEach(func() {
			values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
loadBalancerAndIngressService:
  config:
    ako_settings:
      disable_static_route_sync: "false"
    controller_settings:
      service_engine_group_name: SEG1
      cloud_name: Default-Cloud
      controller_ip: 192.168.10.1
    l7_settings:
      disable_ingress_class: false
      default_ing_controller: true
      shard_vs_size: MEDIUM
`
		})

		When("ConfigMap has been L7 ingress settings of type NodePortLocal", func() {
			BeforeEach(func() {
				values = values + `
      service_type: NodePortLocal
`
			})

			It("should render configmap for ako", func() {
				Expect(yttRenderErr).NotTo(HaveOccurred())
				docs, err := matchers.FindDocsMatchingYAMLPath(output, map[string]string{
					"$.kind":               "ConfigMap",
					"$.metadata.name":      "avi-k8s-config",
					"$.metadata.namespace": "avi-system",
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(docs).To(HaveLen(1))
				cfg := unmarshalConfigMap(output)
				Expect(cfg.Namespace).To(Equal("avi-system"))
				Expect(getConfigKeyFromData(cfg, "disableStaticRouteSync")).To(Equal("false"))
				Expect(getConfigKeyFromData(cfg, "defaultIngController")).To(Equal("true"))
				Expect(getConfigKeyFromData(cfg, "serviceType")).To(Equal("NodePortLocal"))
				Expect(getConfigKeyFromData(cfg, "shardVSSize")).To(Equal("MEDIUM"))
			})
		})

		When("ConfigMap has been L7 ingress settings of type NodePort", func() {
			BeforeEach(func() {
				values = values + `
      service_type: NodePort
`
			})

			It("should render configmap for ako", func() {
				Expect(yttRenderErr).NotTo(HaveOccurred())
				docs, err := matchers.FindDocsMatchingYAMLPath(output, map[string]string{
					"$.kind":               "ConfigMap",
					"$.metadata.name":      "avi-k8s-config",
					"$.metadata.namespace": "avi-system",
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(docs).To(HaveLen(1))
				cfg := unmarshalConfigMap(output)
				Expect(cfg.Namespace).To(Equal("avi-system"))
				Expect(getConfigKeyFromData(cfg, "disableStaticRouteSync")).To(Equal("false"))
				Expect(getConfigKeyFromData(cfg, "defaultIngController")).To(Equal("true"))
				Expect(getConfigKeyFromData(cfg, "serviceType")).To(Equal("NodePort"))
				Expect(getConfigKeyFromData(cfg, "shardVSSize")).To(Equal("MEDIUM"))
			})
		})

		When("ConfigMap has been L7 ingress settings of type ClusterIP", func() {
			BeforeEach(func() {
				values = values + `
      service_type: ClusterIP
`
			})

			It("should render configmap for ako", func() {
				Expect(yttRenderErr).NotTo(HaveOccurred())
				docs, err := matchers.FindDocsMatchingYAMLPath(output, map[string]string{
					"$.kind":               "ConfigMap",
					"$.metadata.name":      "avi-k8s-config",
					"$.metadata.namespace": "avi-system",
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(docs).To(HaveLen(1))
				cfg := unmarshalConfigMap(output)
				Expect(cfg.Namespace).To(Equal("avi-system"))
				Expect(getConfigKeyFromData(cfg, "disableStaticRouteSync")).To(Equal("false"))
				Expect(getConfigKeyFromData(cfg, "defaultIngController")).To(Equal("true"))
				Expect(getConfigKeyFromData(cfg, "serviceType")).To(Equal("ClusterIP"))
				Expect(getConfigKeyFromData(cfg, "shardVSSize")).To(Equal("MEDIUM"))
			})
		})
	})

	Context("L4 config", func() {
		BeforeEach(func() {
			values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
loadBalancerAndIngressService:
  config:
    controller_settings:
      service_engine_group_name: SEG1
      cloud_name: Default-Cloud
      controller_ip: 192.168.10.1
    network_settings:
      vip_network_list: '[{"networkName":"VM Network","cidr":"10.180.112.0/20"}]'
`
		})

		When("ConfigMap has been L4 settings", func() {
			It("should render configmap for ako", func() {
				Expect(yttRenderErr).NotTo(HaveOccurred())
				docs, err := matchers.FindDocsMatchingYAMLPath(output, map[string]string{
					"$.kind":               "ConfigMap",
					"$.metadata.name":      "avi-k8s-config",
					"$.metadata.namespace": "avi-system",
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(docs).To(HaveLen(1))
				cfg := unmarshalConfigMap(output)
				Expect(cfg.Namespace).To(Equal("avi-system"))
				Expect(getConfigKeyFromData(cfg, "controllerIP")).To(Equal("192.168.10.1"))
				Expect(getConfigKeyFromData(cfg, "serviceEngineGroupName")).To(Equal("SEG1"))
				Expect(getConfigKeyFromData(cfg, "cloudName")).To(Equal("Default-Cloud"))
				Expect(getConfigKeyFromData(cfg, "apiServerPort")).To(Equal("8080"))
				Expect(getConfigKeyFromData(cfg, "deleteConfig")).To(Equal("false"))
				Expect(getConfigKeyFromData(cfg, "vipNetworkList")).To(Equal("[{\"networkName\":\"VM Network\",\"cidr\":\"10.180.112.0/20\"}]"))
			})
		})
	})

	Context("AVIInfraSetting CRD test", func() {
		When("Cluster type is management", func() {
			BeforeEach(func() {
				values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
loadBalancerAndIngressService:
  config:
    tkg_cluster_role: "management"
`
			})

			It("AVIInfraSetting CRD should not be installed", func() {
				Expect(yttRenderErr).NotTo(HaveOccurred())
				Expect(AviInfraSettingCRDExists(output)).To(Equal(false))
			})
		})

		When("Cluster type is workload", func() {
			BeforeEach(func() {
				values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
loadBalancerAndIngressService:
  config:
    tkg_cluster_role: "workload"
`
			})

			It("AVIInfraSetting CRD should be installed", func() {
				Expect(yttRenderErr).NotTo(HaveOccurred())
				Expect(AviInfraSettingCRDExists(output)).To(Equal(true))
			})
		})
	})
})

func AviInfraSettingCRDExists(output string) bool {
	docs, err := matchers.FindDocsMatchingYAMLPath(output, map[string]string{
		"$.kind":          "CustomResourceDefinition",
		"$.metadata.name": "aviinfrasettings.ako.vmware.com",
	})

	return err == nil && len(docs) == 1
}

func getConfigKeyFromData(cfg corev1.ConfigMap, key string) string {
	cfg_value, exists := cfg.Data[key]
	Expect(exists).To(Equal(true))
	return cfg_value
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
