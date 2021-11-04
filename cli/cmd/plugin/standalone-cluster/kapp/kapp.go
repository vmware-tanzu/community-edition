// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package kapp installs and manages instances of kapp-controller into Kubernetes clusters.
package kapp

import (
	"context"
	"fmt"
	"log"
	"strings"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	apiRegv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
)

const (
	deploymentKind               = "Deployment"
	metadataKey                  = "metadata"
	nameKey                      = "name"
	namespaceKey                 = "namespace"
	kappControllerDeploymentName = "kapp-controller"
)

// KappClient is a client for interfacing with kapp-controller.
//nolint:golint
type KappClient struct {
	dynClient  dynamic.Interface
	clientSet  kubernetes.Interface
	restMapper meta.RESTMapper
	scheme     *runtime.Scheme
}

// KappInstallOpts contains information about how to install kapp-controller.
//nolint:golint
type KappInstallOpts struct {
	// MergedManifests represents the final (serialized) data for installing kapp-controller and its ancillary
	// resources. This assumes all objects are merged into one byte array. Think of this like a single YAML file
	// containing multiple Kubernetes resources.
	// It assumes the manifest are in their final state, meaning you could kubectl apply them.
	// If template rendering is required (e.g. ytt) this should be done before setting this value.
	MergedManifests []byte
	// Manifests represents the final (serialized) data for installing kapp-controller and its ancillary
	// resources. This assumes each object is contained in its own byte array. Think of this like a many YAML files
	// each containing one Kubernetes resources.
	// It assumes the manifest are in their final state, meaning you could kubectl apply them.
	// If template rendering is required (e.g. ytt) this should be done before setting this value.
	Manifests [][]byte
}

// KappManager defines the interface for performing kapp operations.
//nolint:golint
type KappManager interface {
	// Install installs kapp-controller into the cluster. When successful, it returns the Deployment object that
	// manages the kapp-controller pod.
	Install(opts KappInstallOpts) (*v1.Deployment, error)
	// Status retrieves the pod status for kapp-controller. It expects to be passed the namespace and name for the
	// kapp-controller Deployment object. If it cannot talk to the cluster, that status is reported. If the
	// pod cannot be resolved, a status of not created is reported. Otherwise, the exact status message is returned.
	Status(ns, name string) string
}

// New instantiates a new KappManager.
func New(kubeconfigBytes []byte) (KappManager, error) {
	config, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (*clientcmdapi.Config, error) {
		return clientcmd.Load(kubeconfigBytes)
	})
	// TODO(joshrosso): figure out what to do here
	if err != nil {
		return nil, fmt.Errorf("could not build config using provided kubeconfig: %s", err.Error())
	}

	client, err := dynamic.NewForConfig(config)
	// TODO(joshrosso): figure out what to do here
	if err != nil {
		return nil, fmt.Errorf("could not build dynamic client with config: %s", err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(config)
	// TODO(joshrosso): figure out what to do here
	if err != nil {
		return nil, fmt.Errorf("could not build new client set from config: %s", err.Error())
	}

	// setting up restMapper is an expensive operation, so do it here and re-use it
	// in method invocations
	groupResources, err := restmapper.GetAPIGroupResources(clientSet.Discovery())
	// TODO(joshrosso): figure out what to do here
	if err != nil {
		return nil, fmt.Errorf("could not build restMapper: %s", err.Error())
	}

	rm := restmapper.NewDiscoveryRESTMapper(groupResources)

	sch := runtime.NewScheme()
	// ensure CRD Definitions are detected
	err = apiv1.AddToScheme(sch)
	// TODO(joshrosso): figure out what to do here
	if err != nil {
		return nil, fmt.Errorf("could not detect apiv1 CRD definitions from runtime scheme: %s", err.Error())
	}
	err = v1.AddToScheme(sch)
	// TODO(joshrosso): figure out what to do here
	if err != nil {
		return nil, fmt.Errorf("could not detect v1 CRD definitions from runtime scheme: %s", err.Error())
	}
	err = corev1.AddToScheme(sch)
	// TODO(joshrosso): figure out what to do here
	if err != nil {
		return nil, fmt.Errorf("could not detect corev1 CRD definitions from runtime scheme: %s", err.Error())
	}
	err = rbacv1.AddToScheme(sch)
	// TODO(joshrosso): figure out what to do here
	if err != nil {
		return nil, fmt.Errorf("could not detect rbacv1 CRD definitions from runtime scheme: %s", err.Error())
	}
	err = apiRegv1.AddToScheme(sch)
	// TODO(joshrosso): figure out what to do here
	if err != nil {
		return nil, fmt.Errorf("could not detect apiRegv1 CRD definitions from runtime scheme: %s", err.Error())
	}

	return KappClient{
		dynClient:  client,
		clientSet:  clientSet,
		restMapper: rm,
		scheme:     sch,
	}, nil
}

// Install performs and installation.
func (k KappClient) Install(opts KappInstallOpts) (*v1.Deployment, error) {
	if opts.MergedManifests == nil && opts.Manifests == nil {
		return nil, fmt.Errorf("no objects were provided to install")
	}
	var objects []runtime.Object
	if opts.MergedManifests != nil {
		objects = parseMergedObjects(k.scheme, opts.MergedManifests)
	} else {
		objects = createObjectList(k.scheme, opts.Manifests)
	}

	var kappDeployment v1.Deployment
	for _, obj := range objects {
		createdObj, err := applyObject(k, obj)
		if err != nil {
			return nil, err
		}

		// determine if the returned object is kapp-controller, if so plan to return it
		// TODO(joshrosso): Really ugly embedded ifs, we should make this better
		if createdObj.GetObjectKind().GroupVersionKind().Kind == deploymentKind {
			if createdObj.Object[metadataKey].(map[string]interface{})[nameKey] != nil {
				if createdObj.Object[metadataKey].(map[string]interface{})[nameKey].(string) == kappControllerDeploymentName {
					err = runtime.DefaultUnstructuredConverter.FromUnstructured(createdObj.Object, &kappDeployment)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}
	return &kappDeployment, nil
}

// Status gets the status of a package.
func (k KappClient) Status(ns, name string) string {
	pods, err := k.clientSet.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "failed to talk to cluster"
	}

	var kappPod *corev1.Pod
	// super bad, should not be iterating through pods on every invocation
	for i := range pods.Items {
		// pod name will always prefix with deployment name
		pod := pods.Items[i]
		if strings.HasPrefix(pod.Name, name) {
			kappPod = &pod
		}
	}

	if kappPod == nil {
		return "Not created"
	}

	return string(kappPod.Status.Phase)
}

// applyObject takes a runtime.Object and converts it into an unstructured object. It then
// uses the dynamic client to apply the object to the cluster. If the namespace field is nil,
// it applies the object cluster wide, if it contains a string, it applies it in the
// appropriate namespace.
func applyObject(k KappClient, obj runtime.Object) (*unstructured.Unstructured, error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
	uObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}

	objectBody := &unstructured.Unstructured{Object: uObj}

	mapping, _ := k.restMapper.RESTMapping(gk, gvk.Version)

	nsInterface := uObj[metadataKey].(map[string]interface{})[namespaceKey]
	var createObj *unstructured.Unstructured
	if nsInterface != nil {
		ns := nsInterface.(string)
		createObj, err = k.dynClient.
			Resource(schema.GroupVersionResource{
				Group:    gvk.Group,
				Version:  gvk.Version,
				Resource: mapping.Resource.Resource,
			}).Namespace(ns).
			Create(context.TODO(), objectBody, metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}
	} else {
		createObj, err = k.dynClient.
			Resource(schema.GroupVersionResource{
				Group:    gvk.Group,
				Version:  gvk.Version,
				Resource: mapping.Resource.Resource,
			}).
			Create(context.TODO(), objectBody, metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}
	}

	return createObj, nil
}

// parseMergedObjects takes multiple YAML objects, separated by '---' and returns a list of runtime objects.
func parseMergedObjects(sch *runtime.Scheme, fileR []byte) []runtime.Object {
	fileAsString := string(fileR)
	sepYamlfiles := strings.Split(fileAsString, "---")
	retVal := make([]runtime.Object, 0, len(sepYamlfiles))
	for _, f := range sepYamlfiles {
		if f == "\n" || f == "" {
			// ignore empty cases
			continue
		}

		decode := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode
		obj, _, err := decode([]byte(f), nil, nil)

		if err != nil {
			log.Println(fmt.Sprintf("Error while decoding YAML object. Err was: %s", err))
			continue
		}
		retVal = append(retVal, obj)
	}
	return retVal
}

// createObjectList returns a list of runtime objects based on objects living in a list of byte arrays
func createObjectList(sch *runtime.Scheme, objects [][]byte) []runtime.Object {
	retVal := make([]runtime.Object, 0, len(objects))
	for _, o := range objects {
		decode := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode
		obj, _, err := decode(o, nil, nil)

		if err != nil {
			log.Println(fmt.Sprintf("Error while decoding YAML object. Err was: %s", err))
			continue
		}
		retVal = append(retVal, obj)
	}
	return retVal
}
