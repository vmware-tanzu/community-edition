// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"bytes"
	"fmt"
	"log"

	crashdexec "github.com/vmware-tanzu/crash-diagnostics/exec"
)

// collectManagementDiags runs command/script to collect management cluster diagnostics data
func collectManagementDiags() error {
	if mgmtArgs.skip {
		log.Println("management cluster: skip=true: diagnostics will not be collected")
		return nil
	}
	if mgmtArgs.clusterName == "" {
		return fmt.Errorf("management cluster: clusterName not set")
	}
	if mgmtArgs.kubeconfig == "" {
		return fmt.Errorf("management cluster: kubeconfig is required")
	}
	if mgmtArgs.contextName == "" {
		mgmtArgs.contextName = getDefaultClusterContext(mgmtArgs.clusterName)
	}

	argsMap := crashdexec.ArgMap{
		"workdir":                 commonArgs.workDir,
		"outputdir":               commonArgs.outputDir,
		"management_cluster_name": mgmtArgs.clusterName,
		"management_kubeconfig":   mgmtArgs.kubeconfig,
		"management_context":      mgmtArgs.contextName,
	}

	libScript := libScriptPath
	libData, err := scriptFS.ReadFile(libScript)
	if err != nil {
		return err
	}

	scriptName := mgmtScriptPath
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
