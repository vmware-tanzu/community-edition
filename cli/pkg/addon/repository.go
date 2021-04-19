// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	klog "k8s.io/klog/v2"

	"github.com/vmware-tanzu/tce/cli/utils"
)

var installDefault bool
var repoFilename string

// RepositoryCmd represents the repository command
var RepositoryCmd = &cobra.Command{
	Use:   "repository",
	Short: "Manage repositories for packages",
}

func init() {
	InstallRepoCmd.Flags().BoolVarP(&installDefault, "default", "d", false, "Install the default TCE repository")
	InstallRepoCmd.Flags().StringVarP(&repoFilename, "file", "f", "", "Install a repository based on a provided file")

	DeleteRepoCmd.Flags().StringVarP(&repoFilename, "file", "f", "", "Delete a repository based on a provided file")

	RepositoryCmd.AddCommand(InstallRepoCmd)
	RepositoryCmd.AddCommand(ListRepoCmd)
	RepositoryCmd.AddCommand(DeleteRepoCmd)
}

// InstallRepoCmd adds a repository.
var InstallRepoCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a package repository into the cluster",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		mgr, err = NewManager()
		return err
	},
	RunE: installRepository,
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func installRepository(cmd *cobra.Command, args []string) error {
	// install the default TCE repo
	if installDefault {
		klog.V(2).Infof("installDefault: %t", installDefault)
		err := mgr.kapp.InstallDefaultRepository()
		if err != nil {
			fmt.Printf("InstallDefaultRepository Failed. Err: %v\n", err)
			return err
		}
		fmt.Printf("Install repository succeeded\n")
		return nil
	}

	filename := strings.TrimSpace(repoFilename)
	if filename == "" {
		fmt.Printf("Missing repo name. Example: package repository install --file <filename>\n")
		return ErrMissingParameter
	}
	klog.V(2).Infof("filename: %s", filename)

	err := mgr.kapp.InstallRepositoryFromFile(filename)
	if err != nil {
		fmt.Printf("InstallRepository Failed. Err: %v\n", err)
		return err
	}

	fmt.Printf("Install repository succeeded\n")
	return nil
}

// ListRepoCmd lists all repositories
var ListRepoCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all repositories installed in the cluster",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		mgr, err = NewManager()
		return err
	},
	RunE: listRepository,
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func listRepository(cmd *cobra.Command, args []string) error {
	repos, err := mgr.kapp.ListRepositories()
	if err != nil {
		fmt.Printf("ListRepositories Failed. Err: %v\n", err)
		return err
	}

	writer := utils.NewTableWriter(cmd.OutOrStdout(), "NAME")
	for _, repo := range repos.Items {
		writer.AddRow(repo.ObjectMeta.Name)
	}
	writer.Render()

	return nil
}

// DeleteRepoCmd deletes a repository.
var DeleteRepoCmd = &cobra.Command{
	Use:   "delete [repo name] or --file <filename>",
	Short: "Deletes a repository of packages from the cluster",
	RunE:  deleteRepository,
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		mgr, err = NewManager()
		return err
	},
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func deleteRepository(cmd *cobra.Command, args []string) error {
	if repoFilename == "" {
		err := mgr.kapp.DeleteRepositoryFromFile(repoFilename)
		if err != nil {
			fmt.Printf("DeleteRepository Failed. Err: %v\n", err)
			return err
		}
		fmt.Printf("Delete repository succeeded\n")
		return nil
	}

	if len(args) == 0 {
		fmt.Printf("Missing repo name. Example: package repository delete <name>\n")
		return ErrMissingParameter
	}

	param := strings.TrimSpace(args[0])
	if param == "" {
		fmt.Printf("Missing repo name. Example: package repository delete <filename>\n")
		return ErrMissingParameter
	}
	klog.V(2).Infof("param: %s", param)

	err := mgr.kapp.DeleteRepository(param)
	if err != nil {
		fmt.Printf("DeleteRepository Failed. Err: %v\n", err)
		return err
	}

	fmt.Printf("Delete repository %s succeeded\n", param)
	return nil
}
