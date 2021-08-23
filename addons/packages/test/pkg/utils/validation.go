// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"time"

	"github.com/onsi/gomega"
)

var (
	DeploymentTimeout       = 120 * time.Second
	DeploymentCheckInterval = 5 * time.Second
)

func ValidateDeploymentReady(namespace, name string) {
	gomega.EventuallyWithOffset(1, func() (string, error) {
		return Kubectl(nil, "-n", namespace, "get", "deployment", name, "-o", "jsonpath={.status.conditions[?(@.type == 'Available')].status}")
	}, DeploymentTimeout, DeploymentCheckInterval).Should(gomega.Equal("True"), fmt.Sprintf("%s/%s deployment was never ready", namespace, name))
}

func ValidatePodReady(namespace, name string) {
	gomega.EventuallyWithOffset(1, func() (string, error) {
		return Kubectl(nil, "-n", namespace, "get", "pod", name, "-o", "jsonpath={.status.conditions[?(@.type == 'Ready')].status}")
	}, DeploymentTimeout, DeploymentCheckInterval).Should(gomega.Equal("True"), fmt.Sprintf("%s/%s pod was never ready", namespace, name))
}

func ValidatePackageReady(namespace, name string) {
	gomega.EventuallyWithOffset(1, func() (string, error) {
		return Kubectl(nil, "-n", namespace, "get", "installedpackage", name, "-o", "jsonpath={.status.conditions[?(@.type == 'ReconcileSucceeded')].status}")
	}, DeploymentTimeout, DeploymentCheckInterval).Should(gomega.Equal("True"), fmt.Sprintf("%s/%s installedpackage was never ready", namespace, name))
}

func ValidateDeploymentNotFound(namespace, name string) {
	gomega.Eventually(func() error {
		_, err := Kubectl(nil, "-n", namespace, "get", "deployment", name)
		return err
	}, DeploymentTimeout, DeploymentCheckInterval).Should(gomega.MatchError(gomega.Or(
		gomega.ContainSubstring(fmt.Sprintf(`deployments.apps %q not found`, name)),
		gomega.ContainSubstring(fmt.Sprintf(`namespaces %q not found`, namespace)),
	)), fmt.Sprintf("%s/%s deployment was never deleted", namespace, name))
}

func ValidatePodNotFound(namespace, name string) {
	gomega.EventuallyWithOffset(1, func() error {
		_, err := Kubectl(nil, "-n", namespace, "get", "pod", name)
		return err
	}, DeploymentTimeout, DeploymentCheckInterval).Should(gomega.MatchError(gomega.Or(
		gomega.ContainSubstring(fmt.Sprintf(`pods %q not found`, name)),
		gomega.ContainSubstring(fmt.Sprintf(`namespaces %q not found`, namespace)),
	)), fmt.Sprintf("%s/%s pod was never deleted", namespace, name))
}

func ValidatePackageNotFound(namespace, name string) {
	gomega.EventuallyWithOffset(1, func() error {
		_, err := Kubectl(nil, "-n", namespace, "get", "installedpackage", name)
		return err
	}, DeploymentTimeout, DeploymentCheckInterval).Should(gomega.MatchError(gomega.Or(
		gomega.ContainSubstring(fmt.Sprintf(`installedpackages.install.package.carvel.dev %q not found`, name)),
		gomega.ContainSubstring(fmt.Sprintf(`namespaces %q not found`, namespace)),
	)), fmt.Sprintf("%s/%s installedpackage was never deleted", namespace, name))
}
