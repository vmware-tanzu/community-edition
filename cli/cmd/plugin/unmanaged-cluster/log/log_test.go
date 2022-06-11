package log

import (
	"bufio"
	"bytes"
	"io/ioutil"
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
	logger.Writer(ioutil.Discard)
	logger.Info("This should be discarded")
	writer.Flush()
	lineCount := len(strings.Split(strings.Trim(b.String(), "\n"), "\n"))
	if lineCount != 2 {
		t.Errorf("Expected 2 lines to be captured into buffer but found %d", lineCount)
	}
}
