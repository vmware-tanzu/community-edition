// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package extension

import (
	"errors"
)

var (
	// ErrMissingToken is missing GitHub token
	ErrMissingToken = errors.New("Missing GitHub token")
	// ErrMissingExtensionName is missing extension name
	ErrMissingExtensionName = errors.New("Missing extension name")
)
