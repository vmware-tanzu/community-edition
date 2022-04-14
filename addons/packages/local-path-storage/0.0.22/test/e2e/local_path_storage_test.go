// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"fmt"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("local-path-storage Package E2E Test", func() {

	JustAfterEach(func() {
		if CurrentGinkgoTestDescription().Failed {
			fmt.Fprintf(GinkgoWriter, "\nCollecting diagnostic information just after test failure\n")
			fmt.Fprintf(GinkgoWriter, "\nResources summary:\n")
			utils.Kubectl(nil, "-n", packageInstallNamespace, "get", "all,packageinstalls,apps") // nolint:errcheck

			fmt.Fprintf(GinkgoWriter, "\npackage install status:\n")
			utils.Kubectl(nil, "-n", packageInstallNamespace, "get", "app", packageInstallName, "-o", "jsonpath={.status}") // nolint:errcheck
		}
	})

	It("works", func() {
		_, err := utils.Kubectl(nil, "get", "StorageClass", "local-path")
		Expect(err).NotTo(HaveOccurred())
	})
})
