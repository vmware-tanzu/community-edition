// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmdhelper

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type CmdHelper struct {
	CommandArgs map[string][]string
	Writer      io.Writer
}

var (
	ErrNilCmdHelper        = errors.New("nil cmdHelper")
	ErrCommandNotFound     = errors.New("command not found")
	ErrCommandArgsNotFound = errors.New("command arguments not found")
	ErrCommandsMustBeFed   = errors.New("commands must be provided")
)

func New(cmds map[string][]string, writer io.Writer) (c *CmdHelper, err error) {
	if cmds == nil {
		return nil, ErrCommandsMustBeFed
	}
	c = &CmdHelper{CommandArgs: cmds, Writer: writer}
	return c, nil
}

// Format formats the command array with a given replace arr. It uses spl (special char) parameter to replace
// it will update the CommandArgs value with a formatted array
func (c *CmdHelper) Format(cmdKey, spl string, rarr []string) {
	arr, ok := c.CommandArgs[cmdKey]
	if ok {
		arr = StrArrReplace(spl, arr, rarr)
		c.CommandArgs[cmdKey] = arr
	}
}

// GetFormatted formats the command array with a given replace array
// it uses spl (special char) parameter to replace and then returns an array with all replaced contents
// this will not update/replace the original array from CommandArgs key
func (c *CmdHelper) GetFormatted(cmdKey, spl string, rarr []string) []string {
	arr, ok := c.CommandArgs[cmdKey]
	if ok {
		farr := make([]string, len(arr))
		copy(farr, arr)
		farr = StrArrReplace(spl, farr, rarr)
		return farr
	}
	return nil
}

func (c *CmdHelper) Run(name string, input io.Reader, cmdKey string) (string, error) {
	if c == nil {
		return "", ErrNilCmdHelper
	}
	if strings.Trim(name, " ") == "" {
		return "", ErrCommandNotFound
	}

	arr, ok := c.CommandArgs[cmdKey]
	if !ok {
		return "", ErrCommandArgsNotFound
	}
	return c.CliRunner(name, input, arr...)
}

// StrArrReplace is to replace an array with a replace array based on a special char
// rarr(replace array) contains strings that have to be replaced in an order.
func StrArrReplace(spl string, arr, rarr []string) []string {
	if len(arr) == 0 || len(rarr) == 0 || strings.Trim(spl, " ") == "" {
		return arr
	}
	j := 0
	for i := 0; i < len(arr); i++ {
		if arr[i] == spl && j < len(rarr) {
			arr[i] = rarr[j]
			j++
		}
	}
	return arr
}

func (c *CmdHelper) CliRunner(name string, input io.Reader, args ...string) (string, error) {
	if c.Writer != nil {
		fmt.Fprintf(c.Writer, "+ %s %s\n", name, strings.Join(args, " "))
	}
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdin = input
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		rc := -1
		if ee, ok := err.(*exec.ExitError); ok {
			rc = ee.ExitCode()
		}

		if c.Writer != nil {
			fmt.Fprintln(c.Writer, stderr.String())
		}
		return "", fmt.Errorf("%s\nexit status: %d", stderr.String(), rc)
	}
	// todo : This code has to be removed once tanzu bug is fixed
	// The below is the workaround since there is a bug as tanzu package install always writes to stderr irrespective of the output
	if stdout.String() == "" {
		if c.Writer != nil {
			fmt.Fprintln(c.Writer, stderr.String())
		}
		return stderr.String(), nil
	}
	// workaround ends.

	if c.Writer != nil {
		fmt.Fprintln(c.Writer, stdout.String())
	}
	return stdout.String(), nil
}

// CliRunnerChan is to kill long running commands upon a signal.
// it sends the *exec.Cmd. Receiver channel will receive it and Kill the process based on conditions
func (c *CmdHelper) CliRunnerChan(name string, input io.Reader, singnal chan<- *exec.Cmd, args ...string) (string, error) {

	if c.Writer != nil {
		fmt.Fprintf(c.Writer, "+ %s %s\n", name, strings.Join(args, " "))
	}
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdin = input
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	singnal <- cmd
	err := cmd.Run()
	if err != nil {
		rc := -1
		if ee, ok := err.(*exec.ExitError); ok {
			rc = ee.ExitCode()
		}

		if c.Writer != nil {
			fmt.Fprintln(c.Writer, stderr.String())
		}
		fmt.Errorf("%s\nexit status: %d", stderr.String(), rc)
		return "", fmt.Errorf("%s\nexit status: %d", stderr.String(), rc)
	}

	// todo : This code has to be removed once tanzu bug is fixed
	// The below is the workaround since there is a bug as tanzu package install always writes to stderr irrespective of the output
	if stdout.String() == "" {
		if c.Writer != nil {
			fmt.Fprintln(c.Writer, stderr.String())
		}
		return stderr.String(), nil
	}
	// workaround ends.

	if c.Writer != nil {
		fmt.Fprintln(c.Writer, stdout.String())
	}

	return stdout.String(), nil
}
