package e2e_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("External-dns Addon E2E Test", func() {
	var (
		bindDeployment string
		dnsutilsPod    string
	)

	BeforeEach(func() {
		By("installing bind deployment")
		corednsServiceClusterIP, err := kubectl(nil, "-n", "kube-system", "get", "service", "kube-dns", "-o", "jsonpath={.spec.clusterIP}")
		Expect(err).NotTo(HaveOccurred())

		bindDeployment, err = readFileAndReplaceContents(filepath.Join("fixtures", "bind-deployment.yaml"), map[string]string{
			"docker.io":          dockerhubProxy,
			"COREDNS_CLUSTER_IP": corednsServiceClusterIP,
		})
		Expect(err).NotTo(HaveOccurred())

		_, err = kubectl(bytes.NewBufferString(bindDeployment), "-n", fixtureNamespace, "apply", "-f", "-")
		Expect(err).NotTo(HaveOccurred())

		validateDeploymentReady(fixtureNamespace, "bind")

		bindServiceClusterIP, err := kubectl(nil, "-n", fixtureNamespace, "get", "service", "bind", "-o", "jsonpath={.spec.clusterIP}")
		Expect(err).NotTo(HaveOccurred())

		By("installing kuard deployment")
		_, err = kubectl(nil, "-n", fixtureNamespace, "apply", "-f", filepath.Join("fixtures", "kuard-deployment.yaml"))
		Expect(err).NotTo(HaveOccurred())

		By("installing dnsutils pod")
		dnsutilsPod, err = readFileAndReplaceContents(filepath.Join("fixtures", "dnsutils-pod.yaml"),
			map[string]string{
				"RFC2136_HOST": bindServiceClusterIP,
			},
		)
		Expect(err).NotTo(HaveOccurred())

		_, err = kubectl(bytes.NewBufferString(dnsutilsPod), "-n", fixtureNamespace, "apply", "-f", "-")
		Expect(err).NotTo(HaveOccurred())

		By("installing external-dns addon package")
		externalDNSAddonValuesFilename, err := readFileAndReplaceContentsTempFile(filepath.Join("fixtures", "external-dns.tce.vmware.com-values.yaml"),
			map[string]string{
				"ADDON_NAMESPACE":  addonNamespace,
				"SOURCE_NAMESPACE": fixtureNamespace,
				"RFC2136_HOST":     bindServiceClusterIP,
			},
		)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(externalDNSAddonValuesFilename)

		_, err = tanzu(nil, "package", "install", "external-dns.tce.vmware.com", "-n", packageNamespace, "-g", externalDNSAddonValuesFilename)
		Expect(err).NotTo(HaveOccurred())

		By("validating everything is ready")
		validateDeploymentReady(fixtureNamespace, "kuard")
		validatePodReady(fixtureNamespace, "dnsutils")
		validatePackageReady(packageNamespace, "external-dns.tce.vmware.com")
	})

	JustAfterEach(func() {
		if CurrentGinkgoTestDescription().Failed {
			fmt.Fprintf(GinkgoWriter, "\nCollecting diagnostic information just after test failure\n")
			_, _ = kubectl(nil, "-n", packageNamespace, "get", "all,installedpackage,apps")
			_, _ = kubectl(nil, "-n", addonNamespace, "get", "all")
			_, _ = kubectl(nil, "-n", fixtureNamespace, "get", "all")
			_, _ = kubectl(nil, "-n", fixtureNamespace, "get", "deployment", "bind", "-o", "jsonpath={.status}")
			_, _ = kubectl(nil, "-n", packageNamespace, "get", "app", "external-dns.tce.vmware.com", "-o", "jsonpath={.status}")
			_, _ = kubectl(nil, "-n", addonNamespace, "get", "deployment", "external-dns", "-o", "jsonpath={.status}")
			_, _ = kubectl(nil, "-n", fixtureNamespace, "get", "deployment", "kuard", "-o", "jsonpath={.status}")
		}
	})

	AfterEach(func() {
		By("cleaning up external-dns addon package")
		_, err := tanzu(nil, "package", "delete", "external-dns.tce.vmware.com", "-n", packageNamespace)
		Expect(err).NotTo(HaveOccurred())

		By("cleaning up dnsutils pod")
		_, err = kubectl(nil, "-n", fixtureNamespace, "delete", "-f", filepath.Join("fixtures", "dnsutils-pod.yaml"))
		Expect(err).NotTo(HaveOccurred())

		By("cleaning up kuard deployment")
		_, err = kubectl(nil, "-n", fixtureNamespace, "delete", "-f", filepath.Join("fixtures", "kuard-deployment.yaml"))
		Expect(err).NotTo(HaveOccurred())

		By("cleaning up bind deployment")
		_, err = kubectl(bytes.NewBufferString(bindDeployment), "-n", fixtureNamespace, "delete", "-f", "-")
		Expect(err).NotTo(HaveOccurred())

		validatePodNotFound(fixtureNamespace, "dnsutils")
		validateDeploymentNotFound(fixtureNamespace, "kuard")
		validateDeploymentNotFound(fixtureNamespace, "bind")
		validatePackageNotFound(packageNamespace, "external-dns.tce.vmware.com")
	})

	It("journeys through the external dns addon lifecycle", func() {
		By("sending an HTTP Request to the kuard deployment using the FQDN in a pod within the cluster")
		Eventually(func() (string, error) {
			return kubectl(nil, "-n", fixtureNamespace, "exec", "dnsutils", "--", "wget", "-O", "-", "http://kuard.k8s.example.org")
		}, httpRequestTimeout, httpRequestInterval).Should(ContainSubstring("KUAR Demo"))
	})
})
