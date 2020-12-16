// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package extension

import (
	config "github.com/vmware-tanzu/tce/pkg/common/config"
	gcp "github.com/vmware-tanzu/tce/pkg/common/gcp"
	github "github.com/vmware-tanzu/tce/pkg/common/github"
	kapp "github.com/vmware-tanzu/tce/pkg/common/kapp"
	types "github.com/vmware-tanzu/tce/pkg/common/types"
)

// Manager encapsulates everything about how to manage extensions
type Manager struct {
	// config manager
	cfg *config.Config
	// GitHub manager
	gh *github.Manager
	// GCP Bucket manager
	b *gcp.Bucket
	// kapp manaer
	kapp *kapp.Kapp

	// CACHE
	metadata *types.Metadata
	release  *types.Release
	// CACHE
}
