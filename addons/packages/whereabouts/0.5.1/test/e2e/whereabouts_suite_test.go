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

type addonPackage struct {
	name               string
	version            string
	fullVersion        string
	installedNamespace string
	installedName      string
}

var (
	multusCNIPackage       addonPackage
	whereaboutsPackage     addonPackage
	multusCNIUseConfFile   bool
	whereaboutsUseConfFile bool
	packageVersion         string
)

func init() {
	flag.BoolVar(&multusCNIUseConfFile, "multuscni-use-conf", true, "use configuration file or not")
	flag.BoolVar(&whereaboutsUseConfFile, "whereabouts-use-conf", true, "use configuration file or not")
	flag.StringVar(&packageVersion, "version", "0.5.1", "the version of the package to test")
}

func TestWhereaboutsE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Whereabouts addon E2E Test Suite")
}

var _ = BeforeSuite(func() {
	// As Whereabouts CNI plugin usually works with a meta CNI plugin like multus-cni, we need to install multus-cni package as well.
	// Needs to have package related values predefined.
	multusCNIpackageName := utils.TanzuPackageName("multus-cni")
	multusCNIPackage = addonPackage{
		name:               multusCNIpackageName,
		version:            "3.7.1",
		fullVersion:        utils.TanzuPackageAvailableVersionWithVersionSubString(multusCNIpackageName, "3.7.1"),
		installedNamespace: "default",
		installedName:      "multus-cni-pkg",
	}

	whereaboutsPackageName := utils.TanzuPackageName("whereabouts")
	whereaboutsPackage = addonPackage{
		name:               whereaboutsPackageName,
		version:            "0.5.1",
		fullVersion:        utils.TanzuPackageAvailableVersionWithVersionSubString(whereaboutsPackageName, "0.5.1"),
		installedNamespace: "default",
		installedName:      "whereabouts-pkg",
	}

	if strings.Compare(multusCNIPackage.installedNamespace, "default") != 0 {
		_, err := utils.Kubectl(nil, "create", "ns", multusCNIPackage.installedNamespace)
		Expect(err).NotTo(HaveOccurred())
	}

	if strings.Compare(whereaboutsPackage.installedNamespace, "default") != 0 {
		_, err := utils.Kubectl(nil, "create", "ns", whereaboutsPackage.installedNamespace)
		Expect(err).NotTo(HaveOccurred())
	}

})

var _ = AfterSuite(func() {
	if strings.Compare(multusCNIPackage.installedNamespace, "default") != 0 {
		_, err := utils.Kubectl(nil, "delete", "ns", multusCNIPackage.installedNamespace)
		Expect(err).NotTo(HaveOccurred())
	}

	if strings.Compare(whereaboutsPackage.installedNamespace, "default") != 0 {
		_, err := utils.Kubectl(nil, "delete", "ns", whereaboutsPackage.installedNamespace)
		Expect(err).NotTo(HaveOccurred())
	}
})
