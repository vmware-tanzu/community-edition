// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/config"
	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/log"
)

const configureDesc = `
Generate a configuration file that can be used when running:
tanzu unmanaged-cluster create -f <config-file-name>.yaml. Configure generates
a config file injected with default values. When flags are specified
(e.g. --cni) the flag value is respected in the overridden config.
`

// ConfigureCmd creates an unmanaged workload cluster.
var ConfigureCmd = &cobra.Command{
	Use:     "configure <cluster name>",
	Aliases: []string{"config", "conf"},
	Short:   "Generate a config file to be used in cluster creation",
	Long:    configureDesc,
	RunE:    configure,
	Args:    checkSingleClusterArg,
}

//nolint:dupl
func init() {
	ConfigureCmd.Flags().StringVarP(&co.clusterConfigFile, "config", "f", "", "Configuration file for unmanaged cluster creation")
	ConfigureCmd.Flags().StringVar(&co.kubeconfigPath, "kubeconfig-path", "", "File path to where the kubeconfig will be persisted. Defaults to global user kubeconfig")
	ConfigureCmd.Flags().StringVarP(&co.existingClusterKubeconfig, "existing-cluster-kubeconfig", "e", "", "Use an existing kubeconfig to tanzu-ify a cluster")
	ConfigureCmd.Flags().StringVar(&co.nodeImage, "node-image", "", "The host OS image to use for kubernetes nodes")
	ConfigureCmd.Flags().StringVar(&co.infrastructureProvider, "provider", "", "The infrastructure provider for cluster creation; default is kind")
	ConfigureCmd.Flags().StringVarP(&co.tkrLocation, "tkr", "t", "", "The URL to the image or path to local file containing a Tanzu Kubernetes release")
	ConfigureCmd.Flags().StringSliceVar(&co.additionalRepo, "additional-repo", []string{}, "Addresses for additional package repositories to install")
	ConfigureCmd.Flags().StringVarP(&co.cni, "cni", "c", "", "The CNI to deploy; default is calico")
	ConfigureCmd.Flags().StringVar(&co.podcidr, "pod-cidr", "", "The CIDR for Pod IP allocation; default is 10.244.0.0/16")
	ConfigureCmd.Flags().StringVar(&co.servicecidr, "service-cidr", "", "The CIDR for Service IP allocation; default is 10.96.0.0/16")
	ConfigureCmd.Flags().StringSliceVarP(&co.portMapping, "port-map", "p", []string{}, "Ports to map between container node and the host (format: '127.0.0.1:80:80/tcp', '80:80/tcp', '80:80', or just '80')")
	ConfigureCmd.Flags().Bool("tty-disable", false, "Disable log stylization and emojis")
	ConfigureCmd.Flags().BoolVar(&co.skipPreflightChecks, "skip-preflight", false, "Skip the preflight checks; default is false")
	ConfigureCmd.Flags().StringVar(&co.numContPlanes, "control-plane-node-count", "", "The number of control plane nodes to deploy; default is 1")
	ConfigureCmd.Flags().StringVar(&co.numWorkers, "worker-node-count", "", "The number of worker nodes to deploy; default is 0")
	ConfigureCmd.Flags().StringSliceVar(&co.installPackage, "install-package", []string{}, "(experimental) A package to install on bootstrapping. May be specified multiple times. install-package mappings supported - package-name:package-version:package-config-file. package-name should be the fully qualified package name or a prefix to a package name found in an installed package repository. package-version is optional and resolves to the latest semantic versioned package if not specified or latest is entered. package-config-file is optional and should be the path to a values yaml file in order to configure the package.")
}

func configure(cmd *cobra.Command, args []string) error {
	// args have already been checked by ConfigureCmd.Args()
	clusterName := args[0]

	log := logger.NewLogger(TtySetting(cmd.Flags()), LoggingVerbosity(cmd))

	portMaps, err := config.ParsePortMappings(co.portMapping)
	if err != nil {
		log.Error(err.Error())
	}

	installPackages, err := config.ParseInstallPackageMappings(co.installPackage)
	if err != nil {
		log.Error(err.Error())
	}

	// Get the log file from the global parent flag
	logFile, err := cmd.Parent().PersistentFlags().GetString("log-file")
	if err != nil {
		log.Errorf("Failed to parse log file string. Error %v\n", err)
	}

	// Determine our configuration to use
	//nolint:dupl
	configArgs := map[string]interface{}{
		config.ClusterConfigFile:         co.clusterConfigFile,
		config.ClusterName:               clusterName,
		config.KubeconfigPath:            co.kubeconfigPath,
		config.ExistingClusterKubeconfig: co.existingClusterKubeconfig,
		config.NodeImage:                 co.nodeImage,
		config.Provider:                  co.infrastructureProvider,
		config.TKRLocation:               co.tkrLocation,
		config.Cni:                       co.cni,
		config.PodCIDR:                   co.podcidr,
		config.ServiceCIDR:               co.servicecidr,
		config.ControlPlaneNodeCount:     co.numContPlanes,
		config.WorkerNodeCount:           co.numWorkers,
		config.AdditionalPackageRepos:    co.additionalRepo,
		config.PortsToForward:            portMaps,
		config.SkipPreflightChecks:       co.skipPreflightChecks,
		config.InstallPackages:           installPackages,
		config.LogFile:                   logFile,
	}

	scConfig, err := config.InitializeConfiguration(configArgs)
	if err != nil {
		log.Errorf("Failed to initialize configuration. Error: %s\n", err.Error())
		return nil
	}
	fileName := fmt.Sprintf("%s.yaml", clusterName)

	err = config.RenderConfigToFile(fileName, scConfig)
	if err != nil {
		log.Errorf("Failed to write configuration file: %s\n", err.Error())
		return nil
	}
	log.Infof("Wrote configuration file to: %s\n", fileName)

	return nil
}
