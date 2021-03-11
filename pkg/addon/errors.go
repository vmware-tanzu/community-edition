// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	"errors"
)

var (
	// ErrMissingToken is missing GitHub token
	ErrMissingToken = errors.New("Missing GitHub token")
	// ErrMissingPackageName is missing extension name
	ErrMissingPackageName = errors.New("Missing package name")
)
