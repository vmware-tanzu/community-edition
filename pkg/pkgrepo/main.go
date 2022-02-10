package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

//TODO: make this an interface with the fields as methods, as well as Name
type Repo struct {
	Root, ImgpkgPath, PackagesPath string
}

// CreateRepo creates a repository on the local filesystem based on local package definitions.
// packagesDir holds the package files that should be processed.
// outputPath defines the root output path.
// name defines the directory representing a specific repository (such as `core` and `user`).
func CreateRepo(packagesDir, outputPath, name string) error {
	r, err := GenerateRepoFileSystem(outputPath, name)
	if err != nil {
		return err
	}
	allpkgs := ReadPackages(packagesDir)
	os.WriteFile(filepath.Join(r.PackagesPath, "packages.yaml"), allpkgs, 0755)
	// WritePackages
	// Write lock?
	return nil
}

// GenerateRepoFileSystem creates the proper repository file system structure at a given path.
// outputPath is the level at which new repositories will be created.
// name is the child of outputPath that will be populated
func GenerateRepoFileSystem(outputPath, name string) (Repo, error) {
	fullPath := filepath.Join(outputPath, name)
	imagesPath := filepath.Join(fullPath, "images")
	pkgsPath := filepath.Join(fullPath, "pkgs")

	for _, d := range []string{fullPath, imagesPath, pkgsPath} {
		if err := os.MkdirAll(d, 0755); err != nil {
			return Repo{}, err
		}
	}
	return Repo{Root: fullPath, ImgpkgPath: imagesPath, PackagesPath: pkgsPath}, nil
}

// ReadPackages combines all package.yaml files into document in a byte array.
// packagesDir defines the upper level where packages are located.
// It is assumed the overall structure is packagesDir/<package>/<version>/packages.yaml.
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
			// TODO: make this a log statement instead
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

func main() {
	// Inputs:
	// directory where packages live
	// directory where lock files will go
	// name/channel
	// bonus: input file?

	if err := CreateRepo(os.Args[1], os.Args[2], os.Args[3]); err != nil {
		panic(err)
	}
	// 1: bring all the YAML files together into 1 file.
	// 2: Use kbld lock to get image lock files.
}
