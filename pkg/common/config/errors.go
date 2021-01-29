// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"errors"
)

var (
	// ErrDatawriterFailed is Datawriter is empty
	ErrDatawriterFailed = errors.New("Datawriter is empty")
	// ErrDatareaderFailed is Datareader is empty
	ErrDatareaderFailed = errors.New("Datareader is empty")
)
