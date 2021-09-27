// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utils

func ListClusters() (string, error) {
	return Tanzu(nil, "cluster", "list")
}

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

func GetBootstrapClusterDiagnostics() (string, error) {
	return Tanzu(nil, "diagnostics", "collect", "--output-dir", "/tmp/")
}

func GetManagementClusterDiagnostics(clusterName string) (string, error) {
	return Tanzu(nil, "diagnostics", "collect", "--management-cluster-name", clusterName, "--output-dir", "/tmp/")
}

func GetWorkloadClusterDiagnostics(clusterName string) (string, error) {
	return Tanzu(nil, "diagnostics", "collect", "--workload-cluster-name", clusterName, "--output-dir", "/tmp/")
}
