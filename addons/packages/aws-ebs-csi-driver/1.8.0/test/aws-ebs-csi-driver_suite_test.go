// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package awsebscsidrivertest

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAwsEBSCSIDriver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AWS EBS CSI Driver Addons Templates Suite")
}
