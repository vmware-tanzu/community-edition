// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package testdata

import (
	"fmt"

	u "github.com/vmware-tanzu/community-edition/test/e2e/utils"
)

func InstallMetallb() error {
	_, err := u.Kubectl(nil, "apply", "-f", u.WorkingDir+"/testdata/metal-lb/namespace.yaml")
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}

	_, err = u.Kubectl(nil, "apply", "-f", u.WorkingDir+"/testdata/metal-lb/metallb.yaml")
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}

	_, err = u.Kubectl(nil, "apply", "-f", u.WorkingDir+"/testdata/metal-lb/metallb_cm.yaml")
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}

	return nil
}
func InstallVelero(version string) error {

	isVeleroInstall, _ := u.Kubectl(nil, "get", "apps", "velero", "-o=jsonpath={..status.conditions[0].status}")
	if isVeleroInstall != "True" {
		_, err := u.Tanzu(nil, "package", "install", "velero", "--package-name", "velero.community.tanzu.vmware.com", "--version", version, "--values-file", u.WorkingDir+"/testdata/velero/velero_values.yaml")
		if err != nil {
			fmt.Printf("%s", err)
			return err
		}
	}

	return nil
}

func UnsinstallVelero() error {
	_, err := u.Tanzu(nil, "package", "installed", "delete", "velero", "-y")
	if err != nil {
		fmt.Printf("%s", err)
	}
	return err
}

func UninstallMetallb() error {
	_, err := u.Kubectl(nil, "delete", "-f", u.WorkingDir+"/testdata/metal-lb/namespace.yaml")
	if err != nil {
		fmt.Printf("%s", err)
	}
	return err
}
