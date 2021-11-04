// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"os"
	"runtime"
	"strconv"

	"github.com/spf13/pflag"
)

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

	if flags.Changed("tty") {
		// User has explicitly set the flag, use that value
		result, _ = flags.GetBool("tty")
	} else if tty := os.Getenv("TANZU_TTY"); tty != "" {
		// Not explicitly provided, but there is an env setting
		val, err := strconv.ParseBool(tty)
		if err == nil {
			result = val
		}
	}
	return result
}
