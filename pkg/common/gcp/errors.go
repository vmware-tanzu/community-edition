// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package gcp

import (
	"errors"
)

var (
	// ErrBucketConnect is returned when the provided username is empty.
	ErrBucketConnect = errors.New("Could not connect to repository")
	// ErrBucketDownload is returned when we cannot download the file
	ErrBucketDownload = errors.New("Could not fetch artifact from repository")
	// ErrBucketObject is bucket object failed
	ErrBucketObject = errors.New("Bucket object failed")
)
