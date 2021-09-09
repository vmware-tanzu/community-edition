// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	conformance "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/conformance/pkg"

	klog "k8s.io/klog/v2"

	cliv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/command/plugin"
)

var descriptor = cliv1alpha1.PluginDescriptor{
	Name:        "conformance",
	Description: "Run Sonobuoy conformance tests against clusters",
	Group:       cliv1alpha1.RunCmdGroup,
}

var logLevel int32

func main() {
	// plugin!
	p, err := plugin.NewPlugin(&descriptor)
	if err != nil {
		klog.Fatalf("%v", err)
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
		klog.Fatalf("%v", err)
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
