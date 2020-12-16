// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package download

import (
	klog "k8s.io/klog/v2"
)

// NewManager generates a Manager object
func NewManager(byConfig []byte) (*Manager, error) {

	cfg, err := InitDownloadConfig(byConfig)
	if err != nil {
		klog.Errorf("InitGitHubConfig failed. Err: %v", err)
		return nil, err
	}

	m := &Manager{
		cfg: cfg,
	}

	return m, nil
}
