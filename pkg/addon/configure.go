// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	//"os"

	"github.com/spf13/cobra"
	ctlconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	"github.com/vmware-tanzu/carvel-vendir/pkg/vendir/fetch/imgpkgbundle"
	"k8s.io/klog/v2"
)

// ConfigureCmd represents the tanzu addon configure command. It resolves the desired
// package and downloads the imgpkg bundle from the OCI repository. It then unpacks
// the values failes, which a user can modify.
var ConfigureCmd = &cobra.Command{
	Use:   "configure <addon name>",
	Short: "Configure addon",
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

func configure(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		fmt.Printf("Please provide addon name\n")
		return ErrMissingExtensionName
	}
	name := args[0]

	pkg, err := mgr.kapp.ResolvePackage(name, inputAppCrd.Version)
	if err != nil {
		klog.Errorf("Failed to resolve package %s. error: %s", name, err.Error())
	}
	pkgBundleLocation, err := mgr.kapp.ResolvePackageBundleLocation(*pkg)
	if err != nil {
		klog.Errorf("Failed to resolve package %s. error: %s", name, err.Error())
	}

	klog.Infof("pkgbundle location resolved to: %s", pkgBundleLocation)
	err = fetchConfig(pkgBundleLocation)
	if err != nil {
		klog.Errorf("Falied to fetch pkgbundle. error: %s", err.Error())
	}
	return nil
}

// fetchConfig fetches the remote OCI bundle and saves it in a temp directory.
// it then extracts and returns the values file to the current directory.
func fetchConfig(imageURL string) error {
	klog.Infoln("Downloading addon")
	dir, err := ioutil.TempDir("/tmp/", "tce-addon-")
	if err != nil {
		return err
	}
	//defer os.RemoveAll(dir)

	conf := ctlconf.DirectoryContentsImgpkgBundle{
		Image: imageURL,
	}

	sync := imgpkgbundle.NewSync(conf, nil)

	klog.Infoln("dir is: " + dir)
	_, err = sync.Sync(dir)

	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	// location of the values file
	valuesFile := fmt.Sprintf(dir + string(os.PathSeparator) + "config" + string(os.PathSeparator) + "values.yaml")
	// copy file to current directory

	sourceFileStat, err := os.Stat(valuesFile)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", valuesFile)
	}

	s, err := os.Open(valuesFile)
	if err != nil {
		return fmt.Errorf("Failed to open file %s. error: %s", valuesFile, err.Error())
	}
	defer s.Close()

	valuesFileNew := "test.yaml"
	d, err := os.Create(valuesFileNew)
	if err != nil {
		return err
	}
	defer d.Close()
	_, err = io.Copy(d, s)
	if err != nil {
		return fmt.Errorf("Failed to copy values file. Error: %s", err.Error())
	}

	return nil
}
