package main

import (
	"fmt"
	"io/ioutil"
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
	Packages []Package `packages`
}

func main() {
	var OciRegistry = "projects.registry.vmware.com/tce"
	var PackagesDirectoryPath = filepath.Join("addons", "packages")
	var RepoDirectoryPath = filepath.Join("addons", "repos")
	var GeneratedRepoDirectoryPath = filepath.Join(RepoDirectoryPath, "generated")
	var repository Repository

	channel := os.Args[1]
	channelDir := filepath.Join(GeneratedRepoDirectoryPath, channel)
	imgpkgDir := filepath.Join(channelDir, ".imgpkg")
	packagesDir := filepath.Join(channelDir, "packages")

	// Remove any existing generated files
	os.RemoveAll(fmt.Sprintf(channelDir))
	err := os.MkdirAll(imgpkgDir, 0755)
	check(err)

	err = os.MkdirAll(packagesDir, 0755)
	check(err)

	targetChannelFilename := filepath.Join(RepoDirectoryPath, channel+".yaml")
	source, err := ioutil.ReadFile(targetChannelFilename)
	check(err)

	err = yaml.Unmarshal(source, &repository)
	check(err)

	var outputPackageYaml = filepath.Join(packagesDir, "packages.yaml")
	outputFile, err := os.Create(outputPackageYaml)
	defer func() {
		if err := outputFile.Close(); err != nil {
			panic(err)
		}
	}()

	for _, p := range repository.Packages {
		metadataFilepath := filepath.Join(PackagesDirectoryPath, p.Name, "metadata.yaml")
		copyYaml(metadataFilepath, outputFile)

		for _, version := range p.Versions {
			packageFilepath := filepath.Join(PackagesDirectoryPath, p.Name, version, "package.yaml")
			copyYaml(packageFilepath, outputFile)
		}
	}

	imagesLockFile := filepath.Join(imgpkgDir, "images.yml")
	execCommand("kbld", []string{"--file", packagesDir, "--imgpkg-lock-output", imagesLockFile})

	registryPathAndTag := OciRegistry + "/" + channel + ":latest"
	execCommand("imgpkg", []string{"push", "--bundle", registryPathAndTag, "--file", channelDir})
}

func execCommand(command string, commandArgs []string) {
	out, err := exec.Command(command, commandArgs...).CombinedOutput()
	fmt.Println(string(out))
	check(err)
}

func copyYaml(packageFilepath string, outputFile *os.File) {
	source, err := ioutil.ReadFile(packageFilepath)
	check(err)

	var slice = source[0:3]
	if !strings.HasPrefix(string(slice), "---") {
		outputFile.WriteString("---\n")
	}
	outputFile.Write(source)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
