package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

//TODO: make this an interface with the fields as methods, as well as Name
type Repo struct {
	Root, Name, InputPath, ImgpkgPath, PkgsPath string
}

func newRepo(inputPath, root, name string) *Repo {
	return &Repo{
		Root:      root,
		Name:      name,
		InputPath: inputPath,
	}
}

func (r *Repo) FullPath() string {
	return filepath.Join(r.Root, r.Name)
}
func (r *Repo) CreateRepo() error {
	err := r.GenerateFileSystem()
	if err != nil {
		return err
	}
	allpkgs := r.ReadPackages()
	os.WriteFile(filepath.Join(r.PkgsPath, "packages.yaml"), allpkgs, 0755)
	// WritePackages
	// Write lock?
	return nil
}

// GenerateRepoFileSystem creates the proper repository file system structure at a given path.
// outputPath is the level at which new repositories will be created.
// name is the child of outputPath that will be populated
func (r *Repo) GenerateFileSystem() error {
	r.ImgpkgPath = filepath.Join(r.FullPath(), ".imgpkg")
	r.PkgsPath = filepath.Join(r.FullPath(), "pkgs")

	for _, d := range []string{r.FullPath(), r.ImgpkgPath, r.PkgsPath} {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}
	return nil
}

// ReadPackages combines all package.yaml files into document in a byte array.
// packagesDir defines the upper level where packages are located.
// It is assumed the overall structure is packagesDir/<package>/<version>/packages.yaml.
func (r *Repo) ReadPackages() []byte {
	packages, err := os.ReadDir(r.InputPath)
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
		pkgpath := filepath.Join(r.InputPath, p.Name())
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
			f := filepath.Join(r.InputPath, p.Name(), v.Name(), "package.yaml")
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

	r := newRepo(os.Args[1], os.Args[2], os.Args[3])
	if err := r.CreateRepo(); err != nil {
		panic(err)
	}
	// 1: bring all the YAML files together into 1 file.
	// 2: Use kbld lock to get image lock files.
}
