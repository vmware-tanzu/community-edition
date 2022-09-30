// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kube_vip_cloud_provider_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKubeVipCloudProivder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "kube-vip CloudProvider Addons Templates Suite")
}
