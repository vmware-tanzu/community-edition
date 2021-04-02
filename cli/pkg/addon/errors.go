// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	"errors"
)

var (
	// ErrMissingToken is missing GitHub token
	ErrMissingToken = errors.New("missing GitHub token")
	// ErrMissingPackageName is missing extension name
	ErrMissingPackageName = errors.New("missing package name")

	// ErrMissingOperation is missing operation
	ErrMissingOperation = errors.New("missing sub operation or command")
	// ErrMissingParameter is missing a required parameter
	ErrMissingParameter = errors.New("missing a required parameter")
)
