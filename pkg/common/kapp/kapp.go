// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kapp

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	yaml "github.com/ghodss/yaml"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	klog "k8s.io/klog/v2"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	kappctrl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"

	types "github.com/vmware-tanzu/tce/pkg/common/types"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = kappctrl.AddToScheme(scheme)
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

func duration(s string) time.Duration {
	t, _ := time.ParseDuration(s)
	return t
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

func (k *Kapp) installServiceAccount(client *client.Client, extensionName string) error {

	klog.V(2).Infof("installServiceAccount(%s)", extensionName)

	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      extensionName + k.config.ExtensionServiceAccountPostfix,
			Namespace: k.config.ExtensionNamespace,
		},
	}
	klog.V(6).Infof("serviceAccount.Name = %s", serviceAccount.ObjectMeta.Name)
	klog.V(6).Infof("sa.Namespace = %s", serviceAccount.ObjectMeta.Namespace)

	_, err := controllerutil.CreateOrPatch(context.TODO(), *client, serviceAccount, nil)
	if err != nil {
		klog.Errorf("Error creating or patching addon service account. Err: %v", err)
		return err
	}

	klog.V(2).Infof("installServiceAccount(%s) succeeded", extensionName)
	return nil
}

func (k *Kapp) installRoleBinding(client *client.Client, extensionName string) error {

	klog.V(2).Infof("installRoleBinding(%s)", extensionName)

	roleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: extensionName + k.config.ExtensionRoleBindingPostfix,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      extensionName + k.config.ExtensionServiceAccountPostfix,
				Namespace: k.config.ExtensionNamespace,
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

	klog.V(2).Infof("installRoleBinding(%s) succeeded", extensionName)
	return nil
}

// installAppCrd install app crd
func (k *Kapp) installAppCrd(client *client.Client, extensionName string) error {

	klog.V(2).Infof("InstallAppCrd(%s)", extensionName)

	workingExtensionDir := filepath.Join(k.localWorkingDirectory, extensionName)
	fullFilename := filepath.Join(workingExtensionDir, types.DefaultAppCrdFilename)
	klog.V(2).Infof("workingExtensionDir = %s", workingExtensionDir)
	klog.V(2).Infof("fullFilename = %s", fullFilename)

	// read the contents of the provided file
	byFile, err := ioutil.ReadFile(fullFilename)
	if err != nil {
		klog.Errorf("Open failed. Err:", err)
		return err
	}

	klog.V(6).Infof("Data:\n")
	klog.V(6).Infof("%s\n\n", string(byFile))

	// unmarshal kappctrl.App
	app := &kappctrl.App{}

	err = yaml.Unmarshal(byFile, &app)
	if err != nil {
		klog.Errorf("Unmarshal failed. Err: ", err)
		return err
	}

	klog.V(4).Infof("Unmarshal succeeded\n")
	klog.V(6).Infof("app.ObjectMeta.Name = %s\n", app.ObjectMeta.Name)
	klog.V(6).Infof("app.ObjectMeta.Namespace = %s\n", app.ObjectMeta.Namespace)
	klog.V(6).Infof("app.Spec.ServiceAccountName = %s\n", app.Spec.ServiceAccountName)
	for _, f := range app.Spec.Fetch {
		if f.Image != nil {
			klog.V(6).Infof("Image.URL = %s\n", f.Image.URL)
		}
		if f.ImgpkgBundle != nil {
			klog.V(6).Infof("ImgpkgBundle.Image = %s\n", f.ImgpkgBundle.Image)
		}
	}
	for _, t := range app.Spec.Template {
		if t.Ytt == nil {
			continue
		}
		if t.Ytt.Inline == nil {
			continue
		}
		for _, p := range t.Ytt.Inline.Paths {
			klog.V(6).Infof("Path = %s\n", p)
		}
	}
	klog.V(4).Info("\n%v\n\n", app)

	appMutateFn := func() error {
		if app.ObjectMeta.Annotations == nil {
			app.ObjectMeta.Annotations = make(map[string]string)
		}

		app.Spec.Deploy = []kappctrl.AppDeploy{
			{
				Kapp: &kappctrl.AppDeployKapp{},
			},
		}

		return nil
	}

	_, err = controllerutil.CreateOrPatch(context.TODO(), *client, app, appMutateFn)
	if err != nil {
		klog.Errorf("CreateOrPatch failed. Err: %v", err)
		return err
	}

	klog.V(2).Infof("InstallFromFile(%s) succeeded", extensionName)
	return nil
}

// InstallFromFile install extension from file
func (k *Kapp) InstallFromFile(input *AppCrdInput) error {
	client, err := k.createClient()
	if err != nil {
		klog.Errorf("createClient failed. Err: %v", err)
		return err
	}

	err = k.installServiceAccount(client, input.Name)
	if err != nil {
		klog.Errorf("installServiceAccount failed. Err: %v", err)
		return err
	}

	err = k.installRoleBinding(client, input.Name)
	if err != nil {
		klog.Errorf("installRoleBinding failed. Err: %v", err)
		return err
	}

	err = k.installAppCrd(client, input.Name)
	if err != nil {
		klog.Errorf("installAppCrd failed. Err: %v", err)
		return err
	}

	return nil
}

func (k *Kapp) deleteServiceAccount(client *client.Client, extensionName string) error {

	klog.V(2).Infof("deleteServiceAccount(%s)", extensionName)

	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      extensionName + k.config.ExtensionServiceAccountPostfix,
			Namespace: k.config.ExtensionNamespace,
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

	klog.V(2).Infof("deleteServiceAccount(%s) succeeded", extensionName)
	return nil
}

func (k *Kapp) deleteRoleBinding(client *client.Client, extensionName string) error {

	klog.V(2).Infof("deleteRoleBinding(%s)", extensionName)

	roleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: extensionName + k.config.ExtensionRoleBindingPostfix,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      extensionName + k.config.ExtensionServiceAccountPostfix,
				Namespace: k.config.ExtensionNamespace,
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

	klog.V(2).Infof("deleteRoleBinding(%s) succeeded", extensionName)
	return nil
}

func (k *Kapp) deleteAppCrd(client *client.Client, extensionName string, force bool) error {

	klog.V(2).Infof("deleteAppCrd(%s)", extensionName)

	workingExtensionDir := filepath.Join(k.localWorkingDirectory, extensionName)
	fullFilename := filepath.Join(workingExtensionDir, types.DefaultAppCrdFilename)
	klog.V(2).Infof("workingExtensionDir = %s", workingExtensionDir)
	klog.V(2).Infof("fullFilename = %s", fullFilename)

	// read the contents of the provided file
	byFile, err := ioutil.ReadFile(fullFilename)
	if err != nil {
		klog.Errorf("Open failed. Err:", err)
		return err
	}

	klog.V(6).Infof("Data:\n")
	klog.V(6).Infof("%s\n\n", string(byFile))

	// unmarshal kappctrl.App
	app := &kappctrl.App{}

	err = yaml.Unmarshal(byFile, &app)
	if err != nil {
		klog.Errorf("Unmarshal failed. Err: ", err)
		return err
	}

	klog.V(4).Infof("Unmarshal succeeded\n")
	klog.V(6).Infof("app.ObjectMeta.Name = %s\n", app.ObjectMeta.Name)
	klog.V(6).Infof("app.ObjectMeta.Namespace = %s\n", app.ObjectMeta.Namespace)
	klog.V(6).Infof("app.Spec.ServiceAccountName = %s\n", app.Spec.ServiceAccountName)
	for _, f := range app.Spec.Fetch {
		if f.Image != nil {
			klog.V(6).Infof("Image.URL = %s\n", f.Image.URL)
		}
		if f.ImgpkgBundle != nil {
			klog.V(6).Infof("ImgpkgBundle.Image = %s\n", f.ImgpkgBundle.Image)
		}
	}
	for _, t := range app.Spec.Template {
		if t.Ytt == nil {
			continue
		}
		if t.Ytt.Inline == nil {
			continue
		}
		for _, p := range t.Ytt.Inline.Paths {
			klog.V(6).Infof("Path = %s\n", p)
		}
	}
	klog.V(4).Info("\n%v\n\n", app)

	var errRet error
	if errRet = (*client).Delete(context.TODO(), app); err != nil {
		klog.Errorf("Error deleting App CRD. Err: %v", err)
		if apierrors.IsNotFound(err) {
			klog.Warningf("App CRD is not present/installed")
			errRet = ErrAppNotPresentOrInstalled
		}
	}

	if force {
		app = &kappctrl.App{}

		err = yaml.Unmarshal(byFile, &app)
		if err != nil {
			klog.Errorf("Unmarshal failed. Err: ", err)
			return err
		}

		app.ObjectMeta.Finalizers = []string{}

		klog.V(4).Infof("Unmarshal succeeded\n")
		klog.V(6).Infof("app.ObjectMeta.Name = %s\n", app.ObjectMeta.Name)
		klog.V(6).Infof("app.ObjectMeta.Namespace = %s\n", app.ObjectMeta.Namespace)
		klog.V(6).Infof("app.Spec.ServiceAccountName = %s\n", app.Spec.ServiceAccountName)
		for _, f := range app.Spec.Fetch {
			if f.Image != nil {
				klog.V(6).Infof("Image.URL = %s\n", f.Image.URL)
			}
			if f.ImgpkgBundle != nil {
				klog.V(6).Infof("ImgpkgBundle.Image = %s\n", f.ImgpkgBundle.Image)
			}
		}
		for _, t := range app.Spec.Template {
			if t.Ytt == nil {
				continue
			}
			if t.Ytt.Inline == nil {
				continue
			}
			for _, p := range t.Ytt.Inline.Paths {
				klog.V(6).Infof("Path = %s\n", p)
			}
		}
		klog.V(4).Info("\n%v\n\n", app)

		appMutateFn := func() error {
			app.ObjectMeta.Finalizers = []string{}
			return nil
		}

		_, err := controllerutil.CreateOrPatch(context.TODO(), *client, app, appMutateFn)
		if err != nil {
			klog.Errorf("Error creating or patching addon data values secret. Err: %v", err)
			return err
		}
	}

	if errRet == nil {
		klog.V(2).Infof("deleteAppCrd(%s) succeeded", extensionName)
	} else {
		klog.V(2).Infof("deleteAppCrd(%s) failed. Err: %v", extensionName, errRet)
	}
	return errRet
}

// DeleteFromFile delete extension from file
func (k *Kapp) DeleteFromFile(input *AppCrdInput) error {
	client, err := k.createClient()
	if err != nil {
		klog.Errorf("createClient failed. Err: %v", err)
		return err
	}

	err = k.deleteAppCrd(client, input.Name, input.Force)
	if err != nil {
		klog.Errorf("installServiceAccount failed. Err: %v", err)
		return err
	}

	if input.Teardown {
		err = k.deleteRoleBinding(client, input.Name)
		if err != nil {
			klog.Errorf("installRoleBinding failed. Err: %v", err)
			return err
		}

		err = k.deleteServiceAccount(client, input.Name)
		if err != nil {
			klog.Errorf("installAppCrd failed. Err: %v", err)
			return err
		}
	}

	return nil
}

// Next version...
/*
// InstallFromSecret install from secret
func (k *Kapp) InstallFromSecret(appCrd *AppCrdInput) error {

	klog.V(2).Infof("extension = %s", appCrd.Name)

	// create k8s client
	config, err := clientcmd.BuildConfigFromFlags("", k.config.Kubeconfig)
	if err != nil {
		klog.Errorf("BuildConfigFromFlags failed. Err: %v", err)
		return err
	}
	client, err := client.New(config, client.Options{})
	if err != nil {
		klog.Errorf("client.New failed. Err: %v", err)
		return err
	}

	// get secret
	cluster := &clusterapiv1alpha3.Cluster{
		Name:      appCrd.ClusterName,
		Namespace: appCrd.Namespace,
	}

	addonSecret, err := util.GetAddonSecretsForCluster(context.TODO(), clusterClient, cluster)
	if err != nil {
		klog.Errorf("GetAddonSecretsForCluster failed. Err: %v", err)
		return err
	}

	// populate from secret
	addonName := util.GetAddonNameFromAddonSecret(addonSecret)

	app := &kappctrl.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      util.GenerateAppNameFromAddonSecret(addonSecret),
			Namespace: util.GenerateAppNamespaceFromAddonSecret(addonSecret),
		},
	}

	appMutateFn := func() error {
		if app.ObjectMeta.Annotations == nil {
			app.ObjectMeta.Annotations = make(map[string]string)
		}

		// app.ObjectMeta.Annotations[addontypes.AddonTypeAnnotation] = fmt.Sprintf("%s/%s", addonConfig.Category, addonName)
		// app.ObjectMeta.Annotations[addontypes.AddonNameAnnotation] = addonSecret.Name
		// app.ObjectMeta.Annotations[addontypes.AddonNamespaceAnnotation] = addonSecret.Namespace

		// remoteApp means App is not present on local workload cluster. It is present in the remote management cluster.
		// workload clusters kubeconfig details need to be added for remote App so that kapp-controller on management
		// cluster can reconcile and push the addon/app to the workload cluster
		//
		// if remoteApp {
		// 	clusterKubeconfigDetails := util.GetClusterKubeconfigSecretDetails(remoteCluster)

		// 	app.Spec.Cluster = &kappctrl.AppCluster{
		// 		KubeconfigSecretRef: &kappctrl.AppClusterKubeconfigSecretRef{
		// 			Name: clusterKubeconfigDetails.Name,
		// 			Key:  clusterKubeconfigDetails.Key,
		// 		},
		// 	}
		// } else {
		app.Spec.ServiceAccountName = addonconstants.TKGAddonsAppServiceAccount
		// }

		app.Spec.Fetch = []kappctrl.AppFetch{
			{
				Image: &kappctrl.AppFetchImage{
					URL: fmt.Sprintf("%s/%s:%s", imageRepository, addonConfig.TemplatesImagePath, addonConfig.TemplatesImageTag),
				},
			},
		}

		app.Spec.Template = []kappctrl.AppTemplate{
			{
				Ytt: &kappctrl.AppTemplateYtt{
					IgnoreUnknownComments: true,
					Strict:                false,
					Inline: &kappctrl.AppFetchInline{
						PathsFrom: []kappctrl.AppFetchInlineSource{
							{
								SecretRef: &kappctrl.AppFetchInlineSourceRef{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: util.GenerateAppSecretNameFromAddonSecret(addonSecret),
									},
								},
							},
						},
					},
				},
			},
		}

		app.Spec.Deploy = []kappctrl.AppDeploy{
			{
				Kapp: &kappctrl.AppDeployKapp{},
			},
		}

		return nil
	}

	_, err = controllerutil.CreateOrPatch(context.TODO(), client, app, appMutateFn)
	if err != nil {
		klog.Errorf("CreateOrPatch failed. Err: %v", err)
		return err
	}

	klog.V(2).Info("CreateOrPatch succeeded")

	return nil
}
*/

// Next version...
/*
// InstallFromUser install from user defined values
func (k *Kapp) InstallFromUser(appCrd *AppCrdInput) error {

	klog.V(2).Infof("extension = %s", appCrd.Name)

	// create k8s client
	config, err := clientcmd.BuildConfigFromFlags("", k.config.Kubeconfig)
	if err != nil {
		klog.Errorf("BuildConfigFromFlags failed. Err: %v", err)
		return err
	}
	client, err := client.New(config, client.Options{})
	if err != nil {
		klog.Errorf("client.New failed. Err: %v", err)
		return err
	}

	// build out kappctrl.App
	app := &kappctrl.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appCrd.Name,
			Namespace: appCrd.Namespace,
		},
	}

	appMutateFn := func() error {
		if app.ObjectMeta.Annotations == nil {
			app.ObjectMeta.Annotations = make(map[string]string)
		}

		// app.ObjectMeta.Annotations[addontypes.AddonTypeAnnotation] = fmt.Sprintf("%s/%s", addonConfig.Category, addonName)
		// app.ObjectMeta.Annotations[addontypes.AddonNameAnnotation] = addonSecret.Name
		// app.ObjectMeta.Annotations[addontypes.AddonNamespaceAnnotation] = addonSecret.Namespace

		// remoteApp means App is not present on local workload cluster. It is present in the remote management cluster.
		// workload clusters kubeconfig details need to be added for remote App so that kapp-controller on management
		// cluster can reconcile and push the addon/app to the workload cluster
		//
		// if remoteApp {
		// 	clusterKubeconfigDetails := util.GetClusterKubeconfigSecretDetails(remoteCluster)

		// 	app.Spec.Cluster = &kappctrl.AppCluster{
		// 		KubeconfigSecretRef: &kappctrl.AppClusterKubeconfigSecretRef{
		// 			Name: clusterKubeconfigDetails.Name,
		// 			Key:  clusterKubeconfigDetails.Key,
		// 		},
		// 	}
		// } else {
		app.Spec.ServiceAccountName = appCrd.Name + DefaultAppCrdServiceAccountPostfix
		// }

		app.Spec.Fetch = []kappctrl.AppFetch{
			{
				Image: &kappctrl.AppFetchImage{
					URL: appCrd.URL,
				},
			},
		}

		app.Spec.Template = []kappctrl.AppTemplate{
			{
				Ytt: &kappctrl.AppTemplateYtt{
					IgnoreUnknownComments: true,
					Strict:                false,
					Inline: &kappctrl.AppFetchInline{
						Paths: appCrd.Paths,
					},
				},
			},
		}

		app.Spec.Deploy = []kappctrl.AppDeploy{
			{
				Kapp: &kappctrl.AppDeployKapp{},
			},
		}

		return nil
	}

	_, err = controllerutil.CreateOrPatch(context.TODO(), client, app, appMutateFn)
	if err != nil {
		klog.Errorf("CreateOrPatch failed. Err: %v", err)
		return err
	}

	klog.V(2).Info("CreateOrPatch succeeded")

	return err
}
*/
