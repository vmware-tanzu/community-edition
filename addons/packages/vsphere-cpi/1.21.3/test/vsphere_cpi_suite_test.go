// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package vspherecpi_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVsphereCPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "vSphere CPI Addons Templates Suite")
}
