// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package main is a utility program used to generate TCE addon manifests
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Package struct {
	Name     string   `yaml:"name"`
	Versions []string `yaml:"versions"`
}

type Repository struct {
	Packages []Package `yaml:"packages"`
}

type BundleLock struct {
	APIVersion string    `json:"apiVersion"`
	Kind       string    `json:"kind"`
	Bundle     BundleRef `json:"bundle"` // This generated yaml, but due to lib we need to use `json`
}

type BundleRef struct {
	Image string `json:"image,omitempty"` // This generated yaml, but due to lib we need to use `json`
	Tag   string `json:"tag,omitempty"`   // This generated yaml, but due to lib we need to use `json`
}

func main() {
	var PackagesDirectoryPath = filepath.Join("..", "..", "addons", "packages")
	var RepoDirectoryPath = filepath.Join("..", "..", "addons", "repos")
	var GeneratedRepoDirectoryPath = filepath.Join(RepoDirectoryPath, "generated")
	var repository Repository

	channel := os.Args[1]
	tag := os.Args[2]
	ociRegistry := os.Args[3]
	channelDir := filepath.Join(GeneratedRepoDirectoryPath, channel)
	imgpkgDir := filepath.Join(channelDir, ".imgpkg")
	packagesDir := filepath.Join(channelDir, "packages")

	// Remove any existing generated files
	os.RemoveAll(channelDir)
	err := os.MkdirAll(imgpkgDir, 0755)
	check(err)

	err = os.MkdirAll(packagesDir, 0755)
	check(err)

	targetChannelFilename := filepath.Join(RepoDirectoryPath, channel+".yaml")
	source, err := os.ReadFile(targetChannelFilename)
	check(err)

	err = yaml.Unmarshal(source, &repository)
	check(err)

	var outputPackageYaml = filepath.Join(packagesDir, "packages.yaml")
	outputFile, err := os.Create(outputPackageYaml)
	check(err)

	defer func() {
		if err := outputFile.Close(); err != nil {
			panic(err)
		}
	}()

	for _, p := range repository.Packages {
		metadataFilepath := getYamlFilepath(filepath.Join(PackagesDirectoryPath, p.Name, "metadata"))
		copyYaml(metadataFilepath, outputFile)

		for _, version := range p.Versions {
			packageFilepath := getYamlFilepath(filepath.Join(PackagesDirectoryPath, p.Name, version, "package"))
			copyYaml(packageFilepath, outputFile)
		}
	}

	imagesLockFile := filepath.Join(imgpkgDir, "images.yml")
	execCommand("kbld", []string{"--file", packagesDir, "--imgpkg-lock-output", imagesLockFile})

	bundleLockFilename := "output.yaml"
	registryPathAndTag := ociRegistry + "/" + channel + ":" + tag
	execCommand("imgpkg", []string{"push", "--tty", "--bundle", registryPathAndTag, "--file", channelDir, "--lock-output", bundleLockFilename})

	bundleLockYamlFile, err := os.ReadFile(bundleLockFilename)
	check(err)

	var bundleLock BundleLock
	err = yaml.Unmarshal(bundleLockYamlFile, &bundleLock)
	check(err)

	fmt.Println(bundleLock.Bundle.Image)
	os.RemoveAll(bundleLockFilename)
}

func execCommand(command string, commandArgs []string) {
	output, err := exec.Command(command, commandArgs...).CombinedOutput()
	if err != nil {
		log.Fatal(string(output))
	}
}

func copyYaml(packageFilepath string, outputFile *os.File) {
	source, err := os.ReadFile(packageFilepath)
	check(err)

	var slice = source[0:3]
	if !strings.HasPrefix(string(slice), "---") {
		if _, err := outputFile.WriteString("---\n"); err != nil {
			panic(err)
		}
	}

	_, err = outputFile.Write(source)
	check(err)

	slice = source[len(source)-1:]
	if string(slice) != "\n" {
		if _, err := outputFile.WriteString("\n"); err != nil {
			panic(err)
		}
	}
}

func getYamlFilepath(filename string) string {
	yamlFilepath := filename + ".yaml"
	_, err := os.Stat(yamlFilepath)
	if err != nil {
		yamlFilepath = filename + ".yml"
		_, err = os.Stat(yamlFilepath)
		check(err)
	}
	return yamlFilepath
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
