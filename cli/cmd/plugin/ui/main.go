// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/api"
)

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
	staticFiles := "web"

	// Add our command line options
	p.Cmd.Flags().StringVarP(&bindAddress, "bind", "b", bindAddress, "Specify the IP and port to bind the Kickstart UI against (e.g. 127.0.0.1:8080).")
	p.Cmd.Flags().StringVar(&browser, "browser", "", "Specify the browser to open the Kickstart UI on. Use 'none' for no browser. Defaults to OS default browser. Supported: ['chrome', 'firefox', 'safari', 'ie', 'edge', 'none']")
	p.Cmd.Flags().StringVarP(&staticFiles, "web-files", "f", staticFiles, "Specify the directory path for static HTML files to serve.")
	p.Cmd.PersistentFlags().Int32VarP(&logLevel, "verbose", "v", 0, "Number for the log level verbosity (0-9)")

	p.Cmd.Run = func(cmd *cobra.Command, args []string) {
		launch(bindAddress, browser, staticFiles)
	}

	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}

func launch(bindAddress, browser, staticFiles string) {
	workingDir, _ := os.Getwd()
	staticFiles, err := filepath.Abs(filepath.Join(workingDir, staticFiles))
	if err != nil {
		fmt.Printf("Error getting static directory path: %s\n", err.Error())
		os.Exit(1)
	}

	router := api.NewRouter()
	router.PathPrefix("/ui").Handler(http.StripPrefix("/ui", api.Logger(http.FileServer(http.Dir(staticFiles)), "ui")))

	if logLevel > 3 {
		if err := api.PrintRoutes(router); err != nil {
			fmt.Printf("Failed to print registered routes: %s\n", err.Error())
		}
	}

	fmt.Printf("Serving from %s\n", staticFiles)
	fmt.Printf("http://%s/ui/ browser: %s\n", bindAddress, browser)
	if err := http.ListenAndServe(bindAddress, router); err != nil {
		fmt.Printf("Error starting web server: %v", err)
	}
}
