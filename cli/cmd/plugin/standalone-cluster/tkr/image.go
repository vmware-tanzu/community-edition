// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkr

import (
	"fmt"
	"os"
	"path/filepath"

	goUi "github.com/cppforlife/go-cli-ui/ui"
	"github.com/k14s/imgpkg/pkg/imgpkg/cmd"
	"github.com/k14s/ytt/pkg/cmd/template"
	"github.com/k14s/ytt/pkg/cmd/ui"
	"github.com/k14s/ytt/pkg/files"
)

//nolint:golint
type TkrImage struct {
	RegistryURL  string
	DownloadPath string
	ConfigPath   string

	YttValuesFiles   []string
	YttKVsFromYAML   []string
	YttRenderedBytes [][]byte
}

// TkrImageReader enables operations on indivdual image bundles that are referenced from the TKR bom.
//nolint:golint
type TkrImageReader interface {
	// GetRegistryURL returns the bundle's registry URL
	GetRegistryURL() string

	// DownloadBundleImage downloads the OCI image bundle using imgpkg libraries (use DownloadImage for a regular image).
	// The unpacked bundle's files are stored in a temporary directory
	DownloadBundleImage() error

	// DownloadImage downloads a plain, regular image using imgpkg libraries (use DownloadBundleImage for a bundle image).
	// The unpacked image file is stored in a temporary directory
	DownloadImage() error

	// GetDownloadPath returns the path to the local filesystem where the OCI image is/will be downloaded
	GetDownloadPath() string

	// SetRelativeConfigPath sets the _relative_ path for the YTT config bundle in the downloaded OCI image.
	// Example: kapp controller stores it's YTT bundle under "config/" in it's bundle.
	//          So therefore, this function should be called with "config/" as an argument
	SetRelativeConfigPath(string)

	// SetYttYamlValuesBytes adds the files to use as values YAML when rendering ytt for the intended bundle.
	// This method may be called multiple times to add multiple byte slice chunks.
	// For each call, this method will create a temporary values.yaml file that gets piped into YTT.
	AddYttYamlValuesBytes([]byte) error

	// AddYttKVsFromYAML sets the key value pairings parsed as YAML when rendering ytt for the intended bundle.
	// This method may be called multiple times to add multiple string slices.
	// Expected format: all.key1.subkey=true
	AddYttKVsFromYAML([]string)

	// RenderYaml renders the OCI bundle using ytt libraries. The returned slice of byte slices contain the rendered yaml
	// Each byte slice represents a "file" that has been rendered. Typically, this is one single chunk that's been rendered
	// from a directory
	RenderYaml() ([][]byte, error)
}

// NewTkrImageReader provides a new TkrImageReader through the TkrImage struct
// and automatically populates a temporary directory to download the OCI image
func NewTkrImageReader(imagePath string) (TkrImageReader, error) {
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(tempDir)

	t := &TkrImage{
		RegistryURL:  imagePath,
		DownloadPath: tempDir,
	}

	return t, nil
}

func (t *TkrImage) GetRegistryURL() string {
	return t.RegistryURL
}

func (t *TkrImage) DownloadBundleImage() error {
	po := cmd.NewPullOptions(goUi.NewNoopUI())
	po.BundleFlags = cmd.BundleFlags{
		Bundle: t.RegistryURL,
	}
	po.BundleRecursiveFlags = cmd.BundleRecursiveFlags{
		Recursive: true,
	}
	po.OutputPath = t.DownloadPath

	err := po.Run()
	if err != nil {
		return err
	}

	return nil
}

func (t *TkrImage) DownloadImage() error {
	po := cmd.NewPullOptions(goUi.NewNoopUI())
	po.ImageFlags = cmd.ImageFlags{
		Image: t.RegistryURL,
	}
	po.BundleRecursiveFlags = cmd.BundleRecursiveFlags{
		Recursive: true,
	}
	po.OutputPath = t.DownloadPath

	err := po.Run()
	if err != nil {
		return err
	}

	return nil
}

func (t *TkrImage) GetDownloadPath() string {
	return t.DownloadPath
}

func (t *TkrImage) SetRelativeConfigPath(configPath string) {
	t.ConfigPath = filepath.Join(t.DownloadPath, configPath)
}

func (t *TkrImage) AddYttYamlValuesBytes(b []byte) error {
	file, err := os.CreateTemp("", "ytt-values")
	if err != nil {
		return err
	}

	// TODO (jpmcb) - if we attempt to `defer os.Remove(file.Name())` here,
	// the file will be removed once we leave this function scope. We should think about how we can
	// clean up for the user. Or not. These are all landing in the /tmp dir anyways.

	_, err = file.Write(b)
	if err != nil {
		return err
	}

	t.YttValuesFiles = append(t.YttValuesFiles, file.Name())

	return nil
}

func (t *TkrImage) AddYttKVsFromYAML(kvs []string) {
	t.YttKVsFromYAML = append(t.YttKVsFromYAML, kvs...)
}

func (t *TkrImage) RenderYaml() ([][]byte, error) {
	filesToProcess, err := files.NewSortedFilesFromPaths([]string{t.ConfigPath}, files.SymlinkAllowOpts{})
	if err != nil {
		return nil, err
	}

	o := template.NewOptions()

	if t.YttValuesFiles != nil {
		o.DataValuesFlags.FromFiles = t.YttValuesFiles
	}

	if t.YttKVsFromYAML != nil {
		o.DataValuesFlags.KVsFromYAML = t.YttKVsFromYAML
	}

	out := o.RunWithFiles(template.Input{Files: filesToProcess}, ui.NewTTY(false))
	if out.Err != nil {
		return nil, fmt.Errorf("evaluating ytt: %s", out.Err)
	}

	if len(out.Files) == 0 {
		return nil, fmt.Errorf("expected to find yaml files but saw zero files after ytt processing")
	}

	processedBytes := [][]byte{}
	for _, file := range out.Files {
		processedBytes = append(processedBytes, file.Bytes())
	}

	// This sets the in the image reader itself so they may be referenced elsewhere.
	// These bytes are also returned from this function call
	t.YttRenderedBytes = processedBytes

	return processedBytes, nil
}
