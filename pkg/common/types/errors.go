// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"errors"
)

var (
	// ErrExtensionNotFound is extension not found
	ErrExtensionNotFound = errors.New("Extension not found")
	// ErrVersionNotFound is version not found
	ErrVersionNotFound = errors.New("Version not found")
)
