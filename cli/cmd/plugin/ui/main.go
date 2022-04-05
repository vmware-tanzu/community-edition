// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server"
)

//nolint:typecheck // When the lint job runs when the site hasn't been built it complains
//go:embed web/tanzu-ui/build
var content embed.FS

var descriptor = plugin.PluginDescriptor{
	Name:        "ui",
	Description: "Launch the Tanzu user interface.",
	Group:       plugin.ManageCmdGroup,
	Version:     plugin.Version,
}

var (
	// logLevel is the verbosity to print
	logLevel int32
)

func main() {
	p, err := plugin.NewPlugin(&descriptor)
	if err != nil {
		log.Fatal(err, "unable to initialize new plugin")
	}

	bindAddress := "0.0.0.0:8080"
	browser := ""

	// Add our command line options
	p.Cmd.Flags().StringVarP(&bindAddress, "bind", "b", bindAddress, "Specify the IP and port to on which to server the UI (e.g. 127.0.0.1:8080).")
	p.Cmd.Flags().StringVar(&browser, "browser", "", "Specify the browser to use to automatically launch the UI. Use 'none' for no browser. Defaults to OS default browser. Supported: ['chrome', 'firefox', 'safari', 'ie', 'edge', 'none']")
	p.Cmd.PersistentFlags().Int32VarP(&logLevel, "verbose", "v", 0, "Number for the log level verbosity (0-9)")

	p.Cmd.Run = func(cmd *cobra.Command, args []string) {
		launch(bindAddress, browser)
	}

	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}

func launch(bindAddress, browser string) {
	// get static content from go embed
	server.Content = content

	fmt.Printf("http://%s/ui/\n", bindAddress)
	err := server.Serve(bindAddress, browser)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
