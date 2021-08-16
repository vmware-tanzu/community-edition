package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

func main() {
	// step-1 : List out all files what may have docker images
	var extensions []string
	var ImageLints []ImageLint
	cimageLint := make(chan ImageLint, 0)
	extensions = append(extensions, ".sh", ".yaml", ".yml")

	go func() {
		for cil := range cimageLint {
			ImageLints = append(ImageLints, cil)
		}
	}()

	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if contains(extensions, filepath.Ext(path)) {
				fmt.Println(path)
				// TODO stuff here --?
				// Start Reading line by line
				go ReadFile(path, cimageLint)
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}
	for _, img := range ImageLints {

		fmt.Println(img.Path, ">>>>>><<<<<", img.Line)
	}
	fmt.Scanln()
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
			imageLint <- ImageLint{Path: path, Line: strings.Trim(s.Text(), " "), Position: Position{Row: count, Col: 0}}
		}
		count++
	}
	err = s.Err()
	if err != nil {
		log.Println(err)
	}
}
