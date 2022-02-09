package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

//func CreateRepo(packagesDir, outputPath, name string) {
// Create repo FS
// ReadPackages
// WritePackages
// Write lock?
//}

func ReadPackages(packagesDir string) []byte {
	packages, err := os.ReadDir(packagesDir)
	if err != nil {
		panic(err)
	}

	realPackages := filterPackageDirs(packages)

	// Make a slice of strings to hold all the YAML we'll pull out of packages.
	// Assume initial capacity of 2 versions per package
	allpkgs := make([][]byte, 0, len(packages)*2)
	for _, p := range realPackages {
		// TODO: Make this a log statement instead
		fmt.Println("Processing package", p.Name())
		pkgpath := filepath.Join(packagesDir, p.Name())
		versions, err := os.ReadDir(pkgpath)
		if err != nil {
			panic(err)
		}
		for _, v := range versions {
			if !isVersion(v) {
				break
			}
			fmt.Println("  Processing version", v.Name())
			f := filepath.Join(packagesDir, p.Name(), v.Name(), "package.yaml")
			bin, err := os.ReadFile(f)
			if err != nil {
				panic(err)
			}
			// Normalize entries to remove whitespace at the end and beginning
			bin = bytes.TrimSpace(bin)

			allpkgs = append(allpkgs, bin)
		}
	}

	// Join YAML documents with "\n---\n"
	return bytes.Join(allpkgs, []byte{'\n', '-', '-', '-', '\n'})
}

// isVersion checks to see if an os.DirEntry is a directory conforming to version expectations.
func isVersion(e os.DirEntry) bool {
	ok := true
	if !e.IsDir() {
		ok = false
	}
	if e.Name() == "hack" || e.Name() == "pkg" {
		ok = false
	}
	return ok
}

// filterPackageDirs filters a list of os.DirEntrys for only those that are directories.
func filterPackageDirs(entries []os.DirEntry) []os.DirEntry {
	new := make([]os.DirEntry, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			new = append(new, e)
		}
	}
	return new
}

//func GenerateRepoFilesystem(outPath string) {
//}

func main() {
	// Inputs:
	// directory where packages live
	// directory where lock files will go
	// name/channel
	// bonus: input file?

	allpkgs := ReadPackages(os.Args[1])

	// 1: bring all the YAML files together into 1 file.
	// 2: Use kbld lock to get image lock files.
	os.WriteFile("packages.yaml", allpkgs, 0755)
}
