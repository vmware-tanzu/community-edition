// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/instrumenta/kubeval/kubeval"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/repo"
)

var _ = Describe("SR-IOV CNI Template Test", func() {
	var (
		err            error
		renderedOutput string
		expected       []byte
		templatePath   []string

		valuesFileDir     = filepath.Join(repo.RootDir(), "addons/packages/sriov-cni/2.6.1/test/unit/sriov-conf")
		expectedResultDir = filepath.Join(repo.RootDir(), "addons/packages/sriov-cni/2.6.1/test/unit/expected")
		baseTemplatesDir  = filepath.Join(repo.RootDir(), "addons/packages/sriov-cni/2.6.1/bundle/config")
		testFiles         = map[string]string{
			"base":           "base.yaml",
			"args.yaml":      "args-expected.yaml",
			"resources.yaml": "resources-expected.yaml",
			"namespace.yaml": "namespace-expected.yaml",
		}
	)
	It("Should be the same with expected file", func() {
		for k, v := range testFiles {
			By("Check generated templates content", func() {
				fmt.Fprintf(GinkgoWriter, fmt.Sprintf("Compare %s with %s\n", k, v))
			})
			// Render templates with values
			if k == "base" {
				templatePath = []string{baseTemplatesDir}
			} else {
				templatePath = []string{baseTemplatesDir, filepath.Join(valuesFileDir, k)}
			}
			renderedOutput, err = yttCli.RenderTemplate(templatePath, nil)
			Expect(err).NotTo(HaveOccurred())

			// Check if same as expected results
			expected, err = ioutil.ReadFile(filepath.Join(expectedResultDir, v))
			Expect(err).NotTo(HaveOccurred())
			Expect(strings.Compare(renderedOutput, string(expected))).Should(BeNumerically("==", 0))

			// Check if resouces can be created successfully.
			// Just a kubeval
			results, err := kubeval.Validate([]byte(renderedOutput), kubevalConfig)
			Expect(err).NotTo(HaveOccurred())

			for _, result := range results {
				Expect(len(result.Errors)).Should(BeNumerically("==", 0))
				fmt.Fprintf(GinkgoWriter,
					fmt.Sprintf("resource %s/%s %s in %s is ok to create as defined in file %s\n", result.APIVersion, result.Kind, result.ResourceName, result.ResourceNamespace, v))
			}
		}

	})
})
