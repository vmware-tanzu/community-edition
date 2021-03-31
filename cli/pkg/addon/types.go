// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	kapp "github.com/vmware-tanzu/tce/cli/pkg/common/kapp"
)

// Manager encapsulates everything about how to manage extensions
type Manager struct {
	// kapp manaer
	kapp *kapp.Kapp
}
