// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/ytt"
)

var (
	// Global ytt commands
	yttCli *ytt.Command
)

func TestSRIOVDPE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SR-IOV NETWORK DEVICE PLUGIN Test Suite")
}

var _ = BeforeSuite(func() {
	// Initialize ytt command
	options := ytt.CommandOptions{
		FailOnUnknownComments:  false,
		Strict:                 false,
		DangerousAllowSymlinks: false,
	}

	yttCli = ytt.NewYttCommand(options)
})
