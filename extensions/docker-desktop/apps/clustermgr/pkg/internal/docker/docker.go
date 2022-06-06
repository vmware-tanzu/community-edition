// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/config"
)

func GetDockerInfo() (types.Info, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return types.Info{}, err
	}

	info, err := cli.Info(context.Background())
	if err != nil {
		return types.Info{}, err
	}
	return info, nil
}

func GetAllTCEContainers() ([]types.Container, error) {
	f := filters.NewArgs()
	f.Add("name", config.GetTCEContainerName())

	// TODO: Extract init of Docker cli to common init function
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: f})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

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

func ForceStopAndDeleteCluster() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
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
