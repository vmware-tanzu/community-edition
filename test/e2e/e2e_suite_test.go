// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"log"
	"testing"

	"github.com/vmware-tanzu/community-edition/test/e2e"
	"github.com/vmware-tanzu/community-edition/test/e2e/testdata"
	"github.com/vmware-tanzu/community-edition/test/e2e/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2e Suite")
}

var _ = BeforeSuite(func() {
	Describe("Check if cluster is up and running else install if required", func() {
		e2e.Initialize()

		// cluster context must be set to guest cluster context
		clusterContext := utils.GetClusterContext(e2e.ConfigVal.GuestClusterName)
		if !e2e.ConfigVal.ClusterInstallRequired && !e2e.ConfigVal.ClusterCleanupRequired {
			By("Check Cluster health")
			err := e2e.CheckClusterHealth(clusterContext)
			if err != nil {
				log.Println("error while checking cluster health, deleting cluster...")
				err = e2e.DeleteCluster()
			}
			Expect(err).NotTo(HaveOccurred())
		}
	})
})

var _ = AfterSuite(func() {
	Describe("Delete dependency package installed and delete cluster...", func() {
		if e2e.MetallbInstalled {
			err := testdata.UninstallMetallb()
			Expect(err).NotTo(HaveOccurred())
		}

		if e2e.VeleroInstalled {
			err := testdata.UnsinstallVelero()
			Expect(err).NotTo(HaveOccurred())
		}

		// delete the cluster
		if e2e.ConfigVal.ClusterCleanupRequired {
			log.Println("Deleting the cluster created")
			err := e2e.DeleteCluster()
			Expect(err).NotTo(HaveOccurred())
		}
	})
})
