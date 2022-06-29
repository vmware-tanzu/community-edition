// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// package e2e_test implements running the external DNS end to end tests
package e2e_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	version = "0.11.0"
)

var _ = Describe("External-dns Addon E2E Test", func() {
	var (
		bindDeployment       string
		dnsutilsPod          string
		packageName          string
		bindServiceClusterIP string
		packageInstallName   = "external-dns"
	)

	BeforeEach(func() {
		packageName = utils.TanzuPackageName("external-dns")

		By("installing bind deployment")
		corednsServiceClusterIP, err := utils.Kubectl(nil, "-n", "kube-system", "get", "service", "kube-dns", "-o", "jsonpath={.spec.clusterIP}")
		Expect(err).NotTo(HaveOccurred())

		bindDeployment, err = utils.ReadFileAndReplaceContents(filepath.Join("fixtures", "bind-deployment.yaml"),
			map[string]string{
				"docker.io":          dockerhubProxy,
				"COREDNS_CLUSTER_IP": corednsServiceClusterIP,
			},
		)
		Expect(err).NotTo(HaveOccurred())

		_, err = utils.Kubectl(bytes.NewBufferString(bindDeployment), "-n", fixtureNamespace, "apply", "-f", "-")
		Expect(err).NotTo(HaveOccurred())

		utils.ValidateDeploymentReady(fixtureNamespace, "bind")

		bindServiceClusterIP, err = utils.Kubectl(nil, "-n", fixtureNamespace, "get", "service", "bind", "-o", "jsonpath={.spec.clusterIP}")
		Expect(err).NotTo(HaveOccurred())

		By("installing kuard deployment")
		_, err = utils.Kubectl(nil, "-n", fixtureNamespace, "apply", "-f", filepath.Join("fixtures", "kuard-deployment.yaml"))
		Expect(err).NotTo(HaveOccurred())

		By("installing dnsutils pod")
		dnsutilsPod, err = utils.ReadFileAndReplaceContents(filepath.Join("fixtures", "dnsutils-pod.yaml"),
			map[string]string{
				"BIND_SERVER_ADDRESS": bindServiceClusterIP,
			},
		)
		Expect(err).NotTo(HaveOccurred())

		_, err = utils.Kubectl(bytes.NewBufferString(dnsutilsPod), "-n", fixtureNamespace, "apply", "-f", "-")
		Expect(err).NotTo(HaveOccurred())

		By("installing external-dns addon package")
		valuesFilename, err := utils.ReadFileAndReplaceContentsTempFile(filepath.Join("fixtures", "external-dns-values-fixture.yaml"),
			map[string]string{
				"PACKAGE_COMPONENTS_NAMESPACE":   packageComponentsNamespace,
				"EXTERNAL_DNS_SOURCES_NAMESPACE": fixtureNamespace,
				"BIND_SERVER_ADDRESS":            bindServiceClusterIP,
			},
		)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(valuesFilename)

		_, err = utils.Tanzu(nil, "package", "install", packageInstallName,
			"--namespace", packageInstallNamespace,
			"--package-name", packageName,
			"--version", utils.TanzuPackageAvailableVersionWithVersionSubString(packageName, version),
			"--values-file", valuesFilename)
		Expect(err).NotTo(HaveOccurred())

		By("validating everything is ready")
		utils.ValidateDeploymentReady(fixtureNamespace, "kuard")
		utils.ValidateLoadBalancerReady(fixtureNamespace, "kuard")
		utils.ValidatePodReady(fixtureNamespace, "dnsutils")
		utils.ValidatePackageInstallReady(packageInstallNamespace, packageInstallName)
	})

	JustAfterEach(func() {
		if CurrentGinkgoTestDescription().Failed {
			fmt.Fprintf(GinkgoWriter, "\nCollecting diagnostic information just after test failure\n")
			fmt.Fprintf(GinkgoWriter, "\nResources summary:\n")
			_, _ = utils.Kubectl(nil, "-n", packageInstallNamespace, "get", "all,packageinstalls,apps")
			_, _ = utils.Kubectl(nil, "-n", packageComponentsNamespace, "get", "all")
			_, _ = utils.Kubectl(nil, "-n", fixtureNamespace, "get", "all")

			fmt.Fprintf(GinkgoWriter, "\nbind deployment status:\n")
			_, _ = utils.Kubectl(nil, "-n", fixtureNamespace, "get", "deployment", "bind", "-o", "jsonpath={.status}")

			fmt.Fprintf(GinkgoWriter, "\nkuard deployment status:\n")
			_, _ = utils.Kubectl(nil, "-n", fixtureNamespace, "get", "deployment", "kuard", "-o", "jsonpath={.status}")

			fmt.Fprintf(GinkgoWriter, "\npackage install status:\n")
			_, _ = utils.Kubectl(nil, "-n", packageInstallNamespace, "get", "app", packageInstallName, "-o", "jsonpath={.status}")

			fmt.Fprintf(GinkgoWriter, "\npackage components status:\n")
			_, _ = utils.Kubectl(nil, "-n", packageComponentsNamespace, "get", "deployment", "external-dns", "-o", "jsonpath={.status}")

			fmt.Fprintf(GinkgoWriter, "\nexternal-dns logs:\n")
			_, _ = utils.Kubectl(nil, "-n", packageComponentsNamespace, "logs", "-l", "app=external-dns")
		}
	})

	AfterEach(func() {
		By("cleaning up external-dns addon package")
		_, err := utils.Tanzu(nil, "package", "installed", "delete", packageInstallName, "--namespace", packageInstallNamespace, "--yes")
		Expect(err).NotTo(HaveOccurred())

		By("cleaning up dnsutils pod")
		_, err = utils.Kubectl(nil, "-n", fixtureNamespace, "delete", "-f", filepath.Join("fixtures", "dnsutils-pod.yaml"))
		Expect(err).NotTo(HaveOccurred())

		By("cleaning up kuard deployment")
		_, err = utils.Kubectl(nil, "-n", fixtureNamespace, "delete", "-f", filepath.Join("fixtures", "kuard-deployment.yaml"))
		Expect(err).NotTo(HaveOccurred())

		By("cleaning up bind deployment")
		_, err = utils.Kubectl(bytes.NewBufferString(bindDeployment), "-n", fixtureNamespace, "delete", "-f", "-")
		Expect(err).NotTo(HaveOccurred())

		By("validating that dnsutils, kuard, and bind no longer exist")
		utils.ValidatePodNotFound(fixtureNamespace, "dnsutils")
		utils.ValidateDeploymentNotFound(fixtureNamespace, "kuard")
		utils.ValidateDeploymentNotFound(fixtureNamespace, "bind")

		By("validating that package install no longer exists")
		utils.ValidatePackageInstallNotFound(packageInstallNamespace, packageInstallName)
	})

	It("journeys through the external dns addon lifecycle", func() {
		By("sending an HTTP Request to the kuard deployment using the FQDN in a pod within the cluster")
		Eventually(func() (string, error) {
			return utils.Kubectl(nil, "-n", fixtureNamespace, "exec", "dnsutils", "--", "wget", "-O", "-", "http://kuard.k8s.example.org")
		}, httpRequestTimeout, httpRequestInterval).Should(ContainSubstring("KUAR Demo"))
	})

	Context("when applying a custom overlay", func() {
		BeforeEach(func() {
			By("applying custom overlay and custom schema secret")
			customOverlayAndSchemaSecret, err := utils.ReadFileAndReplaceContentsTempFile(filepath.Join("fixtures", "custom-overlay-and-schema-secret.yaml"),
				map[string]string{
					"PACKAGE_INSTALL_NAMESPACE": packageInstallNamespace,
				},
			)

			Expect(err).NotTo(HaveOccurred())
			defer os.Remove(customOverlayAndSchemaSecret)

			_, err = utils.Kubectl(nil, "apply", "-f", customOverlayAndSchemaSecret)
			Expect(err).NotTo(HaveOccurred())

			By("configure package install to use a custom overlay")
			_, err = utils.Kubectl(nil, "annotate", "packageinstall", "-n", packageInstallNamespace, packageInstallName,
				"ext.packaging.carvel.dev/ytt-paths-from-secret-name.0=external-dns-overlay-secret",
			)
			Expect(err).NotTo(HaveOccurred())

			By("updating the installed package with additional schema and data value")
			customValuesFilename, err := utils.ReadFileAndReplaceContentsTempFile(filepath.Join("fixtures", "custom-data-values.yaml"),
				map[string]string{
					"PACKAGE_COMPONENTS_NAMESPACE":   packageComponentsNamespace,
					"EXTERNAL_DNS_SOURCES_NAMESPACE": fixtureNamespace,
					"BIND_SERVER_ADDRESS":            bindServiceClusterIP,
				},
			)
			Expect(err).NotTo(HaveOccurred())
			defer os.Remove(customValuesFilename)

			Eventually(func() error {
				_, err = utils.Tanzu(nil, "package", "installed", "update", packageInstallName,
					"--namespace", packageInstallNamespace,
					"--package-name", packageName,
					"--version", utils.TanzuPackageAvailableVersionWithVersionSubString(packageName, version),
					"--values-file", customValuesFilename)
				return err
			}, "30s", "5s").Should(Succeed())

			By("validating package install is ready")
			utils.ValidatePackageInstallReady(packageInstallNamespace, packageInstallName)
		})
		It("applies the customizations", func() {
			Eventually(func() (string, error) {
				return utils.Kubectl(nil, "-n", fixtureNamespace, "exec", "dnsutils", "--", "wget", "-O", "-", "http://kuard.k8s.example.org")
			}, httpRequestTimeout, httpRequestInterval).Should(ContainSubstring("KUAR Demo"))

			Eventually(func() (string, error) {
				return utils.Kubectl(nil, "get", "deployment", "-n", packageComponentsNamespace, packageInstallName, "-o", "jsonpath={.metadata.labels.custom-label}")
			}).Should(Equal("customized-label"))
		})
	})
})
