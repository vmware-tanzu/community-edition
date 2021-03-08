// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	klog "k8s.io/klog/v2"
)

// InstallCmd represents the install command
var InstallCmd = &cobra.Command{
	Use:   "install <extension name>",
	Short: "Install extension",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		mgr, err = NewManager()
		return err
	},
	RunE: install,
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func init() {
	// common between secret and user-defined
	InstallCmd.Flags().StringVarP(&inputAppCrd.Namespace, "namespace", "n", "tanzu-extensions", "Namespace to deploy too")

	// secret
	InstallCmd.Flags().StringVarP(&inputAppCrd.ClusterName, "cluster", "c", "", "Cluster name which corresponds to a secret")

	// user defined
	InstallCmd.Flags().StringVarP(&inputAppCrd.URL, "url", "u", "", "URL to image")
	InstallCmd.Flags().StringToStringVarP(&inputAppCrd.Paths, "paths", "p", nil, "User defined paths for kapp template")
}

func install(cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		fmt.Printf("Please provide extension name\n")
		return ErrMissingExtensionName
	}
	inputAppCrd.Name = args[0]
	klog.V(2).Infof("install(extension) = %s", inputAppCrd.Name)
	if inputAppCrd.Name == "" {
		fmt.Printf("Please provide extension name\n")
		return ErrMissingExtensionName
	}

	extensionWorkingDir := filepath.Join(mgr.kapp.GetWorkingDirectory(), inputAppCrd.Name)
	klog.V(5).Infof("extensionWorkingDir = %s", extensionWorkingDir)
	klog.V(3).Infof("installName = %s", inputAppCrd.Name)

	err := mgr.gh.StageFiles(extensionWorkingDir, inputAppCrd.Name)
	if err != nil {
		fmt.Printf("StageFiles failed. Err: %v\n", err)
		return err
	}

	// TODO next release???
	/*
		if inputAppCrd.ClusterName != "" && inputAppCrd.Namespace != "" {
			err = mgr.kapp.InstallFromSecret(inputAppCrd)
			if err != nil {
				fmt.Printf("kclient install failed. Err: %v\n", err)
				return err
			}
		} else if inputAppCrd.URL != "" && inputAppCrd.Paths != nil {
			err = mgr.kapp.InstallFromUser(inputAppCrd)
			if err != nil {
				fmt.Printf("kclient install failed. Err: %v\n", err)
				return err
			}
		}
	*/

	err = mgr.kapp.InstallFromFile(inputAppCrd)
	if err != nil {
		fmt.Printf("InstallFromFile failed. Err: %v\n", err)
		return err
	}

	fmt.Printf("%s install extension succeeded\n", inputAppCrd.Name)

	return nil
}
