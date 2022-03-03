// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const docsDesc = `
Provides documentation and usage for unmanaged-cluster in multiple formats.
Gives all commands, their flags, and related commands.
This can then be used as a local documentation resource when working with unmanaged-clusters.

Users may generate man pages that can then be zipped and placed on the system.
The 'man' program expects the zipped bundle to be in a location defined by 'manpath'
and the user must run 'mandb' to inject the manpages onto their system.
`

var (
	docsFormat    = ""
	docsOutputDir = ""
	noStdout      = false
)

// ListCmd returns a list of existing clusters.
var DocsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generated documentation & usage of unmanaged-cluster",
	Long:  docsDesc,
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: startDocsGen,
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func init() {
	DocsCmd.Flags().StringVar(&docsFormat, "format", "markdown", "Output format (markdown|manpage|rest|yaml); default is markdown")
	DocsCmd.Flags().BoolVar(&noStdout, "no-stdout", false, "Ignores printing to stdout and inspects 'output-dir' for location to dump files")
	DocsCmd.Flags().StringVar(&docsOutputDir, "output-dir", "", "Location to output the doc page files; providing this prevents printing to stdout. Defaults to the current working dir")
}

// list outputs a list of all unmanaged clusters on the system.
func startDocsGen(cmd *cobra.Command, args []string) error {
	var err error

	if len(args) != 0 {
		return fmt.Errorf("docs generation does not support any arguments")
	}

	// Valid formats
	if (docsFormat != "markdown" && docsFormat != "manpage") && (docsFormat != "rest" && docsFormat != "yaml") {
		return fmt.Errorf("unsupported format: %s - Valid formats: (markdown|manpage|rest|yaml)", docsFormat)
	}

	if noStdout || cmd.Flags().Changed("output-dir") {
		err = docsToFiles(cmd)
	} else {
		err = docsToStdout(cmd)
	}

	return err
}

func docsToFiles(cmd *cobra.Command) error {
	if docsOutputDir == "" {
		p, err := os.Getwd()
		if err != nil {
			return err
		}
		docsOutputDir = p
	}

	switch docsFormat {
	case "markdown":
		err := doc.GenMarkdownTree(cmd.Parent().Root(), docsOutputDir)
		if err != nil {
			return err
		}
	case "manpage":
		manHeader := &doc.GenManHeader{
			Title:   "Tanzu Community Edition",
			Section: "1",
		}
		err := doc.GenManTree(cmd.Parent().Root(), manHeader, docsOutputDir)
		if err != nil {
			return err
		}
	case "rest":
		err := doc.GenReSTTree(cmd.Parent().Root(), docsOutputDir)
		if err != nil {
			return err
		}
	case "yaml":
		err := doc.GenYamlTree(cmd.Parent().Root(), docsOutputDir)
		if err != nil {
			return err
		}
	}

	return nil

}

func docsToStdout(cmd *cobra.Command) error {
	buf := new(bytes.Buffer)

	switch docsFormat {
	case "markdown":
		if err := doc.GenMarkdown(cmd.Root(), buf); err != nil {
			return err
		}

		for _, c := range cmd.Root().Commands() {
			if err := doc.GenMarkdown(c, buf); err != nil {
				return err
			}
		}
	case "manpage":
		manHeader := &doc.GenManHeader{
			Title:   "Tanzu Community Edition",
			Section: "1",
		}

		if err := doc.GenMan(cmd.Root(), manHeader, buf); err != nil {
			return err
		}

		for _, c := range cmd.Root().Commands() {
			if err := doc.GenMan(c, manHeader, buf); err != nil {
				return err
			}
		}
	case "rest":
		if err := doc.GenReST(cmd.Root(), buf); err != nil {
			return err
		}

		if err := doc.GenReST(cmd.Root(), buf); err != nil {
			return err
		}

		for _, c := range cmd.Root().Commands() {
			if err := doc.GenReST(c, buf); err != nil {
				return err
			}
		}
	case "yaml":
		if err := doc.GenYaml(cmd.Root(), buf); err != nil {
			return err
		}
		for _, c := range cmd.Root().Commands() {
			if err := doc.GenYaml(c, buf); err != nil {
				return err
			}
		}
	}

	fmt.Print(buf)

	return nil
}
