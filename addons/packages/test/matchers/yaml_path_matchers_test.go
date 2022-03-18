// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package matchers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FindDocsMatchingYAMLPath", func() {
	It("finds the doc matching the supplied path", func() {
		docs := `---
foo: bar
---
baz: bing
---
blarg: thing
`
		nodes, err := FindDocsMatchingYAMLPath(docs, map[string]string{
			"$.baz": "bing",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(nodes).To(HaveLen(1))
		Expect(nodes[0]).To(Equal(`---
baz: bing
`))
	})

	It("finds multiple docs matching the supplied path", func() {
		docs := `---
foo: bar
---
baz: bing
---
foo: bar
`
		nodes, err := FindDocsMatchingYAMLPath(docs, map[string]string{
			"$.foo": "bar",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(nodes).To(HaveLen(2))
		Expect(nodes[0]).To(Equal(`---
foo: bar
`))
		Expect(nodes[1]).To(Equal(`---
foo: bar
`))
	})

	It("finds doc matching multiple supplied paths", func() {
		docs := `---
foo: bar
namespace: a
---
baz: bing
---
foo: bar
namespace: b
`
		nodes, err := FindDocsMatchingYAMLPath(docs, map[string]string{
			"$.foo":       "bar",
			"$.namespace": "b",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(nodes).To(HaveLen(1))
		Expect(nodes[0]).To(Equal(`---
foo: bar
namespace: b
`))
	})

	It("returns an error when given an invalid yaml path", func() {
		docs := `---
foo: bar
namespace: a
`
		_, err := FindDocsMatchingYAMLPath(docs, map[string]string{
			"foo bar": "bar",
		})
		Expect(err).To(MatchError("invalid yaml path \"foo bar\" err: invalid character ' ' at position 3, following \"foo\""))
	})

	Describe("a path describes multiple nodes", func() {
		When("all found nodes match value", func() {
			It("returns the doc", func() {
				docs := `---
a:
  foo: bar
b:
  foo: bar
`
				nodes, err := FindDocsMatchingYAMLPath(docs, map[string]string{
					"$..foo": "bar",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(nodes).To(HaveLen(1))
				Expect(nodes[0]).To(Equal(`---
a:
  foo: bar
b:
  foo: bar
`))
			})
		})

		When("all found nodes do not match value", func() {
			It("does not return the doc", func() {
				docs := `---
a:
  foo: bar
b:
  foo: baz
`
				nodes, err := FindDocsMatchingYAMLPath(docs, map[string]string{
					"$..foo": "bar",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(nodes).To(HaveLen(0))
			})
		})
	})
})
