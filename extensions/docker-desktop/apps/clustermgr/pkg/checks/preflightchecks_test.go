// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package checks

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/docker/docker/api/types"
)

func TestLocalPorts(t *testing.T) {
	// TODO: Test for Port < 1024 needs to have special permission to simulate binds
	// need to find out what works in the CI/CD.
	tests := []struct {
		name         string
		portToBind   int
		bindType     string
		portsToCheck []string
		err          bool
	}{
		{
			name:         "One of the two port is in use",
			portToBind:   45664,
			bindType:     "tcp4",
			portsToCheck: []string{"45664", "45667"},
			err:          true,
		},
		{
			name:         "none of the ports are bound",
			portToBind:   0,
			bindType:     "tcp4",
			portsToCheck: []string{"45664", "45667"},
			err:          false,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.portToBind != 0 {
				l, err := net.Listen(tc.bindType, fmt.Sprintf("127.0.0.1:%d", tc.portToBind))
				if err != nil {
					t.Errorf("Did not expect an error trying to bind a port. Got: %v", err)
				}
				defer l.Close()
			}
			err := testLocalPorts(tc.portsToCheck)
			if tc.err && err == nil {
				t.Errorf("testLocalPorts(%v) = %v. Expected an error", tc.portsToCheck, err)
			}
			if !tc.err && err != nil {
				t.Errorf("testLocalPorts(%v) = %v. Did not an error", tc.portsToCheck, err)
			}
		})
	}
}

func TestPreflightChecks(t *testing.T) {
	tests := []struct {
		name       string
		portToBind int
		bindType   string
		dockerInfo types.Info
		infoError  error
		err        bool
	}{
		{
			name:       "Resources are not matching",
			portToBind: 0,
			bindType:   "",
			dockerInfo: types.Info{
				NCPU:     1,
				MemTotal: 1 * 1024 * 1024 * 1024,
			},
			infoError: nil,
			err:       true,
		},
		{
			name:       "Failed to fetch resources from Docker",
			portToBind: 0,
			bindType:   "",
			dockerInfo: types.Info{},
			infoError:  errors.New("Simulated error with docker interaction"),
			err:        true,
		},
		{
			name:       "All healthy",
			portToBind: 0,
			bindType:   "",
			dockerInfo: types.Info{
				NCPU:     3,
				MemTotal: 3 * 1024 * 1024 * 1024,
			},
			infoError: nil,
			err:       false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dockerInfoHandler = func() (types.Info, error) {
				return tc.dockerInfo, tc.infoError
			}
			err := PreflightChecks()
			if tc.err && err == nil {
				t.Errorf("PreflightChecks() = %v. Expected an error", err)
			}
			if !tc.err && err != nil {
				t.Errorf("testLocalPorts() = %v. Did not an error", err)
			}
		})
	}
}
