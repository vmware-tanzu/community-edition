// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package checks

import (
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
)

func TestGetContainerClusterStatus(t *testing.T) {
	tests := []struct {
		name             string
		mockedContainers []types.Container
		mockedErr        error
		state            int
		message          string
	}{
		{
			name:             "Failed to fetch container information",
			mockedContainers: []types.Container{},
			mockedErr:        errors.New("Fake error from Docker Interaction"),
			state:            Error,
			message:          "Error executing docker command (Fake error from Docker Interaction)",
		},
		{
			name:             "No Containers exists",
			mockedContainers: []types.Container{},
			mockedErr:        nil,
			state:            NotExist,
			message:          "Cluster does not exist",
		},
		{
			name: "Exists but not running",
			mockedContainers: []types.Container{
				{
					ID: "12344",
					Names: []string{
						"tanzu-commuinity-edition-control-plane",
					},
					State: "exited",
				},
			},
			mockedErr: nil,
			state:     NotRunning,
			message:   "Cluster exists but the container is exited",
		},
		{
			name: "Healthy state",
			mockedContainers: []types.Container{
				{
					ID: "12344",
					Names: []string{
						"tanzu-commuinity-edition-control-plane",
					},
					State: "running",
				},
			},
			mockedErr: nil,
			state:     Running,
			message:   "Cluster is running",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			getAllTCEContainerHandler = func() ([]types.Container, error) {
				return tc.mockedContainers, tc.mockedErr
			}
			state, msg := GetContainerClusterStatus()
			if tc.state != state {
				t.Errorf("State Check Expected: %v, got: %v", tc.state, state)
			}

			if tc.message != msg {
				t.Errorf("Message check Expected: %s, got: %s", tc.message, msg)
			}
		})
	}
}

func TestIsClusterUpAndRunning(t *testing.T) {
	tests := []struct {
		name             string
		mockedContainers []types.Container
		mockedErr        error
		ok               bool
		canBeCreated     bool
		err              error
	}{
		{
			name:             "Failed to fetch container information",
			mockedContainers: []types.Container{},
			mockedErr:        errors.New("Fake error from Docker Interaction"),
			ok:               false,
			canBeCreated:     false,
			err:              errors.New("error executing docker command (Fake error from Docker Interaction)"),
		},
		{
			name:             "No Containers exists",
			mockedContainers: []types.Container{},
			mockedErr:        nil,
			ok:               false,
			canBeCreated:     true,
			err:              errors.New("Cluster does not exist"),
		},
		{
			name: "Exists but not running",
			mockedContainers: []types.Container{
				{
					ID: "12344",
					Names: []string{
						"tanzu-commuinity-edition-control-plane",
					},
					State: "exited",
				},
			},
			mockedErr:    nil,
			ok:           false,
			canBeCreated: false,
			err:          errors.New("Cluster exists but the container is exited"),
		},
		{
			name: "Healthy state",
			mockedContainers: []types.Container{
				{
					ID: "12344",
					Names: []string{
						"tanzu-commuinity-edition-control-plane",
					},
					State: "running",
				},
			},
			mockedErr:    nil,
			ok:           true,
			canBeCreated: false,
			err:          nil,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			getAllTCEContainerHandler = func() ([]types.Container, error) {
				return tc.mockedContainers, tc.mockedErr
			}
			ok, exists, err := IsClusterUpAndRunning()
			if tc.ok != ok {
				t.Errorf("OK Check Expected: %v, got: %v", tc.ok, ok)
			}

			if tc.canBeCreated != exists {
				t.Errorf("CanBeCreated Check Expected: %v, got: %v", tc.canBeCreated, exists)
			}

			if tc.err != nil && err != nil {
				if tc.err.Error() != err.Error() {
					t.Errorf("Error Check Expected: %v, got: %v", tc.err, err)
				}
			} else {
				if tc.err != err {
					t.Errorf("Error Check Expected: %v, got: %v", tc.err, err)
				}
			}
		})
	}
}
