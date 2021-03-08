// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	"fmt"

	"github.com/spf13/cobra"
	klog "k8s.io/klog/v2"
)

var clearToken bool

// TokenCmd represents the token command
var TokenCmd = &cobra.Command{
	Use:   "token <token>",
	Short: "GitHub token",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		mgr, err = NewManager()
		return err
	},
	RunE: token,
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func init() {
	TokenCmd.Flags().BoolVarP(&clearToken, "clear", "c", false, "Clear token")
}

func token(cmd *cobra.Command, args []string) error {

	var token string
	if !clearToken {
		if len(args) == 0 {
			fmt.Printf("Please provide GitHub token\n")
			return ErrMissingToken
		}
		token = args[0]
		klog.V(6).Infof("token = %s", token)
		if len(token) < TokenMinLength {
			fmt.Printf("Please provide GitHub token\n")
			return ErrMissingToken
		}
	} else {
		token = ""
	}

	err := mgr.cfg.UpdateToken(token)
	if err != nil {
		fmt.Printf("Failed to update GitHub token\n")
	}

	fmt.Printf("Updated GitHub token\n")
	return nil
}
