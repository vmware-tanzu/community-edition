// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin"
)

var pluginName = "package"

var descriptor = plugin.NewTestFor(pluginName)

func main() {
	p, err := plugin.NewPlugin(descriptor)
	if err != nil {
		log.Fatal(err)
	}
	p.Cmd.RunE = test
	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}

func test(c *cobra.Command, _ []string) error {
	m := plugin.NewMain(pluginName, c, Cleanup)
	defer m.Finish()
	// TODO: setup test

	err := m.RunTest(
		"list package",
		"package list -o json",
		func(t *plugin.Test) error {
			// TODO: do some work...
			return nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// Cleanup the test.
func Cleanup() error {
	return nil
}
