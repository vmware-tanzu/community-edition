package main

import (
	"fmt"
	"os"
	"context"
	"bufio"
	"path/filepath"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	yaml "github.com/ghodss/yaml"
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
	Stable   string    `json:"stable"`
	Date     string    `json:"date"`
} 

func getTags(token string) ([]*Version, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// tagsGH, _, err := client.Repositories.ListTags(context.Background(), "vmware-tanzu", "tce", nil)
	// if err != nil {
	// 	fmt.Printf("Repositories.ListTags() returned error: %v", err)
	// 	return nil, err
	// }

	opt := &github.ListOptions{}
	tagsGH, _, err := client.Repositories.ListReleases(ctx, "vmware-tanzu", "tce", opt)
	if err != nil {
		fmt.Printf("Repositories.ListReleases returned error: %v\n", err)
		return nil, err
	}

	var tags []*Version
	for _, tag := range tagsGH {
		fmt.Printf("Found: %s\n", *tag.TagName)

		tags = append(tags, &Version{
			Version: *tag.TagName,
			Date: (*tag.PublishedAt).Format(layoutISO),
		})
	}

	return tags, nil
}

func main() {

	var token string
	if v := os.Getenv("GH_ACCESS_TOKEN"); v != "" {
		token = v
	}

	if token == "" {
		fmt.Printf("token is empty\n")
		return
	}

	list, err := getTags(token)
	if err != nil {
		fmt.Printf("getTags failed: %v\n", err)
		return
	}

	releases := &Release{}

	first := true
	for _, item := range list {

		if first {
			releases.Stable = item.Version
			releases.Date = item.Date
			first = false
		}

		releases.Versions = append(releases.Versions, item)
	}
	//fmt.Printf("DUMP:\n\n")
	//fmt.Printf("%v\n", releases)

	byRaw, err := yaml.Marshal(releases)
	if err != nil {
		fmt.Printf("yaml.Marshal error. Err: %v\n", err)
		return
	}
	fmt.Printf("BYTES:\n\n")
	fmt.Printf("%s\n", string(byRaw))

	// make dir
	dirToMake := filepath.Join(MetadataDirectory)
	// err = os.RemoveAll(dirToMake)
	// if err != nil {
	// 	fmt.Printf("RemoveAll failed. Err: %v", err)
	// 	return
	// }
	err = os.MkdirAll(dirToMake, 0755)
	if err != nil {
		fmt.Printf("MkdirAll failed. Err: %v", err)
		return
	}
	
	// write the file
	fileToWrite := filepath.Join(dirToMake, ReleaseFilename)
	fileWrite, err := os.OpenFile(fileToWrite, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("Open Config for write failed. Err: %v\n", err)
		return
	}

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

	// close everything
	err = fileWrite.Close()
	if err != nil {
		fmt.Printf("fileWrite.Close failed. Err: %v", err)
		return
	}

	fmt.Printf("Succeeded\n")
}
