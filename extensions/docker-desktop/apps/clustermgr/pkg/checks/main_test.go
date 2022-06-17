// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package checks

import (
	"os"
	"testing"

	"github.com/docker/docker/client"

	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/internal/docker"
)

func TestMain(m *testing.M) {
	docker.NewClientWithOpts = func(ops ...client.Opt) (*client.Client, error) {
		return &client.Client{}, nil
	}
	os.Exit(m.Run())
}
