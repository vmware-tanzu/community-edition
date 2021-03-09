// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	"fmt"

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
	// common between secret and user-defined
	//InstallCmd.Flags().StringVarP(&inputAppCrd.Namespace, "namespace", "n", "tanzu-extensions", "Namespace to deploy too")

	//// secret
	//InstallCmd.Flags().StringVarP(&inputAppCrd.ClusterName, "cluster", "c", "", "Cluster name which corresponds to a secret")

	//// user defined
	//InstallCmd.Flags().StringVarP(&inputAppCrd.URL, "url", "u", "", "URL to image")
	//InstallCmd.Flags().StringToStringVarP(&inputAppCrd.Paths, "paths", "p", nil, "User defined paths for kapp template")
}

func configure(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		fmt.Printf("Please provide addon name\n")
		return ErrMissingExtensionName
	}

	_, err := mgr.kapp.ResolvePackage("contour-operator", "1.11.0-vmware0")
	if err != nil {
		klog.Errorf("Failed to resolve package. error: %s", err.Error())
	}
	return nil
}

// downloadBundle fetches the remote OCI bundle and saves it in a temp directory
func downloadBundle(imageURL string) {
	klog.Infoln("Downloading addon")

	conf := ctlconf.DirectoryContentsImgpkgBundle{
		Image: "projects.registry.vmware.com/tce/knative-serving-extension-templates:dev",
	}

	sync := imgpkgbundle.NewSync(conf, nil)

	_, err := sync.Sync("/tmp/contents.tar")

	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	fmt.Println("done")
}
