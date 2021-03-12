// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List packages available in the cluster",
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

func list(cmd *cobra.Command, args []string) error {

	pkgs, err := mgr.kapp.RetrievePackages()
	if err != nil {
		return err
	}

	// wrapping in anonymouse fun to ensure flush of writer occurs
	func() {
		// setup tab writer to pretty print output
		w := new(tabwriter.Writer)
		// minwidth, tabwidth, padding, padchar, flags
		w.Init(os.Stdout, 8, 8, 0, '\t', 0)
		defer w.Flush()

		// header for output
		fmt.Fprintf(w, " %s\t%s\t%s\t", "NAME", "VERSION", "DESCRIPTION")

		// list all packages known in the cluster
		for _, pkg := range pkgs {
			fmt.Fprintf(w, "\n %s\t%s\t%s\t", pkg.Spec.PublicName, pkg.Spec.Version, pkg.Spec.Description)
		}
	}()

	// ensures a break line after we flush the tabwriter
	fmt.Println()

	return nil
}
