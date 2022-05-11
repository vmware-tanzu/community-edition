// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"os"
	"runtime"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

const defaultLogVerbosity = 2

// TtySetting gets the setting to use for formatted TTY output based on whether
// the user explicitly set it with a command line argument, or if not, whether
// there is an environment variable set. If neither of these things, it will
// default to whether or not we detect we are running in a terminal that allows
// tty formatting.
func TtySetting(flags *pflag.FlagSet) bool {
	var result bool

	// See if we are running in a tty enabled terminal
	if runtime.GOOS == "windows" {
		// The newer Windows Terminal supports unicode, cmd and powershell do
		// not. Currently the only way to tell if you are running in WinTerm is
		// by the presence of a "WT_SESSION" environment variable.
		result = os.Getenv("WT_SESSION") != ""
	} else {
		// For Mac and Linux we can interrogate the terminal
		fileInfo, _ := os.Stdout.Stat()
		result = (fileInfo.Mode() & os.ModeCharDevice) != 0
	}

	if flags.Changed("tty-disable") {
		// User has explicitly set the flag, use that value
		disableTTY, err := flags.GetBool("tty-disable")
		if err == nil {
			result = !disableTTY
		}
	} else if tty := os.Getenv("TANZU_TTY_DISABLE"); tty != "" {
		// Not explicitly provided, but there is an env setting
		disableTTY, err := strconv.ParseBool(tty)
		if err == nil {
			result = !disableTTY
		}
	}
	return result
}

// SetupRootCommand ensures the root cobra command is setup with our customization
func SetupRootCommand(rootCmd *cobra.Command) {
	cobra.AddTemplateFunc("wrappedFlagUsages", wrappedFlagUsages)

	rootCmd.SetUsageTemplate(usageTemplate)
	rootCmd.SetHelpTemplate(helpTemplate)
}

// LoggingVerbosity will get the configured logging level for this command.
func LoggingVerbosity(cmd *cobra.Command) int {
	level, err := cmd.Flags().GetInt32("verbose")
	if err != nil {
		// Flag missing or read failure, use default
		return defaultLogVerbosity
	}

	return int(level)
}

// Uses the users terminal size or width of 80 if cannot determine users width
func wrappedFlagUsages(cmd *pflag.FlagSet) string {
	fd := int(os.Stdout.Fd())
	width := 80

	// Get the terminal width and dynamically set
	termWidth, _, err := term.GetSize(fd)
	if err == nil {
		width = termWidth
	}

	return cmd.FlagUsagesWrapped(width - 1)
}

// Identical to the default cobra usage template,
// but utilizes wrappedFlagUsages to ensure flag usages don't wrap around
var usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{wrappedFlagUsages .LocalFlags | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{wrappedFlagUsages .InheritedFlags | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

var helpTemplate = `
{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
