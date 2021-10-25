// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/config"
	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/log"
)

const yamlIndent = 2

// ConfigureCmd creates a standalone workload cluster.
var ConfigureCmd = &cobra.Command{
	Use:   "configure <cluster name>",
	Short: "create a configuration file for a future cluster",
	RunE:  configure,
}

//nolint:dupl
func init() {
	ConfigureCmd.Flags().StringVarP(&co.clusterConfigFile, "config", "f", "", "Configuration file for local cluster creation")
	ConfigureCmd.Flags().StringVarP(&co.infrastructureProvider, "provider", "p", "", "The infrastructure provider to use for cluster creation. Default is 'kind'")
	ConfigureCmd.Flags().StringVarP(&co.tkrLocation, "tkr", "t", "", "The Tanzu Kubernetes Release location.")
	ConfigureCmd.Flags().StringVarP(&co.cni, "cni", "c", "", "The CNI to deploy. Default is 'antrea'")
	ConfigureCmd.Flags().StringVar(&co.podcidr, "pod-cidr", "", "The CIDR to use for Pod IP addresses. Default and format is '10.244.0.0/16'")
	ConfigureCmd.Flags().StringVar(&co.servicecidr, "service-cidr", "", "The CIDR to use for Service IP addresses. Default and format is '10.96.0.0/16'")
	ConfigureCmd.Flags().BoolVar(&co.tty, "tty", true, "Specify whether terminal is tty. Set to false to disable styled output; default: true")
}

func configure(cmd *cobra.Command, args []string) error {
	var clusterName string

	// validate a cluster name was passed
	if len(args) < 1 {
		return fmt.Errorf("cluster name not specified")
	} else if len(args) == 1 {
		clusterName = args[0]
	}

	log := logger.NewLogger(co.tty, 0)

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
	lcConfig, err := config.InitializeConfiguration(configArgs)
	if err != nil {
		log.Errorf("Failed to initialize configuration. Error: %s\n", err.Error())
	}

	var rawConfig bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&rawConfig)
	yamlEncoder.SetIndent(yamlIndent)
	err = yamlEncoder.Encode(*lcConfig)
	if err != nil {
		log.Errorf("Failed to render rawConfig file. Error: %s\n", err.Error())
		return nil
	}

	fileName := fmt.Sprintf("%s.yaml", clusterName)
	err = os.WriteFile(fileName, rawConfig.Bytes(), 0644)
	if err != nil {
		log.Errorf("Failed to write rawConfig file. Error: %s\n", err.Error())
		return nil
	}
	log.Infof("Wrote configuration file to: %s\n", fileName)

	return nil
}
