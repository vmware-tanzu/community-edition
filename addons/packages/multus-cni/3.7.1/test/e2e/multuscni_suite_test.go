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
	// To avoid rate-limiting since multus-cni image is currently placed under docker.
	pkgContainerRegistry string
	pkgContainerName     string
	pkgContainerTag      string

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
	var err error
	pkgContainerRegistry = "ghcr.io/k8snetworkplumbingwg"
	pkgContainerName = "multus-cni"
	pkgContainerTag = "v" + packageVersion

	packageName = utils.TanzuPackageName("multus-cni")
	Expect(err).NotTo(HaveOccurred())

	packageFullVersion = utils.TanzuPackageAvailableVersion(packageName)
	Expect(err).NotTo(HaveOccurred())
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
