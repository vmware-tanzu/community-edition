// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	ctlconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	"github.com/vmware-tanzu/carvel-vendir/pkg/vendir/fetch/imgpkgbundle"
	"k8s.io/klog/v2"
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
	pkgBundleLocation, err := mgr.kapp.ResolvePackageBundleLocation(*pkg)
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
func fetchConfig(imageURL string, addonName string) (*string, error) {

	// create a temp directory to store the OCI bundle contents in
	// this directory will be deleted on function return
	dir, err := ioutil.TempDir("", "tce-package-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	// use vendir to sync the OCI bundle down
	conf := ctlconf.DirectoryContentsImgpkgBundle{
		Image: imageURL,
	}
	sync := imgpkgbundle.NewSync(conf, nil)
	klog.V(3).Infof("storing bundle in temp directory %s" + dir)
	_, err = sync.Sync(dir)
	if err != nil {
		return nil, err
	}

	// location of the values file
	valuesFile := dir + string(os.PathSeparator) + "config" + string(os.PathSeparator) + "values.yaml"

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
		return nil, fmt.Errorf("Failed to open file %s. error: %s", valuesFile, err.Error())
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
		return nil, fmt.Errorf("Failed to copy values file. Error: %s", err.Error())
	}

	return &valuesFileNew, nil
}
