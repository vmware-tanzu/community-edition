// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kapp

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "github.com/ghodss/yaml"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	klog "k8s.io/klog/v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	instpkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/installpackage/v1alpha1"
	kappctrl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

func (k *Kapp) generateRepoFromFile(filename string) (*instpkg.PackageRepository, error) {
	path, err := os.Getwd()
	if err != nil {
		klog.Errorf("Getwd failed. Err: ", err)
		return nil, err
	}

	repoFile, err := filepath.Abs(filepath.Join(path, filename))
	if err != nil {
		klog.Errorf("filepath.Abs failed. Err: ", err)
		return nil, err
	}
	klog.V(2).Infof("repoFile: %s", repoFile)

	// read the contents of the provided file
	byFile, err := ioutil.ReadFile(repoFile)
	if err != nil {
		klog.Errorf("Open failed. Err:", err)
		return nil, err
	}

	klog.V(6).Infof("Data:\n")
	klog.V(6).Infof("%s\n\n", string(byFile))

	// unmarshal instpkg.PackageRepository
	repo := &instpkg.PackageRepository{}

	err = yaml.Unmarshal(byFile, &repo)
	if err != nil {
		klog.Errorf("Unmarshal failed. Err: ", err)
		return nil, err
	}

	klog.V(2).Infof("generateRepoFromFile() succeeded")
	return repo, nil
}

// InstallDefaultRepository uses TCE
func (k *Kapp) InstallDefaultRepository() error {
	return k.InstallRepository(DefaultRepositoryName, DefaultRepositoryImage)
}

// InstallRepository installs a generic repo
func (k *Kapp) InstallRepository(name, url string) error {
	klog.V(2).Infof("InstallRepository()")
	klog.V(2).Infof("name: %s", name)
	klog.V(2).Infof("url: %s", url)

	client, err := k.createClient()
	if err != nil {
		klog.Errorf("createClient failed. Err: %v", err)
		return err
	}

	// unmarshal instpkg.App
	repo := &instpkg.PackageRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: instpkg.PackageRepositorySpec{
			Fetch: &instpkg.PackageRepositoryFetch{
				Image: &kappctrl.AppFetchImage{
					URL: url,
				},
			},
		},
	}

	_, err = controllerutil.CreateOrPatch(context.TODO(), client, repo, nil)
	if err != nil {
		klog.Errorf("Create failed. Err: %v", err)
		return err
	}

	klog.V(2).Infof("InstallRepository %s succeeded", name)
	return nil
}

// InstallRepositoryFromFile does it from a file
func (k *Kapp) InstallRepositoryFromFile(filename string) error {
	klog.V(2).Infof("InstallRepositoryFromFile()")
	klog.V(2).Infof("filename: %s", filename)

	client, err := k.createClient()
	if err != nil {
		klog.Errorf("createClient failed. Err: %v", err)
		return err
	}

	repo, err := k.generateRepoFromFile(filename)
	if err != nil {
		klog.Errorf("generateRepoFromFile failed. Err: %v", err)
		return err
	}

	_, err = controllerutil.CreateOrPatch(context.TODO(), client, repo, nil)
	if err != nil {
		klog.Errorf("Create failed. Err: %v", err)
		return err
	}

	klog.V(2).Infof("InstallRepositoryFromFile %s succeeded", filename)
	return nil
}

// ListRepositories lists all repos
func (k *Kapp) ListRepositories() (*instpkg.PackageRepositoryList, error) {
	klog.V(2).Infof("ListRepositories()")

	client, err := k.createClient()
	if err != nil {
		klog.Errorf("createClient failed. Err: %v", err)
		return nil, err
	}

	list := &instpkg.PackageRepositoryList{}

	err = (client).List(context.TODO(), list)
	if err != nil {
		klog.Errorf("List failed. Err: %v", err)
		return nil, err
	}

	return list, nil
}

// DeleteRepository deletes a repo
func (k *Kapp) DeleteRepository(name string) error {
	klog.V(2).Infof("DeleteRepository()")
	klog.V(2).Infof("name: %s", name)

	client, err := k.createClient()
	if err != nil {
		klog.Errorf("createClient failed. Err: %v", err)
		return err
	}

	repo := &instpkg.PackageRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	if err := (client).Delete(context.TODO(), repo); err != nil {
		if apierrors.IsNotFound(err) {
			klog.V(2).Info("Repository not found")
			return nil
		}
		klog.Errorf("Error deleting repository. Err: %v", err)
		return err
	}

	klog.V(2).Infof("DeleteRepository %s succeeded", name)
	return nil
}

// DeleteRepositoryFromFile does it from a file
func (k *Kapp) DeleteRepositoryFromFile(filename string) error {
	klog.V(2).Infof("InstallRepositoryFromFile()")
	klog.V(2).Infof("filename: %s", filename)

	client, err := k.createClient()
	if err != nil {
		klog.Errorf("createClient failed. Err: %v", err)
		return err
	}

	repo, err := k.generateRepoFromFile(filename)
	if err != nil {
		klog.Errorf("generateRepoFromFile failed. Err: %v", err)
		return err
	}

	if err := (client).Delete(context.TODO(), repo); err != nil {
		if apierrors.IsNotFound(err) {
			klog.V(2).Info("Repository not found")
			return nil
		}
		klog.Errorf("Error deleting repository. Err: %v", err)
		return err
	}

	klog.V(2).Infof("DeleteRepositoryFromFile %s succeeded", filename)
	return nil
}
