// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	klog "k8s.io/klog/v2"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List extensions",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		mgr, err = NewManager()
		return err
	},
	RunE: list,
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func init() {
	ListCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Print metadata in format (yaml|json)")
}

func printMetadata(cmd *cobra.Command, args []string) error {
	byMeta, err := mgr.RawMetadata()
	if err != nil {
		klog.Errorf("RawMetadata failed. Err: %v", err)
		return err
	}

	// Print some results...
	switch {
	case outputFormat == "json":
		y, err := yaml.YAMLToJSON(byMeta)
		if err != nil {
			klog.Errorf("err: %v\n", err)
		} else {
			fmt.Printf("%s\n\n", string(y))
		}
	case outputFormat == "yaml":
		y, err := yaml.JSONToYAML(byMeta)
		if err != nil {
			klog.Errorf("err: %v\n", err)
		} else {
			fmt.Printf("%s\n\n", string(y))
		}
	default:
		return fmt.Errorf("unknown output format %v", outputFormat)
	}

	return nil
}

func list(cmd *cobra.Command, args []string) error {

	meta, err := mgr.InitMetadata()
	if err != nil {
		klog.Errorf("InitMetadata failed. Err: %v", err)
		return err
	}

	if outputFormat != "" {
		err = printMetadata(cmd, args)
		if err != nil {
			klog.Errorf("printMetadata failed. Err: %v", err)
		}
	}

	for _, extension := range meta.Extensions {
		fmt.Printf("Extension: %s\n", extension.Name)
	}

	return nil
}
