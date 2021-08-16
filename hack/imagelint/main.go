package main

import (
	"flag"
	"fmt"

	imglint "github.com/vmware-tanzu/tce/hack/imagelint/pkg/lint"
)

func main() {
	var pathFlag = flag.String("path", "../../../", "path to lint")
	ilc, _ := imglint.New("config.json")
	err := ilc.Init(*pathFlag)
	if err != nil {
		fmt.Println(err)
	}
	for _, img := range ilc.ImageLints {
		//fmt.Println(img.Path, ">>>>>><<<<<", img.Line)
		fmt.Println(img.Line)
	}
	fmt.Scanln()
}
