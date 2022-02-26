// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package cmd contains the CLI-level constructs that call into the tanzu package.
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/config"
	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/tanzu"
)

type createUnmanagedOpts struct {
	skipPreflightChecks       bool
	clusterConfigFile         string
	existingClusterKubeconfig string
	infrastructureProvider    string
	tkrLocation               string
	cni                       string
	podcidr                   string
	servicecidr               string
	portMapping               []string
	numContPlanes             string
	numWorkers                string
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
4. defaults (least respected)`

// CreateCmd creates an unmanaged workload cluster.
var CreateCmd = &cobra.Command{
	Use:   "create <cluster name>",
	Short: "Create an unmanaged Tanzu cluster",
	Long:  createDesc,
	RunE:  create,
	Args:  cobra.MaximumNArgs(1),
}

var co = createUnmanagedOpts{}

func init() {
	CreateCmd.Flags().StringVarP(&co.clusterConfigFile, "config", "f", "", "A config file describing how to create the Tanzu environment")
	CreateCmd.Flags().StringVarP(&co.existingClusterKubeconfig, "existing-cluster-kubeconfig", "e", "", "Use an existing kubeconfig to tanzu-ify a cluster")
	CreateCmd.Flags().StringVar(&co.infrastructureProvider, "provider", "", "The infrastructure provider for cluster creation; default is kind")
	CreateCmd.Flags().StringVarP(&co.tkrLocation, "tkr", "t", "", "The URL to the image containing a Tanzu Kubernetes release")
	CreateCmd.Flags().StringVarP(&co.cni, "cni", "c", "", "The CNI to deploy; default is antrea")
	CreateCmd.Flags().StringVar(&co.podcidr, "pod-cidr", "", "The CIDR for Pod IP allocation; default is 10.244.0.0/16")
	CreateCmd.Flags().StringVar(&co.servicecidr, "service-cidr", "", "The CIDR for Service IP allocation; default is 10.96.0.0/16")
	CreateCmd.Flags().StringSliceVarP(&co.portMapping, "port-map", "p", []string{}, "Ports to map between container node and the host (format: '80:80/tcp' or just '80')")
	CreateCmd.Flags().Bool("tty-disable", false, "Disable log stylization and emojis")
	CreateCmd.Flags().BoolVar(&co.skipPreflightChecks, "skip-preflight", false, "Skip the preflight checks; default is false")
	CreateCmd.Flags().StringVar(&co.numContPlanes, "control-plane-node-count", "", "The number of control plane nodes to deploy; default is 1")
	CreateCmd.Flags().StringVar(&co.numWorkers, "worker-node-count", "", "The number of worker nodes to deploy; default is 0")
}

func create(cmd *cobra.Command, args []string) error {
	var clusterName string

	// Set the cluster name if it was provided, otherwise read from config file
	if len(args) == 1 {
		clusterName = args[0]
	}

	// initial logger, needed for logging if something goes wrong
	log := logger.NewLogger(TtySetting(cmd.Flags()), 0)

	// Determine our configuration to use
	configArgs := map[string]string{
		config.ClusterConfigFile:         co.clusterConfigFile,
		config.ExistingClusterKubeconfig: co.existingClusterKubeconfig,
		config.ClusterName:               clusterName,
		config.Provider:                  co.infrastructureProvider,
		config.TKRLocation:               co.tkrLocation,
		config.Cni:                       co.cni,
		config.PodCIDR:                   co.podcidr,
		config.ServiceCIDR:               co.servicecidr,
		config.ControlPlaneNodeCount:     co.numContPlanes,
		config.WorkerNodeCount:           co.numWorkers,
	}
	clusterConfig, err := config.InitializeConfiguration(configArgs)
	if err != nil {
		log.Errorf("Failed to initialize configuration. Error %v\n", err)
		return nil
	}
	clusterConfig.SkipPreflightChecks = co.skipPreflightChecks

	// TODO(stmcginnis): For now, we are only supporting port maps from command
	// line arguments. At some point we need to add env variable and config file
	// support.
	for i := range co.portMapping {
		mapping, err := config.ParsePortMap(co.portMapping[i])
		if err != nil {
			log.Warn(err.Error())
			continue
		}
		clusterConfig.PortsToForward = append(clusterConfig.PortsToForward, mapping)
	}

	tm := tanzu.New(log)
	err = tm.Deploy(clusterConfig)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	return nil
}
