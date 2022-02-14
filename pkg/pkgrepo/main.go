package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	goUi "github.com/cppforlife/go-cli-ui/ui"
	kbld "github.com/vmware-tanzu/carvel-kbld/pkg/kbld/cmd"
	kbldLogger "github.com/vmware-tanzu/carvel-kbld/pkg/kbld/logger"
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
	// Write lock? This'll be invoking kbld or imgpkg most likelyk
	err = r.LockImages()
	if err != nil {
		return err
	}
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
		pkgpath := filepath.Join(r.InputPath, p.Name())

		metadata, err := os.ReadFile(filepath.Join(pkgpath, "metadata.yaml"))
		if err != nil {
			// No metadata.yaml file, this isn't actually a package.
			if os.IsNotExist(err) {
				break
			}
			// TODO: handle this error better
			panic(err)
		}
		// TODO: Make this a log statement instead
		fmt.Println("Processing package", p.Name())

		entries, err := os.ReadDir(pkgpath)
		// TODO: handle this error better
		if err != nil {
			panic(err)
		}

		fmt.Println("  Processing metadata")
		allpkgs = append(allpkgs, bytes.TrimSpace(metadata))

		for _, e := range entries {
			if !isVersion(e) {
				break
			}
			// TODO: make this a log statement instead
			fmt.Println("  Processing version", e.Name())
			f := filepath.Join(r.InputPath, p.Name(), e.Name(), "package.yaml")
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
	return bytes.Join(allpkgs, []byte("\n---\n"))
}

func (r *Repo) LockImages() error {
	imagesLockFile := filepath.Join(r.ImgpkgPath, "images.yaml")

	kbldOpts := kbld.NewResolveOptions(goUi.NewNoopUI())
	kbldOpts.FileFlags.Files = []string{r.PkgsPath}
	kbldOpts.ImgpkgLockOutput = imagesLockFile
	kbldOpts.FileFlags.Recursive = false

	// Kbld default registry flag options
	kbldOpts.RegistryFlags.Insecure = true

	// Kbld default resolve options
	kbldOpts.AllowedToBuild = true
	kbldOpts.BuildConcurrency = 4
	kbldOpts.ImagesAnnotation = true

	// Discard its output for our own
	logger := kbldLogger.NewLogger(io.Discard)
	pLogger := logger.NewPrefixedWriter("")

	// TODO: make this a log statement instead
	fmt.Printf("Locking images in %s\n", imagesLockFile)
	_, err := kbldOpts.ResolveResources(&logger, pLogger)
	if err != nil {
		return err
	}
	return nil
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
