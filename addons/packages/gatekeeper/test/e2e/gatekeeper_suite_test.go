package e2e_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/vmware-tanzu/tce/addons/packages/gatekeeper/test/e2e"
	"github.com/vmware-tanzu/tce/test/pkg/cmdhelper"

	. "github.com/onsi/ginkgo"

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
	cmdHelperUp, err = cmdhelper.New(e2e.GetAllUpCmds(), os.Stdin)
	Expect(err).NotTo(HaveOccurred())

	// delete gatekeeper if at all already installed
	cmdHelperUp.Run("tanzu", nil, "tanzu-package-delete-gatekeeper")
	result, err := cmdHelperUp.Run("tanzu", nil, "tanzu-package-install-gatekeeper")
	Expect(err).NotTo(HaveOccurred())
	Expect(result).Should(ContainSubstring("Installed package in default/gatekeeper.tce.vmware.com"))

	// to ensure gatekeeper audit pod is ready
	EventuallyWithOffset(1, func() (string, error) {
		return cmdHelperUp.CliRunner("kubectl", nil, cmdHelperUp.GetFormatted("kubeclt-check-pod-ready-status", "$", []string{"gatekeeper.sh/operation=audit", "gatekeeper-system"})...)
	}, DeploymentTimeout, DeploymentCheckInterval).Should(Equal("True"), fmt.Sprintln("pod was not ready"))

	// to ensure gatekeeper webhook pod is ready
	EventuallyWithOffset(1, func() (string, error) {
		return cmdHelperUp.CliRunner("kubectl", nil, cmdHelperUp.GetFormatted("kubeclt-check-pod-ready-status", "$", []string{"gatekeeper.sh/operation=webhook", "gatekeeper-system"})...)
	}, DeploymentTimeout, DeploymentCheckInterval).Should(Equal("True"), fmt.Sprintln("pod was not ready"))
})
var _ = AfterSuite(func() {
	var err error
	cmdHelperDown, err = cmdhelper.New(e2e.GetTearDownCmds(), os.Stdin)
	Expect(err).NotTo(HaveOccurred())

	// delete the gatekeeper package
	cmdHelperDown.Run("tanzu", nil, "tanzu-package-delete-gatekeeper")
	// delete the namespace
	cmdHelperDown.CliRunner("kubectl", nil, cmdHelperDown.GetFormatted("kubectl-delete-ns", "$", []string{"test-again"})...)
	cmdHelperDown.CliRunner("kubectl", nil, cmdHelperDown.GetFormatted("kubectl-delete-ns", "$", []string{"test"})...)
})
