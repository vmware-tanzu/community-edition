// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// package e2e_test implements running the cert-manager end-to-end tests
package e2e_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cert-manager Package E2E Test Suite")
}

const (
	packagePollInterval = "10s"
	packagePollTimeout  = "20m"
)

var (
	// packageInstallNamespace is the namespace where the package is installed
	packageInstallNamespace string

	// packageInstallName is the app name of the installed package
	packageInstallName string

	// installedPackages record the packages installed by the test
	installedPackages []string
)

var _ = BeforeSuite(func() {
	packageInstallNamespace = "default"

	packageInstallName = "cert-manager"

	By("installing cert-manager package")

	packageName := utils.TanzuPackageName(packageInstallName)

	version := findPackageAvailableVersion(packageName, "1.7.2")

	valuesFilename := filepath.Join("fixtures", "values.yaml")
	installPackage(packageInstallName, packageName, version, valuesFilename)

	By("validating cert-manager package is reconciled")
	utils.ValidatePackageInstallReady(packageInstallNamespace, packageInstallName)

	By("wait for cert-manager to be ready")
	waitForCertManagerToBeReady()
})

var _ = AfterSuite(func() {
	for _, installedPackage := range installedPackages {
		By(fmt.Sprintf("cleaning up %s package", installedPackage))
		_, err := utils.Tanzu(nil, "package", "installed", "delete", installedPackage,
			"--poll-interval", packagePollInterval,
			"--poll-timeout", packagePollTimeout,
			"--namespace", packageInstallNamespace, "--yes")
		Expect(err).NotTo(HaveOccurred())
	}

	By("validating the cert-manager package install no longer exists")
	utils.ValidatePackageInstallNotFound(packageInstallNamespace, packageInstallName)

	By(fmt.Sprintf("cleaning up %s namespace", packageInstallName))
	utils.Kubectl(nil, "delete", "namespace", packageInstallName) // nolint:errcheck
})

func findPackageAvailableVersion(packageName string, versionSubstr string) string {
	packageVersionJSON, err := utils.Tanzu(nil, "package", "available", "list", packageName, "-o", "json")
	Expect(err).NotTo(HaveOccurred())
	versions := []map[string]string{}

	err = json.Unmarshal([]byte(packageVersionJSON), &versions)
	Expect(err).NotTo(HaveOccurred())
	Expect(len(versions)).To(BeNumerically(">", 0))

	var matchedVersions []string
	for _, v := range versions {
		if versionSubstr == "" || strings.Contains(v["version"], versionSubstr) {
			matchedVersions = append(matchedVersions, v["version"])
		}
	}

	Expect(len(matchedVersions)).To(BeNumerically(">", 0), fmt.Sprintf("version contains %s for package %s not found", versionSubstr, packageName))

	return matchedVersions[len(matchedVersions)-1]
}

func waitForCertManagerToBeReady() {
	EventuallyWithOffset(1, func() (string, error) {
		kubectl, err := utils.Kubectl(nil, "get", "validatingwebhookconfigurations.admissionregistration.k8s.io", "cert-manager-webhook", "-o", "jsonpath={.webhooks[].clientConfig.caBundle}")
		fmt.Println(kubectl)
		return kubectl, err
	}, packagePollTimeout, packagePollInterval).Should(Not(BeEmpty()), "waiting for cert-manager to be ready")
}

func installPackage(name, packageName, version, valuesFilename string) {
	installedPackages = append([]string{name}, installedPackages...)

	args := []string{
		"package", "install", name,
		"--poll-interval", packagePollInterval,
		"--poll-timeout", packagePollTimeout,
		"--namespace", packageInstallNamespace,
		"--package-name", packageName,
		"--version", version,
	}

	if valuesFilename != "" {
		args = append(args, "--values-file", valuesFilename)
	}

	_, err := utils.Tanzu(nil, args...)
	Expect(err).NotTo(HaveOccurred())
}
