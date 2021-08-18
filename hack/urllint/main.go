package main

import (
	"flag"
	"log"
	"os"

	urllint "github.com/vmware-tanzu/tce/hack/urllint/pkg/lint"
)

func main() {
	// get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	var pathFlag = flag.String("path", wd, "path to be provided")                                               // default is current working directory
	var configPathFlag = flag.String("config", "config.json", "path for the configuration file to be provided") // default config is the config.json file that is there in the urllint path
	flag.Parse()
	llint, err := urllint.New(*configPathFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = llint.Init(*pathFlag)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("The following is the path that lint is working on: ", *pathFlag)
	llint.LintAll()
}
