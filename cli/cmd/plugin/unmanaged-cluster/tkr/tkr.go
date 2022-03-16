// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package tkr contains functions for working with Tanzu Kubernetes Release information.
package tkr

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ImagePackage represents information for an image.
type ImagePackage struct {
	ImagePath  string `yaml:"imagePath"`
	Tag        string `yaml:"tag"`
	Repository string `yaml:"repository"`
}

// PackageData contains metadata about a package.
type PackageData struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Arch    string `yaml:"arch"`
}

// PackageInfo contains information about a package.
type PackageInfo struct {
	Category     string   `yaml:"category"`
	ClusterTypes []string `yaml:"clusterTypes"`
	PackageName  string   `yaml:"packageName"`
	Repository   string   `yaml:"repository"`
}

type Bom struct {
	APIVersion string `yaml:"apiVersion"`
	Release    struct {
		Version string `yaml:"version"`
	} `yaml:"release"`
	Components struct {
		AkoOperator []struct {
			Version string `yaml:"version"`
			Images  struct {
				AkoOperatorImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"akoOperatorImage"`
			} `yaml:"images"`
		} `yaml:"ako-operator"`
		Antrea []struct {
			Version string `yaml:"version"`
			Images  struct {
				AntreaImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"antreaImage"`
			} `yaml:"images"`
		} `yaml:"antrea"`
		CalicoAll []struct {
			Version string `yaml:"version"`
			Images  struct {
				CalicoCniImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"calicoCniImage"`
				CalicoKubecontrollerImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"calicoKubecontrollerImage"`
				CalicoNodeImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"calicoNodeImage"`
				CalicoPodDaemonImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"calicoPodDaemonImage"`
			} `yaml:"images"`
		} `yaml:"calico_all"`
		CloudProviderVsphere []struct {
			Version string `yaml:"version"`
			Images  struct {
				CcmControllerImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"ccmControllerImage"`
			} `yaml:"images"`
		} `yaml:"cloud_provider_vsphere"`
		CniPlugins []struct {
			Version string `yaml:"version"`
		} `yaml:"cni_plugins"`
		Containerd []struct {
			Version string `yaml:"version"`
		} `yaml:"containerd"`
		Coredns []struct {
			Version string `yaml:"version"`
			Images  struct {
				Coredns struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"coredns"`
			} `yaml:"images"`
		} `yaml:"coredns"`
		CriTools []struct {
			Version string `yaml:"version"`
		} `yaml:"cri_tools"`
		CsiAttacher []struct {
			Version string `yaml:"version"`
			Images  struct {
				CsiAttacherImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"csiAttacherImage"`
			} `yaml:"images"`
		} `yaml:"csi_attacher"`
		CsiLivenessprobe []struct {
			Version string `yaml:"version"`
			Images  struct {
				CsiLivenessProbeImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"csiLivenessProbeImage"`
			} `yaml:"images"`
		} `yaml:"csi_livenessprobe"`
		CsiNodeDriverRegistrar []struct {
			Version string `yaml:"version"`
			Images  struct {
				CsiNodeDriverRegistrarImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"csiNodeDriverRegistrarImage"`
			} `yaml:"images"`
		} `yaml:"csi_node_driver_registrar"`
		CsiProvisioner []struct {
			Version string `yaml:"version"`
			Images  struct {
				CsiProvisonerImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"csiProvisonerImage"`
			} `yaml:"images"`
		} `yaml:"csi_provisioner"`
		Dex []struct {
			Version string `yaml:"version"`
			Images  struct {
				DexImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"dexImage"`
			} `yaml:"images"`
		} `yaml:"dex"`
		Etcd []struct {
			Version string `yaml:"version"`
			Images  struct {
				Etcd struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"etcd"`
			} `yaml:"images"`
		} `yaml:"etcd"`
		KappController []struct {
			Version string `yaml:"version"`
			Images  struct {
				KappControllerImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"kappControllerImage"`
			} `yaml:"images"`
		} `yaml:"kapp-controller"`
		Kubernetes []struct {
			Version string `yaml:"version"`
			Images  struct {
				KubeAPIServer struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"kubeAPIServer"`
				KubeControllerManager struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"kubeControllerManager"`
				KubeE2E struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"kubeE2e"`
				KubeProxy struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"kubeProxy"`
				KubeScheduler struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"kubeScheduler"`
				Pause struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"pause"`
				PauseWindows1809 struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"pause_windows_1809"`
				PauseWindows1903 struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"pause_windows_1903"`
				PauseWindows1909 struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"pause_windows_1909"`
				PauseWindows2004 struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"pause_windows_2004"`
			} `yaml:"images"`
		} `yaml:"kubernetes"`
		KubernetesCsiExternalResizer []struct {
			Version string `yaml:"version"`
			Images  struct {
				CsiExternalResizer struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"csiExternalResizer"`
			} `yaml:"images"`
		} `yaml:"kubernetes-csi_external-resizer"`
		KubernetesSigsKind []struct {
			Version string `yaml:"version"`
			Images  struct {
				KindNodeImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"kindNodeImage"`
			} `yaml:"images"`
		} `yaml:"kubernetes-sigs_kind"`
		LoadBalancerAndIngressService []struct {
			Version string `yaml:"version"`
			Images  struct {
				LoadBalancerAndIngressServiceImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"loadBalancerAndIngressServiceImage"`
			} `yaml:"images"`
		} `yaml:"load-balancer-and-ingress-service"`
		MetricsServer []struct {
			Version string `yaml:"version"`
			Images  struct {
				MetricsServerImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"metricsServerImage"`
			} `yaml:"images"`
		} `yaml:"metrics-server"`
		Pinniped []struct {
			Version string `yaml:"version"`
			Images  struct {
				PinnipedImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"pinnipedImage"`
			} `yaml:"images"`
		} `yaml:"pinniped"`
		TanzuFrameworkAddons []struct {
			Version string `yaml:"version"`
			Images  struct {
				TanzuAddonsManagerImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"tanzuAddonsManagerImage"`
				TkgPinnipedPostDeployImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"tkgPinnipedPostDeployImage"`
			} `yaml:"images"`
		} `yaml:"tanzu-framework-addons"`
		TkgCorePackages []struct {
			Version string `yaml:"version"`
			Images  struct {
				AddonsManagerTanzuVmwareCom struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"addons-manager.tanzu.vmware.com"`
				AkoOperatorTanzuVmwareCom struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"ako-operator.tanzu.vmware.com"`
				AntreaTanzuVmwareCom struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"antrea.tanzu.vmware.com"`
				CalicoTanzuVmwareCom struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"calico.tanzu.vmware.com"`
				KappControllerTanzuVmwareCom struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"kapp-controller.tanzu.vmware.com"`
				LoadBalancerAndIngressServiceTanzuVmwareCom struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"load-balancer-and-ingress-service.tanzu.vmware.com"`
				MetricsServerTanzuVmwareCom struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"metrics-server.tanzu.vmware.com"`
				PinnipedTanzuVmwareCom struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"pinniped.tanzu.vmware.com"`
				TanzuCorePackageRepositoryImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"tanzuCorePackageRepositoryImage"`
				VsphereCpiTanzuVmwareCom struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"vsphere-cpi.tanzu.vmware.com"`
				VsphereCsiTanzuVmwareCom struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"vsphere-csi.tanzu.vmware.com"`
			} `yaml:"images"`
		} `yaml:"tkg-core-packages"`
		VsphereCsiDriver []struct {
			Version string `yaml:"version"`
			Images  struct {
				CsiControllerImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"csiControllerImage"`
				CsiMetaDataSyncerImage struct {
					ImagePackage `yaml:",inline"`
				} `yaml:"csiMetaDataSyncerImage"`
			} `yaml:"images"`
		} `yaml:"vsphere_csi_driver"`
	} `yaml:"components"`
	KubeadmConfigSpec struct {
		APIVersion        string `yaml:"apiVersion"`
		Kind              string `yaml:"kind"`
		ImageRepository   string `yaml:"imageRepository"`
		KubernetesVersion string `yaml:"kubernetesVersion"`
		Etcd              struct {
			Local struct {
				DataDir         string `yaml:"dataDir"`
				ImageRepository string `yaml:"imageRepository"`
				ImageTag        string `yaml:"imageTag"`
			} `yaml:"local"`
		} `yaml:"etcd"`
		DNS struct {
			Type            string `yaml:"type"`
			ImageRepository string `yaml:"imageRepository"`
			ImageTag        string `yaml:"imageTag"`
		} `yaml:"dns"`
	} `yaml:"kubeadmConfigSpec"`
	Ova []struct {
		Name   string `yaml:"name"`
		Osinfo struct {
			PackageData `yaml:",inline"`
		} `yaml:"osinfo"`
		Version string `yaml:"version"`
	} `yaml:"ova"`
	Ami struct {
		ApNortheast1 []struct {
			ID     string `yaml:"id"`
			Osinfo struct {
				PackageData `yaml:",inline"`
			} `yaml:"osinfo"`
		} `yaml:"ap-northeast-1"`
		ApNortheast2 []struct {
			ID     string `yaml:"id"`
			Osinfo struct {
				PackageData `yaml:",inline"`
			} `yaml:"osinfo"`
		} `yaml:"ap-northeast-2"`
		ApSouth1 []struct {
			ID     string `yaml:"id"`
			Osinfo struct {
				PackageData `yaml:",inline"`
			} `yaml:"osinfo"`
		} `yaml:"ap-south-1"`
		ApSoutheast1 []struct {
			ID     string `yaml:"id"`
			Osinfo struct {
				PackageData `yaml:",inline"`
			} `yaml:"osinfo"`
		} `yaml:"ap-southeast-1"`
		ApSoutheast2 []struct {
			ID     string `yaml:"id"`
			Osinfo struct {
				PackageData `yaml:",inline"`
			} `yaml:"osinfo"`
		} `yaml:"ap-southeast-2"`
		EuCentral1 []struct {
			ID     string `yaml:"id"`
			Osinfo struct {
				PackageData `yaml:",inline"`
			} `yaml:"osinfo"`
		} `yaml:"eu-central-1"`
		EuWest1 []struct {
			ID     string `yaml:"id"`
			Osinfo struct {
				PackageData `yaml:",inline"`
			} `yaml:"osinfo"`
		} `yaml:"eu-west-1"`
		EuWest2 []struct {
			ID     string `yaml:"id"`
			Osinfo struct {
				PackageData `yaml:",inline"`
			} `yaml:"osinfo"`
		} `yaml:"eu-west-2"`
		EuWest3 []struct {
			ID     string `yaml:"id"`
			Osinfo struct {
				PackageData `yaml:",inline"`
			} `yaml:"osinfo"`
		} `yaml:"eu-west-3"`
		SaEast1 []struct {
			ID     string `yaml:"id"`
			Osinfo struct {
				PackageData `yaml:",inline"`
			} `yaml:"osinfo"`
		} `yaml:"sa-east-1"`
		UsEast1 []struct {
			ID     string `yaml:"id"`
			Osinfo struct {
				PackageData `yaml:",inline"`
			} `yaml:"osinfo"`
		} `yaml:"us-east-1"`
		UsEast2 []struct {
			ID     string `yaml:"id"`
			Osinfo struct {
				PackageData `yaml:",inline"`
			} `yaml:"osinfo"`
		} `yaml:"us-east-2"`
		UsGovEast1 []struct {
			ID     string `yaml:"id"`
			Osinfo struct {
				PackageData `yaml:",inline"`
			} `yaml:"osinfo"`
		} `yaml:"us-gov-east-1"`
		UsGovWest1 []struct {
			ID     string `yaml:"id"`
			Osinfo struct {
				PackageData `yaml:",inline"`
			} `yaml:"osinfo"`
		} `yaml:"us-gov-west-1"`
		UsWest2 []struct {
			ID     string `yaml:"id"`
			Osinfo struct {
				PackageData `yaml:",inline"`
			} `yaml:"osinfo"`
		} `yaml:"us-west-2"`
	} `yaml:"ami"`
	Azure []struct {
		Sku             string `yaml:"sku"`
		Publisher       string `yaml:"publisher"`
		Offer           string `yaml:"offer"`
		Version         string `yaml:"version"`
		ThirdPartyImage bool   `yaml:"thirdPartyImage"`
		Osinfo          struct {
			PackageData `yaml:",inline"`
		} `yaml:"osinfo"`
	} `yaml:"azure"`
	ImageConfig struct {
		ImageRepository string `yaml:"imageRepository"`
	} `yaml:"imageConfig"`
	Addons struct {
		AkoOperator struct {
			PackageInfo `yaml:",inline"`
		} `yaml:"ako-operator"`
		Antrea struct {
			PackageInfo `yaml:",inline"`
		} `yaml:"antrea"`
		Calico struct {
			PackageInfo `yaml:",inline"`
		} `yaml:"calico"`
		KappController struct {
			PackageInfo `yaml:",inline"`
		} `yaml:"kapp-controller"`
		LoadBalancerAndIngressService struct {
			PackageInfo `yaml:",inline"`
		} `yaml:"load-balancer-and-ingress-service"`
		MetricsServer struct {
			PackageInfo `yaml:",inline"`
		} `yaml:"metrics-server"`
		Pinniped struct {
			PackageInfo `yaml:",inline"`
		} `yaml:"pinniped"`
		TanzuAddonsManager struct {
			PackageInfo `yaml:",inline"`
		} `yaml:"tanzu-addons-manager"`
		VsphereCpi struct {
			PackageInfo `yaml:",inline"`
		} `yaml:"vsphere-cpi"`
		VsphereCsi struct {
			PackageInfo `yaml:",inline"`
		} `yaml:"vsphere-csi"`
	} `yaml:"addons"`
}

func ReadTKRBom(filePath string) (*Bom, error) {
	bom := &Bom{}
	rawBom, err := os.ReadFile(filePath)
	if err != nil {
		return bom, err
	}
	err = yaml.Unmarshal(rawBom, bom)
	if err != nil {
		return bom, err
	}

	return bom, nil
}

func (tkr *Bom) getTKRRegistry() string {
	return tkr.ImageConfig.ImageRepository
}

func (tkr *Bom) GetTKRNodeImage() string {
	repo := tkr.getTKRNodeRepository()
	path := tkr.Components.KubernetesSigsKind[0].Images.KindNodeImage.ImagePath
	tag := tkr.Components.KubernetesSigsKind[0].Images.KindNodeImage.Tag

	return fmt.Sprintf("%s/%s:%s", repo, path, tag)
}

func (tkr *Bom) GetTKRCoreRepoBundlePath() string {
	registry := tkr.getTKRRegistry()
	path := tkr.Components.TkgCorePackages[0].Images.TanzuCorePackageRepositoryImage.ImagePath
	tag := tkr.Components.TkgCorePackages[0].Images.TanzuCorePackageRepositoryImage.Tag

	return fmt.Sprintf("%s/%s:%s", registry, path, tag)
}

func (tkr *Bom) GetTKRKappImage() (ImageReader, error) {
	registry := tkr.getTKRKappRepository()
	path := tkr.Components.TkgCorePackages[0].Images.KappControllerTanzuVmwareCom.ImagePath
	tag := tkr.Components.TkgCorePackages[0].Images.KappControllerTanzuVmwareCom.Tag

	t, err := NewTkrImageReader(fmt.Sprintf("%s/%s:%s", registry, path, tag))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (tkr *Bom) getTKRNodeRepository() string {
	if tkr.Components.KubernetesSigsKind[0].Images.KindNodeImage.Repository == "" {
		return tkr.getTKRRegistry()
	}

	return tkr.Components.KubernetesSigsKind[0].Images.KindNodeImage.Repository
}

func (tkr *Bom) getTKRKappRepository() string {
	if tkr.Components.TkgCorePackages[0].Images.KappControllerTanzuVmwareCom.Repository == "" {
		return tkr.getTKRRegistry()
	}

	return tkr.Components.TkgCorePackages[0].Images.KappControllerTanzuVmwareCom.Repository
}
