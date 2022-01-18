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

		curDir, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		installOptions := []string{"package", "install", packageInstalledName, "--package-name", packageName, "--version", packageFullVersion, "-n", packageInstalledNamespace}
		if useConfFile {
			installOptions = append(installOptions, "-f", filepath.Join(curDir, "multihomed-testfiles/multus-cni-values.yaml"))
		}
		_, err = utils.Tanzu(nil, installOptions...)
		Expect(err).NotTo(HaveOccurred())
		_, err = utils.Kubectl(nil, "apply", "-f", filepath.Join(curDir, "multihomed-testfiles/cni-install.yaml"))
		Expect(err).NotTo(HaveOccurred())

		utils.ValidatePackageInstallReady(packageInstalledNamespace, packageInstalledName)
		utils.ValidateDaemonsetReady("kube-system", "kube-multus-ds-amd64")
		_, err = utils.Kubectl(nil, "apply", "-f", filepath.Join(curDir, "multihomed-testfiles/simple-macvlan.yaml"))
		Expect(err).NotTo(HaveOccurred())

		utils.ValidateDaemonsetReady("kube-system", "install-cni-plugins")
		utils.ValidatePodReady("default", "macvlan-worker")

		// Just do a simple check for which image is using
		_, err = utils.Kubectl(nil, "get", "daemonset", "kube-multus-ds-amd64", "-n", "kube-system", "-o", "jsonpath='{range .spec.template.spec.containers[*]}{.image}{\"\\n\"}{end}'")
		Expect(err).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
		By("Should install jq command")
		_, err := utils.Kubectl(nil, "exec", "macvlan-worker", "--", "yum", "install", "-y", "jq")
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
		curDir, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		_, err = utils.Kubectl(nil, "delete", "-f", filepath.Join(curDir, "multihomed-testfiles/simple-macvlan.yaml"))
		Expect(err).NotTo(HaveOccurred())
		_, err = utils.Kubectl(nil, "delete", "-f", filepath.Join(curDir, "multihomed-testfiles/cni-install.yaml"))
		Expect(err).NotTo(HaveOccurred())
		_, err = utils.Tanzu(nil, "package", "installed", "delete", packageInstalledName, "-n", packageInstalledNamespace, "-y")
		Expect(err).NotTo(HaveOccurred())

		utils.ValidatePackageInstallNotFound(packageInstalledNamespace, packageName)
		utils.ValidateDaemonsetNotFound("kube-system", "kube-multus-ds-amd64")
		utils.ValidateDaemonsetNotFound("kube-system", "install-cni-plugins")
		utils.ValidatePodNotFound("default", "macvlan-worker")

		// Clean up leftover files on nodes.
		_, err = utils.Kubectl(nil, "apply", "-f", filepath.Join(curDir, "multihomed-testfiles/cleanup.yaml"))
		Expect(err).NotTo(HaveOccurred())
		utils.ValidateDaemonsetReady("kube-system", "cleanup-conflists")
		_, err = utils.Kubectl(nil, "delete", "-f", filepath.Join(curDir, "multihomed-testfiles/cleanup.yaml"))
		Expect(err).NotTo(HaveOccurred())
		utils.ValidateDaemonsetNotFound("kube-system", "cleanup-conflists")
	})

	It("Check macvlan interfaces status", func() {
		By("check pod interface: net1")
		_, err := utils.Kubectl(nil, "exec", "macvlan-worker", "--",
			"bash", "-c", checkMacVlanIFaceExistence)
		Expect(err).NotTo(HaveOccurred())

		By("check pod interface status: net1")
		ifstatus, err := utils.Kubectl(nil, "exec", "macvlan-worker", "--",
			"bash", "-c", checkIFaceStatus)
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.Replace(ifstatus, "\n", "", -1)).To(Equal("UP"))

		By("get pod interface address: net1")
		ip, err := utils.Kubectl(nil, "exec", "macvlan-worker", "--",
			"bash", "-c", getMacVlanIP)
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.Replace(ip, "\n", "", -1)).To(Equal("10.1.1.12"))
	})
})
