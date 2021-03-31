// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/joshrosso/image/v5/copy"
	"github.com/joshrosso/image/v5/manifest"
	"github.com/joshrosso/image/v5/signature"
	"github.com/joshrosso/image/v5/transports/alltransports"
	"github.com/spf13/cobra"
)

const (
	limitReadSize   = 2 * 1024 * 1024
	layerDir        = "layer"
	transportPrefix = "docker://"
	parsedImgSize   = 2
)

// ConfigureCmd represents the tanzu package configure command. It resolves the desired
// package and downloads the imgpkg bundle from the OCI repository. It then unpacks
// the values failes, which a user can modify.
var ConfigureCmd = &cobra.Command{
	Use:   "configure <package name>",
	Short: "Configure a package by downloading its configuration",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		mgr, err = NewManager()
		return err
	},
	RunE: configure,
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func init() {
	ConfigureCmd.Flags().StringVarP(&inputAppCrd.Version, "package-version", "t", "", "Version of the package")
}

// configure implements the ConfigureCmd and retrieve configuration by
// 1. resolving the Package CR based on name and/or version
// 2. resolving the imgpkgbundle's repo (OCI registry) location
// 3. downloading the OCI bundle
// 4. extracting the values file for the extension
func configure(cmd *cobra.Command, args []string) error {
	// validate a package name was passed
	if len(args) < 1 {
		fmt.Println("Please provide package name")
		return ErrMissingPackageName
	}
	name := args[0]

	// find the Package CR that corresponds to the name and/or version
	fmt.Printf("Looking up config for package: %s:%s\n", name, inputAppCrd.Version)
	pkg, err := mgr.kapp.ResolvePackage(name, inputAppCrd.Version)
	if err != nil {
		return err
	}

	// extract the OCI bundle's location in a registry
	pkgBundleLocation, err := mgr.kapp.ResolvePackageBundleLocation(pkg)
	if err != nil {
		return err
	}

	// download and extract the values file from the bundle
	configFile, err := fetchConfig(pkgBundleLocation, name)
	if err != nil {
		return err
	}

	fmt.Printf("Values files saved to %s. Configure this file before installing the package.\n", *configFile)
	return nil
}

// fetchConfig fetches the remote OCI bundle and saves it in a temp directory.
// it then extracts and saves the values file to the current directory.
// When successful, the path to the stored values file is returned.
func fetchConfig(imageURL, addonName string) (*string, error) {
	// create a temp directory to store the OCI bundle contents in
	// this directory will be deleted on function return
	dir, err := os.MkdirTemp("", "tce-package-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	_, err = downloadAndUnpackBundle(imageURL, dir)
	if err != nil {
		return nil, err
	}

	// location of the values file
	valuesFile := dir + string(os.PathSeparator) + "layer" + string(os.PathSeparator) + "config" + string(os.PathSeparator) + "values.yaml"

	// copy the values files into the current directory
	sourceFileStat, err := os.Stat(valuesFile)
	if err != nil {
		return nil, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a regular file", valuesFile)
	}
	s, err := os.Open(valuesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s. error: %s", valuesFile, err.Error())
	}
	defer s.Close()
	valuesFileNew := fmt.Sprintf("%s-values.yaml", addonName)
	d, err := os.Create(valuesFileNew)
	if err != nil {
		return nil, err
	}
	defer d.Close()
	_, err = io.Copy(d, s)
	if err != nil {
		return nil, fmt.Errorf("failed to copy values file. Error: %s", err.Error())
	}

	return &valuesFileNew, nil
}

// downloadAndUnpackBundle takes an imageURL and downloads it. It then unpacks the first layer (based on the
// manifest) to a layer directory inside the tempDir directory.
func downloadAndUnpackBundle(imageURL, tempDir string) (*string, error) {
	fullLayerDir := tempDir + string(os.PathSeparator) + layerDir
	ctx := context.Background()
	policy, err := signature.NewPolicyFromBytes([]byte(StaticPolicy))
	if err != nil {
		return nil, err
	}
	pc, err := signature.NewPolicyContext(policy)
	if err != nil {
		return nil, err
	}

	parsedImgName := transportPrefix + imageURL
	srcRef, err := alltransports.ParseImageName(parsedImgName)
	if err != nil {
		return nil, err
	}
	destRef, err := alltransports.ParseImageName("dir:" + fullLayerDir)
	if err != nil {
		return nil, err
	}

	// copy the tar'd image locally
	mf, err := copy.Image(ctx, pc, destRef, srcRef, &copy.Options{})
	if err != nil {
		return nil, err
	}

	// renter the image's manifest into a structured type
	di, err := manifest.FromBlob(mf, manifest.DockerV2Schema2MediaType)
	if err != nil {
		return nil, err
	}

	if len(di.LayerInfos()) < 1 {
		return nil, fmt.Errorf("could not retrieve layers of OCI bundle")
	}
	parsedImg := strings.Split(di.LayerInfos()[0].Digest.String(), "sha256:")
	if len(parsedImg) < parsedImgSize {
		return nil, fmt.Errorf("layer of OCI bundle had invalid digest value")
	}
	sp := parsedImg[1]

	fp := fullLayerDir + string(os.PathSeparator) + sp
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	err = untar(f, fullLayerDir)
	if err != nil {
		return nil, err
	}

	rp := layerDir
	return &rp, nil
}

//gocyclo:ignore
func untar(r io.Reader, dir string) (err error) {
	nFiles := 0
	madeDir := map[string]bool{}
	zr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("requires gzip-compressed body: %v", err)
	}
	tr := tar.NewReader(zr)
	for {
		f, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar error: %v", err)
		}
		if !validRelPath(f.Name) {
			return fmt.Errorf("tar contained invalid name error %q", f.Name)
		}
		rel := filepath.FromSlash(f.Name)
		abs := filepath.Join(dir, rel)

		fi := f.FileInfo()
		mode := fi.Mode()
		switch {
		case mode.IsRegular():
			// Make the directory. This is redundant because it should
			// already be made by a directory entry in the tar
			// beforehand. Thus, don't check for errors; the next
			// write will fail with the same error.
			dir := filepath.Dir(abs)
			if !madeDir[dir] {
				if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
					return err
				}
				madeDir[dir] = true
			}
			wf, err := os.OpenFile(abs, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode.Perm())
			if err != nil {
				return err
			}
			// prevents DoS explosion
			lr := io.LimitReader(tr, limitReadSize)
			n, err := io.Copy(wf, lr)
			if closeErr := wf.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
			if err != nil {
				return fmt.Errorf("error writing to %s: %v", abs, err)
			}
			if n != f.Size {
				return fmt.Errorf("only wrote %d bytes to %s; expected %d", n, abs, f.Size)
			}
			nFiles++
		case mode.IsDir():
			if err := os.MkdirAll(abs, 0755); err != nil {
				return err
			}
			madeDir[abs] = true
		default:
			return fmt.Errorf("tar file entry %s contained unsupported file type %v", f.Name, mode)
		}
	}
	return nil
}

func validRelPath(p string) bool {
	if p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../") {
		return false
	}
	return true
}
