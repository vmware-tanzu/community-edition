// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"fmt"
	"path/filepath"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cert-manager Package E2E Test", func() {

	JustAfterEach(func() {
		utils.Kubectl(nil, "delete", "namespace", "cert-manager-test") // nolint:errcheck

		if CurrentGinkgoTestDescription().Failed {
			fmt.Fprintf(GinkgoWriter, "\nCollecting diagnostic information just after test failure\n")
			fmt.Fprintf(GinkgoWriter, "\nResources summary:\n")
			utils.Kubectl(nil, "-n", packageInstallNamespace, "get", "all,packageinstalls,apps") // nolint:errcheck

			fmt.Fprintf(GinkgoWriter, "\npackage install status:\n")
			utils.Kubectl(nil, "-n", packageInstallNamespace, "get", "app", packageInstallName, "-o", "jsonpath={.status}") // nolint:errcheck
		}
	})

	It("works", func() {
		_, err := utils.Kubectl(nil, "apply", "-f", filepath.Join("fixtures", "cert-manager-resources.yaml"))
		Expect(err).NotTo(HaveOccurred())
	})
})
