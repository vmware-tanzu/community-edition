// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package extension

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	klog "k8s.io/klog/v2"
)

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete <extension name>",
	Short: "Delete extension",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		mgr, err = NewManager()
		return err
	},
	RunE: delete,
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func init() {
	// common between secret and user-defined
	DeleteCmd.Flags().StringVarP(&inputAppCrd.Namespace, "namespace", "n", "tanzu-extensions", "Namespace to deploy too")

	// secret
	DeleteCmd.Flags().StringVarP(&inputAppCrd.ClusterName, "cluster", "c", "", "Cluster name which corresponds to a secret")

	// user defined
	DeleteCmd.Flags().StringVarP(&inputAppCrd.URL, "url", "u", "", "URL to image")
	DeleteCmd.Flags().StringToStringVarP(&inputAppCrd.Paths, "paths", "p", nil, "User defined paths for kapp template")

	// delete force
	DeleteCmd.Flags().BoolVarP(&inputAppCrd.Force, "force", "f", false, "Force delete")
	DeleteCmd.Flags().BoolVarP(&inputAppCrd.Teardown, "teardown", "t", false, "Delete associated ServiceAccount and RoleBinding")
}

func delete(cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		fmt.Printf("Please provide extension name\n")
		return ErrMissingExtensionName
	}
	inputAppCrd.Name = args[0]
	if inputAppCrd.Name == "" {
		fmt.Printf("Please provide extension name\n")
		return ErrMissingExtensionName
	}

	extensionWorkingDir := filepath.Join(mgr.kapp.GetWorkingDirectory(), inputAppCrd.Name)
	klog.V(5).Infof("extensionWorkingDir = %s", extensionWorkingDir)
	klog.V(3).Infof("installName = %s", inputAppCrd.Name)

	err := mgr.gh.StageFiles(extensionWorkingDir, inputAppCrd.Name)
	if err != nil {
		fmt.Printf("DownloadExtension failed. Err: %v\n", err)
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

	err = mgr.kapp.DeleteFromFile(inputAppCrd)
	if err != nil {
		fmt.Printf("kclient delete failed. Err: %v\n", err)
		return err
	}

	fmt.Printf("%s delete extension succeeded\n", inputAppCrd.Name)

	return nil
}
