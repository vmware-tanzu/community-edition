// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"os"
	"path/filepath"
	"time"

	"github.com/juju/fslock"

	"github.com/vmware-tanzu/community-edition/errors"
)

const (
	// DefaultLockTimeout is the default time waiting on the filelock
	DefaultLockTimeout = 10 * time.Minute
)

// GetClusterCreateLockFilename gets the full path to our create lock.
func GetClusterCreateLockFilename() string {
	return filepath.Join(os.TempDir(), "cluster-create.lck")
}

// GetClusterDeleteLockFilename gets the full path to our delete lock.
func GetClusterDeleteLockFilename() string {
	return filepath.Join(os.TempDir(), "cluster-delete.lck")
}

// GetFileLockWithDefaultTimeOut returns a file lock with the default timeout of 10 minutes.
func GetFileLockWithDefaultTimeOut(lockPath string) (*fslock.Lock, error) {
	return GetFileLockWithTimeOut(lockPath, DefaultLockTimeout)
}

// GetFileLockWithTimeOut returns a file lock with timeout.
func GetFileLockWithTimeOut(lockPath string, lockDuration time.Duration) (*fslock.Lock, error) {
	lock, err := GetLockForFile(lockPath)

	if err != nil {
		return nil, err
	}

	if err := lock.LockWithTimeout(lockDuration); err != nil {
		return &fslock.Lock{}, errors.NewTimeout(err, "failed to acquire a lock with timeout")
	}
	return lock, nil
}

// GetLockForFile gets a new lock for a given lock file path.
func GetLockForFile(lockPath string) (*fslock.Lock, error) {
	dir := filepath.Dir(lockPath)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return nil, err
		}
	}

	return fslock.New(lockPath), nil
}
