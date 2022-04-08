// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package containerruntime has functions for interacting with the local container
// runtime engine.
package containerruntime

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// RuntimeInfo is information about the current container runtime environment.
type RuntimeInfo struct {
	Name         string `json:"Name"`
	OSType       string `json:"OSType"`
	OSVersion    string `json:"OSVersion"`
	Architecture string `json:"Architecture"`
	CPU          int    `json:"NCPU"`
	Memory       int64  `json:"MemTotal"`
	Containers   int    `json:"Containers"`
}

func GetRuntimeInfo() (*RuntimeInfo, error) {
	// Check presence of docker
	cmd := exec.Command("docker", "ps")
	if err := cmd.Run(); err != nil {
		// Doesn't appear there is a running container engine, as far as we can tell
		return nil, fmt.Errorf("docker is not installed or not reachable; verify it's installed, running, and your user has permissions to interact with it. Error when attempting to run docker ps: %w", err)
	}

	// Get Docker info
	cmd = exec.Command("docker", "info", "--format", "{{ json . }}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("unable to get docker info: %w", err)
	}

	result := &RuntimeInfo{}
	err = json.Unmarshal(output, result)
	if err != nil {
		return nil, fmt.Errorf("unable to parse docker info: %w", err)
	}

	return result, nil
}
