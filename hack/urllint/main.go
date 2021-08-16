package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"

	urllint "github.com/vmware-tanzu/tce/hack/urllint/pkg/lint"
)

func main() {
	var pathFlag = flag.String("path", "../../../", "path to lint")
	flag.Parse()
	llint, _ := urllint.New("config.json")
	err := llint.Init(*pathFlag)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("User given Path", *pathFlag)

	for _, link := range llint.LinkLints {
		if !IsUrl(link.Line) {
			log.Fatalf("file:%s line:%d:%d Link:%s has error", link.Path, link.Position.Row, link.Position.Col, link.Line)
			//log.Printf("file:%s line:%d:%d Link:%s has error", link.Path, link.Position.Row, link.Position.Col, link.Line)
		}
		resp, err := http.Get(link.Line)
		if err != nil {
			log.Fatalf("file:%s line:%d:%d Link:%s has error", link.Path, link.Position.Row, link.Position.Col, link.Line)
			//log.Printf("file:%s line:%d:%d Link:%s has error", link.Path, link.Position.Row, link.Position.Col, link.Line)
		}
		if resp.StatusCode >= 300 {
			log.Fatalf("file:%s line:%d:%d Link:%s has error", link.Path, link.Position.Row, link.Position.Col, link.Line)
			//log.Printf("file:%s line:%d:%d Link:%s has error", link.Path, link.Position.Row, link.Position.Col, link.Line)
		}
		//fmt.Println(link.Path, ":", link.Line)
	}
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
