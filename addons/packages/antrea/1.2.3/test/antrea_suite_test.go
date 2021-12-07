// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package antrea_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAntrea(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Antrea Addons Templates Suite")
}
