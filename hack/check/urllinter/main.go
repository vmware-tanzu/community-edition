// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"

	linturl "github.com/vmware-tanzu/community-edition/hack/urllinter/pkg/lint"
)

var (
	//go:embed config/urllintconfig.yaml
	data string
)

func main() {
	// get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	var pathFlag = flag.String("path", wd, "path to be provided")                                    // default is current working directory
	var configPathFlag = flag.String("config", "", "path for the configuration file to be provided") // default config is the config.json file that is there in the urllint path
	var showSumary = flag.Bool("summary", false, "to get summary pass summary=true;to off either dont pass or summary=false")
	var detailedSummary = flag.String("details", "Fail", "detailed summary can be Fail,Pass")
	flag.Parse()
	var llint *linturl.LinkLintConfig
	if *configPathFlag == "" {
		llint, err = linturl.NewFromContent([]byte(data))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		llint, err = linturl.New(*configPathFlag)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = llint.Init(*pathFlag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The following is the path that lint is working on: ", *pathFlag)

	isFatal := llint.LintAll()

	if *showSumary {
		llint.ShowSummary()
		fmt.Println()
	}
	switch *detailedSummary {
	case "Fail", "fail", "FAIL":
		llint.ShowFailSummary()
	case "Pass", "pass", "PASS":
		llint.ShowPassSummary()
	}

	if isFatal {
		os.Exit(1)
	}
}
