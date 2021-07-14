package e2e_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("fluent-bit Addon E2E Test", func() {

	/*BeforeEach(func() {
		By("installing fluent-bit addon package")
		fluentBitAddonValuesFilename, err := readFileAndReplaceContentsTempFile(filepath.Join("fixtures", "fluent-bit.tce.vmware.com-values.yaml"),
			map[string]string{
				"ADDON_NAMESPACE": addonNamespace,
			},
		)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(fluentBitAddonValuesFilename)

		output, err := tanzu(nil, "package", "install", "fluent-bit.tce.vmware.com", "-n", packageNamespace, "-g", fluentBitAddonValuesFilename)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).Should(ContainSubstring("Installed "))

		By("validating everything is ready")
		validatePackageReady(packageNamespace, "fluent-bit.tce.vmware.com")
	})*/

	It("has to create", func() {
		By("installing fluent-bit addon package")
		output, err := tanzu(nil, "package", "install", "fluent-bit.tce.vmware.com", "-n", packageNamespace)
		Expect(err).NotTo(HaveOccurred())
		//Expect(output).Should(ContainSubstring("Installed "))
		print("Output is" + output)

		By("validating everything is ready")
		validatePackageReady(packageNamespace, "fluent-bit.tce.vmware.com")
	})

	JustAfterEach(func() {
		if CurrentGinkgoTestDescription().Failed {
			fmt.Fprintf(GinkgoWriter, "\nCollecting diagnostic information just after test failure\n")
			_, _ = kubectl(nil, "-n", packageNamespace, "get", "all,installedpackage,apps")
			_, _ = kubectl(nil, "-n", addonNamespace, "get", "all")
			_, _ = kubectl(nil, "-n", packageNamespace, "get", "app", "fluent-bit.tce.vmware.com", "-o", "jsonpath={.status}")
			_, _ = kubectl(nil, "-n", addonNamespace, "get", "deployment", "fluent-bit", "-o", "jsonpath={.status}")
		}
	})

	AfterEach(func() {
		By("cleaning up fluent-bit addon package")
		_, err := tanzu(nil, "package", "delete", "fluent-bit.tce.vmware.com", "-n", packageNamespace)
		Expect(err).NotTo(HaveOccurred())
		validatePackageNotFound(packageNamespace, "fluent-bit.tce.vmware.com")
	})

	/*It("journeys through the external dns addon lifecycle", func() {
		By("sending an HTTP Request to the kuard deployment using the FQDN in a pod within the cluster")
		Eventually(func() (string, error) {
			return kubectl(nil, "-n", fixtureNamespace, "exec", "dnsutils", "--", "wget", "-O", "-", "http://kuard.k8s.example.org")
		}, httpRequestTimeout, httpRequestInterval).Should(ContainSubstring("KUAR Demo"))
	})*/
})
