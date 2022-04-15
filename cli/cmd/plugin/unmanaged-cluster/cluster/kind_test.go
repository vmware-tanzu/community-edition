// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//nolint:goconst
package cluster

import (
	"encoding/json"
	"testing"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/config"

	kindconfig "sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
)

var normalDockerInfoJSON = `{"ID":"SEB7:L67H:GZMX:VPIN:YZ7V:RTRC:DCML:3C7C:PNN3:2DQA:6GD2:ZIWU","Containers":7,"ContainersRunning":1,"ContainersPaused":0,"ContainersStopped":6,"Images":151,"Driver":"overlay2","DriverStatus":[["Backing Filesystem","extfs"],["Supports d_type","true"],["Native Overlay Diff","true"],["userxattr","false"]],"Plugins":{"Volume":["local"],"Network":["bridge","host","ipvlan","macvlan","null","overlay"],"Authorization":null,"Log":["awslogs","fluentd","gcplogs","gelf","journald","json-file","local","logentries","splunk","syslog"]},"MemoryLimit":true,"SwapLimit":true,"KernelMemory":true,"KernelMemoryTCP":true,"CpuCfsPeriod":true,"CpuCfsQuota":true,"CPUShares":true,"CPUSet":true,"PidsLimit":true,"IPv4Forwarding":true,"BridgeNfIptables":true,"BridgeNfIp6tables":true,"Debug":false,"NFd":32,"OomKillDisable":true,"NGoroutines":40,"SystemTime":"2022-01-11T15:43:55.314860422-06:00","LoggingDriver":"json-file","CgroupDriver":"cgroupfs","CgroupVersion":"1","NEventsListener":0,"KernelVersion":"5.11.0-43-generic","OperatingSystem":"Ubuntu 20.04.3 LTS","OSVersion":"20.04","OSType":"linux","Architecture":"x86_64","IndexServerAddress":"https://index.docker.io/v1/","RegistryConfig":{"AllowNondistributableArtifactsCIDRs":[],"AllowNondistributableArtifactsHostnames":[],"InsecureRegistryCIDRs":["127.0.0.0/8"],"IndexConfigs":{"docker.io":{"Name":"docker.io","Mirrors":[],"Secure":true,"Official":true}},"Mirrors":[]},"NCPU":16,"MemTotal":33613119488,"GenericResources":null,"DockerRootDir":"/var/lib/docker","HttpProxy":"","HttpsProxy":"","NoProxy":"","Name":"sm-workstation","Labels":[],"ExperimentalBuild":false,"ServerVersion":"20.10.12","Runtimes":{"io.containerd.runc.v2":{"path":"runc"},"io.containerd.runtime.v1.linux":{"path":"runc"},"runc":{"path":"runc"}},"DefaultRuntime":"runc","Swarm":{"NodeID":"","NodeAddr":"","LocalNodeState":"inactive","ControlAvailable":false,"Error":"","RemoteManagers":null},"LiveRestoreEnabled":false,"Isolation":"","InitBinary":"docker-init","ContainerdCommit":{"ID":"7b11cfaabd73bb80907dd23182b9347b4245eb5d","Expected":"7b11cfaabd73bb80907dd23182b9347b4245eb5d"},"RuncCommit":{"ID":"v1.0.2-0-g52b36a2","Expected":"v1.0.2-0-g52b36a2"},"InitCommit":{"ID":"de40ad0","Expected":"de40ad0"},"SecurityOptions":["name=apparmor","name=seccomp,profile=default"],"Warnings":null,"ClientInfo":{"Debug":false,"Context":"default","Plugins":[{"SchemaVersion":"0.1.0","Vendor":"Docker Inc.","Version":"v0.9.1-beta3","ShortDescription":"Docker App","Experimental":true,"Name":"app","Path":"/usr/libexec/docker/cli-plugins/docker-app"},{"SchemaVersion":"0.1.0","Vendor":"Docker Inc.","Version":"v0.7.1-docker","ShortDescription":"Docker Buildx","Name":"buildx","Path":"/usr/libexec/docker/cli-plugins/docker-buildx"},{"SchemaVersion":"0.1.0","Vendor":"Docker Inc.","Version":"v0.12.0","ShortDescription":"Docker Scan","Name":"scan","Path":"/usr/libexec/docker/cli-plugins/docker-scan"}],"Warnings":null}}`

func TestValidateDockerInfoNoIssues(t *testing.T) {
	// Test the full real response from docker info, we'll use a truncated response below
	warnings, errs := validateDockerInfo([]byte(normalDockerInfoJSON))
	if len(warnings) > 0 {
		t.Errorf("no warnings should be detected but %d returned", len(warnings))
	}
	if len(errs) > 0 {
		t.Errorf("no errors should be detected but %d returned", len(errs))
	}
}

func TestValidateDockerInfoNotEnoughMemory(t *testing.T) {
	testInfo := dockerInfo{
		CPUs:         1,
		Memory:       1073741824,
		Architecture: "x86_64",
	}
	output, _ := json.Marshal(testInfo)
	warnings, errs := validateDockerInfo(output)
	if len(warnings) > 0 {
		t.Errorf("no warnings should be detected but %d returned", len(warnings))
	}
	if len(errs) != 1 {
		t.Errorf("expected 1 error but %d returned", len(errs))
	}
}

func TestValidateDockerInfoNotEnoughCPU(t *testing.T) {
	testInfo := dockerInfo{
		CPUs:         0,
		Memory:       2147483648,
		Architecture: "x86_64",
	}
	output, _ := json.Marshal(testInfo)
	warnings, errs := validateDockerInfo(output)
	if len(warnings) > 0 {
		t.Errorf("no warnings should be detected but %d returned", len(warnings))
	}
	if len(errs) != 1 {
		t.Errorf("expected 1 error but %d returned", len(errs))
	}
}

func TestValidateDockerInfoArchitectureNotSupported(t *testing.T) {
	testInfo := dockerInfo{
		CPUs:         16,
		Memory:       2147483648,
		Architecture: "arm",
	}
	output, _ := json.Marshal(testInfo)
	warnings, errs := validateDockerInfo(output)
	if len(warnings) > 0 {
		t.Errorf("no warnings should be detected but %d returned", len(warnings))
	}
	if len(errs) != 1 {
		t.Errorf("expected 1 error but %d returned", len(errs))
	}
}

func TestValidateDockerInfoArchitectureARM64(t *testing.T) {
	testInfo := dockerInfo{
		CPUs:         16,
		Memory:       2147483648,
		Architecture: "aarch64",
	}
	output, _ := json.Marshal(testInfo)
	warnings, errs := validateDockerInfo(output)
	if len(warnings) != 1 {
		t.Errorf("warnings should be detected but %d returned", len(warnings))
	}
	if len(errs) != 0 {
		t.Errorf("no errors expected but %d returned", len(errs))
	}
}

func TestValidateDockerInfoBadData(t *testing.T) {
	warnings, errs := validateDockerInfo([]byte{240, 159, 146, 169})
	if len(warnings) > 0 {
		t.Errorf("no warnings should be detected but %d returned", len(warnings))
	}
	if len(errs) != 1 {
		t.Errorf("expected 1 error but %d returned", len(errs))
	}
}

func TestDefaultKindConfigCluster(t *testing.T) {
	defaultKindconfigCluster, err := kindConfigFromClusterConfig(config.GenerateDefaultConfig())

	if err != nil {
		t.Errorf("expected getting default kindconfig.Cluster to pass. Error: %s\n", err.Error())
	}

	if defaultKindconfigCluster.Kind != KindTypedataCluster {
		t.Errorf("expected expected default kindconfig.Kind to be `Cluster`. Actual: %s", defaultKindconfigCluster.Kind)
	}

	if defaultKindconfigCluster.APIVersion != KindTypedataAPIVersion {
		t.Errorf("expected expected default kindconfig.APIVersion to be `kind.x-k8s.io/v1alpha4`. Actual: %s", defaultKindconfigCluster.APIVersion)
	}
}

// TypeKind in global config takes precedence when not default
func TestMergeConfigType(t *testing.T) {
	globalConfig := kindconfig.Cluster{
		TypeMeta: kindconfig.TypeMeta{
			Kind:       "global-kind",
			APIVersion: "1.2.3",
		},
	}

	providerConfig := kindconfig.Cluster{
		TypeMeta: kindconfig.TypeMeta{
			Kind:       "provider-kind",
			APIVersion: "7.8.9",
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if globalConfig.Kind != "global-kind" {
		t.Errorf("expected merged kindconfig.Kind to have kept value. Actual: %s", globalConfig.Kind)
	}

	if globalConfig.APIVersion != "1.2.3" {
		t.Errorf("expected merged kindconfig.APIVersion to have kept value. Actual: %s", globalConfig.APIVersion)
	}
}

// Custom config for TypeKind in provider config takes precedence
func TestMergeConfigTypeWhenDefault(t *testing.T) {
	defaultKindconfigCluster, _ := kindConfigFromClusterConfig(config.GenerateDefaultConfig())

	globalConfig := kindconfig.Cluster{
		TypeMeta: kindconfig.TypeMeta{
			Kind:       defaultKindconfigCluster.Kind,
			APIVersion: defaultKindconfigCluster.APIVersion,
		},
	}

	providerConfig := kindconfig.Cluster{
		TypeMeta: kindconfig.TypeMeta{
			Kind:       "custom-kind",
			APIVersion: "custom-version",
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if globalConfig.Kind != "custom-kind" {
		t.Errorf("expected merged kindconfig.Kind to have merged from provider config. Actual: %s", globalConfig.Kind)
	}

	if globalConfig.APIVersion != "custom-version" {
		t.Errorf("expected merged kindconfig.APIVersion to have merged from provider config. Actual: %s", globalConfig.APIVersion)
	}
}

// Cluster name in global config takes precedence
func TestMergeConfigName(t *testing.T) {
	globalConfig := kindconfig.Cluster{
		Name: "name-via-arg",
	}

	providerConfig := kindconfig.Cluster{
		Name: "name-via-config-file",
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if globalConfig.Name != "name-via-arg" {
		t.Errorf("expected merged kindconfig.Name to have kept global config value. Actual: %s", globalConfig.Name)
	}
}

// Cluster name in provider config is used when default name detected
func TestMergeConfigNameDefault(t *testing.T) {
	defaultKindconfigCluster, _ := kindConfigFromClusterConfig(config.GenerateDefaultConfig())

	globalConfig := kindconfig.Cluster{
		Name: defaultKindconfigCluster.Name,
	}

	providerConfig := kindconfig.Cluster{
		Name: "name-via-config-file",
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if globalConfig.Name != "name-via-config-file" {
		t.Errorf("expected merged kindconfig.Name to have merged from provider config. Actual: %s", globalConfig.Name)
	}
}

// Takes the global node config and image
// when none was given in provider config
func TestMergeConfigNodesUsesGlobal(t *testing.T) {
	globalConfig := kindconfig.Cluster{
		Nodes: []kindconfig.Node{
			{
				Role:  "control-plane",
				Image: "global.image.url",
			},
			{
				Role:  "worker",
				Image: "global.image.url",
			},
		},
	}

	providerConfig := kindconfig.Cluster{
		Nodes: []kindconfig.Node{
			{
				Role: "control-plane",
			},
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if len(globalConfig.Nodes) != 2 {
		t.Errorf("expected two nodes from global config. Actual: %v", globalConfig.Nodes)
	}

	if globalConfig.Nodes[0].Role != "control-plane" {
		t.Errorf("expected first node Role to have kept global config value. Actual: %s", globalConfig.Nodes[0].Role)
	}

	if globalConfig.Nodes[1].Role != "worker" {
		t.Errorf("expected second node Role to have kept global config value. Actual: %s", globalConfig.Nodes[1].Role)
	}

	if globalConfig.Nodes[0].Image != "global.image.url" {
		t.Errorf("expected first node Image to have kept global config value. Actual: %s", globalConfig.Nodes[0].Image)
	}

	if globalConfig.Nodes[1].Image != "global.image.url" {
		t.Errorf("expected second node Image to have kept global config value. Actual: %s", globalConfig.Nodes[1].Image)
	}
}

// Keeps the default node image
// when user gives number of nodes via cli flag
// and provider config doesn't have any defined via provider config
func TestMergeConfigNodesUsesDefaultImage(t *testing.T) {
	defaultKindconfigCluster, _ := kindConfigFromClusterConfig(config.GenerateDefaultConfig())

	globalConfig := kindconfig.Cluster{
		Nodes: []kindconfig.Node{
			{
				Role:  "control-plane",
				Image: defaultKindconfigCluster.Nodes[0].Image,
			},
			{
				Role:  "worker",
				Image: defaultKindconfigCluster.Nodes[0].Image,
			},
		},
	}

	providerConfig := kindconfig.Cluster{
		Nodes: []kindconfig.Node{
			{
				Role: "control-plane",
			},
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if len(globalConfig.Nodes) != 2 {
		t.Errorf("expected two nodes from global config. Actual: %v", globalConfig.Nodes)
	}

	if globalConfig.Nodes[0].Role != "control-plane" {
		t.Errorf("expected first node Role to have kept global config value. Actual: %s", globalConfig.Nodes[0].Role)
	}

	if globalConfig.Nodes[1].Role != "worker" {
		t.Errorf("expected second node Role to have kept global config value. Actual: %s", globalConfig.Nodes[1].Role)
	}

	if globalConfig.Nodes[0].Image != defaultKindconfigCluster.Nodes[0].Image {
		t.Errorf("expected first node Image to have kept default config value. Actual: %s", globalConfig.Nodes[0].Image)
	}

	if globalConfig.Nodes[1].Image != defaultKindconfigCluster.Nodes[0].Image {
		t.Errorf("expected second node Image to have kept default config value. Actual: %s", globalConfig.Nodes[1].Image)
	}
}

// Keeps the default node image for second node
// when user gives number of nodes via cli flag
// and provider config has fewer nodes but defines a node image via provider config
func TestMergeConfigNodesUsesProviderImageDiffNumbers(t *testing.T) {
	defaultKindconfigCluster, _ := kindConfigFromClusterConfig(config.GenerateDefaultConfig())
	providerNodeImage := "provider.node.image.url"

	globalConfig := kindconfig.Cluster{
		Nodes: []kindconfig.Node{
			{
				Role:  "control-plane",
				Image: defaultKindconfigCluster.Nodes[0].Image,
			},
			{
				Role:  "worker",
				Image: defaultKindconfigCluster.Nodes[0].Image,
			},
		},
	}

	providerConfig := kindconfig.Cluster{
		Nodes: []kindconfig.Node{
			{
				Role:  "control-plane",
				Image: providerNodeImage,
			},
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if len(globalConfig.Nodes) != 2 {
		t.Errorf("expected two nodes from global config. Actual: %v", globalConfig.Nodes)
	}

	if globalConfig.Nodes[0].Role != "control-plane" {
		t.Errorf("expected first node Role to have kept global config value. Actual: %s", globalConfig.Nodes[0].Role)
	}

	if globalConfig.Nodes[1].Role != "worker" {
		t.Errorf("expected second node Role to have kept global config value. Actual: %s", globalConfig.Nodes[1].Role)
	}

	if globalConfig.Nodes[0].Image != providerNodeImage {
		t.Errorf("expected first node Image to used non-default value from provider config. Actual: %s", globalConfig.Nodes[0].Image)
	}

	if globalConfig.Nodes[1].Image != defaultKindconfigCluster.Nodes[0].Image {
		t.Errorf("expected second node Image to have kept default config value. Actual: %s", globalConfig.Nodes[1].Image)
	}
}

// Uses default node image for provider config nodes
// when users doesn't define number of nodes via CLI flags
// and provider config has nodes without a custom image defined
func TestMergeConfigNodesUsesProviderNodes(t *testing.T) {
	defaultKindconfigCluster, _ := kindConfigFromClusterConfig(config.GenerateDefaultConfig())

	globalConfig := kindconfig.Cluster{
		Nodes: []kindconfig.Node{
			{
				Role:  "control-plane",
				Image: defaultKindconfigCluster.Nodes[0].Image,
			},
		},
	}

	providerConfig := kindconfig.Cluster{
		Nodes: []kindconfig.Node{
			{
				Role: "control-plane",
			},
			{
				Role: "control-plane",
			},
			{
				Role: "worker",
			},
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if len(globalConfig.Nodes) != 3 {
		t.Errorf("expected three provider nodes from global config. Actual: %v", globalConfig.Nodes)
	}

	if globalConfig.Nodes[0].Role != "control-plane" {
		t.Errorf("expected first node Role to have used provider config value. Actual: %s", globalConfig.Nodes[0].Role)
	}

	if globalConfig.Nodes[1].Role != "control-plane" {
		t.Errorf("expected second node Role to have used provider config value. Actual: %s", globalConfig.Nodes[1].Role)
	}

	if globalConfig.Nodes[2].Role != "worker" {
		t.Errorf("expected second node Role to have used provider config value. Actual: %s", globalConfig.Nodes[1].Role)
	}

	if globalConfig.Nodes[0].Image != defaultKindconfigCluster.Nodes[0].Image {
		t.Errorf("expected first node Image to have used default config value. Actual: %s", globalConfig.Nodes[0].Image)
	}

	if globalConfig.Nodes[1].Image != defaultKindconfigCluster.Nodes[0].Image {
		t.Errorf("expected second node Image to have used default config value. Actual: %s", globalConfig.Nodes[1].Image)
	}

	if globalConfig.Nodes[2].Image != defaultKindconfigCluster.Nodes[0].Image {
		t.Errorf("expected second node Image to have used default config value. Actual: %s", globalConfig.Nodes[1].Image)
	}
}

// Uses node image from provider config
// when user gives cli flag for number of nodes
// and provider config has matching number of nodes with custom image defined
func TestMergeConfigNodesUsesProviderImage(t *testing.T) {
	defaultKindconfigCluster, _ := kindConfigFromClusterConfig(config.GenerateDefaultConfig())
	providerNodeImage := "provider.node.image.url"

	globalConfig := kindconfig.Cluster{
		Nodes: []kindconfig.Node{
			{
				Role:  "control-plane",
				Image: defaultKindconfigCluster.Nodes[0].Image,
			},
			{
				Role:  "worker",
				Image: defaultKindconfigCluster.Nodes[0].Image,
			},
		},
	}

	providerConfig := kindconfig.Cluster{
		Nodes: []kindconfig.Node{
			{
				Role:  "control-plane",
				Image: providerNodeImage,
			},
			{
				Role:  "worker",
				Image: providerNodeImage,
			},
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if len(globalConfig.Nodes) != 2 {
		t.Errorf("expected two nodes from global config. Actual: %v", globalConfig.Nodes)
	}

	if globalConfig.Nodes[0].Role != "control-plane" {
		t.Errorf("expected first node Role to have kept global config value. Actual: %s", globalConfig.Nodes[0].Role)
	}

	if globalConfig.Nodes[1].Role != "worker" {
		t.Errorf("expected second node Role to have kept global config value. Actual: %s", globalConfig.Nodes[1].Role)
	}

	if globalConfig.Nodes[0].Image != providerNodeImage {
		t.Errorf("expected first node Image to have used provider config value. Actual: %s", globalConfig.Nodes[0].Image)
	}

	if globalConfig.Nodes[1].Image != providerNodeImage {
		t.Errorf("expected second node Image to have used provider config value. Actual: %s", globalConfig.Nodes[1].Image)
	}
}

// When the user has a config file with port mappings
// and defines number of nodes via CLI flags,
// port mappings merge into global config nodes
func TestMergeConfigNodesUsesProviderConfigPortMappings(t *testing.T) {
	globalConfig := kindconfig.Cluster{
		Nodes: []kindconfig.Node{
			{
				Role: "control-plane",
			},
			{
				Role: "worker",
			},
		},
	}

	providerConfig := kindconfig.Cluster{
		Nodes: []kindconfig.Node{
			{
				Role: "control-plane",
				ExtraPortMappings: []kindconfig.PortMapping{
					{
						ContainerPort: 777,
						HostPort:      1,
						ListenAddress: "1.2.3.4",
						Protocol:      "TCP",
					},
					{
						ContainerPort: 888,
						HostPort:      2,
						ListenAddress: "2.3.4.5",
						Protocol:      "TCP",
					},
				},
			},
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if len(globalConfig.Nodes) != 2 {
		t.Errorf("expected two nodes from global config. Actual: %v", globalConfig.Nodes)
	}

	if len(globalConfig.Nodes[0].ExtraPortMappings) != 2 {
		t.Errorf("expected first node from to have 2 extra port mappings. Actual: %v", globalConfig.Nodes[0])
	}

	if len(globalConfig.Nodes[1].ExtraPortMappings) != 0 {
		t.Errorf("expected second node to have 0 extra port mappings. Actual: %v", globalConfig.Nodes[1])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[0].ContainerPort != 777 {
		t.Errorf("expected first node's first extra port mapping to have container port from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[0])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[1].ContainerPort != 888 {
		t.Errorf("expected first node's second extra port mapping to have container port from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[0].HostPort != 1 {
		t.Errorf("expected first node's first extra port mapping to have host port from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[0])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[1].HostPort != 2 {
		t.Errorf("expected first node's second extra port mapping to have host port from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[0].ListenAddress != "1.2.3.4" {
		t.Errorf("expected first node's first extra port mapping to have listen address from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[0])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[1].ListenAddress != "2.3.4.5" {
		t.Errorf("expected first node's second extra port mapping to have listen address from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[0].Protocol != "TCP" {
		t.Errorf("expected first node's first extra port mapping to have protocol from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[0])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[1].Protocol != "TCP" {
		t.Errorf("expected first node's second extra port mapping to have protocol from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}
}

// When the user has a config file with port mappings
// and defines number of nodes via CLI flags and port forwarding via CLI flags,
// port mappings merge into global config nodes
//nolint:funlen,gocyclo
func TestMergeConfigNodesMergesPortMappings(t *testing.T) {
	globalConfig := kindconfig.Cluster{
		Nodes: []kindconfig.Node{
			{
				Role: "control-plane",
				ExtraPortMappings: []kindconfig.PortMapping{
					{
						ContainerPort: 111,
						HostPort:      9,
						ListenAddress: "9.9.9.9",
						Protocol:      "TCP",
					},
				},
			},
			{
				Role: "worker",
			},
		},
	}

	providerConfig := kindconfig.Cluster{
		Nodes: []kindconfig.Node{
			{
				Role: "control-plane",
				ExtraPortMappings: []kindconfig.PortMapping{
					{
						ContainerPort: 777,
						HostPort:      1,
						ListenAddress: "1.2.3.4",
						Protocol:      "TCP",
					},
					{
						ContainerPort: 888,
						HostPort:      2,
						ListenAddress: "2.3.4.5",
						Protocol:      "TCP",
					},
				},
			},
			{
				Role: "worker",
				ExtraPortMappings: []kindconfig.PortMapping{
					{
						ContainerPort: 999,
						HostPort:      3,
						ListenAddress: "3.4.5.6",
						Protocol:      "TCP",
					},
				},
			},
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if len(globalConfig.Nodes) != 2 {
		t.Errorf("expected two nodes from global config. Actual: %v", globalConfig.Nodes)
	}

	if len(globalConfig.Nodes[0].ExtraPortMappings) != 3 {
		t.Errorf("expected first node to have 3 extra port mappings. Actual: %v", globalConfig.Nodes[0])
	}

	if len(globalConfig.Nodes[1].ExtraPortMappings) != 1 {
		t.Errorf("expected second node to have 1 extra port mappings. Actual: %v", globalConfig.Nodes[1])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[0].ContainerPort != 111 {
		t.Errorf("expected first node's first extra port mapping to keep container port from global config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[0])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[1].ContainerPort != 777 {
		t.Errorf("expected first node's second extra port mapping to have container port from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[0])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[2].ContainerPort != 888 {
		t.Errorf("expected first node's third extra port mapping to have container port from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[0].HostPort != 9 {
		t.Errorf("expected first node's first extra port mapping to have host port from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[0])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[1].HostPort != 1 {
		t.Errorf("expected first node's second extra port mapping to have host port from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[2].HostPort != 2 {
		t.Errorf("expected first node's second extra port mapping to have host port from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[0].ListenAddress != "9.9.9.9" {
		t.Errorf("expected first node's first extra port mapping to have listen address from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[0])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[1].ListenAddress != "1.2.3.4" {
		t.Errorf("expected first node's second extra port mapping to have listen address from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[2].ListenAddress != "2.3.4.5" {
		t.Errorf("expected first node's second extra port mapping to have listen address from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[0].Protocol != "TCP" {
		t.Errorf("expected first node's first extra port mapping to keep protocol from global config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[0])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[1].Protocol != "TCP" {
		t.Errorf("expected first node's second extra port mapping to have protocol from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[2].Protocol != "TCP" {
		t.Errorf("expected first node's third extra port mapping to have protocol from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}

	if globalConfig.Nodes[1].ExtraPortMappings[0].ContainerPort != 999 {
		t.Errorf("expected second node's extra port mapping to have container port from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}

	if globalConfig.Nodes[1].ExtraPortMappings[0].HostPort != 3 {
		t.Errorf("expected second node's extra port mapping to have host port from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}

	if globalConfig.Nodes[1].ExtraPortMappings[0].ListenAddress != "3.4.5.6" {
		t.Errorf("expected second node's extra port mapping to have listen address from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}

	if globalConfig.Nodes[1].ExtraPortMappings[0].Protocol != "TCP" {
		t.Errorf("expected second node's extra port mapping to have protocol from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}
}

// When the user provides a port mapping but no other flags
// and defines nodes and other mappings via provider config file
func TestMergeConfigNodesUsesCliPortMapping(t *testing.T) {
	globalConfig := kindconfig.Cluster{
		Nodes: []kindconfig.Node{
			{
				Role: "control-plane",
				ExtraPortMappings: []kindconfig.PortMapping{
					{
						ContainerPort: 111,
						HostPort:      9,
					},
				},
			},
		},
	}

	providerConfig := kindconfig.Cluster{
		Nodes: []kindconfig.Node{
			{
				Role: "control-plane",
				ExtraPortMappings: []kindconfig.PortMapping{
					{
						ContainerPort: 777,
						HostPort:      1,
					},
				},
			},
			{
				Role: "worker",
				ExtraPortMappings: []kindconfig.PortMapping{
					{
						ContainerPort: 888,
						HostPort:      2,
					},
					{
						ContainerPort: 999,
						HostPort:      3,
					},
				},
			},
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if len(globalConfig.Nodes) != 2 {
		t.Errorf("expected two nodes from global config. Actual: %v", globalConfig.Nodes)
	}

	if len(globalConfig.Nodes[0].ExtraPortMappings) != 2 {
		t.Errorf("expected first node to have 2 extra port mappings. Actual: %v", globalConfig.Nodes[0])
	}

	if len(globalConfig.Nodes[1].ExtraPortMappings) != 2 {
		t.Errorf("expected second node to have 2 extra port mappings. Actual: %v", globalConfig.Nodes[1])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[0].ContainerPort != 777 {
		t.Errorf("expected first node's first extra port mapping to keep container port from global config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[0])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[1].ContainerPort != 111 {
		t.Errorf("expected first node's second extra port mapping to have container port from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[0])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[0].HostPort != 1 {
		t.Errorf("expected first node's first extra port mapping to have host port from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[0])
	}

	if globalConfig.Nodes[0].ExtraPortMappings[1].HostPort != 9 {
		t.Errorf("expected first node's second extra port mapping to have host port from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}

	if globalConfig.Nodes[1].ExtraPortMappings[0].ContainerPort != 888 {
		t.Errorf("expected second node's first extra port mapping to have container port from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}

	if globalConfig.Nodes[1].ExtraPortMappings[1].ContainerPort != 999 {
		t.Errorf("expected second node's second extra port mapping to have container port from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}

	if globalConfig.Nodes[1].ExtraPortMappings[0].HostPort != 2 {
		t.Errorf("expected second node's first extra port mapping to have host port from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}

	if globalConfig.Nodes[1].ExtraPortMappings[1].HostPort != 3 {
		t.Errorf("expected second node's second extra port mapping to have host port from provider config. Actual: %v", globalConfig.Nodes[0].ExtraPortMappings[1])
	}
}

// Use global config for network IP family when not default
func TestMergeConfigNetworkIPFamily(t *testing.T) {
	globalConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			IPFamily: "1.2.3.4",
		},
	}

	providerConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			IPFamily: "6.7.8.9",
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if globalConfig.Networking.IPFamily != "1.2.3.4" {
		t.Errorf("expected merged kindconfig.Networking.IPFamily to have kept global config value. Actual: %s", globalConfig.Networking.IPFamily)
	}
}

// Use provider config for networking IP faimly when given and global config value is empty
func TestMergeConfigNetworkIPFamilyDefault(t *testing.T) {
	defaultKindconfigCluster, _ := kindConfigFromClusterConfig(config.GenerateDefaultConfig())

	globalConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			IPFamily: defaultKindconfigCluster.Networking.IPFamily,
		},
	}

	providerConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			IPFamily: "6.7.8.9",
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if globalConfig.Networking.IPFamily != "6.7.8.9" {
		t.Errorf("expected merged kindconfig.Networking.IPFamily to have merged from provider config. Actual: %s", globalConfig.Networking.IPFamily)
	}
}

// Use global config for networking API server port when given not default
func TestMergeConfigNetworkAPIServerPort(t *testing.T) {
	globalConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			APIServerPort: 123,
		},
	}

	providerConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			APIServerPort: 789,
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if globalConfig.Networking.APIServerPort != 123 {
		t.Errorf("expected merged kindconfig.Networking.APIServerPort to have kept global config value. Actual: %v", globalConfig.Networking.APIServerPort)
	}
}

// Use provider config for networking API server port when given and global config value is default
func TestMergeConfigNetworkAPIServerPortDefault(t *testing.T) {
	defaultKindconfigCluster, _ := kindConfigFromClusterConfig(config.GenerateDefaultConfig())

	globalConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			APIServerPort: defaultKindconfigCluster.Networking.APIServerPort,
		},
	}

	providerConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			APIServerPort: 789,
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if globalConfig.Networking.APIServerPort != 789 {
		t.Errorf("expected merged kindconfig.Networking.APIServerPort to have merged from provider config. Actual: %v", globalConfig.Networking.APIServerPort)
	}
}

// Use global config for networking API address when not default
func TestMergeConfigNetworkAPIAddress(t *testing.T) {
	globalConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			APIServerAddress: "global-address",
		},
	}

	providerConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			APIServerAddress: "provider-address",
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if globalConfig.Networking.APIServerAddress != "global-address" {
		t.Errorf("expected merged kindconfig.Networking.APIServerAddress to have kept global config value. Actual: %v", globalConfig.Networking.APIServerAddress)
	}
}

// Use provider config for networking API address when given and global config value is default
func TestMergeConfigNetworkAPIAddressDefault(t *testing.T) {
	defaultKindconfigCluster, _ := kindConfigFromClusterConfig(config.GenerateDefaultConfig())

	globalConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			APIServerAddress: defaultKindconfigCluster.Networking.APIServerAddress,
		},
	}

	providerConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			APIServerAddress: "provider-address",
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if globalConfig.Networking.APIServerAddress != "provider-address" {
		t.Errorf("expected merged kindconfig.Networking.APIServerAddress to have merged from provider config. Actual: %v", globalConfig.Networking.APIServerAddress)
	}
}

// Use global config for networking pod subnet when not default
func TestMergeConfigNetworkPodSubnet(t *testing.T) {
	globalConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			PodSubnet: "global-podsubnet",
		},
	}

	providerConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			PodSubnet: "provider-podsubnet",
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if globalConfig.Networking.PodSubnet != "global-podsubnet" {
		t.Errorf("expected merged kindconfig.Networking.PodSubnet to have kept global config value. Actual: %v", globalConfig.Networking.PodSubnet)
	}
}

// Use provider config for networking pod subnet when given and global config value is default
func TestMergeConfigNetworkPodSubnetDefault(t *testing.T) {
	defaultKindconfigCluster, _ := kindConfigFromClusterConfig(config.GenerateDefaultConfig())

	globalConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			PodSubnet: defaultKindconfigCluster.Networking.PodSubnet,
		},
	}

	providerConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			PodSubnet: "provider-podsubnet",
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if globalConfig.Networking.PodSubnet != "provider-podsubnet" {
		t.Errorf("expected merged kindconfig.Networking.PodSubnet to have merged from provider config. Actual: %v", globalConfig.Networking.PodSubnet)
	}
}

// Use global config for networking service subnet when not default
func TestMergeConfigNetworkServiceSubnet(t *testing.T) {
	globalConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			ServiceSubnet: "global-servicesubnet",
		},
	}

	providerConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			ServiceSubnet: "provider-servicesubnet",
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if globalConfig.Networking.ServiceSubnet != "global-servicesubnet" {
		t.Errorf("expected merged kindconfig.Networking.ServiceSubnet to have kept global config value. Actual: %v", globalConfig.Networking.ServiceSubnet)
	}
}

// Use provider config for networking service subnet when given and global config value is default
func TestMergeConfigNetworkServiceSubnetDefault(t *testing.T) {
	defaultKindconfigCluster, _ := kindConfigFromClusterConfig(config.GenerateDefaultConfig())

	globalConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			ServiceSubnet: defaultKindconfigCluster.Networking.ServiceSubnet,
		},
	}

	providerConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			ServiceSubnet: "provider-servicesubnet",
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if globalConfig.Networking.ServiceSubnet != "provider-servicesubnet" {
		t.Errorf("expected merged kindconfig.Networking.ServiceSubnet to have merged from provider config. Actual: %v", globalConfig.Networking.ServiceSubnet)
	}
}

// Use global config for networking kubeproxy mode when not default
func TestMergeConfigNetworkKubeproxy(t *testing.T) {
	globalConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			KubeProxyMode: "global-kubeproxymode",
		},
	}

	providerConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			KubeProxyMode: "provider-kubeproxymode",
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if globalConfig.Networking.KubeProxyMode != "global-kubeproxymode" {
		t.Errorf("expected merged kindconfig.Networking.KubeProxyMode to have kept global config value. Actual: %v", globalConfig.Networking.KubeProxyMode)
	}
}

// Use provider config for networking kubeproxy mode when given and global config value is default
func TestMergeConfigNetworkKubeproxyDefault(t *testing.T) {
	defaultKindconfigCluster, _ := kindConfigFromClusterConfig(config.GenerateDefaultConfig())

	globalConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			KubeProxyMode: defaultKindconfigCluster.Networking.KubeProxyMode,
		},
	}

	providerConfig := kindconfig.Cluster{
		Networking: kindconfig.Networking{
			KubeProxyMode: "provider-kubeproxymode",
		},
	}

	err := mergeConfigsLeft(&globalConfig, &providerConfig)
	if err != nil {
		t.Errorf("expected merge to succeed")
	}

	if globalConfig.Networking.KubeProxyMode != "provider-kubeproxymode" {
		t.Errorf("expected merged kindconfig.Networking.KubeProxyMode to have merged from provider config. Actual: %v", globalConfig.Networking.KubeProxyMode)
	}
}
