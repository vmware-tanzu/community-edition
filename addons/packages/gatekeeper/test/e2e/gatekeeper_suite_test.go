package e2e_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGateKeeperE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gatekeeper addon package e2e test suite")
}
