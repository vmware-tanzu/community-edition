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

var emptyConfig = map[string]interface{}{
	ClusterConfigFile:      "",
	ClusterName:            "",
	Tty:                    "",
	TKRLocation:            "",
	Provider:               "",
	Cni:                    "",
	PodCIDR:                "",
	ServiceCIDR:            "",
	ControlPlaneNodeCount:  "",
	WorkerNodeCount:        "",
	AdditionalPackageRepos: []string{},
}

func TestInitializeConfigurationNoName(t *testing.T) {
	_, err := InitializeConfiguration(emptyConfig)
	if err == nil {
		t.Error("initialization should fail if no cluster name provided")
	}
}

func TestInitializeConfigurationDefaults(t *testing.T) {
	args := map[string]interface{}{ClusterName: "test"}
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

	if len(config.AdditionalPackageRepos) != 0 {
		t.Errorf("expected no AdditionalPackageRepos, was: %q", config.AdditionalPackageRepos)
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

	if config.ControlPlaneNodeCount != defaultConfigValues[ControlPlaneNodeCount] {
		t.Errorf("expected default ControlPlaneNodeCount, was: %q", config.ControlPlaneNodeCount)
	}

	if config.WorkerNodeCount != defaultConfigValues[WorkerNodeCount] {
		t.Errorf("expected default WorkerNodeCount, was: %q", config.ControlPlaneNodeCount)
	}
}

func TestInitializeConfigurationEnvVariables(t *testing.T) {
	os.Setenv("TANZU_PROVIDER", "test_provider")
	os.Setenv("TANZU_CLUSTER_NAME", "test2")
	config, err := InitializeConfiguration(emptyConfig)
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

	if len(config.AdditionalPackageRepos) != 0 {
		t.Errorf("expected no AdditionalPackageRepos, was: %q", config.AdditionalPackageRepos)
	}

	if config.PodCidr != defaultConfigValues[PodCIDR] {
		t.Errorf("expected default PodCidr, was: %q", config.PodCidr)
	}

	if config.ServiceCidr != defaultConfigValues[ServiceCIDR] {
		t.Errorf("expected default ServiceCidr, was: %q", config.ServiceCidr)
	}

	if config.ControlPlaneNodeCount != defaultConfigValues[ControlPlaneNodeCount] {
		t.Errorf("expected default ControlPlaneNodeCount value, was: %q", config.ControlPlaneNodeCount)
	}

	if config.WorkerNodeCount != defaultConfigValues[WorkerNodeCount] {
		t.Errorf("expected default WorkerNodeCount value, was: %q", config.WorkerNodeCount)
	}
}

func TestInitializeConfigurationArgsTakePrecedent(t *testing.T) {
	os.Setenv("TANZU_PROVIDER", "test_provider")
	os.Setenv("TANZU_CLUSTER_NAME", "test2")
	args := map[string]interface{}{ClusterName: "test"}
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

	if len(config.AdditionalPackageRepos) != 0 {
		t.Errorf("expected no AdditionalPackageRepos, was: %q", len(config.AdditionalPackageRepos))
	}

	if config.PodCidr != defaultConfigValues[PodCIDR] {
		t.Errorf("expected default PodCidr, was: %q", config.PodCidr)
	}

	if config.ServiceCidr != defaultConfigValues[ServiceCIDR] {
		t.Errorf("expected default ServiceCidr, was: %q", config.ServiceCidr)
	}

	if config.ControlPlaneNodeCount != defaultConfigValues[ControlPlaneNodeCount] {
		t.Errorf("expected default ControlPlaneNodeCount value, was: %q", config.ControlPlaneNodeCount)
	}

	if config.WorkerNodeCount != defaultConfigValues[WorkerNodeCount] {
		t.Errorf("expected default WorkerNodeCount value, was: %q", config.WorkerNodeCount)
	}
}

func TestInitializeConfigurationFromConfigFile(t *testing.T) {
	os.Setenv("TANZU_PROVIDER", "")
	os.Setenv("TANZU_CLUSTER_NAME", "")
	var configData bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&configData)
	yamlEncoder.SetIndent(2)

	if err := yamlEncoder.Encode(UnmanagedClusterConfig{
		ClusterName:            "test3",
		Provider:               "courteous",
		Cni:                    "bongos",
		PodCidr:                "8.8.8.0/24",
		ServiceCidr:            "9.9.9.0/24",
		TkrLocation:            "here",
		AdditionalPackageRepos: []string{"example.registry.com", "another.example.com"},
		ControlPlaneNodeCount:  "99",
		WorkerNodeCount:        "25",
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

	args := map[string]interface{}{ClusterConfigFile: f.Name()}
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

	if config.AdditionalPackageRepos[0] != "example.registry.com" {
		t.Errorf("expected first AdditionalPackageRepos value to be 'example.registry.com', was: %q", config.AdditionalPackageRepos[0])
	}

	if config.AdditionalPackageRepos[1] != "another.example.com" {
		t.Errorf("expected first AdditionalPackageRepos value to be 'another.example.com', was: %q", config.AdditionalPackageRepos[1])
	}

	if config.PodCidr != "8.8.8.0/24" {
		t.Errorf("expected PodCidr to be set to '8.8.8.0/24', was: %q", config.PodCidr)
	}

	if config.ServiceCidr != "9.9.9.0/24" {
		t.Errorf("expected ServiceCidr to be set to '9.9.9.0/24', was: %q", config.ServiceCidr)
	}

	if config.ControlPlaneNodeCount != "99" {
		t.Errorf("expected ControlPlaneNodeCount to be set to '99', was: %q", config.ControlPlaneNodeCount)
	}

	if config.WorkerNodeCount != "25" {
		t.Errorf("expected WorkerNodeCount to be set to '25', was: %q", config.WorkerNodeCount)
	}
}

func TestFieldNameToEnvName(t *testing.T) {
	result := fieldNameToEnvName("SomeCamelCaseVar")
	if result != "TANZU_SOME_CAMEL_CASE_VAR" {
		t.Errorf("Conversion to env var failed, got %s", result)
	}
}

func TestSanatizeKubeconfigPath(t *testing.T) {
	result := sanatizeKubeconfigPath("/path/to/file/kubeconfig.yaml")
	if result != "/path/to/file/kubeconfig.yaml" {
		t.Errorf("Sanatizing kubeconfig path failed, got %s", result)
	}

	result = sanatizeKubeconfigPath("~/path/with/tilda/kubeconfig.yaml")
	home, _ := os.UserHomeDir()
	if result != home+"/path/with/tilda/kubeconfig.yaml" {
		t.Errorf("Sanatizing kubeconfig path failed, got %s", result)
	}
}

func TestParsePortMapFullString(t *testing.T) {
	portMap, err := ParsePortMap("80:8080/tcp")
	if err != nil {
		t.Error("Parsing should pass")
	}

	if portMap.ContainerPort != 80 {
		t.Errorf("Container port should be 80, was %d", portMap.ContainerPort)
	}

	if portMap.HostPort != 8080 {
		t.Errorf("Host port should be 8080, was %d", portMap.HostPort)
	}

	if portMap.Protocol != "tcp" {
		t.Errorf("Protocol should be tcp, was %s", portMap.Protocol)
	}
}

func TestParsePortMapContainerPort(t *testing.T) {
	portMap, err := ParsePortMap("80")
	if err != nil {
		t.Error("Parsing should pass")
	}

	if portMap.ContainerPort != 80 {
		t.Errorf("Container port should be 80, was %d", portMap.ContainerPort)
	}

	if portMap.HostPort != 0 {
		t.Errorf("Host port should be 0, was %d", portMap.HostPort)
	}

	if portMap.Protocol != "" {
		t.Errorf("Protocol should be empty, was %s", portMap.Protocol)
	}
}

func TestParsePortMapContainerPortProtocol(t *testing.T) {
	portMap, err := ParsePortMap("80/UDP")
	if err != nil {
		t.Error("Parsing should pass")
	}

	if portMap.ContainerPort != 80 {
		t.Errorf("Container port should be 80, was %d", portMap.ContainerPort)
	}

	if portMap.HostPort != 0 {
		t.Errorf("Host port should be 0, was %d", portMap.HostPort)
	}

	if portMap.Protocol != "udp" {
		t.Errorf("Protocol should be udp, was %s", portMap.Protocol)
	}
}

func TestParsePortMapInvalid(t *testing.T) {
	_, err := ParsePortMap("http")
	if err == nil {
		t.Error("Parsing should fail")
	}
}
