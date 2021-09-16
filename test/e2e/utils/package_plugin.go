// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utils

func CheckPackageRepositoryList() (string, error) {
	return Tanzu(nil, "package", "repository", "list")
}

func CheckPackageRepositoryListAllNamespaces() (string, error) {
	return Tanzu(nil, "package", "repository", "list", "-A")
}

func CheckPackageAvailableList(namespace string) (string, error) {
	return Tanzu(nil, "package", "available", "list", "-n", namespace)
}

func CheckPackageAvailableListAllNamespaces() (string, error) {
	return Tanzu(nil, "package", "available", "list", "-A")
}

func PackageInstalledList(namespace string) (string, error) {
	return Tanzu(nil, "package", "installed", "list", "-n", namespace)
}

func PackageInstalledListAllNamespaces() (string, error) {
	return Tanzu(nil, "package", "installed", "list", "-A")
}
