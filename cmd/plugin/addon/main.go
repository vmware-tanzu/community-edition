// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	klog "k8s.io/klog/v2"

	"github.com/vmware-tanzu-private/core/pkg/v1/cli"
	"github.com/vmware-tanzu-private/core/pkg/v1/cli/command/plugin"

	addon "github.com/vmware-tanzu/tce/pkg/addon"
)

var descriptor = cli.PluginDescriptor{
	Name:        "package",
	Description: "Package management",
	Version:     cli.BuildVersion,
	BuildSHA:    "",
	Group:       cli.ManageCmdGroup,
}

func main() {

	// plugin!
	p, err := plugin.NewPlugin(&descriptor)
	if err != nil {
		klog.Fatalf("%v", err)
	}

	// p.Cmd.PersistentFlags().Int32VarP(&logLevel, "v", "v", 0, "Number for the log level verbosity(0-9)")
	// p.Cmd.PersistentFlags().StringVar(&logFile, "log_file", "", "Log file path")

	p.AddCommands(
		addon.ListCmd,
		addon.ConfigureCmd,
		addon.InstallCmd,
		addon.DeleteCmd,
		addon.RepositoryCmd,
	)
	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}
