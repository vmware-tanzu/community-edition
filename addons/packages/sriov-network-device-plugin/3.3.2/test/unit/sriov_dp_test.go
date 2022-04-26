// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/instrumenta/kubeval/kubeval"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/repo"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/utils"
)

var _ = Describe("SR-IOV NETWORK DEVICE PLUGIN Template Test", func() {
	var (
		err            error
		renderedOutput string
		expected       []byte
		templatePath   []string
		targetFile     *os.File

		valuesFileDir     = filepath.Join(repo.RootDir(), "addons/packages/sriov-network-device-plugin/3.3.2/test/unit/sriov-conf")
		expectedResultDir = filepath.Join(repo.RootDir(), "addons/packages/sriov-network-device-plugin/3.3.2/test/unit/expected")
		baseTemplatesDir  = filepath.Join(repo.RootDir(), "addons/packages/sriov-network-device-plugin/3.3.2/bundle/config")
		testFiles         = map[string]string{
			"base":            "base.yaml",
			"node-pools.yaml": "node-pools-expected.yaml",
			"resources.yaml":  "resources-expected.yaml",
		}
	)
	It("Should be the same with expected base file", func() {
		for k, v := range testFiles {
			By("Check generated templates content", func() {
				fmt.Fprintf(GinkgoWriter, "Compare %s with %s\n", k, v)
			})
			// Render templates with values
			if strings.Compare(k, "base") == 0 {
				templatePath = []string{baseTemplatesDir}
			} else {
				templatePath = []string{baseTemplatesDir, filepath.Join(valuesFileDir, k)}
			}
			renderedOutput, err = yttCli.RenderTemplate(templatePath, nil)
			Expect(err).NotTo(HaveOccurred())

			// Check if same as expected results
			expected, err = os.ReadFile(filepath.Join(expectedResultDir, v))
			Expect(err).NotTo(HaveOccurred())
			Expect(strings.Compare(renderedOutput, string(expected))).Should(BeNumerically("==", 0))

			// Check if resouces can be created successfully.
			// Just a dry-run and kubeval
			targetFile, err = os.CreateTemp("", "sriov-dp-base-*.yaml")
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				targetFile.Close()
				os.Remove(targetFile.Name())
			}()

			_, err = targetFile.Write([]byte(renderedOutput))
			Expect(err).NotTo(HaveOccurred())

			_, err = utils.Kubectl(nil, "apply", "-f", targetFile.Name(), "--dry-run=client")
			Expect(err).NotTo(HaveOccurred())

			results, err := kubeval.Validate([]byte(renderedOutput), kubevalConfig)
			Expect(err).NotTo(HaveOccurred())

			for _, result := range results {
				Expect(len(result.Errors)).Should(BeNumerically("==", 0))
				fmt.Fprintf(GinkgoWriter,
					"resource %s/%s %s in %s is ok to create as defined in file %s\n", result.APIVersion, result.Kind, result.ResourceName, result.ResourceNamespace, v)
			}
		}

	})
})
