package main

import (
	"fmt"
	"os"
	"flag"
	"strings"
	"context"
	"bufio"
	"path/filepath"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	yaml "github.com/ghodss/yaml"
)

const (
	// FirstRelease first release tag
	FirstRelease string = "v0.1.0"
	// MainBranchKeyword latest
	MainBranchKeyword string = "main"
	// LatestKeyword latest
	LatestKeyword string = "latest"
	// MetadataDirectory filename
	MetadataDirectory string = "metadata"
	// MetadataFilename filename
	MetadataFilename string = "metadata.yaml"
	// AppCrdFilename filename
	AppCrdFilename string = "extension.yaml"
)

type File struct {
	Name        string `json:"filename"`
	Description string `json:"description,omitempty"`
}

// Extension - yep, it's that
type Extension struct {
	Name                   string `json:"name"`
	Description            string `json:"description,omitempty"`
	Version                string `json:"version"`
	KubernetesMinSupported string `json:"minsupported,omitempty"`
	KubernetesMaxSupported string `json:"maxsupported,omitempty"`
	Files                  []*File `json:"files"`
}

// Metadata outer container for metadata
type Metadata struct {
	Extensions      []*Extension `json:"extensions"`
	Version         string      `json:"version"`
	GitHubRepo      string      `json:"repo,omitempty"`
	GitHubBranchTag string      `json:"branch,omitempty"`
}

func fetchDirectoryList(token string) ([]string, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opts := &github.RepositoryContentGetOptions{}
	// if len(branch) > 0 {
	// 	klog.V(6).Infof("Update Ref = %s", branch)
	// 	opts.Ref = branch
	// }
	_, dirGH, _, err := client.Repositories.GetContents(ctx, "vmware-tanzu", "tce", "extensions", opts)
	if err != nil {
		fmt.Printf("client.Repositories failed. Err: %v", err)
		return nil, err
	}

	var extensions []string
	for _, item := range dirGH {
		if strings.EqualFold(*item.Type, "file") {
			fmt.Printf("skip file: %s\n", *item.Name)
			continue
		}
		fmt.Printf("add extension: %s\n", *item.Name)
		extensions = append(extensions, *item.Name)
	}

	return extensions, nil
}

func main() {

	var token string
	if v := os.Getenv("GH_ACCESS_TOKEN"); v != "" {
		token = v
	}

	var tag string
    flag.StringVar(&tag, "tag", "", "The latest tag")
	var release bool
	flag.BoolVar(&release, "release", false, "Is this a release")
	flag.Parse()

	if token == "" {
		fmt.Printf("token is empty\n")
		return
	}
	if tag == "" {
		fmt.Printf("tag is empty\n")
		return
	}

	list, err := fetchDirectoryList(token)
	if err != nil {
		fmt.Printf("fetchDirectoryList failed: %v\n", err)
		return
	}

	metadata := &Metadata{
		Version: tag,
		GitHubBranchTag: tag,
	}
	if !release {
		metadata.GitHubBranchTag = MainBranchKeyword
	}

	for _, item := range list {
		file := &File{
			Name: AppCrdFilename,
		}
		extension := &Extension{
			Name: item,
			Version: tag,
			KubernetesMinSupported: FirstRelease,
			KubernetesMaxSupported: tag,
			Files: []*File{file},
		}

		metadata.Extensions = append(metadata.Extensions, extension)
	}
	fmt.Printf("DUMP:\n\n")
	fmt.Printf("%v\n", metadata)

	byRaw, err := yaml.Marshal(metadata)
	if err != nil {
		fmt.Printf("yaml.Marshal error. Err: %v\n", err)
		return
	}
	fmt.Printf("BYTES:\n\n")
	fmt.Printf("%s\n", string(byRaw))

	// make dir
	dirToMake := filepath.Join(MetadataDirectory, LatestKeyword)
	if release {
		dirToMake = filepath.Join(MetadataDirectory, tag)
	}
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
	fileToWrite := filepath.Join(dirToMake, MetadataFilename)
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
