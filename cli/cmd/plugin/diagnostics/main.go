// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"embed"
	"log"
	"os"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/diagnostics/pkg"
)

var pluginDesc = plugin.PluginDescriptor{
	Name:        "diagnostics",
	Description: "Cluster diagnostics",
	Group:       plugin.RunCmdGroup,
	Aliases:     []string{"diag", "diags", "diagnostic"},
	Version:     plugin.Version,
}

var (
	//go:embed scripts
	scriptFS embed.FS
)

func main() {
	p, err := plugin.NewPlugin(&pluginDesc)
	if err != nil {
		log.Fatal(err)
	}

	p.AddCommands(pkg.CollectCmd(scriptFS))
	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}
