// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package interworking_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	goyaml "gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/repo"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/ytt"
)

var (
	interworkingConfigName         = "antrea-interworking-config"
	interworkingBootstrapConfigmap = "bootstrap-config"

	// data header overwritten
	dataHeader = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---`
)

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
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
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

func marshalAntreaConfig(valuesFile string, settingFunc func(config *AntreaInterworkingConfig)) (string, error) {
	var (
		err    error
		config *AntreaInterworkingConfig
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
func loadAntreaConfig(configFile string) (*AntreaInterworkingConfig, error) {
	var (
		err     error
		content []byte
		config  = AntreaInterworkingConfig{}
	)

	if content, err = os.ReadFile(configFile); err != nil {
		return nil, err
	}
	if err := goyaml.Unmarshal(content, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

var _ = ginkgo.Describe("Antrea-interworking YTT Templates", func() {
	var (
		filePaths []string
		values    string
		output    string
		err       error

		configDir                                  = filepath.Join(repo.RootDir(), "addons/packages/antrea/1.7.1-p1/bundle/config")
		fileAntreaInterworkingYaml                 = filepath.Join(configDir, "upstream/interworking.yaml")
		fileAntreaInterworkingBootstrapYaml        = filepath.Join(configDir, "upstream/bootstrap.yaml")
		fileAntreaInterworkingOverlayYaml          = filepath.Join(configDir, "overlay/interworking-overlay.yaml")
		fileAntreaInterworkingBootstrapOverlayYaml = filepath.Join(configDir, "overlay/interworking-bootstrap-overlay.yaml")
		fileValuesSchema                           = filepath.Join(configDir, "schema.yaml")
		fileValuesYaml                             = filepath.Join(configDir, "values.yaml")
		fileValuesStar                             = filepath.Join(configDir, "values.star")
	)

	ginkgo.BeforeEach(func() {
		values = ""
	})

	ginkgo.JustBeforeEach(func() {
		filePaths = []string{fileValuesSchema, fileValuesYaml, fileValuesStar, fileAntreaInterworkingOverlayYaml,
			fileAntreaInterworkingBootstrapOverlayYaml, fileAntreaInterworkingYaml, fileAntreaInterworkingBootstrapYaml}
		output, err = ytt.RenderYTTTemplate(ytt.CommandOptions{}, filePaths, strings.NewReader(values))
	})

	ginkgo.Context("antrea-interworking default configuration", func() {
		ginkgo.It("mp-adapter configuration", func() {
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			configMap := findConfigMapByName(unmarshalConfigMaps(output), interworkingConfigName)
			gomega.Expect(configMap).NotTo(gomega.BeNil())
			gomega.Expect(configMap.Data["mp-adapter.conf"]).To(gomega.MatchYAML(`---
NSXRemoteAuth: false
NSXClientAuthCertFile: /etc/antrea/nsx-cert/tls.crt
NSXClientAuthKeyFile: /etc/antrea/nsx-cert/tls.key
NSXCAFile: ""
NSXInsecure: true
NSXClientTimeout: 120
InventoryBatchSize: 50
InventoryBatchPeriod: 5
NSXRPCConnType: tnproxy
EnableDebugServer: false
APIServerPort: 16664
DebugServerPort: 16666
NSXRPCDebug: false
#in second
ConditionTimeout: 150`))
			gomega.Expect(configMap.Data["ccp-adapter.conf"]).To(gomega.MatchYAML(`---
EnableDebugServer: false
APIServerPort: 16665
DebugServerPort: 16667
NSXRPCDebug: false
# Time to wait for realization
RealizeTimeoutSeconds: 60
# An interval for regularly report latest realization error in background
RealizeErrorSyncIntervalSeconds: 600
ReconcilerWorkerCount: 8
# Average QPS = ReconcilerWorkerCount * ReconcilerQPS
ReconcilerQPS: 5.0
# Peak QPS =  ReconcilerWorkerCount * ReconcilerBurst
ReconcilerBurst: 10
# 24 Hours
ReconcilerResyncSeconds: 86400`))
		})
	})

	ginkgo.Context("antrea-interworking set antrea_nsx to false", func() {
		ginkgo.BeforeEach(func() {
			values, err = marshalAntreaConfig(fileValuesYaml, func(config *AntreaInterworkingConfig) {
				config.AntreaInterworking.Config.MpAdapterConf.NSXClientTimeout = 100
				config.AntreaNSX.Enable = false
			})
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		ginkgo.It("renders a ConfigMap with MpAdapterConf configuration", func() {
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			configMap := findConfigMapByName(unmarshalConfigMaps(output), interworkingConfigName)
			gomega.Expect(configMap).NotTo(gomega.BeNil())
			gomega.Expect(configMap.Data["mp-adapter.conf"]).To(gomega.ContainSubstring(`NSXClientTimeout: 120`))
		})
	})

	ginkgo.Context("antrea-interworking with MpAdapterConf configuration", func() {
		ginkgo.BeforeEach(func() {
			values, err = marshalAntreaConfig(fileValuesYaml, func(config *AntreaInterworkingConfig) {
				config.AntreaInterworking.Config.MpAdapterConf.NSXClientTimeout = 100
				config.AntreaNSX.Enable = true
			})
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		ginkgo.It("renders a ConfigMap with MpAdapterConf configuration", func() {
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			configMap := findConfigMapByName(unmarshalConfigMaps(output), interworkingConfigName)
			gomega.Expect(configMap).NotTo(gomega.BeNil())
			gomega.Expect(configMap.Data["mp-adapter.conf"]).To(gomega.ContainSubstring(`NSXClientTimeout: 100`))
		})
	})

	ginkgo.Context("antrea-interworking with mp_adapter_conf and ccp_adapter_conf", func() {
		ginkgo.BeforeEach(func() {
			values, err = marshalAntreaConfig(fileValuesYaml, func(config *AntreaInterworkingConfig) {
				config.AntreaInterworking.Config.MpAdapterConf.EnableDebugServer = false
				config.AntreaInterworking.Config.CCPAdapterConf.EnableDebugServer = true
				config.AntreaNSX.Enable = true
			})
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		ginkgo.It("renders a ConfigMap with mp_adapter_conf and ccp_adapter_conf", func() {
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			configMap := findConfigMapByName(unmarshalConfigMaps(output), interworkingConfigName)
			gomega.Expect(configMap).NotTo(gomega.BeNil())
			gomega.Expect(configMap.Data["mp-adapter.conf"]).To(gomega.ContainSubstring(`EnableDebugServer: false`))
			gomega.Expect(configMap.Data["ccp-adapter.conf"]).To(gomega.ContainSubstring(`EnableDebugServer: true`))
		})
	})

	ginkgo.Context("antrea-interworking with mp_adapter_conf and ccp_adapter_conf", func() {
		ginkgo.BeforeEach(func() {
			values, err = marshalAntreaConfig(fileValuesYaml, func(config *AntreaInterworkingConfig) {
				config.AntreaInterworking.Config.MpAdapterConf.EnableDebugServer = false
				config.AntreaInterworking.Config.CCPAdapterConf.EnableDebugServer = true
				config.AntreaNSX.Enable = true
			})
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		ginkgo.It("renders a ConfigMap with mp_adapter_conf and ccp_adapter_conf", func() {
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			configMap := findConfigMapByName(unmarshalConfigMaps(output), interworkingConfigName)
			gomega.Expect(configMap).NotTo(gomega.BeNil())
			gomega.Expect(configMap.Data["mp-adapter.conf"]).To(gomega.ContainSubstring(`EnableDebugServer: false`))
			gomega.Expect(configMap.Data["ccp-adapter.conf"]).To(gomega.ContainSubstring(`EnableDebugServer: true`))
		})
	})

	ginkgo.Context("antrea-interworking default bootstrap configuration", func() {
		ginkgo.It("default bootstrap configuration", func() {
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			configMap := findConfigMapByName(unmarshalConfigMaps(output), interworkingBootstrapConfigmap)
			gomega.Expect(configMap).NotTo(gomega.BeNil())
			gomega.Expect(configMap.Data["bootstrap.conf"]).To(gomega.ContainSubstring(`clusterName: dummyClusterName`))
		})
	})

	ginkgo.Context("antrea-interworking bootstrap configuration", func() {
		ginkgo.BeforeEach(func() {
			values, err = marshalAntreaConfig(fileValuesYaml, func(config *AntreaInterworkingConfig) {
				config.AntreaInterworking.Config.NSXKey = "fake-cert"
				config.AntreaInterworking.Config.ClusterName = "fake-clusterName"
				config.AntreaNSX.Enable = true
			})
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		ginkgo.It("renders a ConfigMap with bootstrap.conf", func() {
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			configMap := findConfigMapByName(unmarshalConfigMaps(output), interworkingBootstrapConfigmap)
			gomega.Expect(configMap).NotTo(gomega.BeNil())
			gomega.Expect(configMap.Data["bootstrap.conf"]).To(gomega.ContainSubstring(`clusterName: fake-clusterName`))
		})
	})
})
