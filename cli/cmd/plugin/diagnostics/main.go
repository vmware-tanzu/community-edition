// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"embed"
	"log"
	"os"

	cliv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/cmd/cli/plugin/diagnostics/pkg"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/command/plugin"
)

var pluginDesc = cliv1alpha1.PluginDescriptor{
	Name:        "diagnostics",
	Description: "Cluster diagnostics",
	Group:       cliv1alpha1.RunCmdGroup,
	Aliases:     []string{"diag", "diags", "diagnostic"},
	Version:     cli.BuildVersion,
}

var (
	//go:embed scripts
	scriptFS embed.FS
	defaultVersion = "v0.0.1-unversioned"
)

func main() {
	if pluginDesc.Version == ""{
		pluginDesc.Version = defaultVersion
	}
	p, err := plugin.NewPlugin(&pluginDesc)
	if err != nil {
		log.Fatal(err)
	}

	p.AddCommands(pkg.CollectCmd(scriptFS))
	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}
