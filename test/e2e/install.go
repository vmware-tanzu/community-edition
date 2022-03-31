// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	u "github.com/vmware-tanzu/community-edition/test/e2e/utils"
)

func installTCE() error {
	if ConfigVal.TceVersion == "" {
		path, err := exec.LookPath("tanzu")
		if err == nil {
			log.Println("tanzu is already exist in the environment", path)
			return nil
		}

		log.Println("start building a tce release and install ....")
		// build make release and install
		err = os.Chdir(u.WorkingDir + "/../../")
		if err != nil {
			log.Println("error while changing directory :", err)
			return err
		}

		start := time.Now()
		runDeployScript("download-or-build-tce.sh", "")
		log.Println("Time taken for build and install TCE release from source code: ", time.Since(start))

		err = os.Chdir(u.WorkingDir)
		if err != nil {
			log.Println("error while changing directory :", err)
			return err
		}
	} else {
		log.Println("installing tce ....")
		log.Println("Provided tce release version to install is :", ConfigVal.TceVersion)
		// fetch tce release from github page and install
		start := time.Now()
		runDeployScript("fetch-tce.sh", ConfigVal.TceVersion)
		log.Println("Time taken for fetch release ", ConfigVal.TceVersion, " and install TCE: ", time.Since(start))
	}

	return nil
}

func installCluster() error {
	// Bring up the cluster
	log.Println("Provisioning cluster...")
	s := time.Now()
	err := DeployTanzuCluster()
	if err != nil {
		return err
	}

	log.Println("Total Time taken for bringing up "+ConfigVal.ClusterType+" cluster : ", time.Since(s))
	return nil
}

func installAWSCli() {
	// Installing AWS CLI
	log.Println("Installing AWS CLI...")
	runDeployScript("ensure-aws-cli.sh", "")
}

func runDeployScript(filename, releaseVersion string) {
	mwriter := io.MultiWriter(os.Stdout)
	cmd := exec.Command("/bin/sh", u.WorkingDir+"/../"+filename, releaseVersion) //nolint:gosec
	cmd.Stderr = mwriter
	cmd.Stdout = mwriter
	err := cmd.Run() // blocks until sub process is complete
	if err != nil {
		log.Fatal(err)
	}
}
