package main

import (
	"fmt"
	"os"
	"flag"
	"strings"
	"context"
	"bufio"
	"path/filepath"
	"errors"
	"io/ioutil"

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

	// ExtensionDirectory filename
	ExtensionDirectory string = "extensions"
	// OfflineDirectory filename
	OfflineDirectory string = "offline"
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
	_, dirGH, _, err := client.Repositories.GetContents(ctx, "vmware-tanzu", "tce", ExtensionDirectory, opts)
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

func saveMetadata(metadataDir string, token string, tag string, release bool) (*Metadata, error) {
	list, err := fetchDirectoryList(token)
	if err != nil {
		fmt.Printf("fetchDirectoryList failed: %v\n", err)
		return nil, err
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
	//fmt.Printf("DUMP:\n\n")
	//fmt.Printf("%v\n", metadata)

	byRaw, err := yaml.Marshal(metadata)
	if err != nil {
		fmt.Printf("yaml.Marshal error. Err: %v\n", err)
		return nil, err
	}
	fmt.Printf("BYTES:\n\n")
	fmt.Printf("%s\n", string(byRaw))
	
	// write the file
	fileToWrite := filepath.Join(metadataDir, MetadataFilename)
	fileWrite, err := os.OpenFile(fileToWrite, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("Open Config for write failed. Err: %v\n", err)
		return nil, err
	}

	datawriter := bufio.NewWriter(fileWrite)
	if datawriter == nil {
		fmt.Printf("Datawriter creation failed\n")
		return nil, errors.New("Datawriter creation failed")
	}

	_, err = datawriter.Write(byRaw)
	if err != nil {
		fmt.Printf("datawriter.Write error. Err: %v\n", err)
		return nil, err
	}
	datawriter.Flush()

	// close everything
	err = fileWrite.Close()
	if err != nil {
		fmt.Printf("fileWrite.Close failed. Err: %v", err)
		return nil, err
	}

	return metadata, nil
}

func copyFile(source, destination string) error {
	input, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destination, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

func saveForOffline(md *Metadata, release bool) error {

	// copy all the extensions
	for _, extension := range md.Extensions {
		fmt.Printf("Saving App CRD Extension: %s\n", extension.Name)

		offlineDir := filepath.Join(OfflineDirectory, LatestKeyword, extension.Name)
		if release {
			offlineDir = filepath.Join(OfflineDirectory, md.Version, extension.Name)
		}

		err := os.MkdirAll(offlineDir, 0755)
		if err != nil {
			fmt.Printf("MkdirAll failed. Err: %v", err)
			return err
		}

		srcCrdToCopy := filepath.Join(ExtensionDirectory, extension.Name, AppCrdFilename)
		dstCrdToCopy := filepath.Join(offlineDir, AppCrdFilename)

		err = copyFile(srcCrdToCopy, dstCrdToCopy)
		if err != nil {
			fmt.Printf("copyFile failed. Err: %v", err)
			return err
		}
	}

	return nil
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

	// make metadata dir
	metadataDir := filepath.Join(MetadataDirectory, LatestKeyword)
	if release {
		metadataDir = filepath.Join(MetadataDirectory, tag)
	}
	err := os.MkdirAll(metadataDir, 0755)
	if err != nil {
		fmt.Printf("MkdirAll failed. Err: %v", err)
		return
	}

	// save metadata
	md, err := saveMetadata(metadataDir, token, tag, release)
	if err != nil {
		fmt.Printf("saveMetadata failed. Err: %v", err)
		return
	}

	// save extensions
	err = saveForOffline(md, release)
	if err != nil {
		fmt.Printf("saveForOffline failed. Err: %v", err)
		return
	}

	fmt.Printf("Succeeded\n")
}
