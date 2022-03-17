// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	// DefaultTagVersion used after tagging a GA release
	DefaultTagVersion string = "dev.1"

	// BuildFullPathFilename filename
	BuildFullPathFilename string = "../../NEW_BUILD_VERSION"

	// DevFullPathFilename filename
	DevFullPathFilename string = "../../DEV_VERSION"

	// TemplateFullPathFilename filename
	TemplateFullPathFilename string = "./release.template"

	// NumberOfSemVerSeparators is 3
	NumberOfSemVerSeparators int = 3
	// NumberOfSeparatorsInDevTag is 2
	NumberOfSeparatorsInDevTag int = 2
	// NumberOfPartsTag is 2
	NumberOfPartsTag int = 2
)

var (
	// ErrInvalidVersionFormat is Invalid version format
	ErrInvalidVersionFormat = errors.New("invalid version format")
	// ErrDataReaderFailed is Datawriter is empty
	ErrDataReaderFailed = errors.New("datareader is empty")
	// ErrDataWriterFailed is Datawriter is empty
	ErrDataWriterFailed = errors.New("datawriter is empty")
)

// update release notes
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

// update DraftRelease
func updateReleaseNotesOnDraft(previous, tag string) error {
	client, err := getClientWithEnvToken()
	if err != nil {
		fmt.Printf("getClientWithEnvToken returned error: %v\n", err)
		return err
	}

	draftRelease, err := getDraftRelease(client, tag)
	if err != nil {
		fmt.Printf("getDraftRelease failed: %v\n", err)
		return err
	}

	notes, err := os.ReadFile(TemplateFullPathFilename)
	if err != nil {
		fmt.Printf("ReadFile returned error: %v\n", err)
		return err
	}

	notes1 := strings.Replace(string(notes), "<PREVIOUS_VERSION>", previous, -1)
	notesFinal := strings.Replace(notes1, "<BUILD_VERSION>", tag, -1)

	draftRelease.Body = &notesFinal
	_, _, err = client.Repositories.EditRelease(context.Background(), "vmware-tanzu", "community-edition", *draftRelease.ID, draftRelease)
	if err != nil {
		fmt.Printf("Repositories.EditRelease returned error: %v\n", err)
		return err
	}

	return nil
}

// helper DraftRelease functions
func getDraftRelease(client *github.Client, tag string) (*github.RepositoryRelease, error) {
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

// update Release version
func newRelease(current string) error {
	fmt.Printf("tag: %s\n", current)

	newVersion, err := incrementRelease(current)
	if err != nil {
		fmt.Printf("incrementRelease failed. Err: %v\n", err)
		return err
	}

	devTag, err := getDevBuild()
	if err != nil {
		fmt.Printf("getDevBuild failed. Err: %v\n", err)
		return err
	}

	newVersionStr := fmt.Sprintf("%s-%s", newVersion, devTag)
	fmt.Printf("incrementRelease: %s\n", newVersionStr)

	err = saveRelease(newVersionStr)
	if err != nil {
		fmt.Printf("saveRelease failed. Err: %v\n", err)
		return err
	}

	return nil
}

func bumpRelease(current string) error {
	fmt.Printf("tag: %s\n", current)

	devTag, err := getDevBuild()
	if err != nil {
		fmt.Printf("getDevBuild failed. Err: %v\n", err)
		return err
	}

	items := strings.Split(current, "-")
	if len(items) != NumberOfPartsTag {
		fmt.Printf("Split version failed\n")
		return ErrInvalidVersionFormat
	}

	properVersion := strings.TrimSpace(items[0])
	newVersionStr := fmt.Sprintf("%s-%s", properVersion, devTag)
	fmt.Printf("incrementRelease: %s\n", newVersionStr)

	err = saveRelease(newVersionStr)
	if err != nil {
		fmt.Printf("saveRelease failed. Err: %v\n", err)
		return err
	}

	return nil
}

// helper Release functions
func incrementRelease(tag string) (string, error) {
	items := strings.Split(tag, ".")
	if len(items) != NumberOfSemVerSeparators {
		fmt.Printf("Split version failed\n")
		return "", ErrInvalidVersionFormat
	}

	iPatch, err := strconv.Atoi(strings.TrimSpace(items[2]))
	if err != nil {
		fmt.Printf("Patch string to int failed\n")
		return "", ErrInvalidVersionFormat
	}

	iMinor, err := strconv.Atoi(strings.TrimSpace(items[1]))
	if err != nil {
		fmt.Printf("Minor string to int failed\n")
		return "", ErrInvalidVersionFormat
	}

	// are we on a release branch (ie vX.Y.[0-9]+)? then increment the patch version
	// otherwise, this is a minor release and increment the minor version
	if iPatch > 0 {
		iPatch++
	} else {
		iMinor++
	}

	oldMajor := items[0]
	newVersionStr := fmt.Sprintf("%s.%d.%d", oldMajor, iMinor, iPatch)
	fmt.Printf("incrementRelease: %s\n", newVersionStr)

	return newVersionStr, nil
}

func saveRelease(version string) error {
	return saveFile(BuildFullPathFilename, version)
}

// update Dev version
func resetDev() error {
	return saveDev(DefaultTagVersion)
}

func bumpDev() error {
	version, err := getDevBuild()
	if err != nil {
		fmt.Printf("getDevBuild failed. Err: %v\n", err)
		return err
	}

	newVersion, err := incrementDev(version)
	if err != nil {
		fmt.Printf("incrementDev failed. Err: %v\n", err)
		return err
	}

	err = saveDev(newVersion)
	if err != nil {
		fmt.Printf("saveDev failed. Err: %v\n", err)
		return err
	}

	return nil
}

// helper Dev functions
func getDevBuild() (string, error) {
	fileRead, err := os.OpenFile(DevFullPathFilename, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Printf("Open for read failed. Err: %v\n", err)
		return "", err
	}
	defer fileRead.Close()

	dataReader := bufio.NewReader(fileRead)
	if dataReader == nil {
		fmt.Printf("Datareader creation failed\n")
		return "", ErrDataReaderFailed
	}

	byFile, err := io.ReadAll(dataReader)
	if err != nil {
		fmt.Printf("ReadAll failed. Err: %v\n", err)
		return "", err
	}

	devBuild := string(byFile)
	fmt.Printf("DEV_BUILD: %s\n", devBuild)

	return devBuild, nil
}

func incrementDev(devBuild string) (string, error) {
	items := strings.Split(devBuild, ".")
	if len(items) != NumberOfSeparatorsInDevTag {
		fmt.Printf("Split version failed\n")
		return "", ErrInvalidVersionFormat
	}

	ver, err := strconv.Atoi(strings.TrimSpace(items[1]))
	if err != nil {
		fmt.Printf("String to int failed\n")
		return "", ErrInvalidVersionFormat
	}

	newVersion := ver + 1
	newVersionStr := fmt.Sprintf("dev.%d", newVersion)
	fmt.Printf("incrementDev: %s\n", newVersionStr)

	return newVersionStr, nil
}

func saveDev(devBuild string) error {
	return saveFile(DevFullPathFilename, devBuild)
}

// help
func saveFile(filename, content string) error {
	fileWrite, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("Open for write failed. Err: %v\n", err)
		return err
	}

	datawriter := bufio.NewWriter(fileWrite)
	if datawriter == nil {
		fmt.Printf("Datawriter creation failed\n")
		return ErrDataWriterFailed
	}

	_, err = datawriter.WriteString(content)
	if err != nil {
		fmt.Printf("datawriter.Write error. Err: %v\n", err)
		return err
	}
	datawriter.Flush()

	err = fileWrite.Close()
	if err != nil {
		fmt.Printf("fileWrite.Close failed. Err: %v\n", err)
		return err
	}

	return nil
}

func main() {
	// flags
	var previous string
	flag.StringVar(&previous, "previous", "", "The previous release tag")

	var tag string
	flag.StringVar(&tag, "tag", "", "The current release tag")

	var skip bool
	flag.BoolVar(&skip, "skip", false, "Skip making changes to draft release")

	flag.Parse()
	// flags

	if previous == "" {
		fmt.Printf("Must supply -previous input\n")
		return
	}
	if tag == "" {
		fmt.Printf("Must supply -tag input\n")
		return
	}

	release := !strings.ContainsAny(tag, "-")

	if release {
		fmt.Printf("Cutting GA release, so resetting\n")
		err := resetDev()
		if err != nil {
			fmt.Printf("resetDev failed. Err: %v\n", err)
			return
		}

		err = newRelease(tag)
		if err != nil {
			fmt.Printf("newRelease failed. Err: %v\n", err)
			return
		}
	} else {
		fmt.Printf("Cutting a Non-GA release, so bumping\n")
		err := bumpDev()
		if err != nil {
			fmt.Printf("bumpDev failed. Err: %v\n", err)
			return
		}

		err = bumpRelease(tag)
		if err != nil {
			fmt.Printf("newRelease failed. Err: %v\n", err)
			return
		}
	}

	if !skip {
		err := updateReleaseNotesOnDraft(previous, tag)
		if err != nil {
			fmt.Printf("updateReleaseNotesOnDraft failed: %v\n", err)
			return
		}
	}

	fmt.Printf("Succeeded\n")
}
