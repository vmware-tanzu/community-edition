// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gatekeeper Package E2E Test", func() {

	JustAfterEach(func() {
		if CurrentGinkgoTestDescription().Failed {
			fmt.Fprintf(GinkgoWriter, "\nCollecting diagnostic information just after test failure\n")
			fmt.Fprintf(GinkgoWriter, "\nResources summary:\n")
			utils.Kubectl(nil, "-n", packageInstallNamespace, "get", "all,packageinstalls,apps") // nolint:errcheck

			fmt.Fprintf(GinkgoWriter, "\npackage install status:\n")
			utils.Kubectl(nil, "-n", packageInstallNamespace, "get", "app", packageInstallName, "-o", "jsonpath={.status}") // nolint:errcheck
		}
	})

	It("can create mutations", func() {
		_, err := utils.Kubectl(nil, "apply", "-f", filepath.Join("fixtures", "mutations.yaml"))
		Expect(err).NotTo(HaveOccurred())
	})

	It("can create a constraint", func() {
		_, err := utils.Kubectl(nil, "apply", "-f", filepath.Join("fixtures", "constraint.yaml"))
		Expect(err).NotTo(HaveOccurred())

		_, err = utils.Kubectl(nil, "apply", "-f", filepath.Join("fixtures", "k8srequiredlabels.yaml"))
		Expect(err).NotTo(HaveOccurred())

		time.Sleep(5 * time.Second) // Need to sleep for a few seconds to allow the constraint to settle in I guess?

		_, err = utils.Kubectl(nil, "apply", "-f", filepath.Join("fixtures", "label_on_namespace.yaml"))
		Expect(err).NotTo(HaveOccurred())

		_, err = utils.Kubectl(nil, "create", "namespace", "ishouldfail")
		Expect(err).NotTo(BeNil())

		_, err = utils.Kubectl(nil, "delete", "namespace", "test")
		Expect(err).NotTo(HaveOccurred())
	})
})
