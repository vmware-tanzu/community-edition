// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	crashdexec "github.com/vmware-tanzu/crash-diagnostics/exec"
)

var (
	scriptFS       embed.FS
	libScriptPath  = "scripts/lib.star"
	bootScriptPath = "scripts/bootstrap_cluster.star"
	mgmtScriptPath = "scripts/management_cluster.star"
	wcScriptPath   = "scripts/workload_cluster.star"
)

var (
	commonArgs = collectCommonArgs{
		workDir:   getDefaultWorkdir(),
		outputDir: getDefaultOutputDir(),
	}

	bootstrapArgs = collectBootsrapArgs{
		skip: false,
	}

	mgmtArgs = collectMgmtArgs{
		skip: false,
	}

	workloadArgs = collectWorkloadArgs{
		standalone: false,
		infra:      "docker",
	}
)

func CollectCmd(fs embed.FS) *cobra.Command {
	scriptFS = fs
	mgmtSvr, _ := getDefaultManagementServer()

	if mgmtSvr != nil {
		mgmtArgs.kubeconfig = mgmtSvr.kubeconfig
		mgmtArgs.clusterName = mgmtSvr.name
		mgmtArgs.contextName = mgmtSvr.kubecontext
	} else {
		workloadArgs.standalone = true
		workloadArgs.kubeconfig = getDefaultKubeconfig()
	}

	cmd := &cobra.Command{
		Use:   "collect",
		Short: "Collect cluster diagnostics for the specified cluster",
		Long:  `Collect cluster diagnostics for the specified cluster`,
	}

	// common args
	cmd.Flags().StringVar(&commonArgs.workDir, "work-dir", commonArgs.workDir, "Working directory for collected data")
	cmd.Flags().StringVar(&commonArgs.outputDir, "output-dir", commonArgs.outputDir, "Output directory for collected bundle")

	// bootstrap args
	cmd.Flags().BoolVar(&bootstrapArgs.skip, "bootstrap-cluster-skip", bootstrapArgs.skip, "If true, skips bootstrap cluster diagnostics")
	cmd.Flags().StringVar(&bootstrapArgs.clusterName, "bootstrap-cluster-name", bootstrapArgs.clusterName, "A specific bootstrap cluster name to diagnose")

	// management
	cmd.Flags().BoolVar(&mgmtArgs.skip, "management-cluster-skip", mgmtArgs.skip, "If true, skips management cluster diagnostics")
	cmd.Flags().StringVar(&mgmtArgs.kubeconfig, "management-cluster-kubeconfig", mgmtArgs.kubeconfig, "The management cluster config file (required)")
	cmd.Flags().StringVar(&mgmtArgs.clusterName, "management-cluster-name", mgmtArgs.clusterName, "The name of the management cluster (required)")
	cmd.Flags().StringVar(&mgmtArgs.contextName, "management-cluster-context", mgmtArgs.contextName, "The context name of the management cluster (required)")

	// workload
	cmd.Flags().BoolVar(&workloadArgs.standalone, "workload-cluster-standalone", workloadArgs.standalone, "If true, workload cluster is treated as standalone")
	cmd.Flags().StringVar(&workloadArgs.infra, "workload-cluster-infra", workloadArgs.infra, "Overrides the infrastructure type for the managed cluster (i.e. aws, azure, vsphere, etc)")
	cmd.Flags().StringVar(&workloadArgs.clusterName, "workload-cluster-name", workloadArgs.clusterName, "The name of the managed cluster for which to collect diagnostics (required)")
	cmd.Flags().StringVar(&workloadArgs.namespace, "workload-cluster-namespace", workloadArgs.namespace, "The namespace where managed workload resources are stored (required)")

	cmd.RunE = collectFunc
	return cmd
}

func collectFunc(cmd *cobra.Command, args []string) error {
	defer os.RemoveAll(commonArgs.workDir)

	if err := collectBoostrapDiags(); err != nil {
		log.Printf("Error: skipping bootstrap cluster diagnostics: %s", err)
	}

	if err := collectManagementDiags(); err != nil {
		log.Printf("Error: skipping management cluster diagnostics: %s", err)
	}

	if err := collectWorkloadDiags(); err != nil {
		log.Printf("Error: skipping workload cluster diagnostics: %s", err)
	}

	return nil
}

func collectBoostrapDiags() error {
	if bootstrapArgs.skip {
		log.Println("bootstrap cluster: skip=true: diagnostics will not be collected")
		return nil
	}

	log.Println("Collecting bootstrap cluster diagnostics")

	libScript := libScriptPath
	libData, err := scriptFS.ReadFile(libScript)
	if err != nil {
		return err
	}

	scriptName := bootScriptPath
	scriptData, err := scriptFS.ReadFile(scriptName)
	if err != nil {
		return err
	}

	argsMap := crashdexec.ArgMap{
		"workdir":                commonArgs.workDir,
		"infra":                  "docker",
		"bootstrap_cluster_name": bootstrapArgs.clusterName,
	}

	return crashdexec.ExecuteWithModules(
		scriptName,
		bytes.NewReader(scriptData),
		argsMap,
		crashdexec.StarlarkModule{Name: libScript, Source: bytes.NewReader(libData)},
	)
}

func collectManagementDiags() error {
	if mgmtArgs.skip {
		log.Println("management cluster: skip=true: diagnostics will not be collected")
		return nil
	}

	if mgmtArgs.clusterName == "" {
		return fmt.Errorf("management cluster: name not set")
	}
	if mgmtArgs.kubeconfig == "" {
		return fmt.Errorf("management cluster: kubeconfig is required")
	}
	if mgmtArgs.contextName == "" {
		mgmtArgs.contextName = getDefaultClusterContext(mgmtArgs.clusterName)
	}

	argsMap := crashdexec.ArgMap{
		"workdir":                 commonArgs.workDir,
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

func collectWorkloadDiags() error {
	if workloadArgs.clusterName == "" {
		return fmt.Errorf("workload cluster: name not set")
	}

	argsMap := crashdexec.ArgMap{
		"workdir":                 commonArgs.workDir,
		"management_cluster_name": mgmtArgs.clusterName,
		"management_kubeconfig":   mgmtArgs.kubeconfig,

		"workload_infra":        workloadArgs.infra,
		"workload_kubeconfig":   workloadArgs.kubeconfig,
		"workload_cluster_name": workloadArgs.clusterName,
		"workload_namespace":    workloadArgs.namespace,
	}

	scriptName := wcScriptPath
	if workloadArgs.standalone {
		argsMap["workload_kubeconfig"] = getDefaultKubeconfig()
		argsMap["workload_context"] = getDefaultClusterContext(workloadArgs.clusterName)
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
