// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
)

// Interface is a stripped down version of the functionalities provided by the Docker Client
// This is setup in order to ease the mocking aspects of Unit Test that is required. As and when a
// new function of the Docker Client needs to be accessed, the same can be added here and invoked
// from respective places
type Interface interface {
	Info(context.Context) (types.Info, error)
	ContainerList(context.Context, types.ContainerListOptions) ([]types.Container, error)
	ContainerStop(context.Context, string, *time.Duration) error
	ContainerRemove(context.Context, string, types.ContainerRemoveOptions) error
}
