// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/repo"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/utils"
)

var _ = Describe("SR-IOV NETWORK DEVICE PLUGIN Test", func() {
	var (
		err            error
		testPath       string
		renderedOutput string
		expected       []byte
		targetFile     *os.File

		valuesFilePath    = filepath.Join(testPath, "sriov-conf/values.yaml")
		baseTemplatesPath = filepath.Join(repo.RootDir(), "addons/packages/sriov-network-device-plugin/3.3.2/bundle/config")
		expectedResultDir = filepath.Join(testPath, "expected")
	)
	BeforeEach(func() {
		testPath, err = os.Getwd()
		Expect(err).NotTo(HaveOccurred())
	})
	AfterEach(func() {
		_ = os.Remove(targetFile.Name())
	})
	Context("Check the template without values", func() {
		It("Should be the same with expected base file", func() {
			renderedOutput, err = yttCli.RenderTemplate([]string{baseTemplatesPath}, nil)
			Expect(err).NotTo(HaveOccurred())

			expected, err = ioutil.ReadFile(filepath.Join(expectedResultDir, "base.yaml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(strings.Compare(renderedOutput, string(expected))).Should(BeNumerically("==", 0))

			targetFile, err = ioutil.TempFile("", "sriov-dp-base-*.yaml")
			Expect(err).NotTo(HaveOccurred())
			defer targetFile.Close()

			_, err = targetFile.Write([]byte(renderedOutput))
			Expect(err).NotTo(HaveOccurred())

			_, err = utils.Kubectl(nil, "apply", "-f", targetFile.Name(), "--dry-run=client")
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Context("Check the template with values", func() {
		It("Should be the same with expected with values file", func() {
			renderedOutput, err = yttCli.RenderTemplate([]string{baseTemplatesPath, valuesFilePath}, nil)
			Expect(err).NotTo(HaveOccurred())

			expected, err = ioutil.ReadFile(filepath.Join(expectedResultDir, "sample.yaml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(strings.Compare(renderedOutput, string(expected))).Should(BeNumerically("==", 0))

			targetFile, err = ioutil.TempFile("", "sriov-dp-with-values-*.yaml")
			Expect(err).NotTo(HaveOccurred())
			defer targetFile.Close()

			_, err = targetFile.Write([]byte(renderedOutput))
			Expect(err).NotTo(HaveOccurred())

			_, err = utils.Kubectl(nil, "apply", "-f", targetFile.Name(), "--dry-run=client")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
