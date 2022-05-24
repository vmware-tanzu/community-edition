// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package contour_test

import (
	"path/filepath"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/repo"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/ytt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Contour Ytt Templates", func() {
	var (
		values string
		output string
		err    error

		configDir = filepath.Join(repo.RootDir(), "addons/packages/contour/1.21.0/bundle/config")
	)

	BeforeEach(func() {
		values = ""
	})

	JustBeforeEach(func() {
		var filePaths []string

		for _, p := range []string{"upstream/*.yaml", "overlays/*.yaml", "*.yaml", "*.star"} {
			matches, err := filepath.Glob(filepath.Join(configDir, p))
			Expect(err).NotTo(HaveOccurred())
			filePaths = append(filePaths, matches...)
		}

		output, err = ytt.RenderYTTTemplate(ytt.CommandOptions{}, filePaths, strings.NewReader(values))
	})

	Context("No data values", func() {
		It("renders without an error", func() {
			_ = output
			Expect(err).NotTo(HaveOccurred())
		})
	})

	// START validation tests
	Context("Invalid instrastructure provider", func() {
		BeforeEach(func() {
			values = invalidInfrastructureProvider
		})

		It("renders with an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err).To(ContainSubstring("infrastructureProvider must be either docker or aws or vsphere or azure"))
		})
	})

	Context("No namespace", func() {
		BeforeEach(func() {
			values = noNamespace
		})

		It("renders with an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err).To(ContainSubstring("namespace must be provided"))
		})
	})

	Context("No contour replicas", func() {
		BeforeEach(func() {
			values = noContourReplicas
		})

		It("renders with an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err).To(ContainSubstring("contour.replicas must be provided"))
		})
	})

	Context("Cert manager enabled, no certificates duration", func() {
		BeforeEach(func() {
			values = noCertificatesDuration
		})

		It("renders with an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err).To(ContainSubstring("certificates.duration must be provided when certificates.useCertManager is true"))
		})
	})

	Context("Cert manager enabled, no certificates renew before", func() {
		BeforeEach(func() {
			values = noCertificatesRenewBefore
		})

		It("renders with an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err).To(ContainSubstring("certificates.renewBefore must be provided when certificates.useCertManager is true"))
		})
	})

	Context("Envoy host ports enabled, no http port", func() {
		BeforeEach(func() {
			values = noHostPortHTTP
		})

		It("renders with an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err).To(ContainSubstring("envoy.hostPorts.http must be provided when envoy.hostPorts.enable is true"))
		})
	})

	Context("Envoy host ports enabled, no https port", func() {
		BeforeEach(func() {
			values = noHostPortHTTPS
		})

		It("renders with an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err).To(ContainSubstring("envoy.hostPorts.https must be provided when envoy.hostPorts.enable is true"))
		})
	})

	Context("Invalid envoy log level", func() {
		BeforeEach(func() {
			values = invalidEnvoyLogLevel
		})

		It("renders with an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err).To(ContainSubstring("envoy.logLevel must be one of trace|debug|info|warning/warn|error|critical|off"))
		})
	})

	Context("No envoy termination grace period seconds", func() {
		BeforeEach(func() {
			values = noTerminationGracePeriodSeconds
		})

		It("renders with an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err).To(ContainSubstring("envoy.terminationGracePeriodSeconds must be provided"))
		})
	})

	Context("Invalid envoy service type", func() {
		BeforeEach(func() {
			values = invalidEnvoyServiceType
		})

		It("renders with an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err).To(ContainSubstring("envoy.service.type must be either LoadBalancer or NodePort or ClusterIP"))
		})
	})

	Context("Invalid envoy service external traffic policy", func() {
		BeforeEach(func() {
			values = invalidEnvoyServiceExternalTrafficPolicy
		})

		It("renders with an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err).To(ContainSubstring("envoy.service.externalTrafficPolicy must be either Cluster or Local"))
		})
	})

	Context("Invalid AWS load balancer type", func() {
		BeforeEach(func() {
			values = invalidAwsLoadBalancerType
		})

		It("renders with an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err).To(ContainSubstring("envoy.service.aws.loadBalancerType must be either classic or nlb when infrastructureProvider is aws"))
		})
	})
	// END validation tests

	// START Envoy service defaulting
	Context("Envoy service defaults when no infrastructure provider is specified", func() {
		BeforeEach(func() {
			values = ""
		})

		It("defaults the Envoy service type to LoadBalancer", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy service not found")

			service := unmarshalService(docs[0])
			Expect(service.Spec.Type).To(Equal(corev1.ServiceTypeLoadBalancer))
		})

		It("defaults the Envoy service external traffic policy to Local", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy service not found")

			service := unmarshalService(docs[0])
			Expect(service.Spec.ExternalTrafficPolicy).To(Equal(corev1.ServiceExternalTrafficPolicyTypeLocal))
		})
	})

	Context("Envoy service defaults for docker", func() {
		BeforeEach(func() {
			values = dockerInfrastructureProvider
		})

		It("defaults the Envoy service type to NodePort", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy service not found")

			service := unmarshalService(docs[0])
			Expect(service.Spec.Type).To(Equal(corev1.ServiceTypeNodePort))
		})

		It("defaults the Envoy service external traffic policy to Local", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy service not found")

			service := unmarshalService(docs[0])
			Expect(service.Spec.ExternalTrafficPolicy).To(Equal(corev1.ServiceExternalTrafficPolicyTypeLocal))
		})
	})

	Context("Envoy service defaults for aws", func() {
		BeforeEach(func() {
			values = awsInfrastructureProvider
		})

		It("defaults the Envoy service type to LoadBalancer", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy service not found")

			service := unmarshalService(docs[0])
			Expect(service.Spec.Type).To(Equal(corev1.ServiceTypeLoadBalancer))
		})

		It("defaults the Envoy service external traffic policy to Local", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy service not found")

			service := unmarshalService(docs[0])
			Expect(service.Spec.ExternalTrafficPolicy).To(Equal(corev1.ServiceExternalTrafficPolicyTypeLocal))
		})

		Context("Envoy service annotation defaults for AWS NLB", func() {
			BeforeEach(func() {
				values = awsLoadBalancerTypeNLB
			})

			It("sets the Envoy service annotations", func() {
				Expect(err).NotTo(HaveOccurred())

				docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
				Expect(docs).To(HaveLen(1), "Envoy service not found")

				service := unmarshalService(docs[0])
				Expect(service.Annotations).To(HaveKeyWithValue("service.beta.kubernetes.io/aws-load-balancer-type", "nlb"))
			})
		})

		Context("Envoy service annotation defaults for AWS ELB", func() {
			BeforeEach(func() {
				values = awsLoadBalancerTypeClassic
			})

			It("sets the Envoy service annotations", func() {
				Expect(err).NotTo(HaveOccurred())

				docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
				Expect(docs).To(HaveLen(1), "Envoy service not found")

				service := unmarshalService(docs[0])
				Expect(service.Annotations).To(HaveKeyWithValue("service.beta.kubernetes.io/aws-load-balancer-backend-protocol", "tcp"))
			})
		})

		Context("Envoy service annotation defaults for AWS ELB with PROXY protocol enabled", func() {
			BeforeEach(func() {
				values = awsLoadBalancerTypeClassicProxyProtocolEnabled
			})

			It("sets the Envoy service annotations", func() {
				Expect(err).NotTo(HaveOccurred())

				docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
				Expect(docs).To(HaveLen(1), "Envoy service not found")

				service := unmarshalService(docs[0])
				Expect(service.Annotations).To(HaveKeyWithValue("service.beta.kubernetes.io/aws-load-balancer-backend-protocol", "tcp"))
				Expect(service.Annotations).To(HaveKeyWithValue("service.beta.kubernetes.io/aws-load-balancer-proxy-protocol", "*"))
			})
		})

	})

	Context("Envoy service defaults for vsphere", func() {
		BeforeEach(func() {
			values = vsphereInfrastructureProvider
		})

		It("defaults the Envoy service type to NodePort", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy service not found")

			service := unmarshalService(docs[0])
			Expect(service.Spec.Type).To(Equal(corev1.ServiceTypeNodePort))
		})

		It("defaults the Envoy service external traffic policy to Cluster", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy service not found")

			service := unmarshalService(docs[0])
			Expect(service.Spec.ExternalTrafficPolicy).To(Equal(corev1.ServiceExternalTrafficPolicyTypeCluster))
		})
	})

	Context("Envoy service defaults for azure", func() {
		BeforeEach(func() {
			values = azureInfrastructureProvider
		})

		It("defaults the Envoy service type to LoadBalancer", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy service not found")

			service := unmarshalService(docs[0])
			Expect(service.Spec.Type).To(Equal(corev1.ServiceTypeLoadBalancer))
		})

		It("defaults the Envoy service external traffic policy to Local", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy service not found")

			service := unmarshalService(docs[0])
			Expect(service.Spec.ExternalTrafficPolicy).To(Equal(corev1.ServiceExternalTrafficPolicyTypeLocal))
		})
	})
	// END Envoy service defaulting

	// START general data values & overlay logic
	Context("Namespace set", func() {
		BeforeEach(func() {
			values = nonDefaultNamespace
		})

		It("changes all resources' namespaces", func() {
			Expect(err).NotTo(HaveOccurred())

			// Have to be fiddly with the whitespace here to avoid matching
			// lines in the config file example and in the ContourConfiguration CRD
			// schema.

			docs := findDocsContainingLines(output, "\n  namespace: projectcontour")
			Expect(docs).To(BeEmpty())

			// Note, the overlay superficially indicates 16 matches (13 resources in the
			// namespace + 3 instances of the namespace in the Cluster/RoleBinding subject
			// references), but we expect 14 documents because the ClusterRoleBinding is
			// found when performing line matching.
			docs = findDocsContainingLines(output, "\n  namespace: non-default-namespace")
			Expect(docs).To(HaveLen(14))

			docs = findDocsContainingLines(output, "kind: Namespace")
			Expect(docs).To(HaveLen(1))
			Expect(docs[0]).To(ContainSubstring("name: non-default-namespace"))
		})
	})

	Context("Contour config file contents set", func() {
		// Note, the data values for the tests here require the "#@overlay/replace"
		// directive to be specified above "configFileContents" in order for these tests to
		// pass. This is *not* required when passing in a data values file to ytt via
		// --data-values-file. I believe it's just due to a difference in how
		// ytt.RenderYTTTemplate is passing the data values in -- i.e. not via --data-values-file.

		BeforeEach(func() {
			values = contourConfigFileContents
		})

		It("sets Contour's config map contents", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: ConfigMap")
			Expect(docs).To(HaveLen(1))

			configmap := unmarshalConfigMap(docs[0])
			Expect(configmap.Data).To(HaveKeyWithValue("contour.yaml", "foo:\n  bar: baz\nboo: bam\n"))
		})
	})

	Context("Contour replicas set", func() {
		BeforeEach(func() {
			values = contourReplicasThree
		})

		It("sets the Contour deployment's replicas", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Deployment", "name: contour")
			Expect(docs).To(HaveLen(1), "Contour deployment not found")

			deployment := unmarshalDeployment(docs[0])
			Expect(*deployment.Spec.Replicas).To(Equal(int32(3)))
		})
	})

	Context("Contour use proxy protocol enabled", func() {
		BeforeEach(func() {
			values = contourUseProxyProtocolEnabled
		})

		It("sets the Contour deployment's --use-proxy-protocol flag", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Deployment", "name: contour")
			Expect(docs).To(HaveLen(1), "Contour deployment not found")

			deployment := unmarshalDeployment(docs[0])
			Expect(deployment.Spec.Template.Spec.Containers[0].Args).To(ContainElement("--use-proxy-protocol"))
		})
	})

	Context("Contour log level set", func() {
		BeforeEach(func() {
			values = contourLogLevelDebug
		})

		It("sets the Contour deployment's --debug flag", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Deployment", "name: contour")
			Expect(docs).To(HaveLen(1), "Contour deployment not found")

			deployment := unmarshalDeployment(docs[0])
			Expect(deployment.Spec.Template.Spec.Containers[0].Args).To(ContainElement("--debug"))
		})
	})

	// TODO(sk) test this for all infrastructure providers?
	Context("Envoy service type explicitly specified", func() {
		BeforeEach(func() {
			values = vsphereClusterIPEnvoyService
		})

		It("sets the Envoy service type to the value specified", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy service not found")

			service := unmarshalService(docs[0])
			Expect(service.Spec.Type).To(Equal(corev1.ServiceType("ClusterIP")))
		})
	})

	Context("Envoy service load balancer IP specified", func() {
		BeforeEach(func() {
			values = envoyServiceLoadBalancerIP
		})

		It("sets the Envoy service load balancer IP to the value specified", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy service not found")

			service := unmarshalService(docs[0])
			Expect(service.Spec.LoadBalancerIP).To(Equal("7.7.7.7"))
		})
	})

	// TODO(sk) test this for all infrastructure providers?
	Context("Envoy service external traffic policy explicitly specified", func() {
		BeforeEach(func() {
			values = vsphereLocalExternalTrafficPolicy
		})

		It("sets the Envoy service external traffic policy to the value specified", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy service not found")

			service := unmarshalService(docs[0])
			Expect(service.Spec.ExternalTrafficPolicy).To(Equal(corev1.ServiceExternalTrafficPolicyTypeLocal))
		})
	})

	Context("Envoy service annotations specified", func() {
		// Note, the data values for the tests here require the "#@overlay/replace"
		// directive to be specified above "annotations" in order for these tests to
		// pass. This is *not* required when passing in a data values file to ytt via
		// --data-values-file. I believe it's just due to a difference in how
		// ytt.RenderYTTTemplate is passing the data values in -- i.e. not via --data-values-file.

		BeforeEach(func() {
			values = envoyServiceAnnotations
		})

		It("sets the Envoy service annotations", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy service not found")

			service := unmarshalService(docs[0])
			Expect(service.Annotations).To(HaveKeyWithValue("foo", "bar"))
			Expect(service.Annotations).To(HaveKeyWithValue("boo", "baz"))

		})

		Context("Envoy service also gets default annotations", func() {
			BeforeEach(func() {
				values = envoyServiceAnnotationsAWSNLB
			})

			It("sets the Envoy service annotations to the union of the default and explicit annotations", func() {
				Expect(err).NotTo(HaveOccurred())

				docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
				Expect(docs).To(HaveLen(1), "Envoy service not found")

				service := unmarshalService(docs[0])
				Expect(service.Annotations).To(HaveKeyWithValue("service.beta.kubernetes.io/aws-load-balancer-type", "nlb"))
				Expect(service.Annotations).To(HaveKeyWithValue("foo", "bar"))
				Expect(service.Annotations).To(HaveKeyWithValue("boo", "baz"))

			})
		})
	})

	Context("Envoy service node ports set", func() {
		BeforeEach(func() {
			values = envoyServiceNodePorts
		})

		It("sets the Envoy service node ports", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Service", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy service not found")

			service := unmarshalService(docs[0])
			Expect(service.Spec.Ports[0].NodePort).To(Equal(int32(30080)))
			Expect(service.Spec.Ports[1].NodePort).To(Equal(int32(30443)))
		})
	})

	Context("Envoy host ports set", func() {
		BeforeEach(func() {
			values = envoyHostPorts
		})

		It("sets Envoy daemonset host ports", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: DaemonSet", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy daemonset not found")

			daemonset := unmarshalDaemonSet(docs[0])
			Expect(daemonset.Spec.Template.Spec.Containers[1].Ports[0].Name).To(Equal("http"))
			Expect(daemonset.Spec.Template.Spec.Containers[1].Ports[0].HostPort).To(Equal(int32(80)))
			Expect(daemonset.Spec.Template.Spec.Containers[1].Ports[1].Name).To(Equal("https"))
			Expect(daemonset.Spec.Template.Spec.Containers[1].Ports[1].HostPort).To(Equal(int32(443)))
		})
	})

	Context("Envoy host networking enabled", func() {
		BeforeEach(func() {
			values = envoyHostNetworkingEnabled
		})

		It("enables host networking for the Envoy daemonset", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: DaemonSet", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy daemonset not found")

			daemonset := unmarshalDaemonSet(docs[0])
			Expect(daemonset.Spec.Template.Spec.HostNetwork).To(BeTrue())
			Expect(daemonset.Spec.Template.Spec.DNSPolicy).To(Equal(corev1.DNSClusterFirstWithHostNet))
		})
	})

	Context("Envoy termination grace period seconds set", func() {
		BeforeEach(func() {
			values = envoyTerminationGracePeriodSeconds
		})

		It("sets termination grace period seconds for the Envoy daemonset", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: DaemonSet", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy daemonset not found")

			daemonset := unmarshalDaemonSet(docs[0])
			Expect(*daemonset.Spec.Template.Spec.TerminationGracePeriodSeconds).To(Equal(int64(777)))
		})
	})

	Context("Envoy log level set", func() {
		BeforeEach(func() {
			values = envoyLogLevelDebug
		})

		It("sets the Envoy daemonset's log level to debug", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: DaemonSet", "name: envoy")
			Expect(docs).To(HaveLen(1), "Envoy daemonset not found")

			daemonset := unmarshalDaemonSet(docs[0])
			Expect(daemonset.Spec.Template.Spec.Containers[1].Args).To(ContainElement("--log-level debug"))
		})
	})

	Context("Use cert manager enabled", func() {
		BeforeEach(func() {
			values = useCertManagerEnabled
		})

		It("uses cert-manager to provision certs", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Job")
			Expect(docs).To(BeEmpty(), "Expected to not find certgen job")

			docs = findDocsContainingLines(output, "apiVersion: cert-manager.io/v1")
			// TODO(sk) this assertion could be improved.
			Expect(docs).To(HaveLen(5), "Expected 5 cert-manager resources")
		})
	})

	Context("Cert manager duration and renew before set", func() {
		BeforeEach(func() {
			values = useCertManagerEnabledDurationRenewBeforeSet
		})

		It("sets the certificate's duration and renew before values", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Certificate")
			Expect(docs).To(HaveLen(3), "Certificates not found")

			for _, cert := range docs {
				// TODO(sk) this could be improved by unmarshaling to a Go
				// struct, requires importing the cert-manager API though.
				Expect(cert).To(ContainSubstring("\n  duration: 777h\n"))
				Expect(cert).To(ContainSubstring("\n  renewBefore: 77h\n"))
			}
		})
	})

	Context("IPv6 enabled", func() {
		BeforeEach(func() {
			values = ipv6Enabled
		})

		It("sets the appropriate contour serve flags", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Deployment")
			Expect(docs).To(HaveLen(1), "Contour Deployment not found")

			deployment := unmarshalDeployment(docs[0])
			Expect(deployment.Spec.Template.Spec.Containers).To(HaveLen(1))
			Expect(deployment.Spec.Template.Spec.Containers[0].Args).To(ContainElements(
				"--xds-address=::",
				"--stats-address=::",
				"--debug-http-address=::1",
				"--http-address=::",
				"--envoy-service-http-address=::",
				"--envoy-service-https-address=::",
				"--health-address=::",
			))
			Expect(deployment.Spec.Template.Spec.Containers[0].Args).ToNot(ContainElements(ContainSubstring("0.0.0.0")))
			Expect(deployment.Spec.Template.Spec.Containers[0].Args).ToNot(ContainElements(ContainSubstring("127.0.0.1")))
		})
	})

	Context("IPv6 disabled", func() {
		BeforeEach(func() {
			values = ipv6Disabled
		})

		It("sets the appropriate contour serve flags", func() {
			Expect(err).NotTo(HaveOccurred())

			docs := findDocsContainingLines(output, "kind: Deployment")
			Expect(docs).To(HaveLen(1), "Contour Deployment not found")

			deployment := unmarshalDeployment(docs[0])
			Expect(deployment.Spec.Template.Spec.Containers).To(HaveLen(1))
			Expect(deployment.Spec.Template.Spec.Containers[0].Args).To(ContainElement("--xds-address=0.0.0.0"))
			Expect(deployment.Spec.Template.Spec.Containers[0].Args).ToNot(ContainElements(ContainSubstring("::")))
		})
	})
	// END general data values & overlay logic
})

func findDocsContainingLines(output string, lines ...string) []string {
	var docs []string
	for _, doc := range strings.Split(output, "---") {
		matches := true

		for _, line := range lines {
			if !strings.Contains(doc, line+"\n") {
				matches = false
				break
			}
		}
		if matches {
			docs = append(docs, doc)
		}
	}

	return docs
}

func unmarshalService(doc string) *corev1.Service {
	service := &corev1.Service{}
	err := yaml.Unmarshal([]byte(doc), service)
	Expect(err).NotTo(HaveOccurred())
	return service
}

func unmarshalDeployment(doc string) *appsv1.Deployment {
	deployment := &appsv1.Deployment{}
	err := yaml.Unmarshal([]byte(doc), deployment)
	Expect(err).NotTo(HaveOccurred())
	return deployment
}

func unmarshalDaemonSet(doc string) *appsv1.DaemonSet {
	daemonset := &appsv1.DaemonSet{}
	err := yaml.Unmarshal([]byte(doc), daemonset)
	Expect(err).NotTo(HaveOccurred())
	return daemonset
}

func unmarshalConfigMap(doc string) *corev1.ConfigMap {
	configmap := &corev1.ConfigMap{}
	err := yaml.Unmarshal([]byte(doc), configmap)
	Expect(err).NotTo(HaveOccurred())
	return configmap
}
