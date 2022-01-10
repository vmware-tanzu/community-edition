// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rs/xid"

	imgwrapper "github.com/vmware-tanzu/community-edition/hack/imagelinter/pkg/imagewrapper"
	imglint "github.com/vmware-tanzu/community-edition/hack/imagelinter/pkg/lint"
)

var (
	//go:embed config/imagelintconfig.yaml
	data                                      string
	pathFlag, configPathFlag, detailedSummary *string
	showSumary                                *bool
	counter                                   int
	isFatal                                   bool
	imc                                       *imglint.ImageLintConfig
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	pathFlag = flag.String("path", wd, "path to be provided")                                    // default is current working directory
	configPathFlag = flag.String("config", "", "path for the configuration file to be provided") // default config is the config.json file that is there in the imagelint path
	showSumary = flag.Bool("summary", false, "to get summary pass summary=true;to off either dont pass or summary=false")
	detailedSummary = flag.String("details", "Fail", "detailed summary can be Fail,Pass,Not-Identified or Pull-Failed")
	flag.Parse()
	if *configPathFlag == "" {
		imc, err = imglint.NewFromContent([]byte(data))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		imc, err = imglint.New(*configPathFlag)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = imc.Init(*pathFlag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User given Path", *pathFlag)
	fmt.Println("Total number of images to process:", len(imc.ImageMap))
	counter = 0
	isFatal = false
}
func main() {
	for key := range imc.ImageMap {
		counter++
		fmt.Println("Currently processing", counter, " out of ", len(imc.ImageMap), "image(s)")
		containerName := xid.New().String()
		wrapper, err := imgwrapper.New(strings.Trim(key, " "), containerName, nil)
		if err != nil {
			log.Fatalln(err)
		}
		cont, err := LintAll(imc, wrapper, key)
		if err != nil {
			fmt.Println(err)
		}
		if cont {
			continue
		}
		imc.OnEvent("Not Identified", "Cound not find Container OS", key)
		_, err = wrapper.DeleteContainer()
		if err != nil {
			fmt.Println("Error deleting container:", err)
		}
	}

	imc.ShowDetailedSummary(*detailedSummary)
	if *showSumary {
		imc.ShowOverallSummary()
		fmt.Println()
	}
	if isFatal {
		os.Exit(1)
	}
}

func IsAlpine(imc *imglint.ImageLintConfig, wrapper *imgwrapper.Wrapper, key string) (fatal, cont bool, err error) {
	osdata, err := os.ReadFile("os-release")
	if err == nil {
		err = os.Remove("os-release")
		if err != nil {
			return false, false, err
		}
		if strings.Contains(string(osdata), "Alpine") {
			imc.OnEvent("Fail", "Alpine Image", key)
			_, err = wrapper.DeleteContainer()
			if err != nil {
				return true, false, err
			}
			return true, true, nil
		}
		imc.OnEvent("Pass", "Not an Alpine image", key)
		_, err = wrapper.DeleteContainer()
		if err != nil {
			fmt.Println("Error deleting container:", err)
		}
		return false, true, nil
	}
	return false, false, nil
}

func HasLicense(imc *imglint.ImageLintConfig, wrapper *imgwrapper.Wrapper, key string) (cont bool, err error) {
	data, err := os.ReadFile("LICENSE")
	if err == nil && len(data) > 0 {
		err = os.Remove("LICENSE")
		if err != nil {
			_, err = wrapper.DeleteContainer()
			if err != nil {
				return false, err
			}
			return false, err
		}
		if strings.Contains(string(data), "Apache License") {
			imc.OnEvent("Pass", "Valid license file found", key)
			_, err = wrapper.DeleteContainer()
			if err != nil {
				return true, err
			}
			return true, nil
		}
	}
	_, err = wrapper.DeleteContainer()
	if err != nil {
		return false, err
	}
	return false, nil
}

func IsBusyBox(imc *imglint.ImageLintConfig, wrapper *imgwrapper.Wrapper, key string) (cont bool, err error) {
	result, err := wrapper.RunCommand("logs", wrapper.Container+"2")
	if err != nil {
		return false, err
	}
	_, err = wrapper.RunCommand("rm", "-f", wrapper.Container+"2")
	if err != nil {
		return false, err
	}
	if strings.Contains(strings.ToUpper(result), "BUSYBOX") && err == nil {
		imc.OnEvent("Pass", "According to the containr bins, this is not an Alpine Image;Its busybox", key)
		return true, nil
	}
	return false, err
}

func InitialChecks(imc *imglint.ImageLintConfig, wrapper *imgwrapper.Wrapper, key string) (cont bool, err error) {
	_, err = wrapper.PullImage()
	if err != nil {
		imc.OnEvent("Pull Failed", "Pulling Image Failed:"+err.Error(), key)
		return true, err
	}
	skip, err := wrapper.Validate(imc.SuccessValidators)
	if skip && err == nil {
		imc.OnEvent("Pass", "According to the image history, this is not an Alpine Image", key)
		return true, nil
	}
	return false, nil
}

func LintAll(imc *imglint.ImageLintConfig, wrapper *imgwrapper.Wrapper, key string) (bool, error) {
	cont, err := InitialChecks(imc, wrapper, key)
	if cont {
		return true, err
	}
	_, err = wrapper.CreateContainer()
	if err != nil {
		return false, err
	}
	if wrapper.IsContainerExists() {
		_, err := wrapper.ContainerCP("/etc/os-release", "./")
		if err == nil {
			fatal, cont, err := IsAlpine(imc, wrapper, key)
			if fatal {
				isFatal = fatal
			}
			if cont {
				return true, err
			}
		}
		_, err = wrapper.ContainerCP("/usr/lib/os-release", "./")
		if err == nil {
			fatal, cont, err := IsAlpine(imc, wrapper, key)
			if fatal {
				isFatal = fatal
			}
			if cont {
				return true, err
			}
		}
		_, err = wrapper.ContainerCP("/licenses/LICENSE", "./")
		if err == nil {
			cont, err := HasLicense(imc, wrapper, key)
			if cont {
				return true, err
			}
		}
		_, err = wrapper.RunCommand("run", `--entrypoint=/bin/busybox`, "--name", wrapper.Container+"a", key)
		if err == nil {
			cont, err := IsBusyBox(imc, wrapper, key)
			if cont {
				return true, err
			}
		}
	}
	return false, nil
}
