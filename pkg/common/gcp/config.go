// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package gcp

import (
	"os"

	yaml "github.com/ghodss/yaml"
	klog "k8s.io/klog/v2"
)

// NewConfig generates a Config object
func NewConfig() *Config {

	cfg := &Config{
		TceBucket:         DefaultTceBucket,
		MetadataDirectory: DefaultMetadataDirectory,
		VersionTag:        DefaultMetadataVersion,
		MetadataFileName:  DefaultMetadataFileName,
		ReleasesFileName:  DefaultReleasesFileName,
	}

	return cfg
}

// FromEnv initializes the configuration object with values
// obtained from environment variables. If an environment variable is set
// for a property that's already initialized, the environment variable's value
// takes precedence.
func (cfg *Config) FromEnv() error {

	if v := os.Getenv("TCE_EXTENSION_GCP_BUCKET"); v != "" {
		cfg.TceBucket = v
	}
	if v := os.Getenv("TCE_EXTENSION_GCP_METADATA_DIRECTORY"); v != "" {
		cfg.MetadataDirectory = v
	}
	if v := os.Getenv("TCE_EXTENSION_GCP_METADATA_FILENAME"); v != "" {
		cfg.MetadataFileName = v
	}
	if v := os.Getenv("TCE_EXTENSION_GCP_METADATA_VERSION"); v != "" {
		cfg.VersionTag = v
	}
	if v := os.Getenv("TCE_EXTENSION_GCP_METADATA_RELEASES"); v != "" {
		cfg.ReleasesFileName = v
	}

	klog.V(2).Infof("TceBucket = %s", cfg.TceBucket)
	klog.V(2).Infof("MetadataDirectory = %s", cfg.MetadataDirectory)
	klog.V(2).Infof("MetadataFileName = %s", cfg.MetadataFileName)
	klog.V(2).Infof("VersionTag = %s", cfg.VersionTag)
	klog.V(2).Infof("ReleasesFileName = %s", cfg.ReleasesFileName)

	return nil
}

// InitBucketConfig inits the Config and also checks Environment variables
func InitBucketConfig(byConfig []byte) (*Config, error) {

	klog.V(4).Infof("Config Data:")
	klog.V(4).Infof("%s", string(byConfig))

	cfg := NewConfig()

	// pull the version from the config.yaml
	err := yaml.Unmarshal(byConfig, &cfg)
	if err != nil {
		klog.Errorf("Unmarshal failed. Err: ", err)
		return nil, err
	}

	if cfg == nil {
		cfg = NewConfig()
	}

	// Env Vars should override config file entries if present
	if err := cfg.FromEnv(); err != nil {
		klog.Errorf("FromEnv failed: %s", err)
		return nil, err
	}

	klog.V(2).Info("Config (GCP) initialized")
	return cfg, nil
}
