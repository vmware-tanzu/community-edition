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
	_, err := InitializeConfiguration(emptyConfig)
	if err == nil {
		t.Error("initialization should fail if no cluster name provided")
	}
}

//nolint:gocyclo
func TestInitializeConfigurationDefaults(t *testing.T) {
	args := map[string]interface{}{ClusterName: "test"}
	config, err := InitializeConfiguration(args)
	if err != nil {
		t.Error("initialization should pass")
	}

	if config.ClusterName != "test" {
		t.Errorf("expected ClusterName to be 'test', was actually: %q", config.ClusterName)
	}

	if config.KubeconfigPath != "" {
		t.Errorf("expected default KubeconfigPath value, was: %q", config.KubeconfigPath)
	}

	if config.ExistingClusterKubeconfig != "" {
		t.Errorf("expected default ExistingClusterKubeconfig value, was: %q", config.ExistingClusterKubeconfig)
	}

	if config.NodeImage != "" {
		t.Errorf("expected default NodeImage value, was: %q", config.NodeImage)
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

	if config.TkrLocation != "" {
		t.Errorf("expected default TkrLocation value, was: %q", config.TkrLocation)
	}

	if len(config.PortsToForward) != 0 {
		t.Errorf("expected default PortsToForward, was: %q", config.PortsToForward)
	}

	if config.SkipPreflightChecks != false {
		t.Errorf("expected default SkipPreflightChecks, was: %v", config.SkipPreflightChecks)
	}

	if config.ControlPlaneNodeCount != defaultConfigValues[ControlPlaneNodeCount] {
		t.Errorf("expected default ControlPlaneNodeCount, was: %q", config.ControlPlaneNodeCount)
	}

	if config.WorkerNodeCount != defaultConfigValues[WorkerNodeCount] {
		t.Errorf("expected default WorkerNodeCount, was: %q", config.WorkerNodeCount)
	}

	if len(config.InstallPackages) != 0 {
		t.Errorf("expected default install-package, was: %q", config.InstallPackages)
	}

	if config.LogFile != "" {
		t.Errorf("expected default LogFile, was: %q", config.LogFile)
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

func TestGenerateDefaultConfig(t *testing.T) {
	config := GenerateDefaultConfig()
	if config.ClusterName != "default-name" {
		t.Errorf("expected ClusterName to be 'default-name', was actually: %q", config.ClusterName)
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

func TestParsePortMapFullStringWithListenAddr(t *testing.T) {
	portMaps, err := ParsePortMappings([]string{"127.0.0.1:80:8080/tcp"})
	if err != nil {
		t.Error("Parsing should pass")
	}

	if len(portMaps) != 1 {
		t.Errorf("Expected one port mapping. Got: %v", portMaps)
	}

	if portMaps[0].ListenAddress != "127.0.0.1" {
		t.Errorf("Listen address should be 127.0.0.1, was %s", portMaps[0].ListenAddress)
	}

	if portMaps[0].ContainerPort != 80 {
		t.Errorf("Container port should be 80, was %d", portMaps[0].ContainerPort)
	}

	if portMaps[0].HostPort != 8080 {
		t.Errorf("Host port should be 8080, was %d", portMaps[0].HostPort)
	}

	if portMaps[0].Protocol != "tcp" {
		t.Errorf("Protocol should be tcp, was %s", portMaps[0].Protocol)
	}
}

func TestParsePortMapFullString(t *testing.T) {
	portMaps, err := ParsePortMappings([]string{"80:8080/tcp"})
	if err != nil {
		t.Error("Parsing should pass")
	}

	if len(portMaps) != 1 {
		t.Errorf("Expected one port mapping. Got: %v", portMaps)
	}

	if portMaps[0].ListenAddress != "" {
		t.Errorf("Listen address should be empty, was %s", portMaps[0].ListenAddress)
	}

	if portMaps[0].ContainerPort != 80 {
		t.Errorf("Container port should be 80, was %d", portMaps[0].ContainerPort)
	}

	if portMaps[0].HostPort != 8080 {
		t.Errorf("Host port should be 8080, was %d", portMaps[0].HostPort)
	}

	if portMaps[0].Protocol != "tcp" {
		t.Errorf("Protocol should be tcp, was %s", portMaps[0].Protocol)
	}
}

func TestParsePortMapContainerPort(t *testing.T) {
	portMaps, err := ParsePortMappings([]string{"80"})
	if err != nil {
		t.Error("Parsing should pass")
	}

	if len(portMaps) != 1 {
		t.Errorf("Expected one port mapping. Got: %v", portMaps)
	}

	if portMaps[0].ListenAddress != "" {
		t.Errorf("Listen address should be empty, was %s", portMaps[0].ListenAddress)
	}

	if portMaps[0].ContainerPort != 80 {
		t.Errorf("Container port should be 80, was %d", portMaps[0].ContainerPort)
	}

	if portMaps[0].HostPort != 0 {
		t.Errorf("Host port should be 0, was %d", portMaps[0].HostPort)
	}

	if portMaps[0].Protocol != "" {
		t.Errorf("Protocol should be empty, was %s", portMaps[0].Protocol)
	}
}

func TestParsePortMapContainerPortProtocol(t *testing.T) {
	portMaps, err := ParsePortMappings([]string{"80/UDP"})
	if err != nil {
		t.Error("Parsing should pass")
	}

	if len(portMaps) != 1 {
		t.Errorf("Expected one port mapping. Got: %v", portMaps)
	}

	if portMaps[0].ContainerPort != 80 {
		t.Errorf("Container port should be 80, was %d", portMaps[0].ContainerPort)
	}

	if portMaps[0].HostPort != 0 {
		t.Errorf("Host port should be 0, was %d", portMaps[0].HostPort)
	}

	if portMaps[0].Protocol != "udp" {
		t.Errorf("Protocol should be udp, was %s", portMaps[0].Protocol)
	}
}

func TestParseMultiplePortMaps(t *testing.T) {
	portMaps, err := ParsePortMappings([]string{"80/UDP", "127.0.0.1:999:999/TCP"})
	if err != nil {
		t.Error("Parsing should pass")
	}

	if len(portMaps) != 2 {
		t.Errorf("Expected two port mapping. Got: %v", portMaps)
	}

	if portMaps[0].ContainerPort != 80 {
		t.Errorf("Container port should be 80, was %d", portMaps[0].ContainerPort)
	}

	if portMaps[0].HostPort != 0 {
		t.Errorf("Host port should be 0, was %d", portMaps[0].HostPort)
	}

	if portMaps[0].Protocol != "udp" {
		t.Errorf("Protocol should be udp, was %s", portMaps[0].Protocol)
	}

	if portMaps[1].ListenAddress != "127.0.0.1" {
		t.Errorf("Listen address should be 127.0.0.1, was %s", portMaps[1].ListenAddress)
	}

	if portMaps[1].ContainerPort != 999 {
		t.Errorf("Container port should be 999, was %d", portMaps[1].ContainerPort)
	}

	if portMaps[1].HostPort != 999 {
		t.Errorf("Host port should be 999, was %d", portMaps[1].HostPort)
	}

	if portMaps[1].Protocol != "tcp" {
		t.Errorf("Protocol should be tcp, was %s", portMaps[1].Protocol)
	}
}

func TestParsePortMapInvalid(t *testing.T) {
	_, err := ParsePortMappings([]string{"http"})
	if err == nil {
		t.Error("Parsing should fail")
	}
}

// When the user provides only a package name:
// --install-package my.package.com
func TestParseInstallPackageMappingsOnlyName(t *testing.T) {
	ipMap, err := ParseInstallPackageMappings(
		[]string{"my.package.com"},
	)

	if err != nil {
		t.Error("Parsing installPackages should pass")
	}

	if len(ipMap) != 1 {
		t.Errorf("expected 1 InstallPackage. Found %v. Actual: %v", len(ipMap), ipMap)
	}

	if ipMap[0].Name != "my.package.com" {
		t.Errorf("expected InstallPackage with name. Found %s Expected: my.package.com", ipMap[0].Name)
	}

	if ipMap[0].Version != "" {
		t.Errorf("expected installPackage with no version. Found %s Expected: empty string", ipMap[0].Version)
	}

	if ipMap[0].Config != "" {
		t.Errorf("expected installPackage with no config. Found %s Expected: empty string", ipMap[0].Config)
	}
}

// When the user provides multiple installPackage flags:
// --install-package my.installPackage.package.com
// --install-package my.other-installPackage.package.com
//
// dequeues values from flags in order they are enqueued despite order of flags
func TestParseInstallPackageMappingsManyNames(t *testing.T) {
	ipMap, err := ParseInstallPackageMappings(
		[]string{"my.installPackage.package.com", "my.other-installPackage.package.com"},
	)

	if err != nil {
		t.Error("Parsing installPackages should pass")
	}

	if len(ipMap) != 2 {
		t.Errorf("expected 2 installPackage. Found %v. Actual: %v", len(ipMap), ipMap)
	}

	if ipMap[0].Name != "my.installPackage.package.com" {
		t.Errorf("expected installPackage with name. Found %s Expected: my.configured.package.com", ipMap[0].Name)
	}

	if ipMap[0].Version != "" {
		t.Errorf("expected installPackage without version. Found %s Expected: empty string", ipMap[0].Version)
	}

	if ipMap[0].Config != "" {
		t.Errorf("expected installPackage with no config. Found %s Expected: empty string", ipMap[0].Config)
	}

	if ipMap[1].Name != "my.other-installPackage.package.com" {
		t.Errorf("expected installPackage with name. Found %s Expected: my.other-installPackage.package.com", ipMap[1].Name)
	}

	if ipMap[1].Version != "" {
		t.Errorf("expected installPackage with no version. Found %s Expected: empty string", ipMap[1].Version)
	}

	if ipMap[1].Config != "" {
		t.Errorf("expected installPackage with no config. Found %s Expected: empty string", ipMap[1].Config)
	}
}

// When the user provides only a installPackage name:
// --install-package my.package.com:1.2.3
func TestParseInstallPackageMappingsNameVersion(t *testing.T) {
	ipMap, err := ParseInstallPackageMappings(
		[]string{"my.package.com:1.2.3"},
	)

	if err != nil {
		t.Error("Parsing installPackages should pass")
	}

	if len(ipMap) != 1 {
		t.Errorf("expected 1 installPackage. Found %v. Actual: %v", len(ipMap), ipMap)
	}

	if ipMap[0].Name != "my.package.com" {
		t.Errorf("expected installPackage with name. Found %s Expected: my.package.com", ipMap[0].Name)
	}

	if ipMap[0].Version != "1.2.3" {
		t.Errorf("expected installPackage with version. Found %s Expected: 1.2.3", ipMap[0].Version)
	}

	if ipMap[0].Config != "" {
		t.Errorf("expected installPackage with no config. Found %s Expected: empty string", ipMap[0].Config)
	}
}

// When the user provides multiple installPackage flags with name and version:
// --install-package my.installPackage.package.com:1.2.3
// --install-package my.other-installPackage.package.com:7.8.9
func TestParseInstallPackageMappingsManyNameVersion(t *testing.T) {
	ipMap, err := ParseInstallPackageMappings(
		[]string{"my.installPackage.package.com:1.2.3", "my.other-installPackage.package.com:7.8.9"},
	)

	if err != nil {
		t.Error("Parsing installPackages should pass")
	}

	if len(ipMap) != 2 {
		t.Errorf("expected 2 installPackage. Found %v. Actual: %v", len(ipMap), ipMap)
	}

	if ipMap[0].Name != "my.installPackage.package.com" {
		t.Errorf("expected installPackage with name. Found %s Expected: my.configured.package.com", ipMap[0].Name)
	}

	if ipMap[0].Version != "1.2.3" {
		t.Errorf("expected installPackage with version. Found %s Expected: 1.2.3", ipMap[0].Version)
	}

	if ipMap[0].Config != "" {
		t.Errorf("expected installPackage with no config. Found %s Expected: empty string", ipMap[0].Config)
	}

	if ipMap[1].Name != "my.other-installPackage.package.com" {
		t.Errorf("expected installPackage with name. Found %s Expected: my.other-installPackage.package.com", ipMap[1].Name)
	}

	if ipMap[1].Version != "7.8.9" {
		t.Errorf("expected installPackage with version. Found %s Expected: 7.8.9", ipMap[1].Version)
	}

	if ipMap[1].Config != "" {
		t.Errorf("expected installPackage with no config. Found %s Expected: empty string", ipMap[1].Config)
	}
}

// When the user provides a full installPackage mapping:
// --install-package my.package.com:4.4.4:woof-path
func TestParseInstallPackageMappingSingle(t *testing.T) {
	ipMap, err := ParseInstallPackageMappings(
		[]string{"my.package.com:4.4.4:woof-path"},
	)

	if err != nil {
		t.Error("Parsing installPackages should pass")
	}

	if len(ipMap) != 1 {
		t.Errorf("expected 1 installPackage. Found %v. Actual: %v", len(ipMap), ipMap)
	}

	if ipMap[0].Name != "my.package.com" {
		t.Errorf("expected installPackage with name. Found %s Expected: my.package.com", ipMap[0].Name)
	}

	if ipMap[0].Version != "4.4.4" {
		t.Errorf("expected installPackage with version. Found %s Expected: 4.4.4", ipMap[0].Version)
	}

	if ipMap[0].Config != "woof-path" {
		t.Errorf("expected installPackage with config. Found %s Expected: woof-path", ipMap[0].Config)
	}
}

// When the user provides only a installPackage name, empty version, and a config path:
// --install-package my.package.com::my-config
func TestParseInstallPackageMappingsEmptyVersion(t *testing.T) {
	ipMap, err := ParseInstallPackageMappings(
		[]string{"my.package.com::my-config"},
	)

	if err != nil {
		t.Error("Parsing installPackages should pass")
	}

	if len(ipMap) != 1 {
		t.Errorf("expected 1 installPackage. Found %v. Actual: %v", len(ipMap), ipMap)
	}

	if ipMap[0].Name != "my.package.com" {
		t.Errorf("expected installPackage with name. Found %s Expected: my.package.com", ipMap[0].Name)
	}

	if ipMap[0].Version != "" {
		t.Errorf("expected installPackage with no version. Found %s Expected: empty string", ipMap[0].Version)
	}

	if ipMap[0].Config != "my-config" {
		t.Errorf("expected installPackage with config. Found %s Expected: my-config", ipMap[0].Config)
	}
}

// When the user provides multiple full installPackage mappings in single install-package flag:
// --install-package my.package.com:4.4.4:woof-path,other.package.com:1.2.3:my-config-path
func TestParseInstallPackageMappingsMultiple(t *testing.T) {
	ipMap, err := ParseInstallPackageMappings(
		[]string{"my.package.com:4.4.4:woof-path,other.package.com:1.2.3:my-config-path"},
	)

	if err != nil {
		t.Error("Parsing installPackages should pass")
	}

	if len(ipMap) != 2 {
		t.Errorf("expected 1 installPackage. Found %v. Actual: %v", len(ipMap), ipMap)
	}

	if ipMap[0].Name != "my.package.com" {
		t.Errorf("expected installPackage with name. Found %s Expected: my.package.com", ipMap[0].Name)
	}

	if ipMap[0].Version != "4.4.4" {
		t.Errorf("expected installPackage with version. Found %s Expected: 4.4.4", ipMap[0].Version)
	}

	if ipMap[0].Config != "woof-path" {
		t.Errorf("expected installPackage with config. Found %s Expected: woof-path", ipMap[0].Config)
	}

	if ipMap[1].Name != "other.package.com" {
		t.Errorf("expected installPackage with name. Found %s Expected: my.package.com", ipMap[0].Name)
	}

	if ipMap[1].Version != "1.2.3" {
		t.Errorf("expected installPackage with version. Found %s Expected: 1.2.3", ipMap[0].Version)
	}

	if ipMap[1].Config != "my-config-path" {
		t.Errorf("expected installPackage with config. Found %s Expected: my-config-path", ipMap[0].Config)
	}
}

// When the user provides multiple full installPackage mappings in multiple install-package flags:
// --install-package my.package.com:4.4.4:woof-path,other.package.com:2.2.2:other-config
// --install-package my.other-package.com:1.2.3:my-path,my.final.package.com:7.8.9:final-config
func TestParseInstallPackageMappingMixedFlags(t *testing.T) {
	ipMap, err := ParseInstallPackageMappings(
		[]string{
			"my.package.com:4.4.4:woof-path,other.package.com:2.2.2:other-config",
			"my-third.package.com:1.2.3:my-path,my-final.package.com:7.8.9:final-config",
		},
	)

	if err != nil {
		t.Error("Parsing installPackages should pass")
	}

	if len(ipMap) != 4 {
		t.Errorf("expected 4 installPackages. Found %v. Actual: %v", len(ipMap), ipMap)
	}

	if ipMap[0].Name != "my.package.com" {
		t.Errorf("expected installPackage with name. Found %s Expected: my.package.com", ipMap[0].Name)
	}

	if ipMap[0].Version != "4.4.4" {
		t.Errorf("expected installPackage with version. Found %s Expected: 4.4.4", ipMap[0].Version)
	}

	if ipMap[0].Config != "woof-path" {
		t.Errorf("expected installPackage with config. Found %s Expected: woof-path", ipMap[0].Config)
	}

	if ipMap[1].Name != "other.package.com" {
		t.Errorf("expected installPackage with name. Found %s Expected: other.package.com", ipMap[1].Name)
	}

	if ipMap[1].Version != "2.2.2" {
		t.Errorf("expected installPackage with version. Found %s Expected: 2.2.2", ipMap[1].Version)
	}

	if ipMap[1].Config != "other-config" {
		t.Errorf("expected installPackage with config. Found %s Expected: other-config", ipMap[1].Config)
	}

	if ipMap[2].Name != "my-third.package.com" {
		t.Errorf("expected installPackage with name. Found %s Expected: third.package.com", ipMap[2].Name)
	}

	if ipMap[2].Version != "1.2.3" {
		t.Errorf("expected installPackage with version. Found %s Expected: 1.2.3", ipMap[2].Version)
	}

	if ipMap[2].Config != "my-path" {
		t.Errorf("expected installPackage with config. Found %s Expected: my-config-path", ipMap[2].Config)
	}

	if ipMap[3].Name != "my-final.package.com" {
		t.Errorf("expected installPackage with name. Found %s Expected: my-final.package.com", ipMap[3].Name)
	}

	if ipMap[3].Version != "7.8.9" {
		t.Errorf("expected installPackage with version. Found %s Expected: 7.8.9", ipMap[3].Version)
	}

	if ipMap[3].Config != "final-config" {
		t.Errorf("expected installPackage with config. Found %s Expected: final-config", ipMap[3].Config)
	}
}

// When the user provides a bad formatting:
// --install-package my.package.com:4.4.4:woof-path:garbage
func TestParseInstallPackageMappingBadFormat(t *testing.T) {
	_, err := ParseInstallPackageMappings(
		[]string{"my.package.com:4.4.4:woof-path:garbage"},
	)

	if err == nil {
		t.Error("Parsing should fail")
	}
}

// Won't errors when nothing is provided:
func TestParseInstallPackageNil(t *testing.T) {
	_, err := ParseInstallPackageMappings(
		[]string{},
	)

	if err != nil {
		t.Error("Parsing shouldn't fail")
	}
}
