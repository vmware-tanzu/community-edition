// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package plugin provides commonly used functionality for implementing Tanzu plugins.
//nolint:revive
package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Date is the date the binary was built.
	// Set by go build -ldflags "-X" flag
	Date string

	// SHA is the git commit SHA the binary was built with.
	// Set by go build -ldflags "-X" flag
	SHA string

	// Version is the version the binary was built with.
	// Set by go build -ldflags "-X" flag
	Version = "dev"
)

// Plugin is a Tanzu CLI plugin.
type Plugin struct {
	Cmd *cobra.Command
}

// CmdGroup is a group of CLI commands.
type CmdGroup string

// PluginCompletionType is the mechanism used for determining command line completion options.
type PluginCompletionType int

const (
	// NativePluginCompletion indicates command line completion is determined using the built in
	// cobra.Command __complete mechanism.
	NativePluginCompletion PluginCompletionType = iota
	// StaticPluginCompletion indicates command line completion will be done by using a statically
	// defined list of options.
	StaticPluginCompletion
	// DynamicPluginCompletion indicates command line completion will be retrieved from the plugin
	// at runtime.
	DynamicPluginCompletion

	// RunCmdGroup are commands associated with Tanzu Run.
	RunCmdGroup CmdGroup = "Run"

	// ManageCmdGroup are commands associated with Tanzu Manage.
	ManageCmdGroup CmdGroup = "Manage"

	// BuildCmdGroup are commands associated with Tanzu Build.
	BuildCmdGroup CmdGroup = "Build"

	// ObserveCmdGroup are commands associated with Tanzu Observe.
	ObserveCmdGroup CmdGroup = "Observe"

	// SystemCmdGroup are system commands.
	SystemCmdGroup CmdGroup = "System"

	// VersionCmdGroup are version commands.
	VersionCmdGroup CmdGroup = "Version"

	// AdminCmdGroup are admin commands.
	AdminCmdGroup CmdGroup = "Admin"

	// TestCmdGroup is the test command group.
	TestCmdGroup CmdGroup = "Test"

	// ExtraCmdGroup is the extra command group.
	ExtraCmdGroup CmdGroup = "Extra"
)

// PluginDescriptor describes a plugin binary.
type PluginDescriptor struct {
	// Name is the name of the plugin.
	Name string `json:"name" yaml:"name"`

	// Description is the plugin's description.
	Description string `json:"description" yaml:"description"`

	// Version of the plugin. Must be a valid semantic version https://semver.org/
	Version string `json:"version" yaml:"version"`

	// BuildSHA is the git commit hash the plugin was built with.
	// TODO(stmcginnis): Update Makefile to set build info with LDFLAG.
	BuildSHA string `json:"buildSHA" yaml:"buildSHA"`

	// Digest is the SHA256 hash of the plugin binary.
	Digest string `json:"digest" yaml:"digest"`

	// Command group for the plugin.
	Group CmdGroup `json:"group" yaml:"group"`

	// DocURL for the plugin.
	DocURL string `json:"docURL,omitempty" yaml:"docURL,omitempty"`

	// Hidden tells whether the plugin should be hidden from the help command.
	Hidden bool `json:"hidden,omitempty" yaml:"hidden,omitempty"`

	// CompletionType determines how command line completion will be determined.
	CompletionType PluginCompletionType `json:"completionType" yaml:"completionType"`

	// CompletionArgs contains the valid command line completion values if `CompletionType`
	// is set to `StaticPluginCompletion`.
	CompletionArgs []string `json:"completionArgs,omitempty" yaml:"completionArgs,omitempty"`

	// CompletionCommand is the command to call from the plugin to retrieve a list of
	// valid completion nouns when `CompletionType` is set to `DynamicPluginCompletion`.
	CompletionCommand string `json:"completionCmd,omitempty" yaml:"completionCmd,omitempty"`

	// Aliases are other text strings used to call this command
	Aliases []string `json:"aliases,omitempty" yaml:"aliases,omitempty"`

	// InstallationPath is a relative installation path for a plugin binary.
	// E.g., cluster/v0.3.2@sha256:...
	InstallationPath string `json:"installationPath"`

	// Discovery is the name of the discovery from where
	// this plugin is discovered.
	Discovery string `json:"discovery"`

	// Scope is the scope of the plugin. Stand-Alone or Context
	Scope string `json:"scope"`
}

// NewPlugin creates an instance of Plugin.
func NewPlugin(descriptor *PluginDescriptor) (*Plugin, error) {
	p := &Plugin{
		Cmd: &cobra.Command{
			Use:     descriptor.Name,
			Short:   descriptor.Description,
			Aliases: descriptor.Aliases,
		},
	}

	// TODO(stmcginnis): There is a docs command we may want to add, need to
	// see if that's necessary.
	p.Cmd.AddCommand(
		newDescribeCmd(descriptor.Description),
		newVersionCmd(descriptor.Version),
		newInfoCmd(descriptor),
		newLintCmd(),
	)
	return p, nil
}

// NewTestFor creates a plugin descriptor for a test plugin.
func NewTestFor(pluginName string) *PluginDescriptor {
	return &PluginDescriptor{
		Name:        fmt.Sprintf("%s-test", pluginName),
		Description: fmt.Sprintf("test for %s", pluginName),
		Version:     "v0.0.1",
		BuildSHA:    SHA,
		Group:       TestCmdGroup,
		Aliases:     []string{fmt.Sprintf("%s-alias", pluginName)},
	}
}

func newDescribeCmd(description string) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "describe",
		Short:  "Describes the plugin",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(description)
			return nil
		},
	}

	return cmd
}

func newVersionCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Short:  "Version the plugin",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(version)
			return nil
		},
	}

	return cmd
}

func newInfoCmd(desc *PluginDescriptor) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "info",
		Short:  "Plugin info",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := json.Marshal(desc)
			if err != nil {
				return err
			}
			fmt.Println(string(b))
			return nil
		},
	}

	return cmd
}

// newLintCmd is used to make sure all commands meet the linting standards.
// TODO(stmcginnis): see if we actually need this.
func newLintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "lint",
		Hidden:       true,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}

// AddCommands adds commands to the plugin.
func (p *Plugin) AddCommands(commands ...*cobra.Command) {
	p.Cmd.AddCommand(commands...)
}

// Execute executes the plugin.
func (p *Plugin) Execute() error {
	return p.Cmd.Execute()
}
