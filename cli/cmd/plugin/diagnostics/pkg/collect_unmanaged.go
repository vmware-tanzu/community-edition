// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"bytes"
	"fmt"

	crashdexec "github.com/vmware-tanzu/crash-diagnostics/exec"
)

// collectUnmanagedDiags runs command/scripts to collect unmanaged cluster diagnostics data
func collectUnmanagedDiags() error {
	if unmanagedArgs.clusterName == "" {
		return fmt.Errorf("unmanaged cluster: cluster name not set")
	}
	if unmanagedArgs.kubeconfig == "" {
		unmanagedArgs.kubeconfig = getDefaultKubeconfig()
	}
	if unmanagedArgs.contextName == "" {
		unmanagedArgs.contextName = getDefaultClusterContext(unmanagedArgs.clusterName)
	}

	scriptName := umScriptPath
	argsMap := crashdexec.ArgMap{
		"workdir":                 commonArgs.workDir,
		"outputdir":               commonArgs.outputDir,
		"unmanaged_cluster_name": unmanagedArgs.clusterName,
		"unmanaged_kubeconfig":   unmanagedArgs.kubeconfig,
		"unmanaged_context":      unmanagedArgs.contextName,
	}

	libScript := libScriptPath
	libData, err := scriptFS.ReadFile(libScript)
	if err != nil {
		return err
	}

	scriptData, err := scriptFS.ReadFile(scriptName)
	if err != nil {
		return err
	}

	return crashdexec.ExecuteWithModules(
		scriptName,
		bytes.NewReader(scriptData),
		argsMap,
		crashdexec.StarlarkModule{Name: libScript, Source: bytes.NewReader(libData)},
	)
}
