package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type PackageRepositoryPackage struct {
	Name     string   `yaml:"name"`
	Versions []string `yaml:"versions"`
}

type PackageRepository struct {
	Packages []PackageRepositoryPackage `yaml:"packages"`
}

type MyPackage struct {
	DisplayName string   `yaml:"displayName"`
	Name        string   `yaml:"name"`
	Versions    []string `yaml:"versions"`
}

type SubItem struct {
	Subpage string `yaml:"subpage"`
	SubUrl  string `yaml:"suburl"`
}

type SubFolderItem struct {
	Page    string    `yaml:"page,omitempty"`
	URL     string    `yaml:"url,omitempty"`
	Package MyPackage `yaml:"package,omitempty"`
	SubItems []SubItem `yaml:"subitems,omitempty"`
}

type TocItem struct {
	Title          string          `yaml:"title"`
	SubFolderItems []SubFolderItem `yaml:"subfolderitems"`
}

type Toc struct {
	Toc []TocItem `yaml:"toc"`
}

func main() {
	docsDir := filepath.Join("..", "..", "..", "docs", "site", "content", "docs", "latest")
	addonsPackagesDir := filepath.Join("..", "..", "..", "addons", "packages")

	// delete any existing package readme files
	files, _ := filepath.Glob(filepath.Join(docsDir, "package-readme-*"))
	for _, file := range files {
		err := os.Remove(file)
		check(err)
	}

	// Create a map of packages and their current versions
	var packageRepository PackageRepository
	packageRepositoryFile := filepath.Join("..", "..", "..", "addons", "repos", "main.yaml")
	source, err := ioutil.ReadFile(packageRepositoryFile)
	check(err)

	err = yaml.Unmarshal(source, &packageRepository)
	check(err)

	var currentPackageVersions = make(map[string][]string)
	for _, pkg := range packageRepository.Packages {
		currentPackageVersions[pkg.Name] = pkg.Versions
	}

	// Copy package readmes to the docs directory
	for packageName, versions := range currentPackageVersions {
		fmt.Println("package:", packageName, "=>", versions)
		for _, version := range versions {
			readmeFilename := findReadme(filepath.Join(addonsPackagesDir, packageName, version))
			source := filepath.Join(addonsPackagesDir, packageName, version, readmeFilename)
			destinationFilename := "package-readme-" + packageName + "-" + version + ".md"
			destination := filepath.Join(docsDir, destinationFilename)
			copyFile(source, destination)
		}
	}

	// load the table of contents yaml
	var latestTocPath = filepath.Join("..", "..", "..", "docs", "site", "data", "docs", "latest-toc.yml")
	var toc Toc

	source, err = ioutil.ReadFile(latestTocPath)
	check(err)

	err = yaml.Unmarshal(source, &toc)
	check(err)

	// Create a new array with the new version information
	newPackageVersions := []SubFolderItem{}
	newPackageVersions = append(newPackageVersions, SubFolderItem{
		Page:    "Packages Introduction",
		URL:     "/packages-intro",
		Package: MyPackage{},
	})
	for key, value := range currentPackageVersions {
		newPackageVersions = append(newPackageVersions, SubFolderItem{
			Package: MyPackage{
				DisplayName: strings.Title(strings.Replace(key, "-", " ", -1)),
				Name:        key,
				Versions:    value,
			},
		})
	}

	sort.Slice(newPackageVersions[:], func(i, j int) bool {
		return newPackageVersions[i].Package.Name < newPackageVersions[j].Package.Name
	})

	// Find the Packages in the TOC and replace it with the new values
	PackageIndex := 0
	for index, item := range toc.Toc {
		if item.Title == "Packages" {
			PackageIndex = index
			break
		}
	}
	toc.Toc[PackageIndex] = TocItem{
		Title:          "Packages",
		SubFolderItems: newPackageVersions,
	}

	// Write out the YAML
	data, err := yaml.Marshal(&toc)
	check(err)

	err = ioutil.WriteFile(latestTocPath, data, 0644)
	check(err)

	fmt.Println("Done!")
}

func findReadme(path string) string {
	var filename string
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			matches, _ := regexp.MatchString(`[rR][eE][aA][dD][mM][eE].md`, info.Name())
			if matches {
				filename = info.Name()
			}
		}
		return nil
	})
	return filename
}

func copyFile(source string, destination string) {
	input, err := ioutil.ReadFile(source)
	check(err)

	err = ioutil.WriteFile(destination, input, 0644)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
