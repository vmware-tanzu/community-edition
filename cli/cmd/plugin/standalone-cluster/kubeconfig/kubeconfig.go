// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package kubeconfig manages the user's kubeconfig. Its primary use is to merge a newly created kubeconfig into the
// default kubeconfig location (paths.Join(os.Home, ".kube", "config") and then switching the in-use context to the
// newly merged one.
package kubeconfig

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	clientcmdapilatest "k8s.io/client-go/tools/clientcmd/api/latest"
)

// KubeConfig contains information about the kubeconfig location.
type KubeConfig struct {
	defaultConfigLocation string
}

// KubeConfigMgr exposes operations that can be to a user's kubeconfig
//nolint:golint
type KubeConfigMgr interface {
	// MergeToDefaultConfig takes a kubeconfig file and merges it into the default kube config location.
	// It does not mutate existing cluster configs, unless there is a record with the same name. In this
	// case, the existing cluster details are overwritten.
	MergeToDefaultConfig(kubeconfigPath string) error
	// SetCurrentContext changes the kubeconfig context (`current-context` value) to the name passed in.
	SetCurrentContext(name string) error
	// TODO(joshrosso): we should considering introducing a backup capability that is called as part of this process.
}

// NewManager returns a KubeConfigMgr implemented by KubeConfig.
// It automatically resolves the default config location for the user's
// machine and stores it locally for future context and merge operations.
func NewManager() KubeConfigMgr {
	kc := &KubeConfig{
		defaultConfigLocation: clientcmd.RecommendedHomeFile,
	}

	return kc
}

// MergeToDefaultConfig merges configuration to the kubeconfig file.
func (kc *KubeConfig) MergeToDefaultConfig(kubeconfigPath string) error {
	rules := clientcmd.ClientConfigLoadingRules{
		Precedence: []string{kc.defaultConfigLocation, kubeconfigPath},
	}
	loadedRules, err := rules.Load()
	if err != nil {
		return err
	}

	output, err := encodeConfig(loadedRules)
	if err != nil {
		return err
	}

	if err := writeKubeConfigFile(kc.defaultConfigLocation, output, 0655); err != nil {
		return err
	}
	return nil
}

// SetCurrentContext sets the current kubeconfig context.
func (kc *KubeConfig) SetCurrentContext(name string) error {
	rules := clientcmd.ClientConfigLoadingRules{
		Precedence: []string{kc.defaultConfigLocation},
	}
	loadedRules, err := rules.Load()
	if err != nil {
		return err
	}
	loadedRules.CurrentContext = name

	output, err := encodeConfig(loadedRules)
	if err != nil {
		return err
	}

	if err := writeKubeConfigFile(kc.defaultConfigLocation, output, 0655); err != nil {
		return err
	}
	return nil
}

// encodeConfig takes the [kube]config struct from the Kubernetes API and returns the YAML
// representation in byte format.
func encodeConfig(config *clientcmdapi.Config) ([]byte, error) {
	var err error
	var output []byte

	encode, err := runtime.Encode(clientcmdapilatest.Codec, config)
	if err != nil {
		return nil, err
	}
	output, err = jSONToYAML(encode)
	return output, err
}

// writeKubeConfigFile writes the encoded ([]byte representation) config into the kubeconfig
// default directory. If this directory does not exist, it is created.
func writeKubeConfigFile(path string, data []byte, perm os.FileMode) error {
	if i := strings.LastIndex(path, "/"); i != -1 {
		kubeDir := path[:i+1]
		if _, err := os.Stat(kubeDir); os.IsNotExist(err) {
			if err := os.MkdirAll(kubeDir, 0744); err != nil {
				return err
			}
		}
	}
	return os.WriteFile(path, data, perm)
}

// jSONTOYAML converts JSON to YAML. It is sourced from the github.com/ghodss/yaml project
// and brought into this package in order to be able to purely import gopkg.in/yaml.v2 rather
// than a wrapper dependency. This is as private function to limit its use to this package.
// If it ever needs to be used more generally for this project, extract it out of this package.
func jSONToYAML(j []byte) ([]byte, error) {
	// Convert the JSON to an object.
	var jsonObj interface{}
	// We are using yaml.Unmarshal here (instead of json.Unmarshal) because the
	// Go JSON library doesn't try to pick the right number type (int, float,
	// etc.) when unmarshalling to interface{}, it just picks float64
	// universally. go-yaml does go through the effort of picking the right
	// number type, so we can preserve number type throughout this process.
	err := yaml.Unmarshal(j, &jsonObj)
	if err != nil {
		return nil, err
	}

	// Marshal this object into YAML.
	return yaml.Marshal(jsonObj)
}
