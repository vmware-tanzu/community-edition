// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/command/plugin"
)

var descriptor = cli.NewTestFor("diagnostics")

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
