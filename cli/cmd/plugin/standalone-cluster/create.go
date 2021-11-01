// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/config"
	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/tanzu"
)

type createStandaloneOpts struct {
	clusterConfigFile      string
	infrastructureProvider string
	tkrLocation            string
	cni                    string
	podcidr                string
	servicecidr            string
	tty                    bool
}

// CreateCmd creates a standalone workload cluster.
var CreateCmd = &cobra.Command{
	Use:   "create <cluster name> -f <configuration location>",
	Short: "create a standalone tanzu cluster",
	RunE:  create,
}

var co = createStandaloneOpts{}

//nolint:dupl
func init() {
	CreateCmd.Flags().StringVarP(&co.clusterConfigFile, "config", "f", "", "Configuration file for standalone cluster creation")
	CreateCmd.Flags().StringVarP(&co.infrastructureProvider, "provider", "p", "", "The infrastructure provider to use for cluster creation. Default is 'kind'")
	CreateCmd.Flags().StringVarP(&co.tkrLocation, "tkr", "t", "", "The Tanzu Kubernetes Release location.")
	CreateCmd.Flags().StringVarP(&co.cni, "cni", "c", "", "The CNI to deploy. Default is 'antrea'")
	CreateCmd.Flags().StringVar(&co.podcidr, "pod-cidr", "", "The CIDR to use for Pod IP addresses. Default and format is '10.244.0.0/16'")
	CreateCmd.Flags().StringVar(&co.servicecidr, "service-cidr", "", "The CIDR to use for Service IP addresses. Default and format is '10.96.0.0/16'")
	CreateCmd.Flags().BoolVar(&co.tty, "tty", true, "Specify whether terminal is tty. Set to false to disable styled output.")
}

func create(cmd *cobra.Command, args []string) error {
	var clusterName string

	// validate a cluster name was passed
	if len(args) < 1 {
		return fmt.Errorf("cluster name not specified")
	} else if len(args) == 1 {
		clusterName = args[0]
	}

	// initial logger, needed for logging if something goes wrong
	log := logger.NewLogger(TtySetting(cmd.Flags()), 0)

	// Determine our configuration to use
	configArgs := map[string]string{
		config.ClusterConfigFile: co.clusterConfigFile,
		config.ClusterName:       clusterName,
		config.Provider:          co.infrastructureProvider,
		config.TKRLocation:       co.tkrLocation,
		config.Cni:               co.cni,
		config.PodCIDR:           co.podcidr,
		config.ServiceCIDR:       co.servicecidr,
	}
	clusterConfig, err := config.InitializeConfiguration(configArgs)
	if err != nil {
		log.Errorf("Failed to initialize configuration. Error %s\n", clusterConfig)
		return nil
	}

	tm := tanzu.New(log)
	err = tm.Deploy(clusterConfig)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	return nil
}
