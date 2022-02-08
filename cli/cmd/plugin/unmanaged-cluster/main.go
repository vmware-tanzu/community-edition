// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"os"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/cmd"
)

var description = `Deploy and manage single-node, static, Tanzu clusters.`

var descriptor = plugin.PluginDescriptor{
	Name:        "unmanaged-cluster",
	Aliases:     []string{"um", "uc", "unmanaged"},
	Description: description,
	Group:       plugin.RunCmdGroup,
	Version:     plugin.Version,
}

var (
	// logLevel is the verbosity to print
	logLevel int32

	// Log file to dump logs to
	logFile string
)

func main() {
	p, err := plugin.NewPlugin(&descriptor)
	if err != nil {
		log.Fatal(err, "unable to initialize new plugin")
	}

	p.Cmd.PersistentFlags().Int32VarP(&logLevel, "verbose", "v", 0, "Number for the log level verbosity(0-9)")
	p.Cmd.PersistentFlags().StringVar(&logFile, "log-file", "", "Log file path")

	p.AddCommands(
		cmd.ConfigureCmd,
		cmd.CreateCmd,
		cmd.DeleteCmd,
		cmd.ListCmd,
	)
	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}
