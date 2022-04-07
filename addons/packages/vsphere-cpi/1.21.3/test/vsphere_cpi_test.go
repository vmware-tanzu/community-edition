// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package vspherecpi_test

import (
	"fmt"
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

var _ = Describe("vSphere CPI Ytt Templates", func() {
	var (
		filePaths    []string
		values       string
		output       string
		yttRenderErr error

		configDir                  = filepath.Join(repo.RootDir(), "addons/packages/vsphere-cpi/1.21.3/bundle/config")
		file01rbac                 = filepath.Join(configDir, "upstream/vsphere-cpi/01-rbac.yaml")
		file02config               = filepath.Join(configDir, "upstream/vsphere-cpi/02-config.yaml")
		file03secret               = filepath.Join(configDir, "upstream/vsphere-cpi/03-secret.yaml")
		file04daemonset            = filepath.Join(configDir, "upstream/vsphere-cpi/04-daemonset.yaml")
		fileOverlayUpdateConfig    = filepath.Join(configDir, "overlays/update-config.yaml")
		fileOverlayAddSecret       = filepath.Join(configDir, "overlays/add-secret.yaml")
		fileOverlayUpdateSecret    = filepath.Join(configDir, "overlays/update-secret.yaml")
		fileOverlayUpdateDaemonset = filepath.Join(configDir, "overlays/update-daemonset.yaml")
		fileValuesYaml             = filepath.Join(configDir, "values.yaml")
		fileValuesStar             = filepath.Join(configDir, "values.star")
		fileVsphereconfLibTxt      = filepath.Join(configDir, "vsphereconf.lib.txt")
	)

	BeforeEach(func() {
		values = `#@data/values
#@overlay/match-child-defaults missing_ok=True
---
vsphereCPI:
  server: fake-server.com
  datacenter: dc0
  username: my-user
  password: my-password
  insecureFlag: True`
	})

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
			fileValuesYaml,
			fileValuesStar,
			fileVsphereconfLibTxt,
		}
		output, yttRenderErr = ytt.RenderYTTTemplate(ytt.CommandOptions{}, filePaths, strings.NewReader(values))
	})

	Context("DaemonSet env vars", func() {
		It("renders a DaemonSet with ENABLE_ALPHA_DUAL_STACK env var feature flag", func() {
			Expect(yttRenderErr).NotTo(HaveOccurred())
			daemonSet := parseDaemonSet(output)
			containerEnvVars := daemonSet.Spec.Template.Spec.Containers[0].Env
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
