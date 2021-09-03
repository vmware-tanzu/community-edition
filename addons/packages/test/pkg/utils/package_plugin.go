// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"encoding/json"

	"github.com/onsi/gomega"
)

func TanzuPackageName(displayName string) string {
	packageVersionJSON, err := Tanzu(nil, "package", "available", "list", "-o", "json")
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	packages := []map[string]string{}

	err = json.Unmarshal([]byte(packageVersionJSON), &packages)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	var packageName string
	for _, pkg := range packages {
		if dn := pkg["display-name"]; dn == displayName {
			packageName = pkg["name"]
			break
		}
	}

	gomega.Expect(packageName).NotTo(gomega.BeEmpty())
	return packageName
}

func TanzuPackageAvailableVersion(packageName string) string {
	packageVersionJSON, err := Tanzu(nil, "package", "available", "list", packageName, "-o", "json")
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	versions := []map[string]string{}

	err = json.Unmarshal([]byte(packageVersionJSON), &versions)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(len(versions)).To(gomega.BeNumerically(">", 0))

	return versions[0]["version"]
}
