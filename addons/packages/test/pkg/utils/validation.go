// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
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
	DaemonsetTimeout        = 120 * time.Second
	DaemonsetCheckInterval  = 5 * time.Second
)

func ValidateDeploymentReady(namespace, name string) {
	gomega.EventuallyWithOffset(1, func() (string, error) {
		return Kubectl(nil, "-n", namespace, "get", "deployment", name, "-o", "jsonpath={.status.conditions[?(@.type == 'Available')].status}")
	}, DeploymentTimeout, DeploymentCheckInterval).Should(gomega.Equal("True"), fmt.Sprintf("%s/%s deployment was never ready", namespace, name))
}

func ValidateDaemonsetReady(namespace, name string) {
	gomega.EventuallyWithOffset(1, func() (bool, error) {
		desiredNumberScheduled, err := Kubectl(nil, "-n", namespace, "get", "daemonset", name, "-o", "jsonpath='{.status.desiredNumberScheduled}'")
		if err != nil {
			return false, err
		}
		numberReady, err := Kubectl(nil, "-n", namespace, "get", "daemonset", name, "-o", "jsonpath='{.status.numberReady}'")
		if err != nil {
			return false, err
		}
		return desiredNumberScheduled == numberReady, nil
	}, DaemonsetTimeout, DaemonsetCheckInterval).Should(gomega.Equal(true), fmt.Sprintf("%s/%s daemonset was never ready", namespace, name))
}

func ValidatePodReady(namespace, name string) {
	gomega.EventuallyWithOffset(1, func() (string, error) {
		return Kubectl(nil, "-n", namespace, "get", "pod", name, "-o", "jsonpath={.status.conditions[?(@.type == 'Ready')].status}")
	}, DeploymentTimeout, DeploymentCheckInterval).Should(gomega.Equal("True"), fmt.Sprintf("%s/%s pod was never ready", namespace, name))
}

func ValidateLoadBalancerReady(namespace, name string) {
	gomega.EventuallyWithOffset(1, func() (string, error) {
		return Kubectl(nil, "-n", namespace, "get", "service", name, "-o", "jsonpath='{.status.loadBalancer.ingress[0].ip}'")
	}, DeploymentTimeout, DeploymentCheckInterval).ShouldNot(gomega.BeEmpty(), fmt.Sprintf("%s/%s load balancer service was never ready", namespace, name))
}

func ValidatePackageInstallReady(namespace, name string) {
	gomega.EventuallyWithOffset(1, func() (string, error) {
		return Kubectl(nil, "-n", namespace, "get", "packageinstalls", name, "-o", "jsonpath={.status.conditions[?(@.type == 'ReconcileSucceeded')].status}")
	}, DeploymentTimeout, DeploymentCheckInterval).Should(gomega.Equal("True"), fmt.Sprintf("%s/%s packageinstalls was never ready", namespace, name))
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

func ValidateDaemonsetNotFound(namespace, name string) {
	gomega.Eventually(func() error {
		_, err := Kubectl(nil, "-n", namespace, "get", "daemonset", name)
		return err
	}, DaemonsetTimeout, DaemonsetCheckInterval).Should(gomega.MatchError(gomega.Or(
		gomega.ContainSubstring(fmt.Sprintf(`daemonsets.apps %q not found`, name)),
		gomega.ContainSubstring(fmt.Sprintf(`namespaces %q not found`, namespace)),
	)), fmt.Sprintf("%s/%s daemonset was never deleted", namespace, name))
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

func ValidatePackageInstallNotFound(namespace, name string) {
	gomega.EventuallyWithOffset(1, func() error {
		_, err := Kubectl(nil, "-n", namespace, "get", "packageinstall", name)
		return err
	}, DeploymentTimeout, DeploymentCheckInterval).Should(gomega.MatchError(gomega.Or(
		gomega.ContainSubstring(fmt.Sprintf(`packageinstalls.packaging.carvel.dev %q not found`, name)),
		gomega.ContainSubstring(fmt.Sprintf(`namespaces %q not found`, namespace)),
	)), fmt.Sprintf("%s/%s packageinstall was never deleted", namespace, name))
}
