// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ytt

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"

	uexec "k8s.io/utils/exec"
)

// Command struct holds ytt command related information.
type Command struct {
	options CommandOptions
	path    string
}

// CommandOptions specifies ytt command options.
type CommandOptions struct {
	FailOnUnknownComments  bool
	Strict                 bool
	DangerousAllowSymlinks bool
}

// NewYttCommand returns a new instance of YTTCommand.
func NewYttCommand(options CommandOptions) *Command {
	return &Command{
		options: options,
		path:    "ytt",
	}
}

// RenderTemplate renders template given a set of file/directory Paths or standard input or both.
func (ytt *Command) RenderTemplate(filePaths []string, input io.Reader) (string, error) {
	var (
		stdout, stderr bytes.Buffer
	)

	args := ytt.buildArgs(filePaths, input)

	cmd := exec.Command(ytt.path, args...) //nolint
	cmd.Stdin = input
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		rc := -1

		if ee, ok := err.(*exec.ExitError); ok {
			rc = ee.ExitCode()
		}

		return stdout.String(), uexec.CodeExitError{
			Err:  fmt.Errorf("error running %v:\nCommand stdout:\n%v\nstderr:\n%v\nerror:\n%v", cmd, cmd.Stdout, cmd.Stderr, err),
			Code: rc,
		}
	}

	return stdout.String(), nil
}

func (ytt *Command) buildArgs(filePaths []string, input io.Reader) []string {
	args := []string{}

	if !ytt.options.FailOnUnknownComments {
		args = append(args, "--ignore-unknown-comments")
	}

	if ytt.options.Strict {
		args = append(args, "--strict")
	}

	if ytt.options.DangerousAllowSymlinks {
		args = append(args, "--dangerous-allow-all-symlink-destinations")
	}

	for _, filePath := range filePaths {
		args = append(args, "-f", filePath)
	}

	if input != nil {
		args = append(args, "-f", "-")
	}

	return args
}

// RenderYTTTemplate is a convenience function to render YTT template.
func RenderYTTTemplate(options CommandOptions, filePaths []string, input io.Reader) (string, error) {
	return NewYttCommand(options).RenderTemplate(filePaths, input)
}
