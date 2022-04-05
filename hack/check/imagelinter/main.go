// Copyright 2021-2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/rs/xid"

	imgwrapper "github.com/vmware-tanzu/community-edition/hack/imagelinter/pkg/imagewrapper"
	imglint "github.com/vmware-tanzu/community-edition/hack/imagelinter/pkg/lint"
)

var (
	//go:embed config/imagelintconfig.yaml
	data                                      string
	pathFlag, configPathFlag, detailedSummary *string
	showSummary                               *bool
	imc                                       *imglint.ImageLintConfig
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	pathFlag = flag.String("path", wd, "path to be provided")                                    // default is current working directory
	configPathFlag = flag.String("config", "", "path for the configuration file to be provided") // default config is the config.json file that is there in the imagelint path
	showSummary = flag.Bool("summary", false, "to get summary pass summary=true; to off either dont pass or summary=false")
	detailedSummary = flag.String("details", "Fail", "detailed summary can be Fail, Pass, Not-Identified, or Pull-Failed")
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
}

type CheckResult struct {
	Status  string
	Message string
	Image   string
	Error   error
}

func main() {
	imageCount := len(imc.ImageMap)
	results := make(chan *CheckResult, 2)
	wg := sync.WaitGroup{}
	guard := make(chan struct{}, 8)

	isFatal := false
	count := 0
	go func() {
		// Loop through our results as they come in and print the results
		for res := range results {
			wg.Done()
			<-guard
			count++
			fmt.Printf("Result %d of %d images\n", count, imageCount)

			if res.Error != nil {
				isFatal = true
				imc.OnEvent(imglint.Fail, res.Error.Error(), res.Image)
				continue
			}

			imc.OnEvent(res.Status, res.Message, res.Image)
		}
	}()

	for key := range imc.ImageMap {
		wg.Add(1)
		guard <- struct{}{}

		go func(key string) {
			imageName := strings.Trim(key, " ")
			wrapper, err := imgwrapper.New(imageName, xid.New().String(), nil)
			if err != nil {
				results <- &CheckResult{Image: imageName, Error: err}
				return
			}

			results <- LintAll(imc, wrapper, imageName)
		}(key)
	}
	wg.Wait()
	close(results)

	imc.ShowDetailedSummary(*detailedSummary)
	if *showSummary {
		imc.ShowOverallSummary()
		fmt.Println()
	}

	if isFatal {
		os.Exit(1)
	}
}

func LintAll(imc *imglint.ImageLintConfig, wrapper *imgwrapper.Wrapper, key string) *CheckResult {
	_, err := wrapper.PullImage()
	if err != nil {
		// Unable to pull image - invalid?
		return &CheckResult{
			Message: fmt.Sprintf("unable to pull image: %s", err.Error()),
			Image:   key,
			Status:  imglint.PullFail,
		}
	}

	// Check the docker history
	skip, err := wrapper.Validate(imc.SuccessValidators)
	if skip && err == nil {
		return &CheckResult{
			Image:   key,
			Status:  imglint.Pass,
			Message: "According to the image history, this is not an Alpine Image",
		}
	}

	_, err = wrapper.CreateContainer()
	if err != nil {
		// May be a local problem, but we can't tell
		return &CheckResult{
			Image:   key,
			Status:  imglint.NotIdentified,
			Message: fmt.Sprintf("unable to create container from image: %s", err.Error()),
		}
	}
	defer wrapper.DeleteContainer() //nolint: errcheck

	// Check if it's alpine based
	err = IsAlpine(wrapper, key)
	if err != nil {
		return &CheckResult{
			Image:  key,
			Status: imglint.Fail,
			Error:  err,
		}
	}

	// See if it declares an unacceptable license
	err = CheckLicense(wrapper, key)
	if err != nil {
		return &CheckResult{
			Image:  key,
			Status: imglint.Fail,
			Error:  err,
		}
	}

	// Check if it uses busybox
	err = IsBusyBox(wrapper)
	if err != nil {
		return &CheckResult{
			Image:  key,
			Status: imglint.Fail,
			Error:  err,
		}
	}

	return &CheckResult{
		Image:   key,
		Status:  imglint.Pass,
		Message: "Not an Alpine image",
	}
}

// IsAlpine checks if the container appears to be alpine based. Returns error if so,
// otherwise returns nil if either it is not alpine or if it could not be determined.
func IsAlpine(wrapper *imgwrapper.Wrapper, key string) error {
	// Find the os-release file in the container
	localName := fmt.Sprintf("./os-release-%s", key)
	_, err := wrapper.ContainerCP("/etc/os-release", localName)
	if err != nil {
		_, _ = wrapper.ContainerCP("/usr/lib/os-release", localName)
	}

	osdata, err := os.ReadFile(localName)
	if err != nil {
		// Unable to read the os-release file, punt
		return nil
	}
	defer os.Remove(localName)

	if strings.Contains(string(osdata), "Alpine") {
		return errors.New("alpine image")
	}
	return nil
}

// CheckLicense checks if the container declares an incompatible license. If so it returns an error, otherwise
// if the license is fine or if it could not be determined it will return nil.
func CheckLicense(wrapper *imgwrapper.Wrapper, key string) error {
	localName := fmt.Sprintf("./LICENSE-%s", key)
	_, err := wrapper.ContainerCP("/licenses/LICENSE", localName)
	if err != nil {
		// No license file to check
		return nil
	}

	license, err := os.ReadFile(localName)
	defer os.Remove(localName)
	if err != nil {
		// Unable to read the LICENSE file, punt
		return nil
	}

	// TODO: Add checks for other known incompatible licenses
	if strings.Contains(string(license), "Alpine") {
		return errors.New("alpine license")
	}
	return nil
}

// IsBusyBox checks if the container image is using busybox. If so, returns an
// error, else if it is not or can't be determined returns nil.
func IsBusyBox(wrapper *imgwrapper.Wrapper) error { //nolint:unparam
	_, err := wrapper.RunCommand("exec", wrapper.Container, "/bin/busybox")
	if err != nil {
		// Successfully ran busybox
		return nil
	}

	// TODO: What the heck are we trying to do here? The original code was checking
	// if it was busybox, and if so just reporting success that it wasn't alpine.
	// Do we care about busybox? What does that have to do with alpine?
	return nil
}
