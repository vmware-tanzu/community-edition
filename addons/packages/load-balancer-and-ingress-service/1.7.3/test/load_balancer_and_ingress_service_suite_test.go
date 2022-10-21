// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package loadbalancer_and_ingress_service_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLoadBalancerAndIngressService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LoadBalancer And Ingress Service Addons Templates Suite")
}
