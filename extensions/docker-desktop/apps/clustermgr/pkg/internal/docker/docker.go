// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"fmt"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/config"
)

var (
	NewClientWithOpts = client.NewClientWithOpts
	clientOnce        sync.Once
	dockerClient      DockerInterface
)

func getDockerClient() DockerInterface {
	clientOnce.Do(func() {
		cli, err := NewClientWithOpts(client.FromEnv)
		if err != nil {
			panic(err)
		}
		dockerClient = cli
	})
	return dockerClient
}

// GetDockerInfo gets the Docker engine runtime info.
func GetDockerInfo() (types.Info, error) {
	cli := getDockerClient()
	info, err := cli.Info(context.Background())
	if err != nil {
		return types.Info{}, err
	}
	return info, nil
}

// GetAllTCEContainers queries the Docker engine for our unmanaged cluster container.
func GetAllTCEContainers() ([]types.Container, error) {
	f := filters.NewArgs()
	f.Add("name", config.GetTCEContainerName())

	cli := getDockerClient()
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: f})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// GetTCEContainerID gets the ID of the cluster container.
func GetTCEContainerID() (string, error) {
	containers, err := GetAllTCEContainers()
	if err != nil {
		return "", err
	}
	if len(containers) != 1 {
		return "", fmt.Errorf("TCE container not found")
	}
	return containers[0].ID, nil
}

// ForceStopAndDeleteCluster will force stopping the cluster container and delete it.
func ForceStopAndDeleteCluster() error {
	cli := getDockerClient()
	containerID, err := GetTCEContainerID()
	if err != nil {
		return err
	}
	err = cli.ContainerStop(context.Background(), containerID, nil)
	if err != nil {
		return err
	}
	err = cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{RemoveVolumes: true, Force: true})
	if err != nil {
		return err
	}
	return nil
}
