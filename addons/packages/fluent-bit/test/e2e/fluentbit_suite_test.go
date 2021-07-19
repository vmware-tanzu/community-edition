// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/vmware-tanzu/tce/addons/packages/fluent-bit/test/e2e"
	"github.com/vmware-tanzu/tce/test/pkg/cmdhelper"

	. "github.com/onsi/ginkgo"

	. "github.com/onsi/gomega"
)

func TestFluentBitE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "fluent-bit addon package e2e test suite")
}

var (
	cmdHelperUp             *cmdhelper.CmdHelper
	cmdHelperDown           *cmdhelper.CmdHelper
	DeploymentTimeout       = 120 * time.Second
	DeploymentCheckInterval = 5 * time.Second
	ApiCallTimeout          = 20 * time.Second
)

var _ = BeforeSuite(func() {
	var err error
	cmdHelperUp, err = cmdhelper.New(e2e.GetAllUpCmds(), os.Stdin)
	Expect(err).NotTo(HaveOccurred())

	cmdHelperDown, err = cmdhelper.New(e2e.GetTearDownCmds(), os.Stdin)
	Expect(err).NotTo(HaveOccurred())

	// delete fluent-bit if at all already installed
	// nothing to worry about the err or result here
	cmdHelperDown.CliRunner("tanzu", nil, cmdHelperDown.GetFormatted("tanzu-package-delete", "$", []string{"fluent-bit.tce.vmware.com"})...)

	result, err := cmdHelperUp.CliRunner("tanzu", nil, cmdHelperUp.GetFormatted("tanzu-package-install", "$", []string{"fluent-bit.tce.vmware.com"})...)
	Expect(err).NotTo(HaveOccurred())
	Expect(result).Should(ContainSubstring("Installed package in default/fluent-bit.tce.vmware.com"))

	// to ensure fluent-bit deamonset
	EventuallyWithOffset(1, func() (string, error) {
		return cmdHelperUp.CliRunner("kubectl", nil, cmdHelperUp.GetFormatted("k8s-daemonset-state", "$", []string{"fluent-bit", `jsonpath={..status.desiredNumberScheduled}`})...)
	}, DeploymentTimeout, DeploymentCheckInterval).Should(Equal("2"), fmt.Sprintln("daemonset was not ready"))
	EventuallyWithOffset(1, func() (string, error) {
		return cmdHelperUp.CliRunner("kubectl", nil, cmdHelperUp.GetFormatted("k8s-daemonset-state", "$", []string{"fluent-bit", `jsonpath={..status.numberAvailable}`})...)
	}, DeploymentTimeout, DeploymentCheckInterval).Should(Equal("2"), fmt.Sprintln("daemonset was not ready"))
})
var _ = AfterSuite(func() {
	// delete the fluent-bit package
	result, err := cmdHelperDown.CliRunner("tanzu", nil, cmdHelperDown.GetFormatted("tanzu-package-delete", "$", []string{"fluent-bit.tce.vmware.com"})...)
	Expect(err).NotTo(HaveOccurred())
	Expect(result).Should(ContainSubstring("Deleted default/fluent-bit.tce.vmware.com"))
})
