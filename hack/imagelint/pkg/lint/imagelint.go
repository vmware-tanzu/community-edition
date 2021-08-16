package lint

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ImageLintConfig struct {
	IncludeExts  []string    `json:"includeExts"`
	IncludeFiles []string    `json:"includeFiles"`
	IncludeLines []string    `json:"includeLines"`
	ExcludeFiles []string    `json:"excludeFiles"`
	ImageLints   []ImageLint `json:"imageLints"`
	chImageLint  chan ImageLint
}

func New(configFile string) (*ImageLintConfig, error) {
	file, err := ioutil.ReadFile(configFile)

	if err != nil {
		return nil, err
	}

	ilc := &ImageLintConfig{}

	err = json.Unmarshal([]byte(file), ilc)
	if err != nil {
		return nil, err
	}
	ilc.chImageLint = make(chan ImageLint)
	return ilc, nil
}

type ImageLint struct {
	Path     string
	Line     string
	Position Position
}
type Position struct {
	Row, Col int
}

func (imc *ImageLintConfig) Init(dir string) error {
	//imc.chImageLint = make(chan ImageLint)
	// todo is chImageLint is nil
	go func() {
		for cil := range imc.chImageLint {
			imc.ImageLints = append(imc.ImageLints, cil)
		}
	}()
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if contains(imc.IncludeExts, filepath.Ext(path)) {
				//fmt.Println(path)
				// TODO stuff here --?
				// Start Reading line by line
				go imc.ReadFile(path)
			}
			return nil
		})
	if err != nil {
		return err
	}
	return err
}

func (imc *ImageLintConfig) ReadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	count := 1
	for s.Scan() {
		// ignore lines
		// if the line is commeted then skip it
		line := strings.Trim(s.Text(), " ")
		if len(line) > 2 && (line[:2] == "//" || line[:2] == "/*") {
			continue
		}
		for _, searchterm := range imc.IncludeLines {
			if strings.Contains(line, searchterm) {
				imc.chImageLint <- ImageLint{Path: path, Line: strings.Trim(s.Text(), " "), Position: Position{Row: count, Col: 0}}
			}
		}
		count++
	}
	err = s.Err()
	if err != nil {
		fmt.Println(err)
		//return err
	}
	return nil
}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}
