// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/tanzu"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// ConfigureCmd creates a standalone workload cluster.
var ConfigureCmd = &cobra.Command{
	Use:   "configure <cluster name>",
	Short: "create a configuration file for a future cluster",
	RunE:  configure,
}

func init() {
	//ConfigureCmd.Flags().BoolVar(&createOpts.tty, "tty", true, "Specify whether terminal is tty;\\nSet to false to disable styled ouput; default: true")
}

func configure(cmd *cobra.Command, args []string) error {
	var clusterName string

	// validate a cluster name was passed when not using the kickstart UI
	if len(args) < 1 && !iso.ui {
		return fmt.Errorf("Cluster name not specified.")
	} else if len(args) == 1 {
		clusterName = args[0]
	}

	log := logger.NewLogger(true, 0)

	var config bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&config)
	yamlEncoder.SetIndent(2)

	err := yamlEncoder.Encode(tanzu.LocalClusterConfig{
		ClusterName: clusterName,
		Provider:    tanzu.DefaultProvider,
		CNI:         tanzu.DefaultCni,
		PodCidr:     tanzu.DefaultPodCidr,
		ServiceCidr: tanzu.DefaultServiceCidr,
		TkrLocation: tanzu.DefaultTkrLocation,
	})

	if err != nil {
		log.Errorf("failed to render config file. Error: %s", err.Error())
		return nil
	}

	fileName := fmt.Sprintf("%s.yaml", clusterName)
	err = ioutil.WriteFile(fileName, config.Bytes(), 0644)
	if err != nil {
		log.Errorf("failed to write config file. Error: %s", err.Error())
		return nil
	}
	log.Infof("Wrote configuration file to: %s", fileName)

	return nil
}
