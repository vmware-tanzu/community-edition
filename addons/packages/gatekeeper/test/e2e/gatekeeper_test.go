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
		filename := filepath.Join("fixtures", "constraint-template.yaml")
		result, err := cmdHelperUp.CliRunner("kubectl", nil, cmdHelperUp.GetFormatted("k8s-apply", "$", []string{filename})...)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).Should(ContainSubstring("constrainttemplate.templates.gatekeeper.sh/k8srequiredlabels "))

		filename = filepath.Join("fixtures", "constraint.yaml")
		result, err = cmdHelperUp.CliRunner("kubectl", nil, cmdHelperUp.GetFormatted("k8s-apply", "$", []string{filename})...)

		Expect(err).NotTo(HaveOccurred())
		Expect(result).Should(ContainSubstring("k8srequiredlabels.constraints.gatekeeper.sh/all-must-have-owner "))
	})

	Specify("check k8srequiredlabels in crd", func() {
		result, err := cmdHelperUp.Run("kubectl", nil, "k8s-get-crds")
		Expect(err).NotTo(HaveOccurred())
		Expect(result).Should(ContainSubstring("k8srequiredlabels"))
	})

	Specify("apply namespace with owner file and delete", func() {
		filename := filepath.Join("fixtures", "test-namespace.yaml")
		result, err := cmdHelperUp.CliRunner("kubectl", nil, cmdHelperUp.GetFormatted("k8s-apply", "$", []string{filename})...)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).Should(ContainSubstring("namespace/test created"))

		result, err = cmdHelperDown.CliRunner("kubectl", nil, cmdHelperDown.GetFormatted("k8s-delete-ns", "$", []string{"test"})...)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).Should(ContainSubstring(`namespace "test" deleted`))
	})

	Specify("create a namespace without owner", func() {
		_, err := cmdHelperUp.CliRunner("kubectl", nil, cmdHelperUp.GetFormatted("k8s-create-ns", "$", []string{"test-again"})...)
		Expect(err.Error()).Should(ContainSubstring("[denied by all-must-have-owner] All namespaces must have an owner label"))
	})

	AfterEach(func() {
		filename := filepath.Join("fixtures", "constraint.yaml")
		result, err := cmdHelperDown.CliRunner("kubectl", nil, cmdHelperDown.GetFormatted("k8s-delete", "$", []string{filename})...)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).Should(ContainSubstring(`k8srequiredlabels.constraints.gatekeeper.sh "all-must-have-owner" deleted`))
	})
})
