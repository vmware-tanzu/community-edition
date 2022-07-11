// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kubeconfig

import (
	"fmt"
)

// KINDClusterKey identifies kind clusters in kubeconfig files
func KINDClusterKey(clusterName string) string {
	return clusterName
}

// checkKubeadmExpectations validates that a kubeadm created KUBECONFIG meets
// our expectations, namely on the number of entries
func checkKubeadmExpectations(cfg *Config) error {
	if len(cfg.Clusters) != 1 {
		return fmt.Errorf("kubeadm KUBECONFIG should have one cluster, but read %d", len(cfg.Clusters))
	}
	if len(cfg.Users) != 1 {
		return fmt.Errorf("kubeadm KUBECONFIG should have one user, but read %d", len(cfg.Users))
	}
	if len(cfg.Contexts) != 1 {
		return fmt.Errorf("kubeadm KUBECONFIG should have one context, but read %d", len(cfg.Contexts))
	}
	return nil
}
