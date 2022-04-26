// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/config"
	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/tanzu"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/wizard"
)

const wizardDesc = `
Starts an interactive bootstrapping wizard`

// WizardCmd starts a stopped cluster.
var WizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "(experimental) Start the interactive terminal bootstrapping wizard",
	Long:  wizardDesc,
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: wizardStart,
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func init() {
	WizardCmd.Flags().Bool("tty-disable", false, "Disable log stylization and emojis")
}

func wizardStart(cmd *cobra.Command, args []string) error {
	log := logger.NewLogger(TtySetting(cmd.Flags()), 0)

	p := tea.NewProgram(wizard.InitalModel())
	m, err := p.StartReturningModel()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	// User hit ctrl-c or <ESC>
	if m.(wizard.Model).Err != nil {
		log.Info("Exiting interactive wizard")
		return nil
	}

	configArgs := map[string]interface{}{
		config.ClusterName:           m.(wizard.Model).Clustername,
		config.ControlPlaneNodeCount: m.(wizard.Model).ControlPlaneNodeCount,
		config.WorkerNodeCount:       m.(wizard.Model).WorkerNodeCount,
	}

	clusterConfig, err := config.InitializeConfiguration(configArgs)
	if err != nil {
		log.Errorf("Failed to initialize configuration. Error: %s\n", err.Error())
		return nil
	}

	clusterConfig.LogFile, err = cmd.Parent().PersistentFlags().GetString("log-file")
	if err != nil {
		log.Errorf("Failed to parse log file string. Error %v\n", err)
		os.Exit(tanzu.InvalidConfig)
	}

	tm := tanzu.New(log)
	exitCode, err := tm.Deploy(clusterConfig)
	if err != nil {
		log.Error(err.Error())
		os.Exit(exitCode)
	}

	return nil
}
