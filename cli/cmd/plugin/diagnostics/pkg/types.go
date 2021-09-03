// Copyright 2021 VMware, Inc. All Rights Reserved.
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
	standalone  bool
	kubeconfig  string
	infra       string
	clusterName string
	namespace   string
	contextName string
	sshUser     string
	sskPkPath   string
}

type collectMgmtArgs struct {
	skip bool
	kubeconfig  string
	contextName string
	clusterName string
}

type managementServer struct {
	name        string
	kubeconfig  string
	kubecontext string
}

