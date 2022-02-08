// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin"
)

var descriptor = plugin.NewTestFor("diagnostics")

func main() {
	retcode := 0

	defer func() { os.Exit(retcode) }()
	defer Cleanup()

	p, err := plugin.NewPlugin(descriptor)
	if err != nil {
		log.Println(err)
		retcode = 1
		return
	}
	p.Cmd.RunE = test
	if err := p.Execute(); err != nil {
		retcode = 1
		return
	}
}

func test(c *cobra.Command, _ []string) error {
	return nil
}

// Cleanup the test.
func Cleanup() {}
