// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	cliv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/command/plugin"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
)

var descriptor = cliv1alpha1.PluginDescriptor{
	Name:        "standalone-cluster",
	Description: `(!!! deprecated - see unmanaged-cluster !!!) Create clusters without a dedicated management cluster`,
	Group:       cliv1alpha1.RunCmdGroup,

	// Since standalone cluster is being deprecated, no changes or feature adds
	// are to be made to it's dependencies or core functionality.
	// Since cli.BuildVersion was deprecated in the build version of tanzu framework v0.10.0,
	// we need to manually set the plugin version in order to compile with hack/builder
	Version: "v0.10.0",
}

var (
	// BuildEdition is the edition the CLI was built for.
	BuildEdition string

	// logLevel is the verbosity to print
	logLevel int32

	// Log file to dump logs to
	logFile string
)

func main() {
	// plugin!
	p, err := plugin.NewPlugin(&descriptor)
	if err != nil {
		log.Fatal(err, "unable to initilize new plugin")
	}

	p.Cmd.PersistentFlags().Int32VarP(&logLevel, "verbose", "v", 0, "Number for the log level verbosity(0-9)")
	p.Cmd.PersistentFlags().StringVar(&logFile, "log-file", "", "Log file path")

	p.AddCommands(
		CreateCmd,
		DeleteCmd,
	)
	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}
