// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package main is a utility program used to copy addon package readmes to the docs
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	SubURL  string `yaml:"suburl"`
}

type SubFolderItem struct {
	Page     string    `yaml:"page,omitempty"`
	URL      string    `yaml:"url,omitempty"`
	Package  MyPackage `yaml:"package,omitempty"`
	SubItems []SubItem `yaml:"subitems,omitempty"`
}

type TocItem struct {
	Title          string          `yaml:"title"`
	SubFolderItems []SubFolderItem `yaml:"subfolderitems"`
}

type Toc struct {
	Toc []TocItem `yaml:"toc"`
}

//nolint:funlen
func main() {
	docsDir := filepath.Join("..", "..", "..", "docs", "site", "content", "docs", "edge")
	imgsDir := filepath.Join(docsDir, "img")
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
	source, err := os.ReadFile(packageRepositoryFile)
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

			// If images directory exists, copy those images to the docs directory
			imagesDirname := findImages(filepath.Join(addonsPackagesDir, packageName, version))
			imgDirSource := filepath.Join(addonsPackagesDir, packageName, version, imagesDirname)
			if imagesDirname != "" {
				entries, err := os.ReadDir(imgDirSource)
				check(err)

				for _, entry := range entries {
					imgSource := filepath.Join(addonsPackagesDir, packageName, version, "images", entry.Name())
					destination := filepath.Join(imgsDir, entry.Name())
					copyFile(imgSource, destination)
				}
			}

			destinationFilename := "package-readme-" + packageName + "-" + version + ".md"
			destination := filepath.Join(docsDir, destinationFilename)
			copyFile(source, destination)

			// If we copied images over, also edit the image references to render correctly in hugo
			if imagesDirname != "" {
				input, err := os.ReadFile(destination)
				check(err)

				fileContents := strings.Replace(string(input), "images/", "../../img/", -1)

				err = os.WriteFile(destination, []byte(fileContents), 0644)
				check(err)
			}
		}
	}

	// load the table of contents yaml
	var latestTocPath = filepath.Join("..", "..", "..", "docs", "site", "data", "docs", "main-toc.yml")
	var toc Toc

	source, err = os.ReadFile(latestTocPath)
	check(err)

	err = yaml.Unmarshal(source, &toc)
	check(err)

	// Create a new array with the new version information
	newPackageVersions := []SubFolderItem{}
	newPackageVersions = append(newPackageVersions, SubFolderItem{
		Page:    "Work with Packages",
		URL:     "/package-management",
		Package: MyPackage{},
	}, SubFolderItem{
		Page:    "Create a Package",
		URL:     "/package-creation-step-by-step",
		Package: MyPackage{},
	})
	for key, value := range currentPackageVersions {
		newPackageVersions = append(newPackageVersions, SubFolderItem{
			Package: MyPackage{
				DisplayName: cases.Title(language.Und).String(strings.Replace(key, "-", " ", -1)),
				Name:        key,
				Versions:    value,
			},
		})
	}

	sort.Slice(newPackageVersions, func(i, j int) bool {
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

	err = os.WriteFile(latestTocPath, data, 0644)
	check(err)

	fmt.Println("Done!")
}

func findReadme(path string) string {
	var filename string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			matches, _ := regexp.MatchString(`[rR][eE][aA][dD][mM][eE].md`, info.Name())
			if matches {
				filename = info.Name()
			}
		}
		return nil
	})

	check(err)
	return filename
}

func findImages(path string) string {
	var dirname string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			matches, _ := regexp.MatchString(`[iI][mM][aA][gG][eE][sS]`, info.Name())
			if matches {
				dirname = info.Name()
			}
		}
		return nil
	})

	check(err)
	return dirname
}

func copyFile(source, destination string) {
	input, err := os.ReadFile(source)
	check(err)

	err = os.WriteFile(destination, input, 0644)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
