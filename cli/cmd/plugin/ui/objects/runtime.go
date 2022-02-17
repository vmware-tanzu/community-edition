// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package objects

type Runtime struct {
	Name string `json:"name,omitempty"`

	OSType string `json:"ostype,omitempty"`

	OSVersion string `json:"osversion,omitempty"`

	CPU int32 `json:"cpu,omitempty"`

	Memory int64 `json:"memory,omitempty"`

	Containers int32 `json:"containers,omitempty"`
}
