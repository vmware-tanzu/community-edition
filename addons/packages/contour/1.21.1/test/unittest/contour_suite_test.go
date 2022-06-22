// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package contour_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestContour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Contour Addon Templates Suite")
}
