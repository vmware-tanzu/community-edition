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

// Forat formats the command array with given replace arr. It uses spl (special char) parameter to replace
func (c *CmdHelper) Format(cmdKey, spl string, rarr []string) {
	arr, ok := c.CommandArgs[cmdKey]
	if ok {
		arr = StrArrReplace(spl, arr, rarr)
		c.CommandArgs[cmdKey] = arr
	}
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
	return c.cliRunner(name, input, arr...)
}

// StrArrReplace is to replace an array with a replace array based on a special charcter
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

func (c *CmdHelper) cliRunner(name string, input io.Reader, args ...string) (string, error) {
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

	if stdout.String() == "" {
		if c.Writer != nil {
			fmt.Fprintln(c.Writer, stderr.String())
		}
		return stderr.String(), nil
	}

	if c.Writer != nil {
		fmt.Fprintln(c.Writer, stdout.String())
	}
	return stdout.String(), nil
}
