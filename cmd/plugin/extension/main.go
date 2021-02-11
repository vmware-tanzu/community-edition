// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	klog "k8s.io/klog/v2"

	"github.com/vmware-tanzu-private/core/pkg/v1/cli"
	"github.com/vmware-tanzu-private/core/pkg/v1/cli/command/plugin"

	extension "github.com/vmware-tanzu/tce/pkg/extension"
)

var descriptor = cli.PluginDescriptor{
	Name:        "extension",
	Description: "Extension management",
	Version:     "v0.2.0-pre-alpha-1", //TODO: Workaround for build breaking. Was: cli.BuildVersion,
	BuildSHA:    "",
	Group:       cli.ManageCmdGroup,
}

// var logLevel int32
// var logFile string

func main() {
	// plugin!
	p, err := plugin.NewPlugin(&descriptor)
	if err != nil {
		klog.Fatal(err)
	}

	// p.Cmd.PersistentFlags().Int32VarP(&logLevel, "v", "v", 0, "Number for the log level verbosity(0-9)")
	// p.Cmd.PersistentFlags().StringVar(&logFile, "log_file", "", "Log file path")

	p.AddCommands(
		extension.ListCmd,
		extension.GetCmd,
		extension.ReleaseCmd,
		extension.InstallCmd,
		extension.DeleteCmd,
		extension.TokenCmd,
		extension.ResetCmd,
	)
	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}
