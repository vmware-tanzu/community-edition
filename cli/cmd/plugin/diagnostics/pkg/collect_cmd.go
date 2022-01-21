// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cluster"
	kindcmd "sigs.k8s.io/kind/pkg/cmd"

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
	cmd.Flags().StringVar(&workloadArgs.infra, "workload-cluster-infra", workloadArgs.infra, "Overrides the infrastructure type for the managed cluster (i.e. aws, azure, vsphere, etc)")
	cmd.Flags().StringVar(&workloadArgs.clusterName, "workload-cluster-name", workloadArgs.clusterName, "The name of the managed cluster for which to collect diagnostics (required)")
	cmd.Flags().StringVar(&workloadArgs.namespace, "workload-cluster-namespace", workloadArgs.namespace, "The namespace where managed workload resources are stored (required)")

	cmd.RunE = collectFunc
	return cmd
}

func collectFunc(_ *cobra.Command, _ []string) error {
	defer os.RemoveAll(commonArgs.workDir)

	if err := collectBoostrapDiags(); err != nil {
		log.Printf("Warn: skipping bootstrap cluster diagnostics: %s", err)
	}

	if err := collectManagementDiags(); err != nil {
		log.Printf("Warn: skipping management cluster diagnostics: %s", err)
	}

	if err := collectWorkloadDiags(); err != nil {
		log.Printf("Warn: skipping workload cluster diagnostics: %s", err)
	}

	return nil
}

func collectBoostrapDiags() error {
	if bootstrapArgs.skip {
		log.Println("bootstrap cluster: skip=true: diagnostics will not be collected")
		return nil
	}

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

	// setup workdir
	if err := os.MkdirAll(commonArgs.workDir, 0744); err != nil && !os.IsExist(err) {
		return fmt.Errorf("bootstrap cluster: %w", err)
	}

	// loop through and collect diags from each cluster
	prov := cluster.NewProvider(cluster.ProviderWithLogger(kindcmd.NewLogger()))
	clusterList, err := prov.List()
	if err != nil {
		return err
	}

	clusters := getTanzuKindClusters(clusterList, bootstrapArgs.clusterName)
	if len(clusters) == 0 {
		return fmt.Errorf("bootstrap cluster: no kind cluster found")
	}

	argsMap := crashdexec.ArgMap{
		"workdir":                commonArgs.workDir,
		"infra":                  "docker",
		"bootstrap_cluster_name": bootstrapArgs.clusterName,
		"outputdir":              commonArgs.outputDir,
	}

	for _, cluster := range clusters {
		cfg, err := prov.KubeConfig(cluster, false)
		if err != nil {
			log.Printf("Warn: failed to get cluster kubeconfig, K8s object not collected: %s: %s", cluster, err)
		}

		path := filepath.Join(commonArgs.workDir, fmt.Sprintf("%s.config", cluster))
		if err := os.WriteFile(path, []byte(cfg), 0644); err != nil {
			return fmt.Errorf("bootstrap diagnostics kubeconfig: %w", err)
		}

		defer func(p string) {
			if err := os.RemoveAll(p); err != nil {
				log.Printf("Warn: bootstrap cluster: failed to remove kubeconfig file: %s", err)
			}
		}(path)

		log.Printf("Collecting bootstrap diagnostics: cluster: %s", cluster)

		argsMap["bootstrap_cluster_name"] = cluster
		argsMap["bootstrap_kubeconfig"] = path

		err = crashdexec.ExecuteWithModules(
			scriptName,
			bytes.NewReader(scriptData),
			argsMap,
			crashdexec.StarlarkModule{Name: libScript, Source: bytes.NewReader(libData)},
		)
		if err != nil {
			log.Printf("Warn: bootstrap script failed, skipping: cluster%s: %s: ", cluster, err)
			continue
		}
	}
	return nil
}

func getTanzuKindClusters(clusters []string, clusterName string) []string {
	var result []string
	for _, cluster := range clusters {
		if clusterName != "" && clusterName == cluster {
			return []string{cluster}
		}

		if strings.HasPrefix(cluster, "tkg-kind") {
			result = append(result, cluster)
		}
	}

	return result
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
		"outputdir":               commonArgs.outputDir,
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

		"outputdir": commonArgs.outputDir,
	}

	scriptName := wcScriptPath
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
