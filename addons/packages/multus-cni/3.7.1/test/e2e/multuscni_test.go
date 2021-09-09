// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/utils"
)

var (
	// The net1 is the default iface name when not appointed.
	checkMacVlanIFaceExistence = "ip a show dev net1"
	checkIFaceStatus           = "ip -j a show net1 | jq -r '.[]|.operstate'"
	getMacVlanIP               = "ip -j a show | jq -r '.[]|select(.ifname ==\"net1\")|.addr_info[]|select(.family==\"inet\").local'"
)

var _ = Describe("Multus CNI Addon E2E Test", func() {
	BeforeEach(func() {
		By("Install Multus CNI Package and tests' resources")
		multusCNIAddonValuesFilename, err := utils.ReadFileAndReplaceContentsTempFile(filepath.Join("multihomed-testfiles", "multus-cni-values.yaml"),
			map[string]string{
				"docker.io/nfvpe": pkgContainerRegistry,
				"name: multus":    "name: " + pkgContainerName,
				"stable":          pkgContainerTag,
			})
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(multusCNIAddonValuesFilename)
		_, err = utils.Tanzu(nil, "package", "install", packageInstalledName, "--package-name", packageName, "--version", packageFullVersion, "-n", packageInstalledNamespace, "-f", multusCNIAddonValuesFilename)
		Expect(err).NotTo(HaveOccurred())

		utils.ValidateDaemonsetReady("kube-system", "kube-multus-ds-amd64")

		macvlanPodsFilename, err := utils.ReadFileAndReplaceContentsTempFile(filepath.Join("multihomed-testfiles", "simple-macvlan1.yaml"),
			map[string]string{
				"docker.io/library": "quay.io/centos",
			})
		defer os.Remove(macvlanPodsFilename)
		Expect(err).NotTo(HaveOccurred())
		_, err = utils.Kubectl(nil, "apply", "-f", macvlanPodsFilename)
		Expect(err).NotTo(HaveOccurred())

		By("Check resources status")
		utils.ValidatePodReady("default", "macvlan1-worker0")
		utils.ValidatePackageInstallReady(packageInstalledNamespace, packageInstalledName)
	})

	JustBeforeEach(func() {
		By("Should install jq command")
		_, err := utils.Kubectl(nil, "exec", "macvlan1-worker0", "--", "yum", "install", "-y", "jq")
		Expect(err).NotTo(HaveOccurred())
	})

	JustAfterEach(func() {
		if CurrentGinkgoTestDescription().Failed {
			fmt.Fprintf(GinkgoWriter, "\nAfter test fails, theses resources are related.\n")
			_, err := utils.Kubectl(nil, "get", "app", packageInstalledName, "-n", packageInstalledNamespace, "-o", "jsonpath={.status}")
			Expect(err).NotTo(HaveOccurred())
			_, err = utils.Kubectl(nil, "get", "daemonset", "kube-multus-ds-amd64", "-n", "kube-system", "-o", "jsonpath={.status}")
			Expect(err).NotTo(HaveOccurred())
		}
	})

	AfterEach(func() {
		By("Should delete Multus CNI test temporary resources")
		_, err := utils.Kubectl(nil, "delete", "-f", filepath.Join("multihomed-testfiles", "simple-macvlan1.yaml"))
		Expect(err).NotTo(HaveOccurred())
		_, err = utils.Tanzu(nil, "package", "installed", "delete", packageInstalledName, "-n", packageInstalledNamespace, "-y")
		Expect(err).NotTo(HaveOccurred())

		utils.ValidatePackageInstallNotFound(packageInstalledNamespace, packageName)
		utils.ValidateDaemonsetNotFound("kube-system", "kube-multus-ds-amd64")
		utils.ValidatePodNotFound("default", "macvlan1-worker0")
	})

	It("Check macvlan interfaces status", func() {
		By("check macvlan1-worker0 interface: net1")
		_, err := utils.Kubectl(nil, "exec", "macvlan1-worker0", "--",
			"bash", "-c", checkMacVlanIFaceExistence)
		Expect(err).NotTo(HaveOccurred())

		By("check macvlan1-worker0 interface status: net1")
		ifstatus, err := utils.Kubectl(nil, "exec", "macvlan1-worker0", "--",
			"bash", "-c", checkIFaceStatus)
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.Replace(ifstatus, "\n", "", -1)).To(Equal("UP"))

		By("get macvlan1-worker0 interface address: net1")
		ip, err := utils.Kubectl(nil, "exec", "macvlan1-worker0", "--",
			"bash", "-c", getMacVlanIP)
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.Replace(ip, "\n", "", -1)).To(Equal("10.1.1.12"))
	})
})
