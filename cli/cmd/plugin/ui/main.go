// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin"
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

	bindAddress := "127.0.0.0:8080"
	browser := ""

	// Add our command line options
	p.Cmd.Flags().StringVarP(&bindAddress, "bind", "b", "127.0.0.1:8080", "Specify the IP and port to bind the Kickstart UI against (e.g. 127.0.0.1:8080).")
	p.Cmd.Flags().StringVarP(&browser, "browser", "", "", "Specify the browser to open the Kickstart UI on. Use 'none' for no browser. Defaults to OS default browser. Supported: ['chrome', 'firefox', 'safari', 'ie', 'edge', 'none']")
	p.Cmd.PersistentFlags().Int32VarP(&logLevel, "verbose", "v", 0, "Number for the log level verbosity (0-9)")

	p.Cmd.Run = func(cmd *cobra.Command, args []string) {
		launch(bindAddress, browser)
	}

	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}

func launch(bindAddress, browser string) {
	fmt.Printf("http://%s browser: %s\n", bindAddress, browser)

	http.HandleFunc("/", hello)
	http.HandleFunc("/headers", headers)

	err := http.ListenAndServe(bindAddress, nil)
	if err != nil {
		fmt.Printf("Error starting web server: %v", err)
	}
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}
