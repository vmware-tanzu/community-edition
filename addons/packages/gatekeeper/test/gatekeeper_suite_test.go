// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"fmt"
	"os"
	"testing"
	"time"

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

var _ = BeforeSuite(func() {
	var err error
	cmdHelperUp, err = cmdhelper.New(GetAllUpCmds(), os.Stdin)
	Expect(err).NotTo(HaveOccurred())

	cmdHelperDown, err = cmdhelper.New(GetTearDownCmds(), os.Stdin)
	Expect(err).NotTo(HaveOccurred())

	// delete gatekeeper if at all already installed
	// nothing to worry about the err or result here
	cmdHelperDown.CliRunner("tanzu", nil, cmdHelperDown.GetFormatted("tanzu-package-delete", "$", []string{"gatekeeper.tce.vmware.com"})...)

	result, err := cmdHelperUp.CliRunner("tanzu", nil, cmdHelperUp.GetFormatted("tanzu-package-install", "$", []string{"gatekeeper.tce.vmware.com"})...)
	Expect(err).NotTo(HaveOccurred())
	Expect(result).Should(ContainSubstring("Installed package in default/gatekeeper.tce.vmware.com"))

	// to ensure gatekeeper audit pod is ready
	EventuallyWithOffset(1, func() (string, error) {
		return cmdHelperUp.CliRunner("kubectl", nil, cmdHelperUp.GetFormatted("k8s-check-pod-ready-status", "$", []string{"gatekeeper.sh/operation=audit", "gatekeeper-system"})...)
	}, DeploymentTimeout, DeploymentCheckInterval).Should(Equal("True"), fmt.Sprintln("pod was not ready"))

	// to ensure gatekeeper webhook pod is ready
	EventuallyWithOffset(1, func() (string, error) {
		return cmdHelperUp.CliRunner("kubectl", nil, cmdHelperUp.GetFormatted("k8s-check-pod-ready-status", "$", []string{"gatekeeper.sh/operation=webhook", "gatekeeper-system"})...)
	}, DeploymentTimeout, DeploymentCheckInterval).Should(Equal("True"), fmt.Sprintln("pod was not ready"))
})
var _ = AfterSuite(func() {
	// delete the gatekeeper package
	result, err := cmdHelperDown.CliRunner("tanzu", nil, cmdHelperDown.GetFormatted("tanzu-package-delete", "$", []string{"gatekeeper.tce.vmware.com"})...)
	Expect(err).NotTo(HaveOccurred())
	Expect(result).Should(ContainSubstring("Deleted default/gatekeeper.tce.vmware.com"))
})
