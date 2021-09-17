// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"flag"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/utils"
)

var (
	// full multus-cni package to test.
	packageName string
	// package version of multus-cni.
	packageVersion string
	// package version with suffix.
	packageFullVersion string
	// package installed namespace.
	packageInstalledNamespace string
	// package installed name.
	packageInstalledName string
)

func init() {
	flag.StringVar(&packageVersion, "version", "3.7.1", "the version of the package to test")
}

func TestMultusCNIE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Multus CNI addon E2E Test Suite")
}

var _ = BeforeSuite(func() {
	// Needs to have package related values predefined.
	packageName = utils.TanzuPackageName("multus-cni")
	packageFullVersion = utils.TanzuPackageAvailableVersion(packageName)
	packageInstalledName = "multus-cni-pkg"
	packageInstalledNamespace = "default"
	if strings.Compare(packageInstalledNamespace, "default") != 0 {
		_, err := utils.Kubectl(nil, "create", "ns", packageInstalledNamespace)
		Expect(err).NotTo(HaveOccurred())
	}
})

var _ = AfterSuite(func() {
	if strings.Compare(packageInstalledNamespace, "default") != 0 {
		_, err := utils.Kubectl(nil, "delete", "ns", packageInstalledNamespace)
		Expect(err).NotTo(HaveOccurred())
	}
})
