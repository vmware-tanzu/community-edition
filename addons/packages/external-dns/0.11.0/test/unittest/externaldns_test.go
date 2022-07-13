// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package externaldns_test

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/vmware-tanzu/community-edition/addons/packages/test/matchers"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/ytt"

	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("External DNS Ytt Templates", func() {
	Context("Providing a minimal configuration", func() {
		var output string

		BeforeEach(func() {
			var err error
			output, err = renderWithDataValuesFixture("minimal-configuration.yaml")
			Expect(err).NotTo(HaveOccurred())
		})

		It("renders upstream yaml documents", func() {
			Expect(FindDocsMatchingYAMLPath(output, map[string]string{
				".kind":          "ClusterRole",
				".metadata.name": "external-dns",
			})).To(HaveLen(1))
			Expect(FindDocsMatchingYAMLPath(output, map[string]string{
				".kind":          "ClusterRoleBinding",
				".metadata.name": "external-dns-viewer",
			})).To(HaveLen(1))
			Expect(FindDocsMatchingYAMLPath(output, map[string]string{
				".kind":          "ServiceAccount",
				".metadata.name": "external-dns",
			})).To(HaveLen(1))
		})

		It("renders the default namespace", func() {
			Expect(FindDocsMatchingYAMLPath(output, map[string]string{
				".kind":          "Namespace",
				".metadata.name": "external-dns",
			})).To(HaveLen(1))
		})

		It("renders the deployment.args", func() {
			deploymentDoc := findDeploymentDoc(output)
			Expect(deploymentDoc).To(HaveYAMLPathWithValue(
				".spec.template.spec.containers[0].args[0]",
				"--source=ingress",
			))
			Expect(deploymentDoc).To(HaveYAMLPathWithValue(
				".spec.template.spec.containers[0].args[1]",
				"--source=contour-httpproxy",
			))
			Expect(deploymentDoc).To(HaveYAMLPathWithValue(
				".spec.template.spec.containers[0].args[2]",
				"--provider=rfc2136",
			))
		})

		It("does not configure optional keys", func() {
			deploymentDoc := findDeploymentDoc(output)
			Expect(deploymentDoc).NotTo(HaveYAMLPath("$.spec.template.spec.containers[0].env"))
			Expect(deploymentDoc).NotTo(HaveYAMLPath("$.spec.template.spec.containers[0].securityContext"))
			Expect(deploymentDoc).NotTo(HaveYAMLPath("$.spec.template.spec.containers[0].volumeMounts"))
			Expect(deploymentDoc).NotTo(HaveYAMLPath("$.spec.template.spec.volumes"))
		})
	})

	Context("Providing a namespace", func() {
		It("renders a setup in a different namespace", func() {
			output, err := renderWithDataValuesFixture("namespace.yaml")
			Expect(err).NotTo(HaveOccurred())
			Expect(FindDocsMatchingYAMLPath(
				output, map[string]string{".metadata.namespace": "custom-external-dns-namespace"},
			)).To(HaveLen(2))
		})
	})

	Context("Providing env vars for the deployment", func() {
		It("renders a deployment with env vars", func() {
			output, err := renderWithDataValuesFixture("deployment-env-vars.yaml")
			Expect(err).NotTo(HaveOccurred())

			deploymentDoc := findDeploymentDoc(output)
			Expect(deploymentDoc).To(HaveYAMLPathWithValue(
				"$.spec.template.spec.containers[0].env[0].name",
				"FOO",
			))
			Expect(deploymentDoc).To(HaveYAMLPathWithValue(
				"$.spec.template.spec.containers[0].env[0].value",
				"bar",
			))
		})
	})

	Context("Providing the security context for the deployment", func() {
		It("renders a deployment with a custom security context", func() {
			output, err := renderWithDataValuesFixture("deployment-security-context.yaml")
			Expect(err).NotTo(HaveOccurred())
			deploymentDoc := findDeploymentDoc(output)
			Expect(deploymentDoc).To(HaveYAMLPathWithValue(
				"$.spec.template.spec.containers[0].securityContext.runAsUser",
				"1000",
			))
			Expect(deploymentDoc).To(HaveYAMLPathWithValue(
				"$.spec.template.spec.containers[0].securityContext.runAsGroup",
				"2000",
			))
			Expect(deploymentDoc).To(HaveYAMLPathWithValue(
				"$.spec.template.spec.containers[0].securityContext.allowPrivilegeEscalation",
				"false",
			))
		})
	})

	Context("Providing volumes and their mounts for the deployment", func() {
		It("renders a deployment with additional volumes mounted", func() {
			output, err := renderWithDataValuesFixture("deployment-volumes.yaml")
			Expect(err).NotTo(HaveOccurred())

			deploymentDoc := findDeploymentDoc(output)
			Expect(deploymentDoc).To(HaveYAMLPathWithValue(
				"$.spec.template.spec.volumes[0].name",
				"additional-volume",
			))
			Expect(deploymentDoc).To(HaveYAMLPathWithValue(
				"$.spec.template.spec.volumes[0].emptyDir",
				"",
			))
			Expect(deploymentDoc).To(HaveYAMLPathWithValue(
				"$.spec.template.spec.containers[0].volumeMounts[0].name",
				"additional-volume",
			))
			Expect(deploymentDoc).To(HaveYAMLPathWithValue(
				"$.spec.template.spec.containers[0].volumeMounts[0].mountPath",
				"/path/in/container",
			))
		})
	})

	Context("Providing annotations for the service account", func() {
		It("renders a service account with annotations", func() {
			output, err := renderWithDataValuesFixture("serviceaccount-annotations.yaml")
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
				_, err := renderWithDataValuesFixture("deployment-args-no-source.yaml")
				Expect(err).To(ContainSubstring("--source is required in deployment.args to query for endpoints"))
			})
		})

		Context("No --provider in deployment.args", func() {
			It("renders with an error", func() {
				_, err := renderWithDataValuesFixture("deployment-args-no-provider.yaml")
				Expect(err).To(ContainSubstring("--provider is required in deployment.args to define a DNS provider where records will be created"))
			})
		})
	})

	Describe("aws secrets", func() {
		Describe("when not supplied", func() {
			var output string
			BeforeEach(func() {
				var err error
				output, err = renderWithDataValuesFixture("minimal-configuration.yaml")
				Expect(err).NotTo(HaveOccurred())
			})
			It("does not render a secret", func() {
				secretDocs, err := FindDocsMatchingYAMLPath(output, map[string]string{".kind": "Secret"})
				Expect(err).NotTo(HaveOccurred())
				Expect(secretDocs).To(BeEmpty())
			})
			It("does not renders env vars on the Deployment with secret refs", func() {
				deploymentDoc := findDeploymentDoc(output)
				Expect(deploymentDoc).NotTo(HaveYAMLPath("$.spec.template.spec.containers[0].envFrom[?(@.secretRef.name=='external-dns-aws-values')]"))
			})
		})
		Describe("when supplied", func() {
			It("renders a secret", func() {
				output, err := renderWithDataValuesFixture("aws-secret.yaml")
				Expect(err).NotTo(HaveOccurred())

				secretDocs, err := FindDocsMatchingYAMLPath(output, map[string]string{".kind": "Secret"})
				Expect(err).NotTo(HaveOccurred())
				Expect(secretDocs).To(HaveLen(1))
				Expect(secretDocs[0]).To(HaveYAMLPathWithValue("$.metadata.name", "external-dns-aws-values"))
				Expect(secretDocs[0]).To(HaveYAMLPathWithValue("$.metadata.namespace", "external-dns-aws"))
				Expect(secretDocs[0]).To(HaveYAMLPathWithValue("$.data.awsAccessKeyID", "YXdzIGFjY2VzcyBrZXk="))
				Expect(secretDocs[0]).To(HaveYAMLPathWithValue("$.data.awsSecretAccessKey", "YXdzIHNlY3JldCBrZXk="))
			})
			It("renders env vars on the Deployment with secret refs", func() {
				output, err := renderWithDataValuesFixture("aws-secret.yaml")
				Expect(err).NotTo(HaveOccurred())

				deploymentDoc := findDeploymentDoc(output)
				Expect(deploymentDoc).To(HaveYAMLPathWithValue("$.spec.template.spec.containers[0].env[?(@.name=='AWS_ACCESS_KEY_ID')].valueFrom.secretKeyRef.key", "awsAccessKeyID"))
				Expect(deploymentDoc).To(HaveYAMLPathWithValue("$.spec.template.spec.containers[0].env[?(@.name=='AWS_ACCESS_KEY_ID')].valueFrom.secretKeyRef.name", "external-dns-aws-values"))
				Expect(deploymentDoc).To(HaveYAMLPathWithValue("$.spec.template.spec.containers[0].env[?(@.name=='AWS_SECRET_ACCESS_KEY')].valueFrom.secretKeyRef.name", "external-dns-aws-values"))
				Expect(deploymentDoc).To(HaveYAMLPathWithValue("$.spec.template.spec.containers[0].env[?(@.name=='AWS_SECRET_ACCESS_KEY')].valueFrom.secretKeyRef.key", "awsSecretAccessKey"))
			})
			It("does not interfere with other env vars", func() {
				output, err := renderWithDataValuesFixture("aws-secret.yaml")
				Expect(err).NotTo(HaveOccurred())

				deploymentDoc := findDeploymentDoc(output)
				Expect(deploymentDoc).To(HaveYAMLPathWithValue("$.spec.template.spec.containers[0].env[?(@.name=='other-key')].value", "other-value"))
			})
			It("requires the provider be `aws`", func() {
				values := []string{
					`#@data/values`,
					`---`,
					`deployment:`,
					`  args:`,
					`    - --source=ingress`,
					`    - --provider=cloudflare`,
					`aws:`,
					`  credentials:`,
					`    accessKey: "the-access-key"`,
					`    secretKey: "the-sekret-key"`,
				}
				valuesYaml := strings.Join(values, "\n")
				_, err := renderWithDataValuesInline(valuesYaml)
				Expect(err).To(MatchError(ContainSubstring("Use of `aws.credentials` requires using the aws provider")))
			})
			It("requires `aws.credentials.secretKey` when accessKey is provided", func() {
				values := []string{
					`#@data/values`,
					`---`,
					`deployment:`,
					`  args:`,
					`    - --source=ingress`,
					`    - --provider=aws`,
					`aws:`,
					`  credentials:`,
					`    accessKey: "the-access-key"`,
				}
				valuesYaml := strings.Join(values, "\n")
				_, err := renderWithDataValuesInline(valuesYaml)
				Expect(err).To(MatchError(
					ContainSubstring("`aws.credentials.accessKey` and `aws.credentials.secretKey` must both be provided")))
			})
			It("requires `aws.credentials.accessKey` when secretKey is provided", func() {
				values := []string{
					`#@data/values`,
					`---`,
					`deployment:`,
					`  args:`,
					`    - --source=ingress`,
					`    - --provider=aws`,
					`aws:`,
					`  credentials:`,
					`    secretKey: "the-secret-key"`,
				}
				valuesYaml := strings.Join(values, "\n")
				_, err := renderWithDataValuesInline(valuesYaml)
				Expect(err).To(MatchError(
					ContainSubstring("`aws.credentials.accessKey` and `aws.credentials.secretKey` must both be provided")))
			})
		})
	})

	Describe("Azure credentials", func() {
		var output string
		Describe("when not provided", func() {
			BeforeEach(func() {
				var err error
				output, err = renderWithDataValuesFixture("minimal-configuration.yaml")
				Expect(err).NotTo(HaveOccurred())
			})
			It("does not add a Secret", func() {
				secretDocs, err := FindDocsMatchingYAMLPath(
					output,
					map[string]string{".kind": "Secret"},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(secretDocs).To(BeEmpty())
			})
			It("does not add Volumes, VolumeMounts", func() {
				deploymentDoc := findDeploymentDoc(output)
				Expect(deploymentDoc).NotTo(HaveYAMLPath("$.spec.template.spec.volumes"))
				Expect(deploymentDoc).NotTo(HaveYAMLPath("$.spec.template.spec.containers[0].volumeMounts"))
			})
		})
		Describe("when provided", func() {
			It("renders a Secret and updates the Deployment to ref the volume", func() {
				var err error
				output, err = renderWithDataValuesFixture("azure-configuration.yaml")
				Expect(err).NotTo(HaveOccurred())
				deploymentDoc := findDeploymentDoc(output)
				Expect(deploymentDoc).To(HaveYAMLPathWithValue(
					"$.spec.template.spec.volumes[?(@.name=='azure-config-file')].secret.secretName",
					"azure-config-file",
				))
				Expect(deploymentDoc).To(HaveYAMLPathWithValue(
					"$.spec.template.spec.containers[0].volumeMounts[?(@.name=='azure-config-file')].mountPath",
					"/etc/kubernetes",
				))
				Expect(deploymentDoc).To(HaveYAMLPathWithValue(
					"$.spec.template.spec.containers[0].volumeMounts[?(@.name=='azure-config-file')].readOnly",
					"true",
				))

				secretDocs, err := FindDocsMatchingYAMLPath(
					output,
					map[string]string{".kind": "Secret"},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(secretDocs).To(HaveLen(1))
				secretName := valueAtYAMLPath(deploymentDoc, "$.spec.template.spec.volumes[?(@.name=='azure-config-file')].secret.secretName")
				Expect(secretDocs[0]).To(HaveYAMLPathWithValue("$.metadata.name", secretName))
				Expect(secretDocs[0]).To(HaveYAMLPathWithValue("$.metadata.namespace", "external-dns-azure"))
				Expect(valueAtYAMLPath(secretDocs[0], "$.stringData['azure.json']")).To(MatchJSON(`{
				  "cloud": "azure-cloud",
				  "tenantId": "azure-tenant-id",
				  "subscriptionId": "azure-subscription-id",
				  "resourceGroup": "azure-resource-group",
				  "aadClientId": "azure-aad-client-id",
				  "aadClientSecret": "azure-aad-client-secret",
				  "useManagedIdentityExtension": false,
				  "userAssignedIdentityID": "azure-user-assigned-identity-id"
				}`))
			})

			When("only the required fields are provided", func() {
				It("renders a Secret that omits the optional fields", func() {
					var err error
					output, err = renderWithDataValuesFixture("azure-minimal-configuration.yaml")
					Expect(err).NotTo(HaveOccurred())
					secretDocs, err := FindDocsMatchingYAMLPath(
						output,
						map[string]string{".kind": "Secret"},
					)
					Expect(err).NotTo(HaveOccurred())
					Expect(secretDocs).To(HaveLen(1))
					Expect(secretDocs[0]).To(HaveYAMLPathWithValue(
						"$.metadata.name",
						"azure-config-file",
					))
					Expect(valueAtYAMLPath(secretDocs[0], "$.stringData['azure.json']")).To(MatchJSON(`{
					  "resourceGroup": "azure-resource-group",
					  "subscriptionId": "azure-subscription-id",
					  "tenantId": "azure-tenant-id",
					  "aadClientId": "azure-aad-client-id",
					  "aadClientSecret": "azure-aad-client-secret"
					}`))
				})
			})
			When("useManagedIdentityExtension is true", func() {
				It("does not require aadClientId and aadClientSecret", func() {
					values := []string{
						`#@data/values`,
						`---`,
						`deployment:`,
						`  args:`,
						`    - --source=ingress`,
						`    - --provider=azure`,
						`azure:`,
						`  subscriptionId: "azure-subscription-id"`,
						`  resourceGroup: "azure-resource-group"`,
						`  tenantId: "azure-tenant-id"`,
						`  useManagedIdentityExtension: true`,
					}
					valuesYaml := strings.Join(values, "\n")
					_, err := renderWithDataValuesInline(valuesYaml)
					Expect(err).NotTo(HaveOccurred())
				})
			})
			When("useManagedIdentityExtension is false or not specified", func() {
				It("requires aadClientId", func() {
					values := []string{
						`#@data/values`,
						`---`,
						`deployment:`,
						`  args:`,
						`    - --source=ingress`,
						`    - --provider=azure`,
						`azure:`,
						`  subscriptionId: "azure-subscription-id"`,
						`  resourceGroup: "azure-resource-group"`,
						`  tenantId: "azure-tenant-id"`,
						`  aadClientSecret: "azure-aad-client-secret"`,
					}
					valuesYaml := strings.Join(values, "\n")
					_, err := renderWithDataValuesInline(valuesYaml)
					Expect(err).To(MatchError(ContainSubstring(
						"aadClientId` must be specified if not using managed identity extension",
					)))
				})
				It("requires aadClientSecret", func() {
					values := []string{
						`#@data/values`,
						`---`,
						`deployment:`,
						`  args:`,
						`    - --source=ingress`,
						`    - --provider=azure`,
						`azure:`,
						`  subscriptionId: "azure-subscription-id"`,
						`  resourceGroup: "azure-resource-group"`,
						`  tenantId: "azure-tenant-id"`,
						`  aadClientId: "azure-aad-client-id"`,
					}
					valuesYaml := strings.Join(values, "\n")
					_, err := renderWithDataValuesInline(valuesYaml)
					Expect(err).To(MatchError(ContainSubstring(
						"aadClientSecret` must be specified if not using managed identity extension",
					)))
				})
			})
			When("the required fields are missing", func() {
				It("returns an error when resourceGroup is not provided", func() {
					values := []string{
						`#@data/values`,
						`---`,
						`deployment:`,
						`  args:`,
						`    - --source=ingress`,
						`    - --provider=azure`,
						`azure:`,
						`  subscriptionId: "azure-subscription-id"`,
						`  tenantId: "azure-tenant-id"`,
					}
					valuesYaml := strings.Join(values, "\n")
					_, err := renderWithDataValuesInline(valuesYaml)
					Expect(err).To(MatchError(ContainSubstring("resourceGroup` must be specified")))
				})
				It("returns an error when tenantId is not provided", func() {
					values := []string{
						`#@data/values`,
						`---`,
						`deployment:`,
						`  args:`,
						`    - --source=ingress`,
						`    - --provider=azure`,
						`azure:`,
						`  resourceGroup: "azure-resource-group"`,
						`  subscriptionId: "azure-subscription-id"`,
					}
					valuesYaml := strings.Join(values, "\n")
					_, err := renderWithDataValuesInline(valuesYaml)
					Expect(err).To(MatchError(ContainSubstring("tenantId` must be specified")))
				})
				It("returns an error when subscriptionID is not provided", func() {
					values := []string{
						`#@data/values`,
						`---`,
						`deployment:`,
						`  args:`,
						`    - --source=ingress`,
						`    - --provider=azure`,
						`azure:`,
						`  resourceGroup: "azure-resource-group"`,
						`  tenantId: "azure-tenant-id"`,
					}
					valuesYaml := strings.Join(values, "\n")
					_, err := renderWithDataValuesInline(valuesYaml)
					Expect(err).To(MatchError(ContainSubstring("subscriptionId` must be specified")))
				})
			})
		})
	})
})

func valueAtYAMLPath(doc, yamlPath string) string {
	var node yaml.Node
	err := yaml.Unmarshal([]byte(doc), &node)
	Expect(err).NotTo(HaveOccurred())

	path, err := yamlpath.NewPath(yamlPath)
	Expect(err).NotTo(HaveOccurred())

	q, err := path.Find(&node)
	Expect(err).NotTo(HaveOccurred())
	Expect(q).To(HaveLen(1))
	return q[0].Value
}

func renderWithDataValuesFixture(filename string) (string, error) {
	workingDir, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())
	path := filepath.Join(workingDir, "fixtures", "values", filename)
	configDir := filepath.Join(workingDir, "..", "..", "bundle", "config")
	return ytt.RenderYTTTemplate(ytt.CommandOptions{}, []string{configDir, path}, nil)
}

func renderWithDataValuesInline(content string) (string, error) {
	workingDir, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	configDir := filepath.Join(workingDir, "..", "..", "bundle", "config")
	return ytt.RenderYTTTemplate(ytt.CommandOptions{}, []string{configDir}, strings.NewReader(content))
}

func findDeploymentDoc(output string) string {
	deploymentDocs, err := FindDocsMatchingYAMLPath(output, map[string]string{
		".kind":          "Deployment",
		".metadata.name": "external-dns",
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(deploymentDocs).To(HaveLen(1))
	return deploymentDocs[0]
}
