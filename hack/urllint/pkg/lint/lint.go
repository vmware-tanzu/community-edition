package lint

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	xurls "mvdan.cc/xurls/v2"
)

type LinkLintConfig struct {
	IncludeExts  []string   `json:"includeExts"`
	ExcludeLinks []string   `json:"excludeLinks"`
	LinkLints    []LinkLint `json:"linkLint"`
}

type LinkLint struct {
	Path     string
	Line     string
	Position Position
}
type Position struct {
	Row, Col int
}

func New(configFile string) (*LinkLintConfig, error) {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	llc := &LinkLintConfig{}
	err = json.Unmarshal([]byte(file), llc)
	if err != nil {
		return nil, err
	}
	return llc, nil
}

func (llc *LinkLintConfig) Init(dir string) error {
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			ext := filepath.Ext(path)
			for _, ex := range llc.IncludeExts {
				if ext == ex {
					llc.ReadFile(path)
				}
			}
			return nil
		})
	if err != nil {
		return err
	}
	return err
}

func (llc *LinkLintConfig) ReadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	count := 1
	rxStrict := xurls.Strict()
	for s.Scan() {
		skip := false
		line := strings.Trim(s.Text(), " ")
		link := rxStrict.FindString(line)
		col := strings.Index(s.Text(), link)
		for _, l := range llc.ExcludeLinks {
			if strings.Contains(link, l) {
				skip = true
				break
			}
		}
		if link != "" && !skip {
			llc.LinkLints = append(llc.LinkLints, LinkLint{Path: path, Line: link, Position: Position{Row: count, Col: col}})
		}
		count++
	}
	err = s.Err()
	if err != nil {
		return err
	}
	return nil
}

func (llc *LinkLintConfig) LintAll() {
	for _, link := range llc.LinkLints {
		if !IsUrl(link.Line) {
			log.Fatalf("file:%s line:%d:%d Link:%s is invalid link", link.Path, link.Position.Row, link.Position.Col, link.Line)
			//log.Printf("file:%s line:%d:%d Link:%s has error", link.Path, link.Position.Row, link.Position.Col, link.Line)
		}
		resp, err := http.Get(link.Line)
		if err != nil {
			log.Fatalf("file:%s line:%d:%d Link:%s is not working.", link.Path, link.Position.Row, link.Position.Col, link.Line)
			//log.Printf("file:%s line:%d:%d Link:%s has error", link.Path, link.Position.Row, link.Position.Col, link.Line)
		}
		if resp.StatusCode >= 300 {
			log.Fatalf("file:%s line:%d:%d Link:%s returns status code %d", link.Path, link.Position.Row, link.Position.Col, link.Line, resp.StatusCode)
			//log.Printf("file:%s line:%d:%d Link:%s has error", link.Path, link.Position.Row, link.Position.Col, link.Line)
		}
		//fmt.Println(link.Path, ":", link.Line)
	}
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
