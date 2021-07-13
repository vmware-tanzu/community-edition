package e2e_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gatekeeper Addon E2E Test", func() {
	//<-signal
	BeforeEach(func() {
		Expect(cmdHelperUp).NotTo(BeNil())
		result, err := cmdHelperUp.Run("tanzu", nil, "tanzu-package-install-gatekeeper")
		Expect(err).NotTo(HaveOccurred())
		Expect(result).Should(ContainSubstring("Installed package in default/gatekeeper.tce.vmware.com"))

	})

	Specify("apply constraint-template file", func() {
		filename := filepath.Join("fixtures", "constraint-template.yaml")
		cmdHelperUp.Format("kubectl-apply-constraint-template", "$", []string{filename})
		result, err := cmdHelperUp.Run("kubectl", nil, "kubectl-apply-constraint-template")
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
		cmdHelperUp.Format("kubectl-apply-constraint", "$", []string{filename})
		result, err := cmdHelperUp.Run("kubectl", nil, "kubectl-apply-constraint")
		Expect(err).NotTo(HaveOccurred())
		Expect(result).Should(ContainSubstring("k8srequiredlabels.constraints.gatekeeper.sh/all-must-have-owner "))
	})

	Specify("create a namespace", func() {
		cmdHelperUp.Format("kubectl-create-ns", "$", []string{"test"})
		result, err := cmdHelperUp.Run("kubectl", nil, "kubectl-create-ns")
		Expect(err).NotTo(HaveOccurred())
		Expect(result).Should(ContainSubstring("namespace/test created"))
	})

	Specify("apply namespace with owner file", func() {
		filename := filepath.Join("fixtures", "test-namespace.yaml")
		cmdHelperUp.Format("kubectl-apply-namespace", "$", []string{filename})
		result, err := cmdHelperUp.Run("kubectl", nil, "kubectl-apply-namespace")
		Expect(err).NotTo(HaveOccurred())
		Expect(result).Should(ContainSubstring("namespace/test"))
	})

	Specify("create a namespace", func() {
		cmdHelperUp.Format("kubectl-create-ns", "$", []string{"test-again"})
		_, err := cmdHelperUp.Run("kubectl", nil, "kubectl-create-ns")
		Expect(err).Should(ContainSubstring("[denied by all-must-have-owner] All namespaces must have an owner label"))

	})

	AfterEach(func() {
		//todo
	})
})
