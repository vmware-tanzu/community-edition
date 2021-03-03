// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package extension

import (
	"fmt"

	"github.com/spf13/cobra"
	klog "k8s.io/klog/v2"
)

var printAlso bool
var getName string
var getForce bool

// GetCmd represents the get command
var GetCmd = &cobra.Command{
	Use:   "get <extension name>",
	Short: "Get extensions",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		mgr, err = NewManager()
		return err
	},
	RunE: get,
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func init() {
	GetCmd.Flags().BoolVarP(&printAlso, "print", "p", false, "Print extension files")
	GetCmd.Flags().BoolVarP(&getForce, "force", "f", false, "Force downloading a new copy of the extension")
}

func getExtension(cmd *cobra.Command, args []string) error {

	klog.V(6).Infof("getExtension = %s", getName)
	klog.V(6).Infof("force download = %t", getForce)

	err := mgr.gh.DownloadExtension(getName, getForce)
	if err != nil {
		fmt.Printf("DownloadExtension failed. Err: %v\n", err)
		return err
	}
	fmt.Printf("Download %s extension succeeded\n", getName)
	return nil
}

func printExtension(cmd *cobra.Command, args []string) error {

	klog.V(6).Infof("printExtension = %s", getName)
	klog.V(6).Infof("force download = %t", getForce)

	err := mgr.gh.PrintExtension(getName, getForce)
	if err != nil {
		fmt.Printf("DownloadExtension failed. Err: %v\n", err)
		return err
	}
	fmt.Printf("Download %s extension succeeded\n", getName)
	return nil
}

func get(cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		fmt.Printf("Please provide extension name\n")
		return ErrMissingExtensionName
	}
	getName = args[0]
	klog.V(2).Infof("get(extension) = %s", getName)
	if getName == "" {
		fmt.Printf("Please provide extension name\n")
		return ErrMissingExtensionName
	}

	if printAlso {
		return printExtension(cmd, args)
	}

	return getExtension(cmd, args)
}
