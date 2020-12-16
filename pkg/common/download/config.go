// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package download

import (
	"os"

	yaml "github.com/ghodss/yaml"
	klog "k8s.io/klog/v2"
)

// NewDownloadConfig generates a Config object
func NewDownloadConfig() *Config {

	cfg := &Config{
		GitHubOrg:  DefaultGitHubOrg,
		GitHubRepo: DefaultGitHubRepo,
	}

	return cfg
}

// FromEnv initializes the configuration object with values
// obtained from environment variables. If an environment variable is set
// for a property that's already initialized, the environment variable's value
// takes precedence.
func (cfg *Config) FromEnv() error {

	if v := os.Getenv("TCE_EXTENSION_DOWNLOAD_ORG"); v != "" {
		cfg.GitHubOrg = v
	}
	if v := os.Getenv("TCE_EXTENSION_DOWNLOAD_REPO"); v != "" {
		cfg.GitHubRepo = v
	}

	klog.V(4).Infof("GitHubOrg = %s", cfg.GitHubOrg)
	klog.V(4).Infof("GitHubRepo = %s", cfg.GitHubRepo)

	return nil
}

// InitDownloadConfig inits the Config and also checks Environment variables
func InitDownloadConfig(byConfig []byte) (*Config, error) {

	klog.V(4).Infof("Config Data:")
	klog.V(4).Infof("%s", string(byConfig))

	cfg := NewDownloadConfig()

	// pull the version from the config.yaml
	err := yaml.Unmarshal(byConfig, &cfg)
	if err != nil {
		klog.Errorf("Unmarshal failed. Err: ", err)
		return nil, err
	}

	// Env Vars should override config file entries if present
	if err := cfg.FromEnv(); err != nil {
		klog.Errorf("FromEnv failed: %s", err)
		return nil, err
	}

	klog.V(2).Info("Config (Download) initialized")
	return cfg, nil
}
