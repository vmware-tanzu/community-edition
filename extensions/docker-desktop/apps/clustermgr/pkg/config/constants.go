// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

const (
	DefaultClusterName              string = "tanzu-community-edition"
	KindTCEImage                    string = "projects.registry.vmware.com/tce/kind:v1.22.7"
	UnmanagedClusterBinary          string = "/backend/tanzu-unmanaged-cluster"
	YttBinary                       string = "ytt"
	KappBinary                      string = "kapp"
	KubectlBinary                   string = "kubectl"
	DefaultPackagesNamespace        string = "tce-packages"
	DefaultHomeKubeConfig           string = "/home/tanzu/.kube/config"
	DefaultHostMountedKubeConfig    string = "/opt/kubeconfig/config"
	ClusterInstallTemplateFile      string = "/backend/templates/cluster-template.yaml"
	ClusterInstallValuesFile        string = "/backend/templates/cluster-values.yaml"
	ClusterPostInstallFile          string = "/backend/templates/cluster-postinstall.yaml"
	IngressTemplatesDir             string = "/backend/templates/apps/ingress/"
	KubeappsTemplatesDir            string = "/backend/templates/apps/kubeapps/"
	KubeappsTargetNamespace         string = "kubeapps"
	GlobalPackagesNamespace         string = "tanzu-package-repo-global"
	ContourPackageName              string = "contour.community.tanzu.vmware.com.1.20.1"
	KubeappsapiDeployName           string = "deploy/kubeapps-internal-kubeappsapis"
	WaitForPackagesAvailableSeconds int    = 120
	ClusterLogFile                  string = "cluster.log"
	ClusterInternalLogFile          string = "cluster-internal.log"
)
