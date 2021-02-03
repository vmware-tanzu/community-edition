// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"path/filepath"

	"github.com/adrg/xdg"
	yaml "github.com/ghodss/yaml"
	klog "k8s.io/klog/v2"
)

// NewConfig generates a Config object
func NewConfig() *Config {
	configFile := filepath.Join(xdg.DataHome, "tanzu-repository", DefaultConfigFile)

	cfg := &Config{
		configFile: configFile,
	}

	return cfg
}

// InitConfig inits the Config and also checks Environment variables
func InitConfig(byConfig []byte) (*Config, error) {

	klog.V(4).Infof("Config Data:")
	klog.V(4).Infof("%s", string(byConfig))

	cfg := NewConfig()

	// pull the version from the config.yaml
	err := yaml.Unmarshal(byConfig, &cfg)
	if err != nil {
		klog.Errorf("Unmarshal failed. Err: ", err)
		return nil, err
	}

	klog.V(4).Infof("releaseVersion = %s", cfg.ReleaseVersion)
	klog.V(6).Infof("githubToken = %s", cfg.GithubToken)
	klog.V(2).Info("Config (Config) initialized")

	return cfg, nil
}
