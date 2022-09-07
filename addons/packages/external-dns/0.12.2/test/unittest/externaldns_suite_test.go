// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package externaldns_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExternalDNS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "External DNS Addons Templates Suite")
}
