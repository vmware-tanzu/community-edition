// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/vmware-tanzu-private/core/pkg/v1/cli"
	"github.com/vmware-tanzu-private/core/pkg/v1/cli/command/plugin"

	extension "github.com/vmware-tanzu/tce/pkg/extension"
)

var descriptor = cli.PluginDescriptor{
	Name:        "extension",
	Description: "Extension management",
	Version:     "v0.1.0",
	Group:       cli.ManageCmdGroup,
}

func main() {
	extensionPlugin, err := plugin.NewPlugin(&descriptor)
	if err != nil {
		fmt.Print(err)
	}
	extensionPlugin.AddCommands(
		extension.ListCmd,
		extension.GetCmd,
		extension.ReleaseCmd,
		extension.InstallCmd,
		extension.DeleteCmd,
		extension.TokenCmd,
		extension.ResetCmd,
	)
	if err := extensionPlugin.Execute(); err != nil {
		os.Exit(1)
	}
}
