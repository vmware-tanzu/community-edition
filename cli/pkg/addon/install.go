// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	"os"

	"github.com/spf13/cobra"
	klog "k8s.io/klog/v2"

	"github.com/vmware-tanzu/tce/cli/pkg/utils"
)

// InstallCmd represents the tanzu package install command. It receives an package name
// and (optional) version. It then looks up the corresponding Package CR to verify
// there is something to install. If the corresponding Package CR resolves, an
// InstalledPacakge CR is create and deployed into the cluster.
var InstallCmd = &cobra.Command{
	Use:   "install <package name>",
	Short: "Install a package into the cluster",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		mgr, err = NewManager()
		return err
	},
	RunE: install,
	Args: cobra.ExactArgs(1),
}

func init() {
	// common between secret and user-defined
	InstallCmd.Flags().StringVarP(&inputAppCrd.Namespace, "namespace", "n", "default", "Namespace to deploy too")

	// secret
	InstallCmd.Flags().StringVarP(&inputAppCrd.ClusterName, "cluster", "c", "", "Cluster name which corresponds to a secret")

	// user defined
	InstallCmd.Flags().StringVarP(&inputAppCrd.URL, "url", "u", "", "URL to image")
	InstallCmd.Flags().StringToStringVarP(&inputAppCrd.Paths, "paths", "p", nil, "User defined paths for kapp template")
	InstallCmd.Flags().StringVarP(&inputAppCrd.Version, "package-version", "o", "", "Version of the package")
	InstallCmd.Flags().StringVarP(&inputAppCrd.ConfigPath, "config", "g", "", "Configuration for the package")
}

func install(cmd *cobra.Command, args []string) error {
	inputAppCrd.Name = args[0]
	klog.V(6).Infof("package name: %s", inputAppCrd.Name)

	// find the Package CR that corresponds to the name and/or version
	cmd.Printf("Looking up package to install: %s:%s\n", inputAppCrd.Name, inputAppCrd.Version)
	pkg, err := mgr.kapp.ResolvePackage(inputAppCrd.Name, inputAppCrd.Version)
	if err != nil {
		return utils.NonUsageError(cmd, err, "unable to resolve package '%s:%s' in namespace '%s'.", inputAppCrd.Name, inputAppCrd.Version, inputAppCrd.Namespace)
	}
	klog.V(6).Infoln(pkg)

	// if the user didn't specify a version, use the version from the resolved package
	if inputAppCrd.Version == "" {
		inputAppCrd.Version = mgr.kapp.ResolvePackageVersion(pkg)
	}

	// if the user specifies a configuration file, load it
	// for later use in the install.
	if inputAppCrd.ConfigPath != "" {
		inputAppCrd.Config, err = os.ReadFile(inputAppCrd.ConfigPath)
		if err != nil {
			return utils.NonUsageError(cmd, err, "package config path '%s' could not be found.", inputAppCrd.ConfigPath)
		}
	}

	// create InstalledPackage CR
	err = mgr.kapp.InstallPackage(inputAppCrd)
	if err != nil {
		return utils.NonUsageError(cmd, err, "error installing package")
	}
	cmd.Printf("Installed package in %s/%s:%s\n", inputAppCrd.Namespace, inputAppCrd.Name, inputAppCrd.Version)

	return nil
}
