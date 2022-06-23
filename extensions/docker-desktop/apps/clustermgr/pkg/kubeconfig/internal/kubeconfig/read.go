// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kubeconfig

import (
	"io"
	"os"

	yaml "gopkg.in/yaml.v3"

	"github.com/vmware-tanzu/community-edition/errors"
)

// KINDFromRawKubeadm returns a kind kubeconfig derived from the raw kubeadm kubeconfig,
// the kind clusterName, and the server.
// server is ignored if unset.
func KINDFromRawKubeadm(rawKubeadmKubeConfig, clusterName, server string) (*Config, error) {
	cfg := &Config{}
	if err := yaml.Unmarshal([]byte(rawKubeadmKubeConfig), cfg); err != nil {
		return nil, err
	}

	// verify assumptions about kubeadm kubeconfigs
	if err := checkKubeadmExpectations(cfg); err != nil {
		return nil, err
	}

	// compute unique kubeconfig key for this cluster
	key := KINDClusterKey(clusterName)

	// use the unique key for all named references
	cfg.Clusters[0].Name = key
	cfg.Users[0].Name = key
	cfg.Contexts[0].Name = key
	cfg.Contexts[0].Context.User = key
	cfg.Contexts[0].Context.Cluster = key
	cfg.CurrentContext = key

	// patch server field if server was set
	if server != "" {
		cfg.Clusters[0].Cluster.Server = server
	}

	return cfg, nil
}

// Read makes a public copy of a config to use outside
func Read(configPath string) (*Config, error) {
	return read(configPath)
}

// read loads a KUBECONFIG file from configPath
func read(configPath string) (*Config, error) {
	// try to open, return default if no such file
	f, err := os.Open(configPath)
	if os.IsNotExist(err) {
		return &Config{}, nil
	} else if err != nil {
		return nil, errors.NewReadFailed(err, "failed to read config file %s", configPath)
	}

	// otherwise read in and deserialize
	cfg := &Config{}
	rawExisting, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.NewIOFailed(err, "failed to read config file data")
	}
	if err := yaml.Unmarshal(rawExisting, cfg); err != nil {
		return nil, errors.NewMarshallingFailed(err, "faled to marshall kubeconfig file %s", configPath)
	}

	return cfg, nil
}
