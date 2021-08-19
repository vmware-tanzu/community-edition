// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/config"
)

func getDefaultWorkdir() string {
	tanzuConfigDir, err := config.LocalDir()
	if err != nil {
		return "./diagnostics"
	}
	return filepath.Join(tanzuConfigDir, "diagnostics")
}

func getDefaultOutputDir() string {
	return "./"
}

func getDefaultKubeconfig() string {
	kcfg := os.Getenv("KUBECONFIG")
	if kcfg == "" {
		kcfg = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}
	return kcfg
}

func getDefaultClusterContext(clusterName string) string {
	return fmt.Sprintf("%s-admin@%s", clusterName, clusterName)
}

func getDefaultManagementServer() (*managementServer, error) {
	svr, err := config.GetCurrentServer()
	if err != nil {
		return nil, err
	}

	if !svr.IsManagementCluster() {
		return nil, fmt.Errorf("current server not management instance")
	}

	return &managementServer{
		name:        svr.Name,
		kubeconfig:  svr.ManagementClusterOpts.Path,
		kubecontext: svr.ManagementClusterOpts.Context,
	}, nil
}

