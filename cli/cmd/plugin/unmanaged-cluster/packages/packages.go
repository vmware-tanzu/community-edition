// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package packages facilitates the creation of Carvel package repositories and package installs
// within the cluster. In order for its operations to be successful, an instance of kapp-controller
// is expected to be running in the cluster.
package packages

import (
	"context"
	"fmt"
	"time"

	kappapis "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	packaging "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	datapackaging "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"

	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	clusterAdminRole       = "cluster-admin"
	clusterRoleKind        = "ClusterRole"
	svcAcctKind            = "ServiceAccount"
	packageRepoResource    = "packagerepositories"
	packageInstallResource = "packageinstalls"
	packageResource        = "packages"
	packageRepoKind        = "PackageRepository"
	packageInstallKind     = "PackageInstall"
	tkgSysNamespace        = "tkg-system"
	rbacAPIGroup           = "rbac.authorization.k8s.io"
	apiBaseURI             = "/apis"
)

// PackageClient implements PackageManager and holds references to both
// clientSet and restClient objects. clientSet is used to interact with native
// Kubernetes objects and restClient is used to interact with CRDs.
type PackageClient struct {
	// used to access packaging CRDs
	// for example, packageinstall
	restClient rest.Interface
	// aggRestClient accessses rsources in the aggregated API server; provided by kapp-controller
	// for example, packages.data.packaging.carvel.dev
	aggRestClient rest.Interface
	// clientSet accesses standard Kubernetes resources
	clientSet kubernetes.Interface
}

type PackageInstallOpts struct {
	// The namespace the PackageInstall object should be created in
	Namespace string
	// The name of the create PackageInstall object
	InstallName string
	// The fully qualified name of the package to install. This and Version are how packages are resolved
	// from a repository.
	FqPkgName string
	// The version of the package to install. This and FqPkgName are how packages are resolved
	// from a repository.
	Version string
	// Optional configuration to be added alongside the package installation. When this value is non-nil, a
	// Secret object is created in the cluster and the package install references it as a values configuration.
	Configuration []byte
	// The ServiceAccount used to facilitate the package install. It must have all privileges required for
	// kapp-controller to create the appropriate objects.
	ServiceAccount string
}

// PackageManager provides operations for doing package management against a cluster.
type PackageManager interface {
	// CreatePackageRepo adds a PackageRepository to the cluster, which in turn makes packages
	// available. If it can successfully add the repository to the cluster, it returns the
	// PackageRepository object created, otherwise an error. It does not wait for a package
	// repository to reconcile. Upon success, it returns the created PackageRepository object.
	CreatePackageRepo(ns, name, url string) (*packaging.PackageRepository, error)
	// CreatePackageInstall adds a PackageInstall object to the cluster. It requires you provide
	// the namespace, install name, fully qualified package name, version, and service account.
	// Configuration may also be passed, if nil, no configuration is added. Configuration is added
	// by injecting a secret object into the cluster and referencing it from the package install.
	// Upon success, it returns the created PackageInstall object.
	CreatePackageInstall(opts *PackageInstallOpts) (*packaging.PackageInstall, error)
	// CreateRootServiceAccount creates a service account in the target namespace with a ClusterRoleBinding
	// referencing the cluster-admin CluterRole. This essentially provides full admin access to anything
	// referencing this service account. Upon success, it returns the created ServiceAccount.
	CreateRootServiceAccount(ns, name string) (*v1.ServiceAccount, error)
	// GetRepositoryStatus outputs the status of a repository based on the namespace and repository name
	// requested. It provides details on kapp-controller process such as "Reconciling" and "Reconcile Succeeded"
	GetRepositoryStatus(ns, name string) (string, error)
	// ListPackagesInNamespace returns a list of packages based on the namespace.
	ListPackagesInNamespace(ns string) ([]datapackaging.Package, error)
}

// NewClient create an instance of a PackageManager, implemented by PackageClient,
// by passing a kubeconfig targeting the cluster. It also sets up both a restClient
// for CRD interaction (Package APIs) and a clientSet for Kubernetes API interaction.
// For the restClient, it registers the packaging APIs to the scheme.
func NewClient(kubeconfigBytes []byte) PackageManager {
	config, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (*clientcmdapi.Config, error) {
		return clientcmd.Load(kubeconfigBytes)
	})

	if err != nil {
		// TODO(joshrosso): do something here
		panic(err.Error())
	}

	// register packaging APIs
	_ = packaging.AddToScheme(scheme.Scheme)
	_ = datapackaging.AddToScheme(scheme.Scheme)
	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &packaging.SchemeGroupVersion
	crdConfig.APIPath = apiBaseURI
	crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	c, err := rest.RESTClientFor(&crdConfig)
	if err != nil {
		// TODO(joshrosso): do something here
		panic(err)
	}

	// registry (aggregated api server) packaging APIs
	aggAPIConfig := *config
	aggAPIConfig.ContentConfig.GroupVersion = &datapackaging.SchemeGroupVersion
	aggAPIConfig.APIPath = apiBaseURI
	aggAPIConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	aggAPIConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	aggRc, err := rest.RESTClientFor(&aggAPIConfig)
	if err != nil {
		panic(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		// TODO(joshrosso): do something here
		panic(err)
	}

	return &PackageClient{
		restClient:    c,
		aggRestClient: aggRc,
		clientSet:     clientSet,
	}
}

func (am *PackageClient) CreatePackageRepo(ns, name, url string) (*packaging.PackageRepository, error) {
	// TODO(joshrosso): do pre-check that url does exist as valid imgpkg bundle

	apiVersion := fmt.Sprintf("%s/%s", packaging.SchemeGroupVersion.Group, packaging.SchemeGroupVersion.Version)
	// create package repository object
	repo := &packaging.PackageRepository{
		TypeMeta: metav1.TypeMeta{
			Kind:       packageRepoKind,
			APIVersion: apiVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: packaging.PackageRepositorySpec{
			SyncPeriod: &metav1.Duration{Duration: 5 * time.Minute},
			Fetch: &packaging.PackageRepositoryFetch{
				ImgpkgBundle: &kappapis.AppFetchImgpkgBundle{
					Image: url,
				},
			},
		},
	}

	createdRepo := &packaging.PackageRepository{}

	// create package repo and store the end state object in an object
	err := am.restClient.
		Post().
		Resource(packageRepoResource).
		Namespace(ns).
		Name(name).
		Body(repo).
		Do(context.TODO()).
		Into(createdRepo)

	if err != nil {
		return nil, err
	}

	return createdRepo, nil
}

func (am *PackageClient) CreatePackageInstall(opts *PackageInstallOpts) (*packaging.PackageInstall, error) {
	// TODO(joshrosso): do pre-check package requesting install resolves in the package repo

	apiVersion := fmt.Sprintf("%s/%s", packaging.SchemeGroupVersion.Group, packaging.SchemeGroupVersion.Version)

	// create package install object
	pkgInstall := &packaging.PackageInstall{
		TypeMeta: metav1.TypeMeta{
			Kind:       packageInstallKind,
			APIVersion: apiVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.InstallName,
			Namespace: opts.Namespace,
		},
		Spec: packaging.PackageInstallSpec{
			ServiceAccountName: opts.ServiceAccount,
			PackageRef: &packaging.PackageRef{
				RefName: opts.FqPkgName,
				VersionSelection: &versions.VersionSelectionSemver{
					Constraints: opts.Version,
				},
			},
			Values: []packaging.PackageInstallValues{
				{},
			},
			SyncPeriod: &metav1.Duration{Duration: 1 * time.Minute},
		},
	}

	if opts.Configuration != nil {
		// create secret based on configuration data
		values := make(map[string]string)
		values["values.yml"] = string(opts.Configuration)
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      opts.InstallName + "-config",
				Namespace: opts.Namespace,
			},
			StringData: values,
		}

		createdSecret, err := am.clientSet.CoreV1().Secrets(tkgSysNamespace).Create(context.TODO(), secret, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("Failed to create secret: %s", err.Error())
			return nil, err
		}

		// set PackageInstall reference to created secret
		pkgInstall.Spec.Values[0].SecretRef = &packaging.PackageInstallValuesSecretRef{
			Name: createdSecret.Name,
		}
	}

	// create package install object in cluster
	createdInstall := &packaging.PackageInstall{}
	err := am.restClient.
		Post().
		Resource(packageInstallResource).
		Namespace(opts.Namespace).
		Name(opts.InstallName).
		Body(pkgInstall).
		Do(context.TODO()).
		Into(createdInstall)
	if err != nil {
		return nil, err
	}

	return createdInstall, nil
}

func (am *PackageClient) GetRepositoryStatus(ns, name string) (string, error) {
	repo := &packaging.PackageRepository{}
	err := am.restClient.
		Get().
		Namespace(ns).
		Name(name).
		Resource(packageRepoResource).
		Do(context.TODO()).
		Into(repo)
	if err != nil {
		return "", err
	}

	return repo.Status.FriendlyDescription, nil
}

func (am *PackageClient) CreateRootServiceAccount(ns, name string) (*v1.ServiceAccount, error) {
	svcAcct := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}

	roleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      svcAcctKind,
				Name:      name,
				Namespace: ns,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacAPIGroup,
			Kind:     clusterRoleKind,
			Name:     clusterAdminRole,
		},
	}

	createdSa, err := am.clientSet.CoreV1().ServiceAccounts(tkgSysNamespace).Create(context.TODO(), svcAcct, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	_, err = am.clientSet.RbacV1().ClusterRoleBindings().Create(context.TODO(), roleBinding, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return createdSa, nil
}

func (am *PackageClient) ListPackagesInNamespace(ns string) ([]datapackaging.Package, error) {
	pkgList := &datapackaging.PackageList{}
	err := am.aggRestClient.Get().
		Resource(packageResource).
		Namespace(ns).
		Do(context.TODO()).
		Into(pkgList)
	if err != nil {
		return nil, err
	}

	return pkgList.Items, nil
}
