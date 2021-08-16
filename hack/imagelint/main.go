package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/vmware-tanzu/tce/hack/imagelint/pkg/cmdhelper"
	imglint "github.com/vmware-tanzu/tce/hack/imagelint/pkg/lint"
)

func main() {
	var pathFlag = flag.String("path", "../../../", "path to lint")
	flag.Parse()
	ilc, _ := imglint.New("config.json")
	err := ilc.Init(*pathFlag)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("User given Path", *pathFlag)
	chelper := &cmdhelper.CmdHelper{Writer: nil}
	for _, img := range ilc.ImageLints {
		log.Println("Image:", strings.Trim(img.Line, " "))
		rnum := rand.Intn(len(ilc.ImageLints))
		// Step-1
		// docker pull

		_, err := chelper.CliRunner("docker", nil, []string{"pull", strings.Trim(img.Line, " ")}...)
		if err != nil {
			log.Println(err)
			// todo some stuff here
		}
		// Step-2
		// Get image history
		history, err := chelper.CliRunner("docker", nil, []string{"history", strings.Trim(img.Line, " "), "--no-trunc"}...)
		if err != nil {
			log.Println(err)
			// todo some stuff here
		}

		// if History has ValidateWords like apt-get etc mentioned in the config
		// then it is sure that the image is not alpine
		skip := false
		for _, item := range ilc.Validators {
			if strings.Contains(history, item) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		// Step-3
		// Create a container
		// Only create container if the above cases fail
		result, err := chelper.CliRunner("docker", nil, []string{"run", "--name", "c" + strconv.Itoa(rnum), strings.Trim(img.Line, " "), "cat", "/etc/os-release"}...)

		if err != nil {
			log.Println(err)

		}
		//fmt.Println(result)
		if strings.Contains(result, "Alpine") {
			fmt.Printf("file:%s line:%d:%d image:%s is Alpine image", img.Path, img.Position.Row, img.Position.Col, img.Line)
		}
	}
}
