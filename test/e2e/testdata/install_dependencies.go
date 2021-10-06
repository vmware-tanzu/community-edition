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

func UninstallMetallb() error {
	_, err := u.Kubectl(nil, "delete", "-f", u.WorkingDir+"/testdata/metal-lb/namespace.yaml")
	if err != nil {
		fmt.Printf("%s", err)
	}
	return err
}
