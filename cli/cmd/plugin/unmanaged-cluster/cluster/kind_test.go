// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cluster

import (
	"encoding/json"
	"testing"
)

var normalDockerInfoJSON = `{"ID":"SEB7:L67H:GZMX:VPIN:YZ7V:RTRC:DCML:3C7C:PNN3:2DQA:6GD2:ZIWU","Containers":7,"ContainersRunning":1,"ContainersPaused":0,"ContainersStopped":6,"Images":151,"Driver":"overlay2","DriverStatus":[["Backing Filesystem","extfs"],["Supports d_type","true"],["Native Overlay Diff","true"],["userxattr","false"]],"Plugins":{"Volume":["local"],"Network":["bridge","host","ipvlan","macvlan","null","overlay"],"Authorization":null,"Log":["awslogs","fluentd","gcplogs","gelf","journald","json-file","local","logentries","splunk","syslog"]},"MemoryLimit":true,"SwapLimit":true,"KernelMemory":true,"KernelMemoryTCP":true,"CpuCfsPeriod":true,"CpuCfsQuota":true,"CPUShares":true,"CPUSet":true,"PidsLimit":true,"IPv4Forwarding":true,"BridgeNfIptables":true,"BridgeNfIp6tables":true,"Debug":false,"NFd":32,"OomKillDisable":true,"NGoroutines":40,"SystemTime":"2022-01-11T15:43:55.314860422-06:00","LoggingDriver":"json-file","CgroupDriver":"cgroupfs","CgroupVersion":"1","NEventsListener":0,"KernelVersion":"5.11.0-43-generic","OperatingSystem":"Ubuntu 20.04.3 LTS","OSVersion":"20.04","OSType":"linux","Architecture":"x86_64","IndexServerAddress":"https://index.docker.io/v1/","RegistryConfig":{"AllowNondistributableArtifactsCIDRs":[],"AllowNondistributableArtifactsHostnames":[],"InsecureRegistryCIDRs":["127.0.0.0/8"],"IndexConfigs":{"docker.io":{"Name":"docker.io","Mirrors":[],"Secure":true,"Official":true}},"Mirrors":[]},"NCPU":16,"MemTotal":33613119488,"GenericResources":null,"DockerRootDir":"/var/lib/docker","HttpProxy":"","HttpsProxy":"","NoProxy":"","Name":"sm-workstation","Labels":[],"ExperimentalBuild":false,"ServerVersion":"20.10.12","Runtimes":{"io.containerd.runc.v2":{"path":"runc"},"io.containerd.runtime.v1.linux":{"path":"runc"},"runc":{"path":"runc"}},"DefaultRuntime":"runc","Swarm":{"NodeID":"","NodeAddr":"","LocalNodeState":"inactive","ControlAvailable":false,"Error":"","RemoteManagers":null},"LiveRestoreEnabled":false,"Isolation":"","InitBinary":"docker-init","ContainerdCommit":{"ID":"7b11cfaabd73bb80907dd23182b9347b4245eb5d","Expected":"7b11cfaabd73bb80907dd23182b9347b4245eb5d"},"RuncCommit":{"ID":"v1.0.2-0-g52b36a2","Expected":"v1.0.2-0-g52b36a2"},"InitCommit":{"ID":"de40ad0","Expected":"de40ad0"},"SecurityOptions":["name=apparmor","name=seccomp,profile=default"],"Warnings":null,"ClientInfo":{"Debug":false,"Context":"default","Plugins":[{"SchemaVersion":"0.1.0","Vendor":"Docker Inc.","Version":"v0.9.1-beta3","ShortDescription":"Docker App","Experimental":true,"Name":"app","Path":"/usr/libexec/docker/cli-plugins/docker-app"},{"SchemaVersion":"0.1.0","Vendor":"Docker Inc.","Version":"v0.7.1-docker","ShortDescription":"Docker Buildx","Name":"buildx","Path":"/usr/libexec/docker/cli-plugins/docker-buildx"},{"SchemaVersion":"0.1.0","Vendor":"Docker Inc.","Version":"v0.12.0","ShortDescription":"Docker Scan","Name":"scan","Path":"/usr/libexec/docker/cli-plugins/docker-scan"}],"Warnings":null}}`

func TestValidateDockerInfoNoIssues(t *testing.T) {
	// Test the full real response from docker info, we'll use a truncated response below
	errs := validateDockerInfo([]byte(normalDockerInfoJSON))
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
	errs := validateDockerInfo(output)
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
	errs := validateDockerInfo(output)
	if len(errs) != 1 {
		t.Errorf("expected 1 error but %d returned", len(errs))
	}
}

func TestValidateDockerInfoArchitectureNotSupported(t *testing.T) {
	testInfo := dockerInfo{
		CPUs:         16,
		Memory:       2147483648,
		Architecture: "arm64",
	}
	output, _ := json.Marshal(testInfo)
	errs := validateDockerInfo(output)
	if len(errs) != 1 {
		t.Errorf("expected 1 error but %d returned", len(errs))
	}
}

func TestValidateDockerInfoBadData(t *testing.T) {
	errs := validateDockerInfo([]byte{240, 159, 146, 169})
	if len(errs) != 1 {
		t.Errorf("expected 1 error but %d returned", len(errs))
	}
}
