// +build tools

// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package builder imports dependencies needed for building and forces `go mod` to see them as dependencies
package builder

import (
	_ "github.com/spf13/cobra"
	_ "github.com/vmware-tanzu/tanzu-framework"
	_ "github.com/vmware-tanzu/tanzu-framework/pkg/v1/builder/command"
	_ "github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli"
	_ "github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/command/core"
	_ "github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/component"
)
