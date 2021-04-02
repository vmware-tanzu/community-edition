// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kapp

import (
	"errors"
)

var (
	// ErrAppNotPresentOrInstalled is Application is not present/installed
	ErrAppNotPresentOrInstalled = errors.New("application is not present/installed")
)
