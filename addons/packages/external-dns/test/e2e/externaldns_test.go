package e2e_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vmware-tanzu/tce/addons/packages/test/pkg/utils"

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

		bindServiceClusterIP, err := utils.Kubectl(nil, "-n", fixtureNamespace, "get", "service", "bind", "-o", "jsonpath={.spec.clusterIP}")
		Expect(err).NotTo(HaveOccurred())

		By("installing kuard deployment")
		_, err = utils.Kubectl(nil, "-n", fixtureNamespace, "apply", "-f", filepath.Join("fixtures", "kuard-deployment.yaml"))
		Expect(err).NotTo(HaveOccurred())

		By("installing dnsutils pod")
		dnsutilsPod, err = utils.ReadFileAndReplaceContents(filepath.Join("fixtures", "dnsutils-pod.yaml"),
			map[string]string{
				"RFC2136_HOST": bindServiceClusterIP,
			},
		)
		Expect(err).NotTo(HaveOccurred())

		_, err = utils.Kubectl(bytes.NewBufferString(dnsutilsPod), "-n", fixtureNamespace, "apply", "-f", "-")
		Expect(err).NotTo(HaveOccurred())

		By("installing external-dns addon package")
		externalDNSAddonValuesFilename, err := utils.ReadFileAndReplaceContentsTempFile(filepath.Join("fixtures", "external-dns.tce.vmware.com-values.yaml"),
			map[string]string{
				"ADDON_NAMESPACE":  addonNamespace,
				"SOURCE_NAMESPACE": fixtureNamespace,
				"RFC2136_HOST":     bindServiceClusterIP,
			},
		)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(externalDNSAddonValuesFilename)

		_, err = utils.Tanzu(nil, "package", "install", "external-dns.tce.vmware.com", "-n", packageNamespace, "-g", externalDNSAddonValuesFilename)
		Expect(err).NotTo(HaveOccurred())

		By("validating everything is ready")
		utils.ValidateDeploymentReady(fixtureNamespace, "kuard")
		utils.ValidatePodReady(fixtureNamespace, "dnsutils")
		utils.ValidatePackageReady(packageNamespace, "external-dns.tce.vmware.com")
	})

	JustAfterEach(func() {
		if CurrentGinkgoTestDescription().Failed {
			fmt.Fprintf(GinkgoWriter, "\nCollecting diagnostic information just after test failure\n")
			_, _ = utils.Kubectl(nil, "-n", packageNamespace, "get", "all,installedpackage,apps")
			_, _ = utils.Kubectl(nil, "-n", addonNamespace, "get", "all")
			_, _ = utils.Kubectl(nil, "-n", fixtureNamespace, "get", "all")
			_, _ = utils.Kubectl(nil, "-n", fixtureNamespace, "get", "deployment", "bind", "-o", "jsonpath={.status}")
			_, _ = utils.Kubectl(nil, "-n", packageNamespace, "get", "app", "external-dns.tce.vmware.com", "-o", "jsonpath={.status}")
			_, _ = utils.Kubectl(nil, "-n", addonNamespace, "get", "deployment", "external-dns", "-o", "jsonpath={.status}")
			_, _ = utils.Kubectl(nil, "-n", fixtureNamespace, "get", "deployment", "kuard", "-o", "jsonpath={.status}")
		}
	})

	AfterEach(func() {
		By("cleaning up external-dns addon package")
		_, err := utils.Tanzu(nil, "package", "delete", "external-dns.tce.vmware.com", "-n", packageNamespace)
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

		utils.ValidatePodNotFound(fixtureNamespace, "dnsutils")
		utils.ValidateDeploymentNotFound(fixtureNamespace, "kuard")
		utils.ValidateDeploymentNotFound(fixtureNamespace, "bind")
		utils.ValidatePackageNotFound(packageNamespace, "external-dns.tce.vmware.com")
	})

	It("journeys through the external dns addon lifecycle", func() {
		By("sending an HTTP Request to the kuard deployment using the FQDN in a pod within the cluster")
		Eventually(func() (string, error) {
			return utils.Kubectl(nil, "-n", fixtureNamespace, "exec", "dnsutils", "--", "wget", "-O", "-", "http://kuard.k8s.example.org")
		}, httpRequestTimeout, httpRequestInterval).Should(ContainSubstring("KUAR Demo"))
	})
})
