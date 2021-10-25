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

func init() {
}

func configure(cmd *cobra.Command, args []string) error {
	var clusterName string

	// validate a cluster name was passed when not using the kickstart UI
	if len(args) < 1 && !co.ui {
		return fmt.Errorf("cluster name not specified")
	} else if len(args) == 1 {
		clusterName = args[0]
	}

	log := logger.NewLogger(true, 0)

	var rawConfig bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&rawConfig)
	yamlEncoder.SetIndent(yamlIndent)

	lcConfig, err := config.InitializeConfiguration(map[string]string{"clustername": clusterName})
	if err != nil {
		log.Errorf("Failed to initialize configuration. Error: %s\n", err.Error())
	}
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
