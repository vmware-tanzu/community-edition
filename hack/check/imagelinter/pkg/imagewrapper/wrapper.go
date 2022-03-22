// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package imagewrapper

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Wrapper struct {
	Image     string
	Container string
	Writer    io.Writer
}

func New(image, contaier string, writer io.Writer) (*Wrapper, error) {
	if image == "" || contaier == "" {
		return nil, errors.New("image or Container parameters cannot be empty")
	}
	return &Wrapper{Image: image, Container: contaier, Writer: writer}, nil
}

func (w *Wrapper) PullImage() (string, error) {
	result, err := w.CliRunner("docker", nil, []string{"pull", w.Image}...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (w *Wrapper) CreateContainer() (string, error) {
	result, err := w.CliRunner("docker", nil, []string{"run", "-d", "--name", w.Container, w.Image}...)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (w *Wrapper) RunCommand(args ...string) (string, error) {
	result, err := w.CliRunner("docker", nil, args...)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (w *Wrapper) IsContainerExists() bool {
	result, _ := w.CliRunner("docker", nil, []string{"ps", "-a", "--format", `table {{.Names}}`}...)
	return strings.Contains(result, w.Container)
}

// src: /etc/os-release	 dst: ./
func (w *Wrapper) ContainerCP(src, dst string) (string, error) {
	result, err := w.CliRunner("docker", nil, []string{"cp", w.Container + ":" + src, dst}...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (w *Wrapper) DeleteContainer() (string, error) {
	result, err := w.CliRunner("docker", nil, []string{"rm", "-f", w.Container}...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (w *Wrapper) Validate(validators []string) (bool, error) {
	history, err := w.CliRunner("docker", nil, []string{"history", w.Image, "--no-trunc"}...)
	if err != nil {
		return false, err
	}
	for _, item := range validators {
		if strings.Contains(history, item) {
			return true, nil
		}
	}
	return false, nil
}

func (w *Wrapper) CliRunner(name string, input io.Reader, args ...string) (string, error) {
	if w.Writer != nil {
		fmt.Fprintf(w.Writer, "+ %s %s\n", name, strings.Join(args, " "))
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

		if w.Writer != nil {
			fmt.Fprintln(w.Writer, stderr.String())
		}
		return "", fmt.Errorf("%s\nexit status: %d", stderr.String(), rc)
	}
	if w.Writer != nil {
		fmt.Fprintln(w.Writer, stdout.String())
	}
	return stdout.String(), nil
}
