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
	// The net1 is the default iface name when not specified.
	checkMacVlanIFaceExistence = "ip a show dev net1"
	checkIFaceStatus           = "ip -j a show net1 | jq -r '.[]|.operstate'"
	getSecondaryIP             = "ip -j a show | jq -r '.[]|select(.ifname ==\"net1\")|.addr_info[]|select(.family==\"inet\").local'"
)

var _ = Describe("Whereabouts Addon E2E Test", func() {
	BeforeEach(func() {
		By("Install Multus CNI and whereabouts Package")

		curDir, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		multusCNIDir := filepath.Join(curDir, "../../../../multus-cni/", multusCNIPackage.version, "/test/e2e/")

		multusCNIInstallOptions := []string{
			"package", "install", multusCNIPackage.installedName,
			"--package-name", multusCNIPackage.name,
			"--version", multusCNIPackage.fullVersion,
			"-n", multusCNIPackage.installedNamespace}
		if multusCNIUseConfFile {
			multusCNIInstallOptions = append(multusCNIInstallOptions, "-f", filepath.Join(multusCNIDir, "multihomed-testfiles/multus-cni-values.yaml"))
		}
		_, err = utils.Tanzu(nil, multusCNIInstallOptions...)
		Expect(err).NotTo(HaveOccurred())
		_, err = utils.Kubectl(nil, "apply", "-f", filepath.Join(multusCNIDir, "multihomed-testfiles/cni-install.yaml"))
		Expect(err).NotTo(HaveOccurred())

		whereaboutsInstallOptions := []string{
			"package", "install", whereaboutsPackage.installedName,
			"--package-name", whereaboutsPackage.name,
			"--version", whereaboutsPackage.fullVersion,
			"-n", whereaboutsPackage.installedNamespace}
		if whereaboutsUseConfFile {
			whereaboutsInstallOptions = append(whereaboutsInstallOptions, "-f", filepath.Join(curDir, "test_assets/whereabouts_values.yaml"))
		}
		_, err = utils.Tanzu(nil, whereaboutsInstallOptions...)
		Expect(err).NotTo(HaveOccurred())

		utils.ValidatePackageInstallReady(multusCNIPackage.installedNamespace, multusCNIPackage.installedName)
		utils.ValidateDaemonsetReady("kube-system", "kube-multus-ds-amd64")
		utils.ValidatePackageInstallReady(whereaboutsPackage.installedNamespace, whereaboutsPackage.installedName)
		utils.ValidateDaemonsetReady("kube-system", "whereabouts")

		_, err = utils.Kubectl(nil, "apply", "-f", filepath.Join(curDir, "test_assets/multi_nic_pod.yaml"))
		Expect(err).NotTo(HaveOccurred())

		// Need to install CNI plugins in case of using environment like TCE on docker
		utils.ValidateDaemonsetReady("kube-system", "install-cni-plugins")
		utils.ValidatePodReady("default", "multi-nic-pod")

	})

	JustBeforeEach(func() {
		By("Install jq command in pod")
		_, err := utils.Kubectl(nil, "exec", "multi-nic-pod", "--", "yum", "install", "-y", "jq")
		Expect(err).NotTo(HaveOccurred())
	})

	JustAfterEach(func() {
		if CurrentGinkgoTestDescription().Failed {
			fmt.Fprintf(GinkgoWriter, "\nTest failed, please check these resources.\n")
			_, err := utils.Kubectl(nil, "get", "app", multusCNIPackage.installedName, "-n", multusCNIPackage.installedNamespace, "-o", "jsonpath={.status}")
			Expect(err).NotTo(HaveOccurred())
			_, err = utils.Kubectl(nil, "get", "daemonset", "kube-multus-ds-amd64", "-n", "kube-system", "-o", "jsonpath={.status}")
			Expect(err).NotTo(HaveOccurred())
		}
	})

	AfterEach(func() {
		By("Delete Multus CNI and whereabouts test resources")
		curDir, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		multusCNIDir := filepath.Join(curDir, "../../../../multus-cni/", multusCNIPackage.version, "/test/e2e/")

		_, err = utils.Kubectl(nil, "delete", "-f", filepath.Join(curDir, "test_assets/multi_nic_pod.yaml"))
		Expect(err).NotTo(HaveOccurred())
		_, err = utils.Kubectl(nil, "delete", "-f", filepath.Join(multusCNIDir, "multihomed-testfiles/cni-install.yaml"))
		Expect(err).NotTo(HaveOccurred())
		_, err = utils.Tanzu(nil, "package", "installed", "delete", multusCNIPackage.installedName, "-n", multusCNIPackage.installedNamespace, "-y")
		Expect(err).NotTo(HaveOccurred())
		_, err = utils.Tanzu(nil, "package", "installed", "delete", whereaboutsPackage.installedName, "-n", whereaboutsPackage.installedNamespace, "-y")
		Expect(err).NotTo(HaveOccurred())

		utils.ValidatePackageInstallNotFound(multusCNIPackage.installedNamespace, multusCNIPackage.name)
		utils.ValidateDaemonsetNotFound("kube-system", "kube-multus-ds-amd64")
		utils.ValidatePackageInstallNotFound(whereaboutsPackage.installedNamespace, whereaboutsPackage.name)
		utils.ValidateDaemonsetNotFound("kube-system", "whereabouts")
		utils.ValidateDaemonsetNotFound("kube-system", "install-cni-plugins")
		utils.ValidatePodNotFound("default", "multi-nic-pod")

		// Some configuration files are left over on nodes after removing multus-cni package, we need to clean them up manually.
		// More details are in https://github.com/vmware-tanzu/community-edition/tree/main/addons/packages/multus-cni/3.7.1#uninstallation-of-multus-cni-package
		_, err = utils.Kubectl(nil, "apply", "-f", filepath.Join(multusCNIDir, "multihomed-testfiles/cleanup.yaml"))
		Expect(err).NotTo(HaveOccurred())
		utils.ValidateDaemonsetReady("kube-system", "cleanup-conflists")
		_, err = utils.Kubectl(nil, "delete", "-f", filepath.Join(multusCNIDir, "multihomed-testfiles/cleanup.yaml"))
		Expect(err).NotTo(HaveOccurred())
		utils.ValidateDaemonsetNotFound("kube-system", "cleanup-conflists")

	})

	It("Check Pod secondary interface", func() {
		By("check pod interface: net1")
		_, err := utils.Kubectl(nil, "exec", "multi-nic-pod", "--",
			"bash", "-c", checkMacVlanIFaceExistence)
		Expect(err).NotTo(HaveOccurred())

		By("check pod interface status: net1")
		ifstatus, err := utils.Kubectl(nil, "exec", "multi-nic-pod", "--",
			"bash", "-c", checkIFaceStatus)
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.Replace(ifstatus, "\n", "", -1)).To(Equal("UP"))

		By("check pod interface address: net1")
		ip, err := utils.Kubectl(nil, "exec", "multi-nic-pod", "--",
			"bash", "-c", getSecondaryIP)
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.HasPrefix(strings.Replace(ip, "\n", "", -1), "192.168.20."))
	})
})
