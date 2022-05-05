// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package checks

import (
	"errors"
	"fmt"

	"github.com/juju/fslock"

	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/internal/docker"
	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/utils"
)

const (
	Running      = 0
	NotExist     = 1
	NotRunning   = 2
	Error        = -1
	runningState = "running"
)

func IsClusterCreating() (bool, error) {
	lock, err := utils.GetLockForFile(utils.GetClusterCreateLockFilename())
	if err != nil {
		return false, err
	}

	if err = lock.TryLock(); err == fslock.ErrLocked {
		return true, nil
	}
	return false, nil
}

func IsClusterDeleting() (bool, error) {
	lock, err := utils.GetLockForFile(utils.GetClusterDeleteLockFilename())
	if err != nil {
		return false, err
	}

	if err = lock.TryLock(); err == fslock.ErrLocked {
		return true, err
	}
	return false, nil
}

// Returns whether it's running, if not running, if can be created, and error message
func IsClusterUpAndRunning() (bool, bool, error) {
	containers, err := docker.GetAllTCEContainers()
	if err != nil {
		return false, false, fmt.Errorf("error executing docker command (%s)", err.Error())
	}

	// Checks: Only one container
	// It must be running
	// It must be of image projects.registry.vmware.com/tce/kind/node:v1.22.5
	// It must be in running state
	if len(containers) < 1 {
		return false, true, errors.New("cluster does not exist")
	}

	tceContainer := containers[0]
	/*
		TODO: See what we do with this check
		if tceContainer.Image != config.DefaultImage {
			return false, false, fmt.Errorf("a cluster is running and using a different image: [%s]", tceContainer.Image)
		}
	*/
	if tceContainer.State != runningState {
		return false, false, fmt.Errorf("Cluster exists but the container is %s", tceContainer.State)
	}

	return true, false, nil
}

// CanTCEContainerBeCreated checks if there are existing cluster containers running.
func CanTCEContainerBeCreated() (bool, error) {
	// OK if no other tce-container exists
	containers, err := docker.GetAllTCEContainers()
	if err != nil {
		return false, fmt.Errorf("error executing docker command (%s)", err.Error())
	}
	if len(containers) == 1 {
		tceContainer := containers[0]
		if tceContainer.State != runningState {
			return false, fmt.Errorf("Cluster exists but the container is %s", tceContainer.State)
		}
		return false, errors.New("Cluster already exist, or a container with same name")
	}
	return true, nil
}

// GetContainerClusterStatus checks the running state of the cluster container.
func GetContainerClusterStatus() (int, string) {
	containers, err := docker.GetAllTCEContainers()
	if err != nil {
		return Error, fmt.Sprintf("Error executing docker command (%s)", err.Error())
	}

	if len(containers) < 1 {
		return NotExist, "Cluster does not exist"
	}

	tceContainer := containers[0]
	/*
		TODO: See what we do with this check
		if tceContainer.Image != config.DefaultImage {
			return false, false, fmt.Errorf("a cluster is running and using a different image: [%s]", tceContainer.Image)
		}
	*/
	if tceContainer.State != runningState {
		return NotRunning, fmt.Sprintf("Cluster exists but the container is %s", tceContainer.State)
	}

	return Running, "Cluster is running"
}
