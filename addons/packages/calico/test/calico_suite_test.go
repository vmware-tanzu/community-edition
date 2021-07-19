// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package calico_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCalico(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Calico Addons Templates Suite")
}
