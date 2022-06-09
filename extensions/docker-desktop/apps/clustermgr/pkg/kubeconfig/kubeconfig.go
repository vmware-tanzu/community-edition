// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package kubeconfig has slightly more generic kubeconfig helpers and
// minimal dependencies on the rest of kind
package kubeconfig

import (
	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/config"
	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/kubeconfig/internal/kubeconfig"
)

// AddConfig adds our cluster kubeconfig to the existing kubeconfig path.
func AddConfig(kubeconfigToAddPath, kubeconfigPath string) error {
	origConf, err := kubeconfig.Read(kubeconfigToAddPath)
	if err != nil {
		return err
	}

	// Modify the kubeconfig
	getTCEConfig(origConf)

	return kubeconfig.WriteMerged(origConf, kubeconfigPath)
}

// RemoveConfig removes our cluster kubeconfig.
func RemoveConfig(kubeconfigToAddPath, kubeconfigPath string) error {
	origConf, err := kubeconfig.Read(kubeconfigToAddPath)
	if err != nil {
		return err
	}
	return kubeconfig.RemoveKIND(origConf.CurrentContext, kubeconfigPath)
}

// RemoveNamedConfig removes a named kubeconfig entry.
func RemoveNamedConfig(context, kubeconfigPath string) error {
	return kubeconfig.RemoveKIND(context, kubeconfigPath)
}

// GetConfig gets the kubeconfig from the given path.
func GetConfig(kubeconfigPath string) ([]byte, error) {
	cfg, err := kubeconfig.Read(kubeconfigPath)
	if err != nil {
		return nil, err
	}
	// Modify the kubeconfig
	getTCEConfig(cfg)
	encoded, err := kubeconfig.Encode(cfg)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

// getTCEConfig gets our cluster kubeconfig.
func getTCEConfig(cfg *kubeconfig.Config) {
	// We change all the config Names in the kubeconfig to be the Cluster name
	if len(cfg.Clusters) > 0 &&
		len(cfg.Users) > 0 &&
		len(cfg.Contexts) > 0 {
		key := config.DefaultClusterName
		cfg.Clusters[0].Name = key
		cfg.Users[0].Name = key
		cfg.Contexts[0].Name = key
		cfg.Contexts[0].Context.User = key
		cfg.Contexts[0].Context.Cluster = key
		cfg.CurrentContext = key
	}
}
