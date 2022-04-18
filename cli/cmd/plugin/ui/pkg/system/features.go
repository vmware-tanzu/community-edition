// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package system has functions to get information about the system
package system

import (
	"github.com/pkg/errors"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
	"github.com/vmware-tanzu/tanzu-framework/apis/config/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/config"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/features"
)

// FeatureFlags gets the configured feature flags for the installation.
func FeatureFlags() (models.Features, error) {
	cfg, err := config.GetClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get client configuration")
	}

	result := models.Features{}
	for pluginName, featureMap := range cfg.ClientOptions.Features {
		result[pluginName] = convertFeatureMap(featureMap)
	}

	return result, nil
}

// convertFeatureMap converts a config file v1alpha1.FeatureMap to payload models.FeatureMap both of which are just hash maps
func convertFeatureMap(featureMap v1alpha1.FeatureMap) models.FeatureMap {
	result := models.FeatureMap{}
	for key, value := range featureMap {
		result[key] = value
	}
	return result
}

// Edition gets the configured edition for the installation.
func Edition() (string, error) {
	featuresClient, err := features.New(GetConfigDir(), "")
	if err != nil {
		return "", err
	}

	tanzuEdition, err := featuresClient.GetFeatureFlag("edition")
	if err != nil {
		return "", err
	}

	return tanzuEdition, nil
}
