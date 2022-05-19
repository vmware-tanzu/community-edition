// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package externaldns_test

import (
	"os"
	"path/filepath"

	. "github.com/vmware-tanzu/community-edition/addons/packages/test/matchers"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/ytt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("External DNS Ytt Templates", func() {
	var (
		configDir  string
		workingDir string
	)

	BeforeEach(func() {
		var err error
		workingDir, err = os.Getwd()
		Expect(err).NotTo(HaveOccurred())

		configDir = filepath.Join(workingDir, "..", "..", "bundle", "config")
	})

	Context("Providing a minimal configuration", func() {
		var output string

		BeforeEach(func() {
			var err error
			output, err = ytt.RenderYTTTemplate(ytt.CommandOptions{}, []string{
				configDir,
				filepath.Join(workingDir, "fixtures", "values", "minimal-configuration.yaml"),
			}, nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("renders upstream yaml documents", func() {
			Expect(FindDocsMatchingYAMLPath(output, map[string]string{".kind": "ClusterRole", ".metadata.name": "external-dns"})).To(HaveLen(1))
			Expect(FindDocsMatchingYAMLPath(output, map[string]string{".kind": "ClusterRoleBinding", ".metadata.name": "external-dns-viewer"})).To(HaveLen(1))
			Expect(FindDocsMatchingYAMLPath(output, map[string]string{".kind": "ServiceAccount", ".metadata.name": "external-dns"})).To(HaveLen(1))
		})

		It("renders the default namespace", func() {
			Expect(FindDocsMatchingYAMLPath(output, map[string]string{".kind": "Namespace", ".metadata.name": "external-dns"})).To(HaveLen(1))
		})

		It("renders the deployment.args", func() {
			deploymentDoc, err := FindDocsMatchingYAMLPath(output, map[string]string{".kind": "Deployment", ".metadata.name": "external-dns"})
			Expect(deploymentDoc).To(HaveLen(1))
			Expect(err).NotTo(HaveOccurred())
			Expect(deploymentDoc[0]).To(HaveYAMLPathWithValue(".spec.template.spec.containers[0].args[0]", "--source=ingress"))
			Expect(deploymentDoc[0]).To(HaveYAMLPathWithValue(".spec.template.spec.containers[0].args[1]", "--source=contour-httpproxy"))
			Expect(deploymentDoc[0]).To(HaveYAMLPathWithValue(".spec.template.spec.containers[0].args[2]", "--provider=rfc2136"))
		})

		It("does not configure optional keys", func() {
			deploymentDoc, err := FindDocsMatchingYAMLPath(output, map[string]string{".kind": "Deployment", ".metadata.name": "external-dns"})
			Expect(deploymentDoc).To(HaveLen(1))
			Expect(err).NotTo(HaveOccurred())
			Expect(deploymentDoc[0]).NotTo(HaveYAMLPath("$.spec.template.spec.containers[0].env"))
			Expect(deploymentDoc[0]).NotTo(HaveYAMLPath("$.spec.template.spec.containers[0].securityContext"))
			Expect(deploymentDoc[0]).NotTo(HaveYAMLPath("$.spec.template.spec.containers[0].volumeMounts"))
			Expect(deploymentDoc[0]).NotTo(HaveYAMLPath("$.spec.template.spec.volumes"))
		})
	})

	Context("Providing a namespace", func() {
		It("renders a setup in a different namespace", func() {
			output, err := ytt.RenderYTTTemplate(ytt.CommandOptions{}, []string{
				configDir,
				filepath.Join(workingDir, "fixtures", "values", "namespace.yaml"),
			}, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(FindDocsMatchingYAMLPath(
				output, map[string]string{".metadata.namespace": "custom-external-dns-namespace"},
			)).To(HaveLen(2))
		})
	})

	Context("Providing env vars for the deployment", func() {
		It("renders a deployment with env vars", func() {
			output, err := ytt.RenderYTTTemplate(ytt.CommandOptions{}, []string{
				configDir,
				filepath.Join(workingDir, "fixtures", "values", "deployment-env-vars.yaml"),
			}, nil)
			Expect(err).NotTo(HaveOccurred())

			deploymentDocs, err := FindDocsMatchingYAMLPath(
				output,
				map[string]string{".kind": "Deployment"},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(deploymentDocs).To(HaveLen(1))
			Expect(deploymentDocs[0]).To(HaveYAMLPathWithValue("$.spec.template.spec.containers[0].env[0].name", "FOO"))
			Expect(deploymentDocs[0]).To(HaveYAMLPathWithValue("$.spec.template.spec.containers[0].env[0].value", "bar"))
		})
	})

	Context("Providing the security context for the deployment", func() {
		It("renders a deployment with a custom security context", func() {
			output, err := ytt.RenderYTTTemplate(ytt.CommandOptions{}, []string{
				configDir,
				filepath.Join(workingDir, "fixtures", "values", "deployment-security-context.yaml"),
			}, nil)
			Expect(err).NotTo(HaveOccurred())

			deploymentDocs, err := FindDocsMatchingYAMLPath(
				output,
				map[string]string{".kind": "Deployment"},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(deploymentDocs).To(HaveLen(1))
			Expect(deploymentDocs[0]).To(HaveYAMLPathWithValue("$.spec.template.spec.containers[0].securityContext.runAsUser", "1000"))
			Expect(deploymentDocs[0]).To(HaveYAMLPathWithValue("$.spec.template.spec.containers[0].securityContext.runAsGroup", "2000"))
			Expect(deploymentDocs[0]).To(HaveYAMLPathWithValue("$.spec.template.spec.containers[0].securityContext.allowPrivilegeEscalation", "false"))
		})
	})

	Context("Providing volumes and their mounts for the deployment", func() {
		It("renders a deployment with additional volumes mounted", func() {
			output, err := ytt.RenderYTTTemplate(ytt.CommandOptions{}, []string{
				configDir,
				filepath.Join(workingDir, "fixtures", "values", "deployment-volumes.yaml"),
			}, nil)
			Expect(err).NotTo(HaveOccurred())

			deploymentDocs, err := FindDocsMatchingYAMLPath(
				output,
				map[string]string{".kind": "Deployment"},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(deploymentDocs).To(HaveLen(1))
			Expect(deploymentDocs[0]).To(HaveYAMLPathWithValue("$.spec.template.spec.volumes[0].name", "additional-volume"))
			Expect(deploymentDocs[0]).To(HaveYAMLPathWithValue("$.spec.template.spec.volumes[0].emptyDir", ""))
			Expect(deploymentDocs[0]).To(HaveYAMLPathWithValue("$.spec.template.spec.containers[0].volumeMounts[0].name", "additional-volume"))
			Expect(deploymentDocs[0]).To(HaveYAMLPathWithValue("$.spec.template.spec.containers[0].volumeMounts[0].mountPath", "/path/in/container"))
		})
	})

	Context("Providing annotations for the service account", func() {
		It("renders a service account with annotations", func() {
			output, err := ytt.RenderYTTTemplate(ytt.CommandOptions{}, []string{
				configDir,
				filepath.Join(workingDir, "fixtures", "values", "serviceaccount-annotations.yaml"),
			}, nil)
			Expect(err).NotTo(HaveOccurred())

			Expect(FindDocsMatchingYAMLPath(
				output,
				map[string]string{
					".kind":                     "ServiceAccount",
					".metadata.annotations.key": "value",
				},
			)).To(HaveLen(1))
		})
	})

	Describe("error validation", func() {
		Context("No --source in deployment.args", func() {
			It("renders with an error", func() {
				_, err := ytt.RenderYTTTemplate(ytt.CommandOptions{}, []string{
					configDir,
					filepath.Join(workingDir, "fixtures", "values", "deployment-args-no-source.yaml"),
				}, nil)
				Expect(err).To(ContainSubstring("--source is required in deployment.args to query for endpoints"))
			})
		})

		Context("No --provider in deployment.args", func() {
			It("renders with an error", func() {
				_, err := ytt.RenderYTTTemplate(ytt.CommandOptions{}, []string{
					configDir,
					filepath.Join(workingDir, "fixtures", "values", "deployment-args-no-provider.yaml"),
				}, nil)
				Expect(err).To(ContainSubstring("--provider is required in deployment.args to define a DNS provider where records will be created"))
			})
		})
	})
})
