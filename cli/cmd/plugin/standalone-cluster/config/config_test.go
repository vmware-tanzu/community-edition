// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//nolint:goconst
package config

import (
	"bytes"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestInitializeConfigurationNoName(t *testing.T) {
	_, err := InitializeConfiguration(make(map[string]string))
	if err == nil {
		t.Error("initialization should fail if no cluster name provided")
	}
}

func TestInitializeConfigurationDefaults(t *testing.T) {
	args := map[string]string{"clustername": "test"}
	config, err := InitializeConfiguration(args)
	if err != nil {
		t.Error("initialization should pass")
	}

	if config.ClusterName != "test" {
		t.Errorf("expected ClusterName to be 'test', was actually: %q", config.ClusterName)
	}

	if config.Cni != defaultConfigValues[Cni] {
		t.Errorf("expected default Cni value, was: %q", config.Cni)
	}

	if config.TkrLocation != defaultConfigValues[TKRLocation] {
		t.Errorf("expected default TkrLocation, was: %q", config.TkrLocation)
	}

	if config.Provider != defaultConfigValues[Provider] {
		t.Errorf("expected default Provider, was: %q", config.Provider)
	}

	if config.PodCidr != defaultConfigValues[PodCIDR] {
		t.Errorf("expected default PodCidr, was: %q", config.PodCidr)
	}

	if config.ServiceCidr != defaultConfigValues[ServiceCIDR] {
		t.Errorf("expected default ServiceCidr, was: %q", config.ServiceCidr)
	}
}

func TestInitializeConfigurationEnvVariables(t *testing.T) {
	os.Setenv("TANZU_PROVIDER", "test_provider")
	os.Setenv("TANZU_CLUSTER_NAME", "test2")
	config, err := InitializeConfiguration(make(map[string]string))
	if err != nil {
		t.Error("initialization should pass")
	}

	if config.ClusterName != "test2" {
		t.Errorf("expected ClusterName to be 'test2', was actually: %q", config.ClusterName)
	}

	if config.Provider != "test_provider" {
		t.Errorf("expected Provider to be set to 'test_provider', was: %q", config.Provider)
	}

	if config.Cni != defaultConfigValues[Cni] {
		t.Errorf("expected default Cni value, was: %q", config.Cni)
	}

	if config.TkrLocation != defaultConfigValues[TKRLocation] {
		t.Errorf("expected default TkrLocation, was: %q", config.TkrLocation)
	}

	if config.PodCidr != defaultConfigValues[PodCIDR] {
		t.Errorf("expected default PodCidr, was: %q", config.PodCidr)
	}

	if config.ServiceCidr != defaultConfigValues[ServiceCIDR] {
		t.Errorf("expected default ServiceCidr, was: %q", config.ServiceCidr)
	}
}

func TestInitializeConfigurationArgsTakePrecedent(t *testing.T) {
	os.Setenv("TANZU_PROVIDER", "test_provider")
	os.Setenv("TANZU_CLUSTER_NAME", "test2")
	args := map[string]string{"clustername": "test"}
	config, err := InitializeConfiguration(args)
	if err != nil {
		t.Error("initialization should pass")
	}

	if config.ClusterName != "test" {
		t.Errorf("expected ClusterName to be 'test', was actually: %q", config.ClusterName)
	}

	if config.Provider != "test_provider" {
		t.Errorf("expected Provider to be set to 'test_provider', was: %q", config.Provider)
	}

	if config.Cni != defaultConfigValues[Cni] {
		t.Errorf("expected default Cni value, was: %q", config.Cni)
	}

	if config.TkrLocation != defaultConfigValues[TKRLocation] {
		t.Errorf("expected default TkrLocation, was: %q", config.TkrLocation)
	}

	if config.PodCidr != defaultConfigValues[PodCIDR] {
		t.Errorf("expected default PodCidr, was: %q", config.PodCidr)
	}

	if config.ServiceCidr != defaultConfigValues[ServiceCIDR] {
		t.Errorf("expected default ServiceCidr, was: %q", config.ServiceCidr)
	}
}

func TestInitializeConfigurationFromConfigFile(t *testing.T) {
	os.Setenv("TANZU_PROVIDER", "")
	os.Setenv("TANZU_CLUSTER_NAME", "")
	var configData bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&configData)
	yamlEncoder.SetIndent(2)

	if err := yamlEncoder.Encode(LocalClusterConfig{
		ClusterName: "test3",
		Provider:    "courteous",
		Cni:         "bongos",
		PodCidr:     "8.8.8.0/24",
		ServiceCidr: "9.9.9.0/24",
		TkrLocation: "here",
	}); err != nil {
		t.Errorf("failed setting up test data")
		return
	}

	f, err := os.CreateTemp("", "test*.yaml")
	if err != nil {
		t.Errorf("failed to create data config file. Error: %s", err.Error())
		return
	}
	_, err = f.Write(configData.Bytes())
	if err != nil {
		t.Errorf("failed to write test data config file. Error: %s", err.Error())
		return
	}

	args := map[string]string{"clusterconfigfile": f.Name()}
	config, err := InitializeConfiguration(args)
	if err != nil {
		t.Error("initialization should pass")
	}

	if config.ClusterName != "test3" {
		t.Errorf("expected ClusterName to be 'test', was actually: %q", config.ClusterName)
	}

	if config.Provider != "courteous" {
		t.Errorf("expected Provider to be set to 'courteous', was: %q", config.Provider)
	}

	if config.Cni != "bongos" {
		t.Errorf("expected Cni to be set to 'bongos', was: %q", config.Cni)
	}

	if config.TkrLocation != "here" {
		t.Errorf("expected TkrLocation to be set to 'here', was: %q", config.TkrLocation)
	}

	if config.PodCidr != "8.8.8.0/24" {
		t.Errorf("expected PodCidr to be set to '8.8.8.0/24', was: %q", config.PodCidr)
	}

	if config.ServiceCidr != "9.9.9.0/24" {
		t.Errorf("expected ServiceCidr to be set to '9.9.9.0/24', was: %q", config.ServiceCidr)
	}
}

func TestFieldNameToEnvName(t *testing.T) {
	result := fieldNameToEnvName("SomeCamelCaseVar")
	if result != "TANZU_SOME_CAMEL_CASE_VAR" {
		t.Errorf("Conversion to env var failed, got %s", result)
	}
}
