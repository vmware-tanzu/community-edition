// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package github

import (
	"os"

	yaml "github.com/ghodss/yaml"
	klog "k8s.io/klog/v2"
)

// NewGitHubConfig generates a Config object
func NewGitHubConfig() *Config {

	cfg := &Config{
		GitHubRepo:         DefaultGitHubRepo,
		GitHubBranchTag:    DefaultGitHubBranchTagMain,
		ExtensionDirectory: DefaultGitHubExtensionsDirectory,
		GitHubURI:          DefaultGitHubRepo + "/" + DefaultGitHubBranchTagMain + "/" + DefaultGitHubExtensionsDirectory,
	}

	return cfg
}

// FromEnv initializes the configuration object with values
// obtained from environment variables. If an environment variable is set
// for a property that's already initialized, the environment variable's value
// takes precedence.
func (cfg *Config) FromEnv() error {

	if v := os.Getenv("TCE_EXTENSION_GITHUB_REPO"); v != "" {
		cfg.GitHubRepo = v
	}
	if v := os.Getenv("TCE_EXTENSION_GITHUB_BRANCH_TAG"); v != "" {
		cfg.GitHubBranchTag = v
	}
	if v := os.Getenv("TCE_EXTENSION_GITHUB_TOKEN"); v != "" {
		cfg.Token = v
	}
	if v := os.Getenv("TCE_EXTENSION_DIRECTORY"); v != "" {
		cfg.ExtensionDirectory = v
	}

	// replace "latest" which is human readable to the GitHub "main" branch
	cfg.originalBranchTag = cfg.GitHubBranchTag
	klog.V(4).Infof("cfg.originalBranchTag = %s", cfg.originalBranchTag)
	klog.V(4).Infof("cfg.GitHubBranchTag = %s", cfg.GitHubBranchTag)
	if cfg.GitHubBranchTag == DefaultGitHubLatest {
		klog.V(4).Infof("Replacing %s with %s", DefaultGitHubLatest, DefaultGitHubBranchTagMain)
		cfg.GitHubBranchTag = DefaultGitHubBranchTagMain
	}

	cfg.GitHubURI = cfg.GitHubRepo + "/" + cfg.GitHubBranchTag + "/" + cfg.ExtensionDirectory

	klog.V(4).Infof("GitHubBranchTag = %s", cfg.GitHubBranchTag)
	klog.V(4).Infof("originalBranchTag = %s", cfg.originalBranchTag)
	klog.V(6).Infof("token = %s", cfg.Token)
	klog.V(4).Infof("GitHubURI = %s", cfg.GitHubURI)

	return nil
}

// InitGitHubConfig inits the Config and also checks Environment variables
func InitGitHubConfig(byConfig []byte) (*Config, error) {

	klog.V(4).Infof("Config Data:")
	klog.V(4).Infof("%s", string(byConfig))

	cfg := NewGitHubConfig()

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

	klog.V(2).Info("Config (GitHub) initialized")
	return cfg, nil
}
