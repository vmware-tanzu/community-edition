// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/mock"
	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/internal/docker/mocks"
)

var (
	fakeDocker *mocks.DockerClient
)

func TestForceStopAndDeleteCluster(t *testing.T) {
	NewClientWithOpts = func(ops ...client.Opt) (*client.Client, error) {
		return &client.Client{}, nil
	}
	clientOnce.Do(func() {})
	fakeDocker = mocks.NewDockerClient(t)
	dockerClient = fakeDocker
	filter := filters.NewArgs()
	filter.Add("name", "tanzu-community-edition-control-plane")
	fakeDocker.On("ContainerList", mock.Anything, mock.Anything).Return(
		[]types.Container{
			{
				ID:    "12333",
				Names: []string{"tanzu-community-edition-control-plane"},
				Image: "image:tag",
			},
		}, nil,
	).Once()

	fakeDocker.On("ContainerStop", mock.Anything, mock.MatchedBy(func(s string) bool {
		return s == "12333"
	}), mock.Anything).Return(nil).Once()

	fakeDocker.On("ContainerRemove", mock.Anything, mock.MatchedBy(func(s string) bool {
		return s == "12333"
	}), mock.Anything).Return(nil)
	if err := ForceStopAndDeleteCluster(); err != nil {
		t.Errorf("Expected err : %v, got %v", nil, err)
	}
}
