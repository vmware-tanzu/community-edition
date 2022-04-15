// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package unit_test

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/repo"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/ytt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gatekeeper Ytt Templates", func() {
	var (
		values string
		output string
		err    error

		configDir = filepath.Join(repo.RootDir(), "addons/packages/gatekeeper/3.7.1/bundle/config")

		ValuesFromFile = func(filename string) string {
			data, err := ioutil.ReadFile(filepath.Join(repo.RootDir(), "addons/packages/gatekeeper/3.7.1/test/unittest/fixtures/values", filename))
			Expect(err).NotTo(HaveOccurred())

			return string(data)
		}

		ExpectOutputEqualToFile = func(filename string) {
			data, err := ioutil.ReadFile(filepath.Join(repo.RootDir(), "addons/packages/gatekeeper/3.7.1/test/unittest/fixtures/expected", filename))
			Expect(err).NotTo(HaveOccurred())

			Expect(output).To(MatchYAML(string(data)))
		}
	)

	BeforeEach(func() {
		values = ""
	})

	JustBeforeEach(func() {
		var filePaths []string

		for _, p := range []string{"upstream/*.yaml", "overlays/*.yaml", "*.yaml"} {
			matches, err := filepath.Glob(filepath.Join(configDir, p))

			Expect(err).NotTo(HaveOccurred())
			filePaths = append(filePaths, matches...)
		}

		filePaths = append(filePaths,
			filepath.Join(repo.RootDir(), "addons/packages/gatekeeper/3.7.1/test/unittest/fixtures/values/default.yaml"),
		)
		output, err = ytt.RenderYTTTemplate(ytt.CommandOptions{}, filePaths, strings.NewReader(values))
	})

	Context("default", func() {
		BeforeEach(func() {
			values = ""
		})

		It("renders with a default configuration", func() {
			Expect(err).NotTo(HaveOccurred())
			ExpectOutputEqualToFile("default.yaml")
		})
	})

	Context("updates the namespace", func() {
		BeforeEach(func() {
			values = ValuesFromFile("custom-namespace.yaml")
		})

		It("renders with a custom namespace", func() {
			Expect(err).NotTo(HaveOccurred())
			ExpectOutputEqualToFile("namespace.yaml")
		})
	})
})
