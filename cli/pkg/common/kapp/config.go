// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kapp

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
)

// NewConfig generates a Config object
func NewConfig() *Config {
	cfg := &Config{
		ExtensionServiceAccountPostfix: DefaultServiceAccountPostfix,
		ExtensionRoleBindingPostfix:    DefaultRoleBindingPostfix,
	}

	return cfg
}

// FromEnv initializes the configuration object with values
// obtained from environment variables. If an environment variable is set
// for a property that's already initialized, the environment variable's value
// takes precedence.
func (cfg *Config) FromEnv() error {
	if v := os.Getenv("TCE_EXTENSION_SERVICEACCOUNT_POSTFIX"); v != "" {
		cfg.ExtensionServiceAccountPostfix = v
	}
	if v := os.Getenv("TCE_EXTENSION_ROLEBINDING_POSTFIX"); v != "" {
		cfg.ExtensionRoleBindingPostfix = v
	}

	home := homedir.HomeDir()
	klog.V(3).Infof("HomeDir = %s", home)

	if v := os.Getenv("TCE_EXTENSION_KUBECONFIG"); v != "" {
		cfg.Kubeconfig = v
	} else if home != "" {
		if _, err := os.Stat(filepath.Join(home, ".kube", "config")); err == nil {
			cfg.Kubeconfig = filepath.Join(home, ".kube", "config")
		} else if _, err := os.Stat(filepath.Join(home, ".kube-tkg", "config")); err == nil {
			cfg.Kubeconfig = filepath.Join(home, ".kube-tkg", "config")
		} else if _, err := os.Stat("/var/run/kubernetes/admin.kubeconfig"); err == nil {
			cfg.Kubeconfig = "/var/run/kubernetes/admin.kubeconfig"
		}
	}

	klog.V(4).Infof("Kubeconfig = %s", cfg.Kubeconfig)
	klog.V(4).Infof("ExtensionServiceAccountPostfix = %s", cfg.ExtensionServiceAccountPostfix)
	klog.V(4).Infof("ExtensionRoleBindingPostfix = %s", cfg.ExtensionRoleBindingPostfix)

	return nil
}

// InitKappConfig inits the Config and also checks Environment variables
func InitKappConfig() (*Config, error) {
	cfg := NewConfig()

	// Env Vars should override config file entries if present
	if err := cfg.FromEnv(); err != nil {
		klog.Errorf("FromEnv failed: %s", err)
		return nil, err
	}

	klog.V(2).Info("Config (Kapp) initialized")
	return cfg, nil
}
