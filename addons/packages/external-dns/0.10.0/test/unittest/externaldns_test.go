// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package externaldns_test

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/repo"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/ytt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("External DNS Ytt Templates", func() {
	var (
		values string
		output string
		err    error

		configDir = filepath.Join(repo.RootDir(), "addons/packages/external-dns/0.10.0/bundle/config")

		ValuesFromFile = func(filename string) string {
			data, err := ioutil.ReadFile(filepath.Join(repo.RootDir(), "addons/packages/external-dns/0.10.0/test/unittest/fixtures/values", filename))
			Expect(err).NotTo(HaveOccurred())

			return string(data)
		}

		ExpectOutputEqualToFile = func(filename string) {
			data, err := ioutil.ReadFile(filepath.Join(repo.RootDir(), "addons/packages/external-dns/0.10.0/test/unittest/fixtures/expected", filename))
			Expect(err).NotTo(HaveOccurred())

			//fmt.Println(output)
			Expect(output).To(BeEquivalentTo(string(data)))
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

		output, err = ytt.RenderYTTTemplate(ytt.CommandOptions{}, filePaths, strings.NewReader(values))
	})

	Context("No configuration", func() {
		It("renders with an error", func() {
			Expect(err).To(ContainSubstring("configuration is required for external-dns"))
		})
	})

	Context("No --source in deployment.args", func() {
		BeforeEach(func() {
			values = ValuesFromFile("deployment-args-no-source.yaml")
		})

		It("renders with an error", func() {
			Expect(err).To(ContainSubstring("--source is required in deployment.args to query for endpoints"))
		})
	})

	Context("No --provider in deployment.args", func() {
		BeforeEach(func() {
			values = ValuesFromFile("deployment-args-no-provider.yaml")
		})

		It("renders with an error", func() {
			Expect(err).To(ContainSubstring("--provider is required in deployment.args to define a DNS provider where records will be created"))
		})
	})

	Context("Providing a minimal configuration", func() {
		BeforeEach(func() {
			values = ValuesFromFile("minimal-configuration.yaml")
		})

		It("renders a working setup", func() {
			Expect(err).NotTo(HaveOccurred())
			ExpectOutputEqualToFile("minimal-configuration.yaml")
		})
	})

	Context("Providing a namespace", func() {
		BeforeEach(func() {
			values = ValuesFromFile("namespace.yaml")
		})

		It("renders a setup in a different namespace", func() {
			Expect(err).NotTo(HaveOccurred())
			ExpectOutputEqualToFile("namespace.yaml")
		})
	})

	Context("Providing env vars for the deployment", func() {
		BeforeEach(func() {
			values = ValuesFromFile("deployment-env-vars.yaml")
		})

		It("renders a deployment with env vars", func() {
			Expect(err).NotTo(HaveOccurred())
			ExpectOutputEqualToFile("deployment-env-vars.yaml")
		})
	})

	Context("Providing the security context for the deployment", func() {
		BeforeEach(func() {
			values = ValuesFromFile("deployment-security-context.yaml")
		})

		It("renders a deployment with a custom security context", func() {
			Expect(err).NotTo(HaveOccurred())
			ExpectOutputEqualToFile("deployment-security-context.yaml")
		})
	})

	Context("Providing volumes and their mounts for the deployment", func() {
		BeforeEach(func() {
			values = ValuesFromFile("deployment-volumes.yaml")
		})

		It("renders a deployment with additional volumes mounted", func() {
			Expect(err).NotTo(HaveOccurred())
			ExpectOutputEqualToFile("deployment-volumes.yaml")
		})
	})

	Context("Providing annotations for the service account", func() {
		BeforeEach(func() {
			values = ValuesFromFile("serviceaccount-annotations.yaml")
		})

		It("renders a service account with annotations", func() {
			Expect(err).NotTo(HaveOccurred())
			ExpectOutputEqualToFile("serviceaccount-annotations.yaml")
		})
	})
})
