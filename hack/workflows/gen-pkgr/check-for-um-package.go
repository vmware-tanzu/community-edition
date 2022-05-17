// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package main is a utility program used to check if the generate-package-repository workflow should run
package main

import (
	"errors"
	"os"
	"path/filepath"
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

func main() {
	err := mainErr()
	if err != nil {
		os.Exit(1)
	}
}

func mainErr() error {
	// The first argument should be a comma separated list of the modified/changed files
	if len(os.Args) <= 1 {
		return errors.New("no files in change set")
	}

	arg1 := os.Args[1]

	// if the repo file, addons/repos/main.yaml, was updated, early exit, we want to generate the repo
	//if strings.Contains(arg1, "addons/repos/main.yaml") {
	if strings.Contains(arg1, "addons/repos/seemiller.yaml") {
		return nil
	}

	changeSet := strings.Split(arg1, ",")

	packageUpdates := make(map[string]string)

	// Go through the change set, look for package or metadata yaml files. Add matched packages to maps.
	for _, file := range changeSet {
		if strings.HasPrefix(file, "addons/packages") {
			fileBits := strings.Split(file, "/")
			if strings.HasSuffix(file, "/metadata.yaml") {
				// split up the string, get the package, the package will be the 3rd arg
				packageMetadata := fileBits[2] + "/metadata"
				packageUpdates[packageMetadata] = packageMetadata
			} else if strings.HasSuffix(file, "/package.yaml") {
				packageVersion := fileBits[2] + "/" + fileBits[3]
				packageUpdates[packageVersion] = packageVersion
			}
		}
	}

	// Read in the package repository
	var packageRepository PackageRepository
	packageRepositoryFile := filepath.Join("..", "..", "..", "addons", "repos", "main.yaml")
	source, err := os.ReadFile(packageRepositoryFile)
	check(err)

	err = yaml.Unmarshal(source, &packageRepository)
	check(err)

	// if a package listed in the main.yaml repo file is in the list of files modified, run the workflow
	for _, pkg := range packageRepository.Packages {
		// check for the metadata file changing
		packageMetadata := pkg.Name + "/metadata"
		if packageUpdates[packageMetadata] != "" {
			// package exists in the repo, need to generate new repository, early exit
			return nil
		} else {
			// check for package versions
			for _, version := range pkg.Versions {
				packageVersion := pkg.Name + "/" + version
				if packageUpdates[packageVersion] != "" {
					// package exists in the repo, need to generate new repository, early exit
					return nil
				}
			}
		}
	}

	return errors.New("Do not generate package repository.")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
