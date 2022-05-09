// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package cmd contains the CLI-level constructs that call into the tanzu package.
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/config"
	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/tanzu"
)

type createUnmanagedOpts struct {
	skipPreflightChecks       bool
	clusterConfigFile         string
	kubeconfigPath            string
	existingClusterKubeconfig string
	nodeImage                 string
	infrastructureProvider    string
	tkrLocation               string
	additionalRepo            []string
	cni                       string
	podcidr                   string
	servicecidr               string
	portMapping               []string
	numContPlanes             string
	numWorkers                string
	installPackage            []string
}

const createDesc = `
Create an unmanaged Tanzu cluster. This sets up a Kubernetes cluster and installs 
Tanzu packages. Once the environment is bootstrapped, your kubectl context is 
automatically switched, enabling you to begin using the kubectl and tanzu CLIs.

No configuration is required for this command. However, you can setup a config
file by running tanzu unmanaged-cluster configure <cluster name>. This config file can
then be referenced using the -f flag.

When create is called, it respects the following precedence for all configuration: 

1. flags (most respected)
2. environment variables
3. config file
4. defaults (least respected)

Exit codes are provided to enhance the automation of bootstrapping and are defined as follows:

0  - Success.
1  - Configuration is invalid.
2  - Could not create local cluster directories.
3  - Unable to get TKR BOM.
4  - Could not render config.
5  - TKR BOM not parseable.
6  - Could not resolve kapp controller bundle.
7  - Unable to create cluster.
8  - Unable to use existing cluster (if provided).
9  - Could not install kapp controller to cluster.
10 - Could not install core package repo to cluster.
11 - Could not install additional package repo
12 - Could not install CNI package.
13 - Failed to merge kubeconfig and set context
14 - Could not install designated installPackage`

// CreateCmd creates an unmanaged workload cluster.
var CreateCmd = &cobra.Command{
	Use:   "create <cluster name>",
	Short: "Create an unmanaged cluster",
	Long:  createDesc,
	Run:   create,
	Args:  cobra.MaximumNArgs(1),
}

var co = createUnmanagedOpts{}

//nolint:dupl
func init() {
	CreateCmd.Flags().StringVarP(&co.clusterConfigFile, "config", "f", "", "A config file describing how to create the Tanzu environment")
	CreateCmd.Flags().StringVar(&co.kubeconfigPath, "kubeconfig-path", "", "File path to where the kubeconfig will be persisted. Defaults to global user kubeconfig")
	CreateCmd.Flags().StringVarP(&co.existingClusterKubeconfig, "existing-cluster-kubeconfig", "e", "", "Use an existing kubeconfig to tanzu-ify a cluster")
	CreateCmd.Flags().StringVar(&co.nodeImage, "node-image", "", "The host OS image to use for kubernetes nodes")
	CreateCmd.Flags().StringVar(&co.infrastructureProvider, "provider", "", "The infrastructure provider for cluster creation; default is kind")
	CreateCmd.Flags().StringVarP(&co.tkrLocation, "tkr", "t", "", "The URL to the image or path to local file containing a Tanzu Kubernetes release")
	CreateCmd.Flags().StringSliceVar(&co.additionalRepo, "additional-repo", []string{}, "Addresses for additional package repositories to install")
	CreateCmd.Flags().StringVarP(&co.cni, "cni", "c", "", "The CNI to deploy; default is calico")
	CreateCmd.Flags().StringVar(&co.podcidr, "pod-cidr", "", "The CIDR for Pod IP allocation; default is 10.244.0.0/16")
	CreateCmd.Flags().StringVar(&co.servicecidr, "service-cidr", "", "The CIDR for Service IP allocation; default is 10.96.0.0/16")
	CreateCmd.Flags().StringSliceVarP(&co.portMapping, "port-map", "p", []string{}, "Ports to map between container node and the host (format: '127.0.0.1:80:80/tcp', '80:80/tcp', '80:80', or just '80')")
	CreateCmd.Flags().Bool("tty-disable", false, "Disable log stylization and emojis")
	CreateCmd.Flags().BoolVar(&co.skipPreflightChecks, "skip-preflight", false, "Skip the preflight checks; default is false")
	CreateCmd.Flags().StringVar(&co.numContPlanes, "control-plane-node-count", "", "The number of control plane nodes to deploy; default is 1")
	CreateCmd.Flags().StringVar(&co.numWorkers, "worker-node-count", "", "The number of worker nodes to deploy; default is 0")
	CreateCmd.Flags().StringSliceVar(&co.installPackage, "install-package", []string{}, "(experimental) A package to install on bootstrapping. May be specified multiple times. install-package mappings supported - package-name:package-version:package-config-file. package-name should be the fully qualified package name or a prefix to a package name found in an installed package repository. package-version is optional and resolves to the latest semantic versioned package if not specified or latest is entered. package-config-file is optional and should be the path to a values yaml file in order to configure the package.")
}

func create(cmd *cobra.Command, args []string) {
	var (
		clusterName string
		err         error
	)

	// Set the cluster name if it was provided, otherwise read from config file
	if len(args) == 1 {
		clusterName = args[0]
	}

	// initial logger, needed for logging if something goes wrong
	log := logger.NewLogger(TtySetting(cmd.Flags()), LoggingVerbosity(cmd))

	// Attempt to read cluster name from provided kubeconfig
	if co.existingClusterKubeconfig != "" {
		clusterName, err = tanzu.ReadClusterContextFromKubeconfig(co.existingClusterKubeconfig)
		if err != nil {
			log.Error(err.Error())
			os.Exit(tanzu.ErrExistingCluster)
		}
	}

	portMaps, err := config.ParsePortMappings(co.portMapping)
	if err != nil {
		log.Error(err.Error())
		os.Exit(tanzu.ErrRenderingConfig)
	}

	installPackages, err := config.ParseInstallPackageMappings(co.installPackage)
	if err != nil {
		log.Error(err.Error())
		os.Exit(tanzu.ErrRenderingConfig)
	}

	// Get the log file from the global parent flag
	logFile, err := cmd.Parent().PersistentFlags().GetString("log-file")
	if err != nil {
		log.Errorf("Failed to parse log file string. Error %v\n", err)
		os.Exit(tanzu.InvalidConfig)
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
	clusterConfig, err := config.InitializeConfiguration(configArgs)
	if err != nil {
		log.Errorf("Failed to initialize configuration. Error %v\n", err)
		os.Exit(tanzu.InvalidConfig)
	}

	tm := tanzu.New(log)
	exitCode, err := tm.Deploy(clusterConfig)
	if err != nil {
		log.Error(err.Error())
		os.Exit(exitCode)
	}
}
