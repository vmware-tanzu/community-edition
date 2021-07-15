// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gatekeeper Addon E2E Test", func() {

	BeforeEach(func() {
		Expect(cmdHelperUp).NotTo(BeNil())
		result, err := cmdHelperUp.Run("tanzu", nil, "tanzu-package-install-gatekeeper")
		Expect(err).NotTo(HaveOccurred())
		Expect(result).Should(ContainSubstring("Installed package in default/gatekeeper.tce.vmware.com"))
	})

	Specify("apply constraint-template file", func() {
		filename := filepath.Join("fixtures", "constraint-template.yaml")
		result, err := cmdHelperUp.CliRunner("kubectl", nil, cmdHelperUp.GetFormatted("kubectl-apply", "$", []string{filename})...)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).Should(ContainSubstring("constrainttemplate.templates.gatekeeper.sh/k8srequiredlabels "))
	})

	Specify("check k8srequiredlabels in crd", func() {
		result, err := cmdHelperUp.Run("kubectl", nil, "kubectl-get-crds-constraint-template")
		Expect(err).NotTo(HaveOccurred())
		Expect(result).Should(ContainSubstring("k8srequiredlabels"))
	})

	Specify("apply constraint file", func() {
		filename := filepath.Join("fixtures", "constraint.yaml")
		result, err := cmdHelperUp.CliRunner("kubectl", nil, cmdHelperUp.GetFormatted("kubectl-apply", "$", []string{filename})...)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).Should(ContainSubstring("k8srequiredlabels.constraints.gatekeeper.sh/all-must-have-owner "))
	})

	Specify("apply namespace with owner file", func() {
		filename := filepath.Join("fixtures", "test-namespace.yaml")
		result, err := cmdHelperUp.CliRunner("kubectl", nil, cmdHelperUp.GetFormatted("kubectl-apply", "$", []string{filename})...)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).Should(ContainSubstring("namespace/test"))
	})

	Specify("create a namespace without owner", func() {
		_, err := cmdHelperUp.CliRunner("kubectl", nil, cmdHelperUp.GetFormatted("kubectl-create-ns", "$", []string{"test-again"})...)
		Expect(err.Error()).Should(ContainSubstring("[denied by all-must-have-owner] All namespaces must have an owner label"))
	})

	AfterEach(func() {
		//todo
	})
})
