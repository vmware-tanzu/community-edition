// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// package e2e_test implements running the external DNS end to end tests
package e2e_test

import (
	"os"
	"testing"
	"time"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExternalDNSE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "External-DNS Addon Package E2E Test Suite")
}

const (
	httpRequestTimeout  = 60 * time.Second
	httpRequestInterval = 5 * time.Second
)

var (
	// dockerhubProxy is an optional configuration option (provided by using
	// DOCKERHUB_PROXY), that allows you to override docker.io with a proxy to
	// docker.io to avoid any potential issues with rate-limiting.
	dockerhubProxy string

	// packageInstallNamespace is the namespace where the external-dns package is
	// installed (i.e this is the namespace tanzu package install is called
	// with)
	packageInstallNamespace string

	// packageComponentsNamespace is the namespace where the external-dns
	// package components are installed  (e.g. the external DNS deployment).
	// This is the namespace passed into the external-dns values.yaml). This
	// namespace is created by the package installation.
	packageComponentsNamespace string

	// fixtureNamespace is the namespace where all test fixtures are created
	// for the purpose of testing the addon (e.g bind, kuard, dnsutils)
	fixtureNamespace string
)

var _ = BeforeSuite(func() {
	dockerhubProxy = os.Getenv("DOCKERHUB_PROXY")
	if dockerhubProxy == "" {
		dockerhubProxy = "docker.io"
	}

	packageComponentsNamespace = "e2e-external-dns-package-components"
	fixtureNamespace = "e2e-external-dns-fixtures"
	packageInstallNamespace = "e2e-external-dns-package"

	_, err := utils.Kubectl(nil, "create", "namespace", packageInstallNamespace)
	Expect(err).NotTo(HaveOccurred())

	_, err = utils.Kubectl(nil, "create", "namespace", fixtureNamespace)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	_, _ = utils.Kubectl(nil, "delete", "namespace", fixtureNamespace)
	_, _ = utils.Kubectl(nil, "delete", "namespace", packageInstallNamespace)
})
