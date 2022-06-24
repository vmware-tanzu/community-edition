// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package checks

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/internal/docker"
)

const minCPUCount = 1
const minMemBytes = 2 * 1024 * 1024 * 1024 // 2 GB

var (
	dockerInfoHandler = docker.GetDockerInfo
)

func testLocalPorts(ports []string) error {
	for _, port := range ports {
		intVar, err := strconv.Atoi(port)
		if err == nil && intVar <= 1024 {
			conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", port), 1*time.Second)
			if err == nil {
				conn.Close()
				return fmt.Errorf("port %s is in Use", port)
			}
		} else {
			dummyListener, err := net.Listen("tcp4", net.JoinHostPort("127.0.0.1", port))
			if err != nil {
				return fmt.Errorf("port %s is in Use (%s)", port, err.Error())
			}
			dummyListener.Close()
		}
	}
	return nil
}

func testCPUandMemory() error {
	info, err := dockerInfoHandler()
	if err != nil {
		return err
	}

	if info.NCPU < minCPUCount {
		return fmt.Errorf("not enough CPUs configured, currently set to %d CPUs", info.NCPU)
	}
	if info.MemTotal < minMemBytes {
		return fmt.Errorf("not enough Memory configured, currently set to %d bytes, %d (%dGB) required", info.MemTotal, minMemBytes, minMemBytes/1024/1024/1024)
	}
	return nil
}

// PreflightChecks performs a set of checks before cluster creation to check
// for minimum requirements or known issues that would prevent successful creation.
func PreflightChecks() error {
	var err error
	// TODO: Get ports (and IP) from config
	err = testLocalPorts([]string{"80", "443"})
	if err != nil {
		// TODO: Log
		return err
	}

	err = testCPUandMemory()
	if err != nil {
		return err
	}

	return nil
}
