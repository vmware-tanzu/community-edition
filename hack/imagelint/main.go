package main

import (
	"fmt"

	imglint "github.com/vmware-tanzu/tce/hack/imagelint/pkg/lint"
)

func main() {

	ilc, _ := imglint.New("config.json")
	err := ilc.Init("../../../")
	fmt.Println(err)
	for _, img := range ilc.ImageLints {
		fmt.Println(img.Path, ">>>>>><<<<<", img.Line)
	}
	fmt.Scanln()
}
