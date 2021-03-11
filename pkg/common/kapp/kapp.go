// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kapp

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/adrg/xdg"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	klog "k8s.io/klog/v2"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	ipkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/installpackage/v1alpha1"
	kappctrl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kapppack "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/package/v1alpha1"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = kappctrl.AddToScheme(scheme)
	_ = kapppack.AddToScheme(scheme)
	_ = ipkg.AddToScheme(scheme)
}

// NewKapp generates a Kapp object
func NewKapp(byConfig []byte) (*Kapp, error) {

	cfg, err := InitKappConfig(byConfig)
	if err != nil {
		klog.Errorf("InitKappConfig failed. Err: %v", err)
		return nil, err
	}

	kapp := &Kapp{
		config:                cfg,
		localWorkingDirectory: filepath.Join(xdg.DataHome, "tanzu-repository", cfg.WorkingDirectory),
	}

	klog.V(4).Infof("localWorkingDirectory = %s", kapp.localWorkingDirectory)

	return kapp, nil
}

// GetWorkingDirectory for ytt
func (k *Kapp) GetWorkingDirectory() string {
	return k.localWorkingDirectory
}

func (k *Kapp) createClient() (*client.Client, error) {
	// create k8s client
	config, err := clientcmd.BuildConfigFromFlags("", k.config.Kubeconfig)
	if err != nil {
		klog.Errorf("BuildConfigFromFlags failed. Err: %v", err)
		return nil, err
	}
	client, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		klog.Errorf("client.New failed. Err: %v", err)
		return nil, err
	}

	return &client, nil
}

// RetrievePackages returns all packages available in the cluster.
func (k *Kapp) RetrievePackages() ([]kapppack.Package, error) {
	cl, err := k.createClient()
	if err != nil {
		return nil, err
	}

	// retrieve a list of all packages
	// these resources are cluster-wide (not namespace scoped)
	pkgs := &kapppack.PackageList{}
	err = (*cl).List(context.Background(), pkgs)
	if err != nil {
		klog.Errorln("Failed to retireve list of Packages from cluster")
		return nil, err
	}

	return pkgs.Items, nil
}

// ResolvePackageBundleLocation takes a Package CR and looks up the associated
// imgpkg bundle. There may only be 1 imgpkg bundle associated with the Package
// CR or else an error is returned.
func (k *Kapp) ResolvePackageBundleLocation(pkg kapppack.Package) (string, error) {

	if len(pkg.Spec.Template.Spec.Fetch) != 1 {
		return "", fmt.Errorf("The package %s's spec can contain only 1 bundle", pkg.Name)
	}

	if pkg.Spec.Template.Spec.Fetch[0].ImgpkgBundle == nil {
		return "", fmt.Errorf("The package %s's spec did not contain an imagepkgbundle", pkg.Name)
	}

	if pkg.Spec.Template.Spec.Fetch[0].ImgpkgBundle.Image == "" {
		return "", fmt.Errorf("The package %s's imagepkgbundle did not contain a valid image", pkg.Name)
	}

	return pkg.Spec.Template.Spec.Fetch[0].ImgpkgBundle.Image, nil
}

// ResolveInstalledPackage takes a package name (publicName) and version and returns the
// contents of that InstalledPackage. When only the name is provided, the newest InstalledPackage
// resolved is returned. If a package cannot be resolved due to the name and/or
// version, an error is returned.
func (k *Kapp) ResolveInstalledPackage(name string, version string, namespace string) (*ipkg.InstalledPackage, error) {

	// create the kubernetes client for retrieving Package CRs
	cl, err := k.createClient()
	if err != nil {
		klog.Errorln("failed to create client")
		return nil, err
	}

	// list all InstalledPackages in the cluster
	//
	// TODO(joshrosso): Listing all InstalledPackges is unideal, but I can't find a way to make
	// field selectors work on CRDs. https://github.com/kubernetes/kubernetes/issues/51046
	packageList := &ipkg.InstalledPackageList{}
	err = (*cl).List(context.Background(), packageList, client.InNamespace(namespace))
	if err != nil {
		klog.Errorf("failed to get package list. error: %s", err.Error())
	}

	// for every package, loop through and resolve the publicName against Name. If no
	// version is specified return the first package. If a version is specified, check
	// resolution against version, it it does not match, continue iterating.
	//
	// TODO(joshrosso): when version is *not* specified, we should resolve the newest
	//                  version and return it.
	var resolvedPackage *ipkg.InstalledPackage
	for _, pkg := range packageList.Items {
		if pkg.Spec.PkgRef.PublicName == name {

			if version == "" {
				resolvedPackage = &pkg
				break
			}

			if pkg.Spec.PkgRef.Version == version {
				resolvedPackage = &pkg
				break
			}

		}
	}

	// when no installedpackage was resolved, return an error
	if resolvedPackage == nil {
		return nil, fmt.Errorf("could not resolve installedpackage %s/%s:%s", namespace, name, version)
	}

	klog.V(6).Infof("Package CR was resolved as: %s", resolvedPackage.Name)
	return resolvedPackage, nil
}

// ResolvePackage takes a package name (publicName) and version and returns the
// contents of that package. When only the name is provided, the newest package
// resolved is returned. If a package cannot be resolved due to the name and/or
// version, an error is returned.
func (k *Kapp) ResolvePackage(name string, version string) (*kapppack.Package, error) {

	// create the kubernetes client for retrieving Package CRs
	client, err := k.createClient()
	if err != nil {
		klog.Errorln("failed to create client")
		return nil, err
	}
	cl := *client

	// list all package in the cluster
	//
	// TODO(joshrosso): Listing all packages is unideal, but I can't find a way to make
	// field selectors work on CRDs. https://github.com/kubernetes/kubernetes/issues/51046
	packageList := &kapppack.PackageList{}
	err = cl.List(context.Background(), packageList)
	if err != nil {
		klog.Errorf("failed to get package list. error: %s", err.Error())
	}

	// for every package, loop through and resolve the publicName against Name. If no
	// version is specified return the first package. If a version is specified, check
	// resolution against version, it it does not match, continue iterating.
	//
	// TODO(joshrosso): when version is *not* specified, we should resolve the newest
	//                  version and return it.
	var resolvedPackage *kapppack.Package
	for _, pkg := range packageList.Items {
		if pkg.Spec.PublicName == name {

			if version == "" {
				resolvedPackage = &pkg
				break
			}

			if pkg.Spec.Version == version {
				resolvedPackage = &pkg
				break
			}

		}
	}

	// when no package was resolved, return an error
	if resolvedPackage == nil {
		return nil, fmt.Errorf("could not resolve package %s with version %s", name, version)
	}

	klog.V(6).Infof("Package CR was resolved as: %s", resolvedPackage.Name)
	return resolvedPackage, nil
}

func (k *Kapp) installServiceAccount(client *client.Client, input *AppCrdInput) error {

	klog.V(2).Infof("installServiceAccount(%s)", input.Name)

	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      input.Name + k.config.ExtensionServiceAccountPostfix,
			Namespace: input.Namespace,
		},
	}
	klog.V(6).Infof("serviceAccount.Name = %s", serviceAccount.ObjectMeta.Name)
	klog.V(6).Infof("sa.Namespace = %s", serviceAccount.ObjectMeta.Namespace)

	_, err := controllerutil.CreateOrPatch(context.TODO(), *client, serviceAccount, nil)
	if err != nil {
		klog.Errorf("Error creating or patching addon service account. Err: %v", err)
		return err
	}

	klog.V(2).Infof("installServiceAccount(%s) succeeded", input.Name)
	return nil
}

func (k *Kapp) installRoleBinding(client *client.Client, input *AppCrdInput) error {

	klog.V(2).Infof("installRoleBinding(%s)", input.Name)

	roleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: input.Name + "-" + input.Namespace + k.config.ExtensionRoleBindingPostfix,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      input.Name + k.config.ExtensionServiceAccountPostfix,
				Namespace: input.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
	}
	klog.V(6).Infof("roleBinding.Name = %s", roleBinding.ObjectMeta.Name)
	klog.V(6).Infof("roleBinding.Subjects.Name = %s", roleBinding.Subjects[0].Name)
	klog.V(6).Infof("roleBinding.Subjects.Namespace = %s", roleBinding.Subjects[0].Namespace)

	_, err := controllerutil.CreateOrPatch(context.TODO(), *client, roleBinding, nil)
	if err != nil {
		klog.Errorf("Error creating or patching addon role binding. Err: %v", err)
		return err
	}

	klog.V(2).Infof("installRoleBinding(%s) succeeded", input.Name)
	return nil
}

// InstallPackage creates the InstalledPackage CR and applies it to the cluster.
func (k *Kapp) InstallPackage(input *AppCrdInput) error {
	client, err := k.createClient()
	if err != nil {
		klog.Errorf("createClient failed. Err: %v", err)
		return err
	}

	err = k.installServiceAccount(client, input)
	if err != nil {
		return err
	}

	err = k.installRoleBinding(client, input)
	if err != nil {
		return err
	}

	// if the configuration data exists, create a secret object
	// and capture its name
	var configName *string
	if len(input.Config) > 0 {
		configName, err = k.installConfigSecret(client, input)
		if err != nil {
			return err
		}
	}

	err = k.installInstalledPackage(client, input, configName)
	if err != nil {
		return err
	}

	return nil
}

// ResolvePackageVersion takes a package an input and returns the Package's version
func (k *Kapp) ResolvePackageVersion(pkg *kapppack.Package) string {
	return pkg.Spec.Version
}

func (k *Kapp) installInstalledPackage(client *client.Client, input *AppCrdInput, configName *string) error {

	// construct the InstalledPackage CR
	ip := &ipkg.InstalledPackage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      input.Name,
			Namespace: input.Namespace,
		},
		Spec: ipkg.InstalledPackageSpec{
			ServiceAccountName: input.Name + k.config.ExtensionServiceAccountPostfix,
			PkgRef: &ipkg.PackageRef{
				PublicName: input.Name,
				VersionSelection: &versions.VersionSelectionSemver{
					Constraints: input.Version,
					Prereleases: &versions.VersionSelectionSemverPrereleases{},
				},
			},
		},
	}

	// if configuration was provided, reference the config (secret) name in
	// the InstalledPackage
	if configName != nil {
		ip.Spec.Values = []ipkg.InstalledPackageValues{
			{
				SecretRef: &ipkg.InstalledPackageValuesSecretRef{
					Name: *configName,
				},
			},
		}
	}

	klog.V(6).Infof("Deploying installed package: %s", ip)
	_, err := controllerutil.CreateOrPatch(context.TODO(), *client, ip, nil)
	if err != nil {
		return err
	}

	return nil
}

func (k *Kapp) deleteServiceAccount(client *client.Client, input *AppCrdInput) error {

	klog.V(2).Infof("deleteServiceAccount(%s)", input.Name)

	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      input.Name + k.config.ExtensionServiceAccountPostfix,
			Namespace: input.Namespace,
		},
	}
	klog.V(6).Infof("serviceAccount.Name = %s", serviceAccount.ObjectMeta.Name)
	klog.V(6).Infof("sa.Namespace = %s", serviceAccount.ObjectMeta.Namespace)

	if err := (*client).Delete(context.TODO(), serviceAccount); err != nil {
		if apierrors.IsNotFound(err) {
			klog.V(2).Info("Service account not found")
			return nil
		}
		klog.Errorf("Error deleting service account. Err: %v", err)
		return err
	}

	klog.V(2).Infof("deleteServiceAccount(%s) succeeded", input.Name)
	return nil
}

func (k *Kapp) deleteInstalledPackage(client *client.Client, input *AppCrdInput) error {

	ipkgToDelete := &ipkg.InstalledPackage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      input.Name,
			Namespace: input.Namespace,
		},
	}

	err := (*client).Delete(context.Background(), ipkgToDelete)
	if err != nil {
		return err
	}

	return nil
}

func (k *Kapp) deleteRoleBinding(client *client.Client, input *AppCrdInput) error {

	klog.V(2).Infof("deleteRoleBinding(%s)", input.Name)

	roleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: input.Name + "-" + input.Namespace + k.config.ExtensionRoleBindingPostfix,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      input.Name + k.config.ExtensionServiceAccountPostfix,
				Namespace: input.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
	}
	klog.V(6).Infof("roleBinding.Name = %s", roleBinding.ObjectMeta.Name)
	klog.V(6).Infof("roleBinding.Subjects.Name = %s", roleBinding.Subjects[0].Name)
	klog.V(6).Infof("roleBinding.Subjects.Namespace = %s", roleBinding.Subjects[0].Namespace)

	if err := (*client).Delete(context.TODO(), roleBinding); err != nil {
		if apierrors.IsNotFound(err) {
			klog.V(2).Info("Role binding not found")
			return nil
		}
		klog.Errorf("Error deleting role binding. Err: %v", err)
		return err
	}

	klog.V(2).Infof("deleteRoleBinding(%s) succeeded", input.Name)

	return nil
}

// DeletePackage removes the InstalledPackage CR and related assets from the cluster.
func (k *Kapp) DeletePackage(input *AppCrdInput) error {
	client, err := k.createClient()
	if err != nil {
		klog.Errorf("createClient failed. Err: %v", err)
		return err
	}

	if input.ConfigPath != "" {
		err = k.deleteConfigSecret(client, input)
		if err != nil {
			return err
		}
	}

	err = k.deleteInstalledPackage(client, input)
	if err != nil {
		return err
	}

	if input.Teardown {
		err = k.deleteServiceAccount(client, input)
		if err != nil {
			return err
		}

		err = k.deleteRoleBinding(client, input)
		if err != nil {
			return err
		}
	}

	return nil
}

// installConfigSecret create a secret object containing the user-provided configuration. It
// returns and errror if it fails to apply. Upon success, it returns the name of the secret
// created.
func (k *Kapp) installConfigSecret(client *client.Client, input *AppCrdInput) (*string, error) {

	configName := input.Name + "-config"

	config := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configName,
			Namespace: input.Namespace,
		},
		Data: map[string][]byte{
			"values.yaml": input.Config,
		},
	}

	_, err := controllerutil.CreateOrPatch(context.TODO(), *client, config, nil)
	if err != nil {
		return nil, err
	}

	return &configName, nil
}

// deleteConfigSecret deletes a secret object containing the user-provided configuration. It
// returns and errror if it fails to delete.
func (k *Kapp) deleteConfigSecret(client *client.Client, input *AppCrdInput) error {

	configName := input.Name + "-config"

	config := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configName,
			Namespace: input.Namespace,
		},
	}

	err := (*client).Delete(context.TODO(), config)
	if err != nil {
		return err
	}

	return nil
}
