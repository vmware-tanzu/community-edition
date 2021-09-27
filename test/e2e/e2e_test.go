// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"log"

	"github.com/vmware-tanzu/community-edition/test/e2e"
	"github.com/vmware-tanzu/community-edition/test/e2e/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E tests", func() {
	Describe("Smoke testing is in progress ....", func() {
		It("Validate tanzu cli", func() {
			_, err := utils.CheckPackageRepositoryListAllNamespaces()
			Expect(err).NotTo(HaveOccurred())

			_, err = utils.CheckPackageAvailableListAllNamespaces()
			Expect(err).NotTo(HaveOccurred())

			_, err = utils.PackageInstalledListAllNamespaces()
			Expect(err).NotTo(HaveOccurred())
			if err != nil {
				// collect diagnostic details for workload cluster
				log.Println("Collecting workload cluster diagnostic details in /tmp folder...")
				_, err = utils.GetWorkloadClusterDiagnostics(e2e.ConfigVal.GuestClusterName)
				if err != nil {
					log.Println("error while collecting logs from cluster", err)
				}

				err = e2e.DeleteCluster()
			}
			Expect(err).NotTo(HaveOccurred())
		})

		It("Run addon package tests", func() {
			err := e2e.RunAddonsTests()
			if err != nil {
				log.Println("error while running package test", err)
				// collect diagnostic details for workload cluster
				log.Println("Collecting workload cluster diagnostic details in /tmp folder...")
				_, err = utils.GetWorkloadClusterDiagnostics(e2e.ConfigVal.GuestClusterName)
				if err != nil {
					log.Println("error while collecting logs from cluster", err)
				}

				err = e2e.DeleteCluster()
			}
			Expect(err).NotTo(HaveOccurred())

		})
	})
})
