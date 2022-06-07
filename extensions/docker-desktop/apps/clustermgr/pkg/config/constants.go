// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

const (
	// DefaultClusterName is the name that will be given to a new cluster.
	DefaultClusterName string = "tanzu-community-edition"
	// UnmanagedClusterBinary is the full path to the unmanaged cluster CLI plugin.
	UnmanagedClusterBinary string = "/backend/tanzu-unmanaged-cluster"
	// YttBinary is the name of the ytt binary to call for ytt operations.
	YttBinary string = "ytt"
	// KubectlBinary is the name of the kubectl command.
	KubectlBinary string = "kubectl"
	// DefaultHomeKubeConfig is the full path to the kubeconfig.
	DefaultHomeKubeConfig string = "/home/tanzu/.kube/config"
	// DefaultHostMountKubeConfig
	DefaultHostMountedKubeConfig string = "/opt/kubeconfig/config"
	// ClusterInstallTemplateFile is the full path to the cluster template yaml file.
	ClusterInstallTemplateFile string = "/backend/templates/cluster-template.yaml"
	// ClusterInstallValuesFile is the full path to the cluster values file.
	ClusterInstallValuesFile string = "/backend/templates/cluster-values.yaml"
	// ClusterConfigFileName is the name of the file containing the cluster config.
	ClusterConfigFileName string = "cluster-config.yaml"
	// ClusterLogFile is the name of the file containing the cluster logs.
	ClusterLogFile string = "cluster.log"
)
