// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"time"

	"github.com/onsi/gomega"
)

var (
	Timeout       = 300 * time.Second
	CheckInterval = 5 * time.Second
)

func ValidateAllDeploymentsReady() {
	gomega.EventuallyWithOffset(1, func() (string, error) {
		return Kubectl(nil, "get", "deployment", "-A", "-o", "jsonpath='{.items[*].status.conditions[?(@.type=='Available')].status}'")
	}, Timeout, CheckInterval).ShouldNot(gomega.ContainSubstring("False"))
}

func ValidateDeploymentReady(namespace, name string) {
	gomega.EventuallyWithOffset(1, func() (string, error) {
		return Kubectl(nil, "-n", namespace, "get", "deployment", name, "-o", "jsonpath={.status.conditions[?(@.type == 'Available')].status}")
	}, Timeout, CheckInterval).Should(gomega.Equal("True"), fmt.Sprintf("%s/%s deployment was never ready", namespace, name))
}
