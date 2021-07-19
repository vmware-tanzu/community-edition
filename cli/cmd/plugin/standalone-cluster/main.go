// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"math/rand"
	"os"
	"time"

	klog "k8s.io/klog/v2"

	cliv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/command/plugin"
)

var descriptor = cliv1alpha1.PluginDescriptor{
	Name:        "standalone-cluster",
	Description: "Create clusters without a dedicated management cluster",
	Group:       cliv1alpha1.RunCmdGroup,
}

var (
	// BuildEdition is the edition the CLI was built for.
	BuildEdition string

	// logLevel is the verbosity to print
	logLevel int32
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	// plugin!
	p, err := plugin.NewPlugin(&descriptor)
	if err != nil {
		klog.Fatalf("%v", err)
	}

	p.Cmd.PersistentFlags().Int32VarP(&logLevel, "verbose", "v", 0, "Number for the log level verbosity(0-9)")

	p.AddCommands(
		CreateCmd,
		DeleteCmd,
	)
	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}
