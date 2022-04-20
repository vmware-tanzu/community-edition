// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkr

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestImageReader(t *testing.T) {
	ir, err := NewTkrImageReader("path-to-image")
	if err != nil {
		t.Errorf(err.Error())
	}

	if ir == nil {
		t.Errorf("expected image struct to not be nil. Was nil")
	}

	if regURL := ir.GetRegistryURL(); regURL != "path-to-image" {
		t.Errorf("expected image path to be `path-to-image`. Actual: %s", regURL)
	}

	if path := ir.GetDownloadPath(); path == "" {
		t.Errorf("expected download path to be not be empty. Was empty")
	}
}

func TestGetTags(t *testing.T) {
	t.Skipf("imgpkg ListTags method libraries not mocked")
}

func TestDownloadBundleImage(t *testing.T) {
	t.Skipf("imgpkg bundle download libraries not mocked")
}

func TestDownloadImage(t *testing.T) {
	t.Skipf("imgpkg image download libraries not mocked")
}

var imgpkgImagesFileContent = `---
apiVersion: imgpkg.carvel.dev/v1alpha1
kind: ImagesLock
images:
- annotations:
    kbld.carvel.dev/id: busybox
  image: test-registry.io/library/busybox@sha256:caa382c432891547782ce7140fb3b7304613d3b0438834dce1cad68896ab110a`

var upstreamDeploymentFileContent = `---
apiVersion: v1
kind: Namespace
metadata:
  name: default-namespace-name

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: default-name
  namespace: default-namespace-name
spec:
  replicas: 1
  selector:
    matchLabels:
      app: default-app-name
  template:
    metadata:
      labels:
        app: default-app-name
    spec:
      containers:
        - name: default-container-name
          image: busybox
          imagePullPolicy: IfNotPresent
          command:
            - echo hello-world`

var overlayNamespaceFileContent = `#@ load("@ytt:data", "data")
#@ load("@ytt:overlay", "overlay")

#@overlay/match by=overlay.subset({"kind":"Namespace", "metadata": {"name": "default-namespace-name"}})
---
metadata:
  name: #@ data.values.namespace

#@overlay/match by=overlay.subset({"metadata":{"namespace": "default-namespace-name"}}), expects=1
---
metadata:
  namespace: #@ data.values.namespace`

var valuesFileContent = `#@data/values

---
namespace: custom-namespace-name`

var expectedBundle = `
---
apiVersion: v1
kind: Namespace
metadata:
  name: custom-namespace-name

---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    kbld.k14s.io/images: |
      - origins:
        - preresolved:
            url: test-registry.io/library/busybox@sha256:caa382c432891547782ce7140fb3b7304613d3b0438834dce1cad68896ab110a
        url: test-registry.io/library/busybox@sha256:caa382c432891547782ce7140fb3b7304613d3b0438834dce1cad68896ab110a
  name: default-name
  namespace: custom-namespace-name
spec:
  replicas: 1
  selector:
    matchLabels:
      app: default-app-name
  template:
    metadata:
      labels:
        app: default-app-name
    spec:
      containers:
      - command:
        - echo hello-world
        image: test-registry.io/library/busybox@sha256:caa382c432891547782ce7140fb3b7304613d3b0438834dce1cad68896ab110a
        imagePullPolicy: IfNotPresent
        name: default-container-name
`

func helperSetupImageLocally() (string, error) {
	// Simulates downloaded image directory
	imageDir, err := ioutil.TempDir("", "tmp-img-bundle-")
	if err != nil {
		log.Fatal(err)
	}

	imgpkgDir := filepath.Join(imageDir, ".imgpkg")
	configDir := filepath.Join(imageDir, "config")
	overlaysDir := filepath.Join(configDir, "overlays")
	upstreamDir := filepath.Join(configDir, "upstream")

	println(imageDir)

	err = os.Mkdir(imgpkgDir, 0755)
	if err != nil {
		return "", err
	}

	err = os.Mkdir(configDir, 0755)
	if err != nil {
		return "", err
	}

	err = os.Mkdir(overlaysDir, 0755)
	if err != nil {
		return "", err
	}

	err = os.Mkdir(upstreamDir, 0755)
	if err != nil {
		return "", err
	}

	// .imgpkg/images.yml file
	imgpkgImagesFile, err := os.Create(filepath.Join(imgpkgDir, "images.yml"))
	if err != nil {
		return "", err
	}

	_, err = imgpkgImagesFile.WriteString(imgpkgImagesFileContent)
	if err != nil {
		return "", err
	}

	// ./config/upstream/deployment.yml file
	deploymentFile, err := os.Create(filepath.Join(upstreamDir, "deployment.yml"))
	if err != nil {
		return "", err
	}

	_, err = deploymentFile.WriteString(upstreamDeploymentFileContent)
	if err != nil {
		return "", err
	}

	// ./config/overlays/overlay-namespace.yml file
	overlayNamespaceFile, err := os.Create(filepath.Join(overlaysDir, "overlay-namespace.yml"))
	if err != nil {
		return "", err
	}

	_, err = overlayNamespaceFile.WriteString(overlayNamespaceFileContent)
	if err != nil {
		return "", err
	}

	// ./config/values.yml file
	valuesFile, err := os.Create(filepath.Join(configDir, "values.yml"))
	if err != nil {
		return "", err
	}

	_, err = valuesFile.WriteString(valuesFileContent)
	if err != nil {
		return "", err
	}

	return imageDir, nil
}

func TestRenderYaml(t *testing.T) {
	p, err := helperSetupImageLocally()
	if err != nil {
		t.Errorf(err.Error())
	}

	ir, err := NewTkrImageReader(p)
	if err != nil {
		t.Errorf(err.Error())
	}

	ir.SetDownloadPath(p)
	ir.SetRelativeConfigPath("config")
	yamls, err := ir.RenderYaml()
	if err != nil {
		t.Errorf(err.Error())
	}

	if string(yamls) != expectedBundle {
		t.Errorf("Rendered yaml is not correct. Was: %s\nExpected: %s", string(yamls), expectedBundle)
	}

	os.RemoveAll(p)
}

func TestAddYttValyesByBytes(t *testing.T) {
	p, err := helperSetupImageLocally()
	if err != nil {
		t.Errorf(err.Error())
	}

	// Remove the values yaml file so we can inject it via our APIs
	os.Remove(filepath.Join(p, "config", "values.yaml"))

	ir, err := NewTkrImageReader(p)
	if err != nil {
		t.Errorf(err.Error())
	}

	ir.SetDownloadPath(p)
	ir.SetRelativeConfigPath("config")

	err = ir.AddYttYamlValuesBytes([]byte("---\nnamespace: custom-namespace-name"))
	if err != nil {
		t.Errorf(err.Error())
	}

	yamls, err := ir.RenderYaml()
	if err != nil {
		t.Errorf(err.Error())
	}

	if string(yamls) != expectedBundle {
		t.Errorf("Rendered yaml is not correct. Was: %s\nExpected: %s", string(yamls), expectedBundle)
	}

	os.RemoveAll(p)
}

func TestAddYttKVsFromYaml(t *testing.T) {
	p, err := helperSetupImageLocally()
	if err != nil {
		t.Errorf(err.Error())
	}

	// Remove the values yaml file so we can inject it via our APIs
	os.Remove(filepath.Join(p, "config", "values.yaml"))

	ir, err := NewTkrImageReader(p)
	if err != nil {
		t.Errorf(err.Error())
	}

	ir.SetDownloadPath(p)
	ir.SetRelativeConfigPath("config")

	ir.AddYttKVsFromYAML([]string{"namespace=custom-namespace-name"})

	yamls, err := ir.RenderYaml()
	if err != nil {
		t.Errorf(err.Error())
	}

	if string(yamls) != expectedBundle {
		t.Errorf("Rendered yaml is not correct. Was: %s\nExpected: %s", string(yamls), expectedBundle)
	}

	os.RemoveAll(p)
}
