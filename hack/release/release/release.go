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

	yaml "github.com/ghodss/yaml"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	// DefaultTagVersion used after tagging a GA release
	DefaultTagVersion string = "dev.1"

	// DevFullPathFilename filename
	DevFullPathFilename string = "../../DEV_BUILD_VERSION.yaml"
	// FakeFullPathFilename filename
	FakeFullPathFilename string = "../../FAKE_BUILD_VERSION.yaml"

	// NewVersionFullPathFilename filename
	NewVersionFullPathFilename string = "../../NEW_BUILD_VERSION"

	// NumberOfSemVerSeparators is 3
	NumberOfSemVerSeparators int = 3
	// NumberOfSeparatorsInDevTag is 2
	NumberOfSeparatorsInDevTag int = 2
)

var (
	// ErrInvalidVersionFormat is Invalid version format
	ErrInvalidVersionFormat = errors.New("invalid version format")
	// ErrDataReaderFailed is Datawriter is empty
	ErrDataReaderFailed = errors.New("datareader is empty")
	// ErrDataWriterFailed is Datawriter is empty
	ErrDataWriterFailed = errors.New("datawriter is empty")
)

type Version struct {
	Version string `json:"version"`
}

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

func updateReleaseNotesOnDraft(release *github.RepositoryRelease, fullPathFilename string) error {
	notes, err := os.ReadFile(fullPathFilename)
	if err != nil {
		fmt.Printf("ReadFile returned error: %v\n", err)
		return err
	}

	client, err := getClientWithEnvToken()
	if err != nil {
		fmt.Printf("getClientWithEnvToken returned error: %v\n", err)
		return err
	}

	strNotes := string(notes)
	release.Body = &strNotes
	_, _, err = client.Repositories.EditRelease(context.Background(), "vmware-tanzu", "community-edition", *release.ID, release)
	if err != nil {
		fmt.Printf("Repositories.EditRelease returned error: %v\n", err)
		return err
	}

	return nil
}

// Release version
func newRelease(current string) error {
	newVersion, err := incrementRelease(current)
	if err != nil {
		fmt.Printf("incrementRelease failed. Err: %v\n", err)
		return err
	}

	err = saveRelease(newVersion)
	if err != nil {
		fmt.Printf("saveDev failed. Err: %v\n", err)
		return err
	}

	return nil
}

func incrementRelease(tag string) (string, error) {
	items := strings.Split(tag, ".")
	if len(items) != NumberOfSemVerSeparators {
		fmt.Printf("Split version failed\n")
		return "", ErrInvalidVersionFormat
	}

	iPatch, err := strconv.Atoi(items[2])
	if err != nil {
		fmt.Printf("Patch string to int failed\n")
		return "", ErrInvalidVersionFormat
	}

	iMinor, err := strconv.Atoi(items[1])
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
	// write the file
	fileWrite, err := os.OpenFile(NewVersionFullPathFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("Open for write failed. Err: %v\n", err)
		return err
	}

	datawriter := bufio.NewWriter(fileWrite)
	if datawriter == nil {
		fmt.Printf("Datawriter creation failed\n")
		return ErrDataWriterFailed
	}

	_, err = datawriter.WriteString(version)
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

// Dev version
func resetDev() error {
	return saveDev(DefaultTagVersion, false)
}

func bumpDev(fake bool) error {
	version, err := getTag(fake)
	if err != nil {
		fmt.Printf("getTag failed. Err: %v\n", err)
		return err
	}

	newVersion, err := incrementDev(version)
	if err != nil {
		fmt.Printf("incrementDev failed. Err: %v\n", err)
		return err
	}

	err = saveDev(newVersion, fake)
	if err != nil {
		fmt.Printf("saveDev failed. Err: %v\n", err)
		return err
	}

	return nil
}

func getTag(fake bool) (string, error) {
	filename := DevFullPathFilename
	if fake {
		filename = FakeFullPathFilename
	}

	fileRead, err := os.OpenFile(filename, os.O_RDONLY, 0755)
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

	version := &Version{}

	err = yaml.Unmarshal(byFile, version)
	if err != nil {
		fmt.Printf("Unmarshal failed. Err: %v\n", err)
		return "", err
	}

	return version.Version, nil
}

func incrementDev(tag string) (string, error) {
	items := strings.Split(tag, ".")
	if len(items) != NumberOfSeparatorsInDevTag {
		fmt.Printf("Split version failed\n")
		return "", ErrInvalidVersionFormat
	}

	ver, err := strconv.Atoi(items[1])
	if err != nil {
		fmt.Printf("String to int failed\n")
		return "", ErrInvalidVersionFormat
	}

	newVersion := ver + 1
	newVersionStr := fmt.Sprintf("dev.%d", newVersion)
	fmt.Printf("incrementDev: %s\n", newVersionStr)

	return newVersionStr, nil
}

func saveDev(tag string, fake bool) error {
	filename := DevFullPathFilename
	if fake {
		filename = FakeFullPathFilename
	}

	version := &Version{
		Version: tag,
	}

	byRaw, err := yaml.Marshal(version)
	if err != nil {
		fmt.Printf("yaml.Marshal error. Err: %v\n", err)
		return err
	}
	fmt.Printf("BYTES:\n\n")
	fmt.Printf("%s\n", string(byRaw))

	// write the file
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

	_, err = datawriter.Write(byRaw)
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
	var tag string
	flag.StringVar(&tag, "tag", "", "The current release tag")

	var notes string
	flag.StringVar(&notes, "notes", "", "The release notes to update with")

	var release bool
	flag.BoolVar(&release, "release", false, "Is this a release")

	flag.Parse()
	// flags

	if tag == "" {
		fmt.Printf("Must supply -tag input\n")
		return
	}
	if notes == "" {
		fmt.Printf("Release notes must be provided\n")
		return
	}

	fake := strings.Contains(tag, "fake")

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
		fmt.Printf("Cutting a release, so bumping\n")
		err := bumpDev(fake)
		if err != nil {
			fmt.Printf("bumpDev failed. Err: %v\n", err)
			return
		}
	}

	draftRelease, err := getDraftRelease(tag)
	if err != nil {
		fmt.Printf("getDraftRelease failed: %v\n", err)
		return
	}

	err = updateReleaseNotesOnDraft(draftRelease, notes)
	if err != nil {
		fmt.Printf("updateReleaseNotesOnDraft failed: %v\n", err)
		return
	}

	fmt.Printf("Succeeded\n")
}
