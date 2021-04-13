// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	"fmt"

	"github.com/spf13/cobra"
	klog "k8s.io/klog/v2"

	"github.com/vmware-tanzu/tce/cli/utils"
)

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete <extension name>",
	Short: "Delete an installed package from the cluster, terminating it",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		mgr, err = NewManager()
		return err
	},
	RunE: deleteCmd,
}

func init() {
	// common between secret and user-defined
	DeleteCmd.Flags().StringVarP(&inputAppCrd.Namespace, "namespace", "n", "default", "Namespace to deploy too")

	// secret
	DeleteCmd.Flags().StringVarP(&inputAppCrd.ClusterName, "cluster", "c", "", "Cluster name which corresponds to a secret")

	// user defined
	DeleteCmd.Flags().StringVarP(&inputAppCrd.URL, "url", "u", "", "URL to image")
	DeleteCmd.Flags().StringToStringVarP(&inputAppCrd.Paths, "paths", "p", nil, "User defined paths for kapp template")
	DeleteCmd.Flags().StringVarP(&inputAppCrd.Version, "package-version", "o", "", "Version of the package")

	// delete force
	DeleteCmd.Flags().BoolVarP(&inputAppCrd.Force, "force", "f", false, "Force delete")
	DeleteCmd.Flags().BoolVarP(&inputAppCrd.Teardown, "teardown", "t", false, "Delete associated ServiceAccount and RoleBinding")
}

func deleteCmd(cmd *cobra.Command, args []string) error {
	// validate a package name was passed
	if len(args) < 1 {
		return ErrMissingPackageName
	}
	inputAppCrd.Name = args[0]
	klog.V(6).Infof("package name: %s", inputAppCrd.Name)

	// lookup an installed package based on the user-provided input. If one cannot be
	// found, return an error as there is nothing to do.
	ipkg, err := mgr.kapp.ResolveInstalledPackage(inputAppCrd.Name, inputAppCrd.Version, inputAppCrd.Namespace)
	if err != nil {
		return utils.NonUsageError(cmd, err, "unable to resolve package '%s:%s' in namespace '%s'.", inputAppCrd.Name, inputAppCrd.Version, inputAppCrd.Namespace)
	}

	// if the version was left empty, fill it with the resolved installedPackage
	if inputAppCrd.Version == "" {
		inputAppCrd.Version = ipkg.Spec.PkgRef.VersionSelection.Constraints
	}
	fmt.Printf("Attempting to delete %s/%s:%s\n", inputAppCrd.Namespace, inputAppCrd.Name, inputAppCrd.Version)

	// if a config secret is referenced in the installedpackage, set its name in configFile
	// so it can be deleted
	if ipkg.Spec.Values != nil {
		inputAppCrd.ConfigPath = ipkg.Spec.Values[0].SecretRef.Name
	}

	err = mgr.kapp.DeletePackage(inputAppCrd)
	if err != nil {
		return utils.NonUsageError(cmd, err, "error deleting package '%s'.", inputAppCrd.Name)
	}

	fmt.Printf("Deleted %s/%s:%s\n", inputAppCrd.Namespace, inputAppCrd.Name, inputAppCrd.Version)
	return nil
}
