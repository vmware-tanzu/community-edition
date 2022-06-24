// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestNewLoggerWithWriter(t *testing.T) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	logger := NewLoggerWithWriter(false, 5, writer)
	logger.Info("Generating Message in Info")
	logger.Info("Generating Message in Info again")
	writer.Flush()
	lineCount := len(strings.Split(strings.Trim(b.String(), "\n"), "\n"))
	if lineCount != 2 {
		t.Errorf("Expected 2 lines to be captured into buffer but found %d", lineCount)
	}
}

func TestNewLogger(t *testing.T) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	logger := NewLogger(false, 5)
	logger.Writer(writer)
	logger.Info("Generating Message in Info")
	logger.Info("Generating Message in Info again")
	logger.Writer(os.Stdout)
	logger.Info("This should not be in the buffer")
	logger.Writer(io.Discard)
	logger.Info("This should be discarded")
	logger.Writer(writer)
	logger.Info("This should be captured in buffer writer")
	writer.Flush()
	lineCount := len(strings.Split(strings.Trim(b.String(), "\n"), "\n"))
	if lineCount != 3 {
		t.Errorf("Expected 2 lines to be captured into buffer but found %d", lineCount)
	}
}
