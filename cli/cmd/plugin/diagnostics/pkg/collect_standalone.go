// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"bytes"
	"fmt"

	crashdexec "github.com/vmware-tanzu/crash-diagnostics/exec"
)

// collectStandaloneDiags runs command/scripts to collect standalone cluster diagnostics data
func collectStandaloneDiags() error {
	if standaloneArgs.clusterName == "" {
		return fmt.Errorf("standalone cluster: cluster name not set")
	}
	if standaloneArgs.kubeconfig == "" {
		standaloneArgs.kubeconfig = getDefaultKubeconfig()
	}
	if standaloneArgs.contextName == "" {
		standaloneArgs.contextName = getDefaultClusterContext(standaloneArgs.clusterName)
	}

	scriptName := saScriptPath
	argsMap := crashdexec.ArgMap{
		"workdir":                 commonArgs.workDir,
		"outputdir":               commonArgs.outputDir,
		"standalone_cluster_name": standaloneArgs.clusterName,
		"standalone_kubeconfig":   standaloneArgs.kubeconfig,
		"standalone_context":      standaloneArgs.contextName,
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
