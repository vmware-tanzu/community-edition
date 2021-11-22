// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"

	"github.com/vladimirvivien/gexe"

	crashdexec "github.com/vmware-tanzu/crash-diagnostics/exec"
)

// collectWorkloadDiags runs command/scripts to collect workload cluster diagnostics data
func collectWorkloadDiags() error {
	if workloadArgs.clusterName == "" {
		return fmt.Errorf("workload cluster: cluster name not set")
	}
	if workloadArgs.infra == "" {
		log.Print("workload cluster: infra not set, setting to 'docker'")
		workloadArgs.infra = "docker"
	}
	if workloadArgs.namespace == "" {
		log.Print("workload cluster: namespace not set, setting to 'default'")
		workloadArgs.namespace = "default"
	}

	// if mgmt server is set, use mgmt cluster to retrieve workload credentials
	if mgmtArgs.clusterName != "" {
		if workloadArgs.kubeconfig == "" {
			workloadArgs.kubeconfig = filepath.Join(commonArgs.workDir, fmt.Sprintf("%s.kubecfg", workloadArgs.clusterName))
			proc := gexe.RunProc(
				fmt.Sprintf(
					"tanzu cluster kubeconfig get %s --admin --namespace=%s --export-file %s",
					workloadArgs.clusterName,
					workloadArgs.namespace,
					workloadArgs.kubeconfig,
				),
			)
			if proc.Err() != nil {
				return fmt.Errorf("workload cluster: generating kubeconfig file: %s: %s", proc.Result(), proc.Err())
			}
		}
	} else {
		workloadArgs.kubeconfig = getDefaultKubeconfig()
	}

	scriptName := wcScriptPath
	argsMap := crashdexec.ArgMap{
		"workdir":               commonArgs.workDir,
		"outputdir":             commonArgs.outputDir,
		"workload_infra":        workloadArgs.infra,
		"workload_kubeconfig":   workloadArgs.kubeconfig,
		"workload_cluster_name": workloadArgs.clusterName,
		"workload_context":      workloadArgs.contextName,
		"workload_namespace":    workloadArgs.namespace,
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
