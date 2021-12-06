// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/kind/pkg/cluster"
	kindcmd "sigs.k8s.io/kind/pkg/cmd"

	crashdexec "github.com/vmware-tanzu/crash-diagnostics/exec"
)

// collectBootstrapDiags runs the command to collect bootstrap diagnostics data
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
		"workdir":   commonArgs.workDir,
		"outputdir": commonArgs.outputDir,
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

// getTanzuKindClusters returns a slice of tanzu kind cluster names
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
