// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/vmware-tanzu/tce/cli/pkg/standalone-cluster"

	klog "k8s.io/klog/v2"

	"github.com/vmware-tanzu-private/core/pkg/v1/cli"
	"github.com/vmware-tanzu-private/core/pkg/v1/cli/command/plugin"
)

var descriptor = cli.PluginDescriptor{
	Name:        "standalone-cluster",
	Description: "Create clusters without a dedicated management cluster",
	Version:     cli.BuildVersion,
	BuildSHA:    "",
	Group:       cli.RunCmdGroup,
}

var logLevel int32

func main() {
	// plugin!
	p, err := plugin.NewPlugin(&descriptor)
	if err != nil {
		klog.Fatalf("%v", err)
	}

	p.Cmd.PersistentFlags().Int32VarP(&logLevel, "verbose", "v", 0, "Number for the log level verbosity(0-9)")

	p.AddCommands(
		standalone.CreateCmd,
	)
	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}
