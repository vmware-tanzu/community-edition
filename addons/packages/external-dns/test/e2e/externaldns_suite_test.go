package e2e_test

import (
	"os"
	"testing"
	"time"

	"github.com/vmware-tanzu/tce/addons/packages/test/pkg/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExternalDNSE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "External-DNS Addon Package E2E Test Suite")
}

const (
	httpRequestTimeout  = 60 * time.Second
	httpRequestInterval = 5 * time.Second
)

var (
	// dockerhubProxy is an optional configuration option (provided by using
	// DOCKERHUB_PROXY), that allows you to override docker.io with a proxy to
	// docker.io to avoid any potential issues with rate-limiting.
	dockerhubProxy string

	// packageNamespace is the namespace where the external-dns package is
	// installed (i.e this is the namespace tanzu package install is called
	// with)
	packageNamespace string

	// addonNamespace is the namespace where the external-dns addon is
	// installed (i.e this is the namespace passed into the external-dns addon
	// values.yaml that is provided to the package). This namespace is created
	// by the package installation.
	addonNamespace string

	// fixtureNamespace is the namespace where all test fixtures are created
	// for the purpose of testing the addon (e.g bind, kuard, dnsutils)
	fixtureNamespace string
)

var _ = BeforeSuite(func() {
	dockerhubProxy = os.Getenv("DOCKERHUB_PROXY")
	if dockerhubProxy == "" {
		dockerhubProxy = "docker.io"
	}

	packageNamespace = "e2e-external-dns-package"
	addonNamespace = "e2e-external-dns-addon"
	fixtureNamespace = "e2e-external-dns-fixtures"

	_, err := utils.Kubectl(nil, "create", "namespace", packageNamespace)
	Expect(err).NotTo(HaveOccurred())

	_, err = utils.Kubectl(nil, "create", "namespace", fixtureNamespace)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	_, err := utils.Kubectl(nil, "delete", "namespace", fixtureNamespace)
	Expect(err).NotTo(HaveOccurred())

	_, err = utils.Kubectl(nil, "delete", "namespace", packageNamespace)
	Expect(err).NotTo(HaveOccurred())
})
