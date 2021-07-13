package e2e_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/vmware-tanzu/tce/test/pkg/cmdhelper"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
)

func TestGateKeeperE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gatekeeper addon package e2e test suite")
}

var (
	cmdHelperUp             *cmdhelper.CmdHelper
	cmdHelperDown           *cmdhelper.CmdHelper
	DeploymentTimeout       = 120 * time.Second
	DeploymentCheckInterval = 5 * time.Second
)

func Init() {
	// todo
}

var _ = BeforeSuite(func() {
	var err error
	cmdHelperUp, err = cmdhelper.New(map[string][]string{
		"tanzu-package-install-gatekeeper":     []string{"package", "install", "gatekeeper.tce.vmware.com"},
		"tanzu-package-delete-gatekeeper":      []string{"package", "delete", "gatekeeper.tce.vmware.com"},
		"kubectl-get-pods-by-namespace":        []string{"get", "pods", "gatekeeper-system"},
		"kubectl-apply-constraint-template":    []string{"apply", "-f", "$"},
		"kubectl-get-crds-constraint-template": []string{"get", "crds"},
		"kubectl-apply-constraint":             []string{"apply", "-f", "$"},
		"kubectl-create-ns":                    []string{"create", "ns", "$"},
		"kubectl-apply-namespace":              []string{"apply", "-f", "$"},
		"kubeclt-check-pod-ready":              []string{"wait", "--for=condition=ready", "pod", "-l", "gatekeeper.sh/operation=audit", "-n", "gatekeeper-system"},
		"kubeclt-check-pod-ready-status":       []string{"get", "pods", "-l", "gatekeeper.sh/operation=audit", "-n", "gatekeeper-system", "-o", `jsonpath={..status.conditions[?(@.type=="Ready")].status}`},
	}, os.Stdin)
	Expect(err).NotTo(HaveOccurred())

	// delete gatekeeper if at all already installed
	cmdHelperUp.Run("tanzu", nil, "tanzu-package-delete-gatekeeper")
	result, err := cmdHelperUp.Run("tanzu", nil, "tanzu-package-install-gatekeeper")
	Expect(err).NotTo(HaveOccurred())
	Expect(result).Should(ContainSubstring("Installed package in default/gatekeeper.tce.vmware.com"))

	gomega.EventuallyWithOffset(1, func() (string, error) {
		return cmdHelperUp.Run("kubectl", nil, "kubeclt-check-pod-ready-status")
	}, DeploymentTimeout, DeploymentCheckInterval).Should(gomega.Equal("True"), fmt.Sprintln("pod was not ready"))

})
var _ = AfterSuite(func() {
	var err error
	cmdHelperDown, err = cmdhelper.New(map[string][]string{
		"tanzu-package-delete-gatekeeper": []string{"package", "delete", "gatekeeper.tce.vmware.com"},
		"kubectl-delete-namespace":        []string{"delete", "ns", "test"},
	}, os.Stdin)
	Expect(err).NotTo(HaveOccurred())

	// delete the gatekeeper package
	cmdHelperDown.Run("tanzu", nil, "tanzu-package-delete-gatekeeper")

	// delete the namespace
	cmdHelperDown.Run("kubectl", nil, "kubectl-delete-namespace")
})
