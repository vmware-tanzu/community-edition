// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/vladimirvivien/gexe"
)

var (
	scriptFS       embed.FS
	libScriptPath  = "scripts/lib.star"
	bootScriptPath = "scripts/bootstrap_cluster.star"
	mgmtScriptPath = "scripts/management_cluster.star"
	wcScriptPath   = "scripts/workload_cluster.star"
)

var (
	commonArgs = &collectCommonArgs{
		workDir:   getDefaultWorkdir(),
		outputDir: getDefaultOutputDir(),
	}

	bootstrapArgs = &collectBootsrapArgs{
		skip: false,
	}

	mgmtArgs = &collectMgmtArgs{
		skip: false,
	}

	workloadArgs = &collectWorkloadArgs{
		infra:     "docker",
		namespace: "default",
	}
)

func CollectCmd(fs embed.FS) *cobra.Command {
	scriptFS = fs
	mgmtSvr, _ := getDefaultManagementServer()

	if mgmtSvr != nil {
		mgmtArgs.kubeconfig = mgmtSvr.kubeconfig
		mgmtArgs.clusterName = mgmtSvr.clusterName
		mgmtArgs.contextName = mgmtSvr.kubecontext
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
	cmd.Flags().StringVar(&bootstrapArgs.clusterName, "bootstrap-cluster-name", bootstrapArgs.clusterName, "A specific bootstrap cluster to diagnose")

	// management
	cmd.Flags().BoolVar(&mgmtArgs.skip, "management-cluster-skip", mgmtArgs.skip, "If true, skips management cluster diagnostics")
	cmd.Flags().StringVar(&mgmtArgs.kubeconfig, "management-cluster-kubeconfig", mgmtArgs.kubeconfig, "The management cluster config file (required)")
	cmd.Flags().StringVar(&mgmtArgs.clusterName, "management-cluster-name", mgmtArgs.clusterName, "The name of the management cluster (required)")
	cmd.Flags().StringVar(&mgmtArgs.contextName, "management-cluster-context", mgmtArgs.contextName, "The context name of the management cluster")

	// workload
	cmd.Flags().StringVar(&workloadArgs.infra, "workload-cluster-infra", workloadArgs.infra, "Overrides the infrastructure type for the managed cluster (i.e. aws, azure, vsphere, etc)")
	cmd.Flags().StringVar(&workloadArgs.kubeconfig, "workload-cluster-kubeconfig", workloadArgs.kubeconfig, "The workload cluster config file")
	cmd.Flags().StringVar(&workloadArgs.clusterName, "workload-cluster-name", workloadArgs.clusterName, "The name of the managed cluster for which to collect diagnostics (required)")
	cmd.Flags().StringVar(&workloadArgs.contextName, "workload-cluster-context", workloadArgs.contextName, "The context name of the workload cluster")
	cmd.Flags().StringVar(&workloadArgs.namespace, "workload-cluster-namespace", workloadArgs.namespace, "The namespace where managed workload resources are stored")

	cmd.RunE = collectFunc
	return cmd
}

func collectFunc(_ *cobra.Command, _ []string) error {
	if gexe.Prog().Avail("tanzu") == "" {
		return fmt.Errorf("program not found: tanzu")
	}

	if gexe.Prog().Avail("docker") == "" {
		return fmt.Errorf("progra not found: docker")
	}

	if commonArgs.workDir == "" {
		log.Printf("workdir empty: setting it to %s", getDefaultWorkdir())
		commonArgs.workDir = getDefaultWorkdir()
	}

	if commonArgs.outputDir == "" {
		log.Printf("output dir empty: setting it to %s", getDefaultOutputDir())
		commonArgs.outputDir = getDefaultOutputDir()
	}

	if err := os.MkdirAll(commonArgs.workDir, 0744); err != nil && !os.IsExist(err) {
		return fmt.Errorf("workdir setup: %w", err)
	}
	defer os.RemoveAll(commonArgs.workDir)

	if err := collectBoostrapDiags(); err != nil {
		log.Printf("Warn: skipping cluster diagnostics: %s", err)
	}

	if err := collectManagementDiags(); err != nil {
		log.Printf("Warn: skipping cluster diagnostics: %s", err)
	}

	if err := collectWorkloadDiags(); err != nil {
		log.Printf("Warn: skipping cluster diagnostics: %s", err)
	}

	return nil
}
