// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"strconv"

	"github.com/spf13/pflag"

	cliv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/command/plugin"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
)

var descriptor = cliv1alpha1.PluginDescriptor{
	Name:        "local",
	Description: "Manage local environments of Tanzu",
	Group:       cliv1alpha1.RunCmdGroup,
}

var (
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

	// TODO(joshrosso): must check if docker daemon is accessible.

	p.AddCommands(
		ConfigureCmd,
		CreateCmd,
		DeleteCmd,
		ListCmd,
	)
	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}

// TtySetting gets the setting to use for formatted TTY output based on whether
// the user explicitly set it with a command line argument, or if not, whether
// there is an environment variable set. If neither of these things, it will
// default to whether or not we detect we are running in a terminal that allows
// tty formatting.
func TtySetting(flags *pflag.FlagSet) bool {
	// See if we are running in a tty enabled terminal
	fileInfo, _ := os.Stdout.Stat()
	result := (fileInfo.Mode() & os.ModeCharDevice) != 0

	if flags.Changed("tty") {
		// User has explicitly set the flag, use that value
		result, _ = flags.GetBool("tty")
	} else if tty := os.Getenv("TANZU_TTY"); tty != "" {
		// Not explicitly provided, but there is an env setting
		val, err := strconv.ParseBool(tty)
		if err == nil {
			result = val
		}
	}
	return result
}
