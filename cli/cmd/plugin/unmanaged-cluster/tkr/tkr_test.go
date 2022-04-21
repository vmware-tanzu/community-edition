// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package tkr

import (
	"io/ioutil"
	"os"
	"testing"
)

var (
	none   = "none"
	tceReg = "projects.registry.vmware.com/tce"
	tkgReg = "projects.registry.vmware.com/tkg"
)

// This is a sample chunk of a TKR from the v0.12.0 bom.
// Since we don't use many of the fields in the bom,
// we can ommit them for testing purposes.
var sampleTkrYamls = `---
apiVersion: run.tanzu.vmware.com/v1alpha2
release:
  version: v1.22.7
components:
  kapp-controller:
  - version: v0.30.1+vmware.1
    images:
      kappControllerImage:
        imagePath: kapp-controller
        tag: v0.30.1_vmware.1
  kubernetes-sigs_kind:
  - version: v1.22.7
    images:
      kindNodeImage:
        imagePath: kind
        tag: v1.22.7
        repository: projects.registry.vmware.com/tce
  tkg-core-packages:
  - version: v1.22.8+vmware.1-tkg.1-tf-v0.11.4
    images:
      kapp-controller.tanzu.vmware.com:
        imagePath: kapp-controller-multi-pkg
        tag: v0.30.1
        repository: projects.registry.vmware.com/tce
      tanzuCorePackageRepositoryImage:
        imagePath: repo-12
        tag: 0.12.0
        repository: projects.registry.vmware.com/tce
      tanzuUserPackageRepositoryImage:
        imagePath: main
        repository: projects.registry.vmware.com/tce
        tag: 0.12.0
imageConfig:
  imageRepository: projects.registry.vmware.com/tkg`

func helperMakeNewBom() (*Bom, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "tkr-test-")
	if err != nil {
		return nil, err
	}

	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(sampleTkrYamls))
	if err != nil {
		return nil, err
	}

	var bom *Bom
	bom, err = ReadTKRBom(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	return bom, nil
}

func TestReadTKRBom(t *testing.T) {
	_, err := helperMakeNewBom()
	if err != nil {
		t.Errorf("expected reading TKR bom to be successful. Error: %s", err.Error())
	}
}

func TestReadTKRBomFails(t *testing.T) {
	_, err := ReadTKRBom("file-doesnt-exist")
	if err == nil {
		t.Error("expected reading TKr file to fail if file doesn't exist")
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "tkr-test-")
	if err != nil {
		t.Errorf(err.Error())
	}

	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(" ' this is bad yamls"))
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = ReadTKRBom(tmpFile.Name())
	if err == nil {
		t.Error("expected unmarshaling TKr file to fail if invalid yaml is provided")
	}
}

func TestGetTKRRegistry(t *testing.T) {
	b, err := helperMakeNewBom()
	if err != nil {
		t.Errorf("expected reading TKR bom to be successful. Error: %s", err.Error())
	}

	repo := b.getTKRRegistry()
	if repo != tkgReg {
		t.Errorf("expected TKR registry to be projects.registry.vmware.com/tkg. Actual: %s", repo)
	}
}

func TestGetTKRNodeImageKind(t *testing.T) {
	b, err := helperMakeNewBom()
	if err != nil {
		t.Errorf("expected reading TKR bom to be successful. Error: %s", err.Error())
	}

	nodeImg := b.GetTKRNodeImage("kind")
	if nodeImg != "projects.registry.vmware.com/tce/kind:v1.22.7" {
		t.Errorf("expected node for kind image to be projects.registry.vmware.com/tce/kind:v1.22.7. Actual: %s", nodeImg)
	}
}

func TestGetTKRNodeImageMinikube(t *testing.T) {
	b, err := helperMakeNewBom()
	if err != nil {
		t.Errorf("expected reading TKR bom to be successful. Error: %s", err.Error())
	}

	nodeImg := b.GetTKRNodeImage("minikube")
	if nodeImg != "v1.22.7" {
		t.Errorf("expected node for minikube image to be v1.22.7. Actual: %s", nodeImg)
	}
}

func TestGetTKRNodeImageNone(t *testing.T) {
	b, err := helperMakeNewBom()
	if err != nil {
		t.Errorf("expected reading TKR bom to be successful. Error: %s", err.Error())
	}

	nodeImg := b.GetTKRNodeImage(none)
	if nodeImg != none {
		t.Errorf("expected node for noop provider image to be `none`. Actual: %s", nodeImg)
	}
}

func TestGetTKRNodeImageUnsupported(t *testing.T) {
	b, err := helperMakeNewBom()
	if err != nil {
		t.Errorf("expected reading TKR bom to be successful. Error: %s", err.Error())
	}

	nodeImg := b.GetTKRNodeImage("bad-provider")
	if nodeImg != "" {
		t.Errorf("expected node for unsupported provider image to be an empty string. Actual: %s", nodeImg)
	}
}

func TestGetTKRCoreRepoBundlePath(t *testing.T) {
	b, err := helperMakeNewBom()
	if err != nil {
		t.Errorf("expected reading TKR bom to be successful. Error: %s", err.Error())
	}

	coreRepo := b.GetTKRCoreRepoBundlePath()
	if coreRepo != "projects.registry.vmware.com/tce/repo-12:0.12.0" {
		t.Errorf("expected core repo path to be projects.registry.vmware.com/tce/repo-12:0.12.0. Actual: %s", coreRepo)
	}
}

func TestGetTKRCoreRepoBundlePathDefaultRegistry(t *testing.T) {
	b, err := helperMakeNewBom()
	if err != nil {
		t.Errorf("expected reading TKR bom to be successful. Error: %s", err.Error())
	}

	b.Components.TkgCorePackages[0].Images.TanzuCorePackageRepositoryImage.Repository = ""

	coreRepo := b.GetTKRCoreRepoBundlePath()
	if coreRepo != tkgReg+"/repo-12:0.12.0" {
		t.Errorf("expected empty core repo path to default to `projects.registry.vmware.com/tkg/repo-12:0.12.0`. Actual: %s", coreRepo)
	}
}

func TestGetTKRUserRepoBundlePath(t *testing.T) {
	b, err := helperMakeNewBom()
	if err != nil {
		t.Errorf("expected reading TKR bom to be successful. Error: %s", err.Error())
	}

	userRepo := b.GetTKRUserRepoBundlePath()
	if userRepo != "projects.registry.vmware.com/tce/main:0.12.0" {
		t.Errorf("expected user repo path to be projects.registry.vmware.com/tce/main:0.12.0. Actual: %s", userRepo)
	}
}

func TestGetTKRUserRepoBundlePathDefaultRegistry(t *testing.T) {
	b, err := helperMakeNewBom()
	if err != nil {
		t.Errorf("expected reading TKR bom to be successful. Error: %s", err.Error())
	}

	b.Components.TkgCorePackages[0].Images.TanzuUserPackageRepositoryImage.Repository = ""

	userRepo := b.GetTKRUserRepoBundlePath()
	if userRepo != tkgReg+"/main:0.12.0" {
		t.Errorf("expected empty core repo path to default to `projects.registry.vmware.com/tkg/main:0.12.0`. Actual: %s", userRepo)
	}
}

func TestGetTKRUserRepoBundlePathFail(t *testing.T) {
	b, err := helperMakeNewBom()
	if err != nil {
		t.Errorf("expected reading TKR bom to be successful. Error: %s", err.Error())
	}

	b.Components.TkgCorePackages[0].Images.TanzuUserPackageRepositoryImage.ImagePath = ""

	userRepo := b.GetTKRUserRepoBundlePath()
	if userRepo != "" {
		t.Errorf("expected empty user repo image path to result in empty full path. Actual: %s", userRepo)
	}

	b.Components.TkgCorePackages[0].Images.TanzuUserPackageRepositoryImage.Tag = ""

	userRepo = b.GetTKRUserRepoBundlePath()
	if userRepo != "" {
		t.Errorf("expected empty user repo tag to result in empty full path. Actual: %s", userRepo)
	}
}

func TestGetTKRKappImage(t *testing.T) {
	b, err := helperMakeNewBom()
	if err != nil {
		t.Errorf("expected reading TKR bom to be successful. Error: %s", err.Error())
	}

	kappImg, err := b.GetTKRKappImage()
	if err != nil {
		t.Errorf("expected building kapp Image reader success. Error: %s", err.Error())
	}

	if kappImg == nil {
		t.Errorf("Expected Kapp image reader from bom to not be nil. Was nil")
	}
}

func TestGetTKRNodeRepository(t *testing.T) {
	b, err := helperMakeNewBom()
	if err != nil {
		t.Errorf("expected reading TKR bom to be successful. Error: %s", err.Error())
	}

	nodeRepo := b.getTKRNodeRepository()
	if nodeRepo != tceReg {
		t.Errorf("expected node repo path to be projects.registry.vmware.com/tce. Actual: %s", nodeRepo)
	}
}

func TestGetTKRNodeRepositoryDefaultRegistry(t *testing.T) {
	b, err := helperMakeNewBom()
	if err != nil {
		t.Errorf("expected reading TKR bom to be successful. Error: %s", err.Error())
	}

	b.Components.KubernetesSigsKind[0].Images.KindNodeImage.Repository = ""

	nodeRepo := b.getTKRNodeRepository()
	if nodeRepo != tkgReg {
		t.Errorf("expected empty node repo path to default to `projects.registry.vmware.com/tkg`. Actual: %s", nodeRepo)
	}
}

func TestGetTKRKappRepository(t *testing.T) {
	b, err := helperMakeNewBom()
	if err != nil {
		t.Errorf("expected reading TKR bom to be successful. Error: %s", err.Error())
	}

	kappRepo := b.getTKRKappRepository()
	if kappRepo != tceReg {
		t.Errorf("expected node repo path to be projects.registry.vmware.com/tce. Actual: %s", kappRepo)
	}
}

func TestGetTKRKappRepositoryDefaultRegistry(t *testing.T) {
	b, err := helperMakeNewBom()
	if err != nil {
		t.Errorf("expected reading TKR bom to be successful. Error: %s", err.Error())
	}

	b.Components.TkgCorePackages[0].Images.KappControllerTanzuVmwareCom.Repository = ""

	kappRepo := b.getTKRKappRepository()
	if kappRepo != tkgReg {
		t.Errorf("expected empty kapp repo path to default to `projects.registry.vmware.com/tkg`. Actual: %s", kappRepo)
	}
}
