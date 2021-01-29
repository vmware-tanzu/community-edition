// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package extension

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var overrideResetPrompt bool

// ResetCmd represents the reset command
var ResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset extension configuration to 'factory'",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		mgr, err = NewManager()
		return err
	},
	RunE: reset,
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func init() {
	ResetCmd.Flags().BoolVarP(&overrideResetPrompt, "yes", "y", false, "Override the verification for deletion prompt")
}

func reset(cmd *cobra.Command, args []string) error {

	if !overrideResetPrompt {
		fmt.Println("Are you sure you want to reset the configuration to factory? [y/N]")
		var response string
		fmt.Scanln(&response)

		if !strings.EqualFold(response, "y") && !strings.EqualFold(response, "yes") {
			fmt.Println("Cancelling operation...")
			return nil
		}
	}

	err := mgr.b.Reset()
	if err != nil {
		fmt.Printf("BitBucket.Reset failed. Err: %v\n", err)
		return err
	}

	err = mgr.gh.Reset(mgr.kapp.GetWorkingDirectory())
	if err != nil {
		fmt.Printf("GitHub.Reset failed. Err: %v\n", err)
		return err
	}

	fmt.Printf("Extension plugin reset successful")
	return nil
}
