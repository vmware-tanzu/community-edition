// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

type collectCommonArgs struct {
	workDir   string
	outputDir string
}

type collectBootsrapArgs struct {
	skip        bool
	clusterName string
}

type collectWorkloadArgs struct {
	infra       string
	kubeconfig  string
	contextName string
	clusterName string
	namespace   string
}

type collectMgmtArgs struct {
	skip        bool
	kubeconfig  string
	contextName string
	clusterName string
}

type collectStandaloneArgs struct {
	kubeconfig  string
	clusterName string
	contextName string
}

type managementServer struct {
	clusterName string
	kubeconfig  string
	kubecontext string
}
