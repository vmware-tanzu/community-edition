// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/adrg/xdg"
	klog "k8s.io/klog/v2"
)

// NewConfig generates a Config object
func NewConfig() *Config {

	configFile := filepath.Join(xdg.DataHome, "tanzu-repository", DefaultConfigFile)
	byConfig, err := ioutil.ReadFile(configFile)
	if err != nil {
		klog.Fatalf("ReadFile(%s) failed. Err:", configFile, err)
	}

	cfg := &Config{
		configFile: configFile,
		byRaw:      byConfig,
	}
	cfg.githubToken = cfg.GetToken()

	return cfg
}

// InitConfig inits the Config and also checks Environment variables
func InitConfig() *Config {

	cfg := NewConfig()

	klog.V(4).Infof("configFile = %s", cfg.configFile)
	klog.V(6).Infof("githubToken = %s", cfg.githubToken)
	klog.V(2).Info("Config (Config) initialized")

	return cfg
}
