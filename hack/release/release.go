// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	yaml "github.com/ghodss/yaml"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	layoutISO = "2006-01-02"

	// MetadataDirectory filename
	MetadataDirectory string = "metadata"
	// ReleaseFilename filename
	ReleaseFilename string = "release-management.yaml"
)

// Version - released version
type Version struct {
	Version string `json:"version"`
	Date    string `json:"date,omitempty"`
}

// Release outer container for metadata
type Release struct {
	Versions []*Version `json:"versions"`
	Stable   string     `json:"stable"`
	Date     string     `json:"date"`
}

func getTags(token string) ([]*Version, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.ListOptions{}
	tagsGH, _, err := client.Repositories.ListReleases(ctx, "vmware-tanzu", "tce", opt)
	if err != nil {
		fmt.Printf("Repositories.ListReleases returned error: %v\n", err)
		return nil, err
	}

	var tags []*Version
	for _, tag := range tagsGH {
		fmt.Printf("Found: %s\n", *tag.TagName)

		if tag.PublishedAt == nil {
			continue
		}

		tags = append(tags, &Version{
			Version: *tag.TagName,
			Date:    tag.PublishedAt.Format(layoutISO),
		})
	}

	return tags, nil
}

func getTagsWithEnvToken() ([]*Version, error) {
	var token string
	if v := os.Getenv("GH_ACCESS_TOKEN"); v != "" {
		token = v
	}

	if token == "" {
		var tags []*Version
		return tags, fmt.Errorf("token is empty")
	}

	return getTags(token)
}

// makeOutputFile ensures the path and file exists for writing out our config.
// It is the responsibility of the caller to close the file if error is not
// returned.
func makeOutputFile() (*os.File, error) {
	// Make sure the directory exists
	dirToMake := filepath.Join(MetadataDirectory)
	err := os.MkdirAll(dirToMake, 0755)
	if err != nil {
		return nil, err
	}

	// Open the file for writing
	fileToWrite := filepath.Join(dirToMake, ReleaseFilename)
	fileWrite, err := os.OpenFile(fileToWrite, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return nil, err
	}

	return fileWrite, nil
}

func main() {
	var tag string
	flag.StringVar(&tag, "tag", "", "The release tag to add")

	var release bool
	flag.BoolVar(&release, "release", false, "Is this a release")

	flag.Parse()

	if release && tag == "" {
		fmt.Printf("If release is set, a tag must be provided\n")
		return
	}

	list, err := getTagsWithEnvToken()
	if err != nil {
		fmt.Printf("getTags failed: %v\n", err)
		return
	}

	first := true
	releases := &Release{}

	if release {
		first = false

		thisRelease := &Version{
			Version: tag,
			Date:    time.Now().Format(layoutISO),
		}

		releases.Stable = thisRelease.Version
		releases.Date = thisRelease.Date

		releases.Versions = append(releases.Versions, thisRelease)
	}

	for _, item := range list {
		if first {
			releases.Stable = item.Version
			releases.Date = item.Date
			first = false
		}

		releases.Versions = append(releases.Versions, item)
	}

	byRaw, err := yaml.Marshal(releases)
	if err != nil {
		fmt.Printf("yaml.Marshal error. Err: %v\n", err)
		return
	}

	// make dir
	fileWrite, err := makeOutputFile()
	if err != nil {
		fmt.Printf("failed to create config file, err: %v\n", err)
		return
	}
	defer fileWrite.Close()

	datawriter := bufio.NewWriter(fileWrite)
	if datawriter == nil {
		fmt.Printf("Datawriter creation failed\n")
		return
	}

	_, err = datawriter.Write(byRaw)
	if err != nil {
		fmt.Printf("datawriter.Write error. Err: %v\n", err)
		return
	}
	datawriter.Flush()

	fmt.Printf("Succeeded\n")
}
