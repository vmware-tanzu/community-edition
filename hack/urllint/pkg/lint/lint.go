package lint

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	xurls "mvdan.cc/xurls/v2"
)

type LinkLintConfig struct {
	IncludeExts []string `json:"includeExts"`
	// IncludeFiles []string  `json:"includeFiles"`
	// IncludeLines []string  `json:"includeLines"`
	ExcludeLinks []string `json:"excludeLinks"`
	// Validators   []string  `json:"validators"`
	LinkLints []LinkLint `json:"linkLint"`
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
		line := strings.Trim(s.Text(), " ")
		link := rxStrict.FindString(line)
		col := strings.Index(s.Text(), link)
		if link != "" {
			//fmt.Println(link)
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
