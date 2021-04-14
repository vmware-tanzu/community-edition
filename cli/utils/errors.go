// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	klog "k8s.io/klog/v2"
)

// NonUsageError returns an error caused by a command that was not due to a usage
// error. In this case, we don't want the CLI to emit the usage help text since
// they presumably called the command correctly. Something else then failed with
// the execution of the command.
func NonUsageError(cmd *cobra.Command, err error, message string, args ...interface{}) error {
	cmd.SilenceUsage = true

	return Error(err, message, args...)
}

// Error is used to have a consistent way to format and log new errors.
func Error(err error, message string, args ...interface{}) error {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	if err != nil {
		// Log the full error details for tracing
		klog.V(3).Infof("Error: %v", err)

		return fmt.Errorf("%s\nCause: %s", message, err.Error())
	}

	return errors.New(message)
}
