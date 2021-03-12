// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	config "github.com/vmware-tanzu/tce/pkg/common/config"
	kapp "github.com/vmware-tanzu/tce/pkg/common/kapp"
)

// Manager encapsulates everything about how to manage extensions
type Manager struct {
	// config manager
	cfg *config.Config
	// kapp manaer
	kapp *kapp.Kapp
}
