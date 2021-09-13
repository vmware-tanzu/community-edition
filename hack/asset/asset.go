// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	// DefaultCheckSumFilename is tce-checksums.txt
	DefaultCheckSumFilename string = "tce-checksums.txt"
)

func getClientWithEnvToken() (*github.Client, error) {
	var token string
	if v := os.Getenv("GITHUB_TOKEN"); v != "" {
		token = v
	}

	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return client, nil
}

func getDraftRelease(tag string) (*github.RepositoryRelease, error) {
	client, err := getClientWithEnvToken()
	if err != nil {
		fmt.Printf("getClientWithEnvToken returned error: %v\n", err)
		return nil, err
	}

	opt := &github.ListOptions{}
	releasesGH, _, err := client.Repositories.ListReleases(context.Background(), "vmware-tanzu", "community-edition", opt)
	if err != nil {
		fmt.Printf("Repositories.ListReleases returned error: %v\n", err)
		return nil, err
	}

	for _, release := range releasesGH {
		fmt.Printf("Check: %s\n", *release.TagName)

		if !strings.EqualFold(tag, *release.TagName) {
			continue
		}

		if release.PublishedAt == nil {
			fmt.Printf("Draft Release Found: %s\n", *release.TagName)
			return release, nil
		}

		fmt.Printf("Release already published: %s\n", *release.TagName)
		return nil, fmt.Errorf("release already published")
	}

	return nil, fmt.Errorf("unable to find a draft release")
}

func uploadToDraftRelease(release *github.RepositoryRelease, fullPathFilename string) error {
	filename := filepath.Base(fullPathFilename)
	fileAsset, err := os.OpenFile(fullPathFilename, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Printf("OpenFile returned error: %v\n", err)
		return err
	}

	client, err := getClientWithEnvToken()
	if err != nil {
		fmt.Printf("getClientWithEnvToken returned error: %v\n", err)
		return err
	}

	opt := &github.UploadOptions{
		Name: filename,
	}
	_, _, err = client.Repositories.UploadReleaseAsset(context.Background(), "vmware-tanzu", "community-edition", *release.ID, opt, fileAsset)
	if err != nil {
		fmt.Printf("Repositories.UploadReleaseAsset returned error: %v\n", err)
		return err
	}

	return nil
}

func main() {
	var tag string
	flag.StringVar(&tag, "tag", "", "The release tag to add")

	flag.Parse()

	if tag == "" {
		fmt.Printf("A tag must be provided\n")
		return
	}

	draftRelease, err := getDraftRelease(tag)
	if err != nil {
		fmt.Printf("getDraftRelease failed: %v\n", err)
		return
	}

	darwinAMD64AssetFilename := fmt.Sprintf("tce-darwin-amd64-%s.tar.gz", tag)
	darwinAMD64Asset := filepath.Join("..", "..", "build", darwinAMD64AssetFilename)
	err = uploadToDraftRelease(draftRelease, darwinAMD64Asset)
	if err != nil {
		fmt.Printf("uploadToDraftRelease(darwin-amd64) failed: %v\n", err)
		return
	}

	// TODO: Uncomment this when cluster creation is supported on Darwin/ARM64
	// darwinARM64AssetFilename := fmt.Sprintf("tce-darwin-arm64-%s.tar.gz", tag)
	// darwinARM64Asset := filepath.Join("..", "..", "build", darwinARM64AssetFilename)
	// err = uploadToDraftRelease(draftRelease, darwinARM64Asset)
	// if err != nil {
	//     fmt.Printf("uploadToDraftRelease(darwin-arm64) failed: %v\n", err)
	// return
	//}

	linuxAssetFilename := fmt.Sprintf("tce-linux-amd64-%s.tar.gz", tag)
	linuxAsset := filepath.Join("..", "..", "build", linuxAssetFilename)
	err = uploadToDraftRelease(draftRelease, linuxAsset)
	if err != nil {
		fmt.Printf("uploadToDraftRelease(linux) failed: %v\n", err)
		return
	}

	windowsAssetFilename := fmt.Sprintf("tce-windows-amd64-%s.zip", tag)
	windowsAsset := filepath.Join("..", "..", "build", windowsAssetFilename)
	err = uploadToDraftRelease(draftRelease, windowsAsset)
	if err != nil {
		fmt.Printf("uploadToDraftRelease(windows) failed: %v\n", err)
		return
	}

	checksumAsset := filepath.Join("..", "..", "build", DefaultCheckSumFilename)
	err = uploadToDraftRelease(draftRelease, checksumAsset)
	if err != nil {
		fmt.Printf("uploadToDraftRelease(checksum) failed: %v\n", err)
		return
	}

	fmt.Printf("Succeeded\n")
}
