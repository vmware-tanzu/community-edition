// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	yaml "github.com/ghodss/yaml"
	klog "k8s.io/klog/v2"

	"github.com/vmware-tanzu-private/core/pkg/v1/cli"
	"github.com/vmware-tanzu-private/core/pkg/v1/cli/command/plugin"

	cfg "github.com/vmware-tanzu/tce/pkg/common/config"
	extension "github.com/vmware-tanzu/tce/pkg/extension"
)

var descriptor = cli.PluginDescriptor{
	Name:        "extension",
	Description: "Extension management",
	Version:     cli.BuildVersion,
	BuildSHA:    "",
	Group:       cli.ManageCmdGroup,
}

// var logLevel int32
// var logFile string

func main() {

	// checks for config file existence
	configFile := filepath.Join(xdg.DataHome, "tanzu-repository", cfg.DefaultConfigFile)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		klog.Error("Config file does not exist")

		c := cfg.NewConfig()
		c.ReleaseVersion = cli.BuildVersion

		fileWrite, err := os.OpenFile(configFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			klog.Fatalf("Open Config for write failed. Err: %v", err)
		}

		datawriter := bufio.NewWriter(fileWrite)
		if datawriter == nil {
			klog.Fatalf("Datawriter creation failed")
		}

		byRaw, err := yaml.Marshal(c)
		if err != nil {
			klog.Fatalf("yaml.Marshal error. Err: %v", err)
		}
		klog.V(6).Infof("byRaw = %v", byRaw)

		_, err = datawriter.Write(byRaw)
		if err != nil {
			klog.Fatalf("datawriter.Write error. Err: %v", err)
		}
		datawriter.Flush()

		err = fileWrite.Close()
		if err != nil {
			klog.Fatalf("fileWrite.Close failed. Err: %v", err)
		}
	}

	// plugin!
	p, err := plugin.NewPlugin(&descriptor)
	if err != nil {
		klog.Fatalf("%v", err)
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
