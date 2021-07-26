// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/vmware-tanzu/tce/addons/packages/prometheus/test/e2e"
	"github.com/vmware-tanzu/tce/addons/packages/test/pkg/cmdhelper"

	. "github.com/onsi/ginkgo"

	. "github.com/onsi/gomega"
)

func TestGrafanaE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "prometheus addon package e2e test suite")
}

var (
	cmdHelperUp             *cmdhelper.CmdHelper
	cmdHelperDown           *cmdhelper.CmdHelper
	DeploymentTimeout       = 120 * time.Second
	DeploymentCheckInterval = 5 * time.Second
	ApiCallTimeout          = 5 * time.Second
)

var _ = BeforeSuite(func() {
	var err error
	cmdHelperUp, err = cmdhelper.New(e2e.GetAllUpCmds(), os.Stdin)
	Expect(err).NotTo(HaveOccurred())

	cmdHelperDown, err = cmdhelper.New(e2e.GetTearDownCmds(), os.Stdin)
	Expect(err).NotTo(HaveOccurred())

	// delete prometheus if at all already installed
	// nothing to worry about the err or result here
	cmdHelperDown.CliRunner("tanzu", nil, cmdHelperDown.GetFormatted("tanzu-package-delete", "$", []string{"prometheus.tce.vmware.com"})...)

	result, err := cmdHelperUp.CliRunner("tanzu", nil, cmdHelperUp.GetFormatted("tanzu-package-install", "$", []string{"prometheus.tce.vmware.com"})...)
	Expect(err).NotTo(HaveOccurred())
	Expect(result).Should(ContainSubstring("Installed package in default/grafana.tce.vmware.com"))

	// to ensure grafana audit pod is ready
	EventuallyWithOffset(1, func() (string, error) {
		return cmdHelperUp.CliRunner("kubectl", nil, cmdHelperUp.GetFormatted("k8s-deployment-ready-status", "$", []string{"app.kubernetes.io/name=grafana", "grafana-addon"})...)
	}, DeploymentTimeout, DeploymentCheckInterval).Should(Equal("True"), fmt.Sprintln("deployment was not ready"))
})
var _ = AfterSuite(func() {
	// delete the grafana package
	result, err := cmdHelperDown.CliRunner("tanzu", nil, cmdHelperDown.GetFormatted("tanzu-package-delete", "$", []string{"grafana.tce.vmware.com"})...)
	Expect(err).NotTo(HaveOccurred())
	Expect(result).Should(ContainSubstring("Deleted default/grafana.tce.vmware.com"))
})
