// +build tools

// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package tools imports things required by build scripts, to force `go mod` to see them as dependencies
package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/shuLhan/go-bindata" // Force load of go-bindata
)
