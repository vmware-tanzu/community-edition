// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package matchers provides custom gomega matchers for tests.
package matchers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"
)

type pathWithValue struct {
	path  *yamlpath.Path
	value string
}

// FindDocsMatchingYAMLPath finds yaml documents that match all paths with values provided.
func FindDocsMatchingYAMLPath(yamlString string, pathsWithValues map[string]string) ([]string, error) {
	docStrings := strings.Split(yamlString, "---")

	var yamlPathWithValues []pathWithValue
	for path, value := range pathsWithValues {
		newPath, err := yamlpath.NewPath(path)
		if err != nil {
			return nil, fmt.Errorf("invalid yaml path %q err: %s", path, err)
		}
		yamlPathWithValues = append(yamlPathWithValues, pathWithValue{newPath, value})
	}

	var docsToReturn []string
	for _, doc := range docStrings {
		var node yaml.Node
		err := yaml.Unmarshal([]byte(doc), &node)
		if err != nil {
			return nil, err
		}

		if matchesAllPaths(&node, yamlPathWithValues) {
			docsToReturn = append(docsToReturn, fmt.Sprintf("---%s", doc))
		}
	}
	return docsToReturn, nil
}

func matchesAllPaths(node *yaml.Node, pathWithValues []pathWithValue) bool {
	for _, pathWithValue := range pathWithValues {
		results, err := pathWithValue.path.Find(node)
		if err != nil {
			panic("Find currently never returns an error, if we're getting this we need to handle these")
		}

		if len(results) == 0 {
			return false
		}

		for _, result := range results {
			if result.Value != pathWithValue.value {
				return false
			}
		}
	}
	return true
}

// HaveYAMLPath searches the *first* document in actual for a path
func HaveYAMLPath(path string) types.GomegaMatcher {
	return &HaveYAMLPathMatcher{
		Path: path,
	}
}

// HaveYAMLPathMatcher is the matcher returned by HaveYAMLPath.
type HaveYAMLPathMatcher struct {
	Path string
}

// Match returns true and no error if actual has a matching path, else it
// returns false.
func (matcher *HaveYAMLPathMatcher) Match(actual interface{}) (bool, error) {
	if actual == nil {
		return false, nil
	}

	if reflect.TypeOf(actual).Kind() != reflect.String {
		return false, fmt.Errorf("matcher HaveYAMLPathWithValue expects a string. Got: %s", format.Object(actual, 1))
	}

	yamlString := reflect.ValueOf(actual).String()

	var node yaml.Node
	err := yaml.Unmarshal([]byte(yamlString), &node)
	if err != nil {
		return false, fmt.Errorf("matcher HaveYAMLPathWithValue failed to unmarshal actual yaml with:\n%s%s", format.Indent, err.Error())
	}

	path, err := yamlpath.NewPath(matcher.Path)
	if err != nil {
		return false, fmt.Errorf("matcher HaveYAMLPathWithValue failed with:\n%s%s", format.Indent, err.Error())
	}

	q, err := path.Find(&node)
	if err != nil { // this should never occur, according to current path.Find documentation
		return false, fmt.Errorf("matcher HaveYAMLPathWithValue failed with:\n%s%s", format.Indent, err.Error())
	}

	return len(q) > 0, nil
}

// FailureMessage returns a human-readable failure message if Match returns false.
func (matcher *HaveYAMLPathMatcher) FailureMessage(actual interface{}) string {
	message := fmt.Sprintf("at path %q to exist\n", matcher.Path)
	return format.Message(actual, message)
}

// NegatedFailureMessage returns a human-readable negated failure message when NotTo Match fails.
func (matcher *HaveYAMLPathMatcher) NegatedFailureMessage(actual interface{}) string {
	message := fmt.Sprintf("at path %q not to exist\n", matcher.Path)
	return format.Message(actual, message)
}

// HaveYAMLPathWithValue searches the *first* document in actual for a path with value.
func HaveYAMLPathWithValue(path string, value interface{}) types.GomegaMatcher {
	return &HaveYAMLPathWithValueMatcher{
		Path:  path,
		Value: value,
	}
}

// HaveYAMLPathWithValueMatcher is the Matcher returned by HaveYAMLPathWithValue.
type HaveYAMLPathWithValueMatcher struct {
	Path        string
	Value       interface{}
	actualValue interface{}
}

// Match returns true and no error if actual has a matching path with value,
// else it returns false.
func (matcher *HaveYAMLPathWithValueMatcher) Match(actual interface{}) (bool, error) {
	if actual == nil {
		return false, nil
	}

	if reflect.TypeOf(actual).Kind() != reflect.String {
		return false, fmt.Errorf("matcher HaveYAMLPathWithValue expects a string. Got: %s", format.Object(actual, 1))
	}

	valueMatcher := &matchers.EqualMatcher{Expected: matcher.Value}

	yamlString := reflect.ValueOf(actual).String()

	var node yaml.Node
	err := yaml.Unmarshal([]byte(yamlString), &node)
	if err != nil {
		return false, fmt.Errorf("matcher HaveYAMLPathWithValue failed to unmarshal actual yaml with:\n%s%s", format.Indent, err.Error())
	}

	path, err := yamlpath.NewPath(matcher.Path)
	if err != nil {
		return false, fmt.Errorf("matcher HaveYAMLPathWithValue failed with:\n%s%s", format.Indent, err.Error())
	}

	q, err := path.Find(&node)
	if err != nil { // this should never occur, according to current path.Find documentation
		return false, fmt.Errorf("matcher HaveYAMLPathWithValue failed with:\n%s%s", format.Indent, err.Error())
	}

	if len(q) != 1 {
		return false, fmt.Errorf("matcher HaveYAMLPathWithValue expected to find one node at path %q, found: %d node(s)", matcher.Path, len(q))
	}
	matcher.actualValue = q[0].Value

	success, err := valueMatcher.Match(q[0].Value)
	if err != nil {
		return false, fmt.Errorf("matcher HaveYAMLPathWithValue failed with:\n%s%s", format.Indent, err.Error())
	}

	return success, nil
}

// FailureMessage returns a human-readable failure message if Match returns false.
func (matcher *HaveYAMLPathWithValueMatcher) FailureMessage(actual interface{}) string {
	message := fmt.Sprintf("at path %q to have value\n%s\nGot:", matcher.Path, format.Object(matcher.Value, 1))
	return format.Message(actual, message, matcher.actualValue)
}

// NegatedFailureMessage returns a human-readable negated failure message when NotTo Match fails.
func (matcher *HaveYAMLPathWithValueMatcher) NegatedFailureMessage(actual interface{}) string {
	message := fmt.Sprintf("at path %q not to have value\n%s\nGot:", matcher.Path, format.Object(matcher.Value, 1))
	return format.Message(actual, message, matcher.actualValue)
}
