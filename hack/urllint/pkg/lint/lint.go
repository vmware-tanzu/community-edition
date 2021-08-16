package lint

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
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

			if contains(llc.IncludeExts, filepath.Ext(path)) {
				// TODO stuff here --?
				// Start Reading line by line
				llc.ReadFile(path)
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
		//skip := false
		if link != "" {
			// for _, k := range llc.ExcludeLinks {
			// 	if strings.Contains(link, k) {
			// 		skip = true
			// 		break
			// 	}
			// }
			// if !skip {
			llc.LinkLints = append(llc.LinkLints, LinkLint{Path: path, Line: link, Position: Position{Row: count, Col: 0}})
			//	}
		}

	}
	count++
	err = s.Err()
	if err != nil {
		//fmt.Println(err)
		return err
	}
	return nil
}
func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}
