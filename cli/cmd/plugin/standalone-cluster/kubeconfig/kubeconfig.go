// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package kubeconfig manages the user's kubeconfig. Its primary use is to merge a newly created kubeconfig into the
// default kubeconfig location (paths.Join(os.Home, ".kube", "config") and then switching the in-use context to the
// newly merged one.
package kubeconfig

import (
	"os"
	"strings"

	"github.com/ghodss/yaml"
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

	output, err = yaml.JSONToYAML(encode)
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
