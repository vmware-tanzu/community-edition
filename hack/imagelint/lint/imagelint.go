package lint

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strings"
)

type ImageLint struct {
	Path     string
	Line     string
	Position Position
}
type Position struct {
	Row, Col int
}

var (
	includeExt  []string
	ignoreFiles []string
	chImageLint chan ImageLint
	ImageLints  []ImageLint
)

func init() {
	chImageLint := make(chan ImageLint, 0)
	go func() {
		for cil := range chImageLint {
			ImageLints = append(ImageLints, cil)
		}
	}()

}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}

func ReadFile(path string, imageLint chan ImageLint) {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Println(err)
		}
	}()
	s := bufio.NewScanner(f)
	count := 1
	for s.Scan() {
		// ignore lines
		// if the line is commeted then skip it

		line := strings.Trim(s.Text(), " ")
		if len(line) > 2 && (line[:2] == "//" || line[:2] == "/*") {
			continue
		}
		if strings.Contains(line, "image:") {
			chImageLint <- ImageLint{Path: path, Line: strings.Trim(s.Text(), " "), Position: Position{Row: count, Col: 0}}
		}
		count++
	}
	err = s.Err()
	if err != nil {
		log.Println(err)
	}
}
