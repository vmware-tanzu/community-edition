// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"os"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin"
	conformance "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/conformance/pkg"
)

var descriptor = plugin.PluginDescriptor{
	Name:        "conformance",
	Description: "Run Sonobuoy conformance tests against clusters",
	Group:       plugin.RunCmdGroup,
	Version:     plugin.Version,
}

var logLevel int32

func main() {
	p, err := plugin.NewPlugin(&descriptor)
	if err != nil {
		log.Fatalf("%v", err)
	}

	p.AddCommands(
		conformance.RunCmd,
		conformance.RetrieveCmd,
		conformance.DeleteCmd,
		conformance.LogsCmd,
		conformance.StatusCmd,
		conformance.ResultsCmd,
		conformance.GenCmd,
	)

	// Remove the generated version command and replace it with ours,
	// so we get more version information.
	c, _, err := p.Cmd.Find([]string{"version"})

	if err != nil {
		log.Fatalf("%v", err)
	}

	p.Cmd.RemoveCommand(c)

	p.AddCommands(
		conformance.VersionCmd,
	)

	p.Cmd.PersistentFlags().Int32VarP(&logLevel, "verbose", "v", 0, "Number for the log level verbosity(0-9)")
	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}
