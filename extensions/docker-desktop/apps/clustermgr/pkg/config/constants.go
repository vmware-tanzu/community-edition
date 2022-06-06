// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

const (
	DefaultClusterName           string = "tanzu-community-edition"
	UnmanagedClusterBinary       string = "/backend/tanzu-unmanaged-cluster"
	YttBinary                    string = "ytt"
	KubectlBinary                string = "kubectl"
	DefaultHomeKubeConfig        string = "/home/tanzu/.kube/config"
	DefaultHostMountedKubeConfig string = "/opt/kubeconfig/config"
	ClusterInstallTemplateFile   string = "/backend/templates/cluster-template.yaml"
	ClusterInstallValuesFile     string = "/backend/templates/cluster-values.yaml"
	ClusterConfigFileName        string = "cluster-config.yaml"
	ClusterLogFile               string = "cluster.log"
)
