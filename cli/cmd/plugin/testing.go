// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	uuid "github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

// Main holds state for multiple command tests.
type Main struct {
	// Name of the main test.
	Name string
	// Tests in this main.
	Tests []*Test
	// Report holds the report for the test.
	Report *Report
	// Cleanup function.
	Cleanup CleanupFunc
	// Whether a report should be printed.
	printReport bool
	// Whether to defer the resource deletion.
	DeferDelete bool
}

// CleanupFunc is executed at the end of the test.
type CleanupFunc func() error

// NoCleanupFunc is cleanup function that just returns nil.
var NoCleanupFunc = func() error { return nil }

const deferDelete = "defer-delete"
const printReportFlagName = "print-report"

// NewMain returns a new CLI test.
func NewMain(name string, c *cobra.Command, cleanup CleanupFunc) *Main {
	b, _ := c.Flags().GetBool(printReportFlagName)
	d, _ := c.Flags().GetBool(deferDelete)
	if d {
		cleanup = NoCleanupFunc
	}

	m := &Main{
		Name: name,
		Report: &Report{
			TestName:  name,
			TimeStart: time.Now(),
		},
		Cleanup:     cleanup,
		printReport: b,
		DeferDelete: d,
	}
	m.PrintStart()
	return m
}

// Report is the report generated from running a CLI test.
type Report struct {
	// Name of the report.
	TestName string `json:"testName"`
	// TimeStart at which the test was started.
	TimeStart time.Time `json:"timeStart"`
	// TimeEnd at which the test was started.
	TimeEnd time.Time `json:"timeEnd"`
	// Pass tells whether all the tests passed.
	Pass bool `json:"pass"`
	// Results of the test.
	Results []*Result `json:"results"`
}

// Result is the result of an individual CLI test.
type Result struct {
	// Command executed.
	Command string `json:"command"`
	// Pass tells if command passed execution.
	Pass bool `json:"pass"`
	// Err holds any error produced.
	Err error `json:"err"`
}

// Error fails a result and reports the error.
func (r *Result) Error(err error) {
	r.Pass = false
	r.Err = err
}

// Success passes a result.
func (r *Result) Success() {
	r.Pass = true
}

// PrintStart prints a main test start message.
func (m *Main) PrintStart() {
	fmt.Println("---")
	log.Printf("testing %s", m.Name)
	fmt.Println("")
}

// PrintSuccess prints a Printful main test message.
func (m *Main) PrintSuccess() {
	fmt.Println("")
	log.Printf("ok: Printfully tested %s", m.Name)
	fmt.Println("")
}

// PrintFailure prints a main test failure message.
func (m *Main) PrintFailure() {
	fmt.Println("")
	log.Printf("FAIL: %s", m.Name)
	fmt.Println("")
}

// ReportResult adds a result to the report.
func (m *Main) ReportResult(res *Result) {
	m.Report.Results = append(m.Report.Results, res)
}

// ReportError adds an error result to the report.
func (m *Main) ReportError(cmd string, err error) {
	res := &Result{Command: cmd, Pass: false, Err: err}
	m.Report.Results = append(m.Report.Results, res)
}

// ReportSuccess adds an error result to the report.
func (m *Main) ReportSuccess(cmd string) {
	res := &Result{Command: cmd, Pass: true}
	m.Report.Results = append(m.Report.Results, res)
}

// ReportTestResult adds an error result to the report.
func (m *Main) ReportTestResult(t *Test) {
	m.Report.Results = append(m.Report.Results, t.Result)
	if t.Result.Pass {
		log.Printf("PASS: %q", t.Name)
	} else {
		log.Printf("FAIL: %q", t.Name)
	}
}

// AddTest adds a test to the main.
func (m *Main) AddTest(t *Test) {
	m.Tests = append(m.Tests, t)
}

// BuildReport will build the report with the current commands.
func (m *Main) BuildReport() {
	pass := true
	for _, t := range m.Tests {
		if !t.Result.Pass {
			pass = false
			break
		}
	}
	m.Report.Pass = pass
	for _, t := range m.Tests {
		m.ReportResult(t.Result)
	}
}

// PrintReport prints the report in json|yaml.
func (m *Main) PrintReport(format string) error {
	switch format {
	case "json":
		b, err := json.Marshal(m.Report)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	case "yaml":
		b, err := yaml.Marshal(m.Report)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		return fmt.Errorf("unknown format %q, can be json|yaml", format)
	}
	return nil
}

// Finish test.
func (m *Main) Finish() {
	m.BuildReport()
	m.Report.TimeEnd = time.Now()
	fmt.Println("")
	log.Println("cleaning up")
	if err := m.Cleanup(); err != nil {
		log.Printf("error cleaning up %s", err)
	}

	if m.printReport {
		if err := m.PrintReport("yaml"); err != nil {
			log.Printf("PrintReport failed %s", err)
		}
	}
	if m.Report.Pass {
		m.PrintSuccess()
	} else {
		m.PrintFailure()
	}
}

// FlagSet returns the default flagset values for cli tests.
func FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("cli-test-flags", pflag.ContinueOnError)
	fs.BoolP(printReportFlagName, "p", false, "print report")
	fs.BoolP(deferDelete, "d", false, "defer resource deletion")
	return fs
}

// NamePrefix is the prefix used in generated names.
const NamePrefix = "cli-test"

// GenerateName returns a name for a cli test.
func GenerateName() string {
	testName := fmt.Sprintf("%s-%s", NamePrefix, uuid.NewString()[:8])
	return testName
}

// Test is a cli test to run.
type Test struct {
	// Name of the test.
	Name string
	// Command to test.
	Command string
	// Run function to wrap any logic exectution.
	run func(t *Test) error
	// Result of command test.
	Result *Result

	stdOut *bytes.Buffer
	stdErr *bytes.Buffer
}

// NewTest returns a new command
func NewTest(name, command string, run func(t *Test) error) *Test {
	log.Println(name)
	return &Test{
		Name:    name,
		Command: command,
		Result: &Result{
			Command: command,
		},
		run: run,
	}
}

// NewTest creates a new Test and adds it to the main.
func (m *Main) NewTest(name, command string, run func(t *Test) error) *Test {
	t := NewTest(name, command, run)
	m.AddTest(t)
	return t
}

// RunTest will create a test and run it.
func (m *Main) RunTest(name, command string, run func(t *Test) error) error {
	t := m.NewTest(name, command, run)
	return t.Run()
}

// PrintSuccess will print a success message.
func (t *Test) PrintSuccess() {
	log.Printf("ok: %s", t.Name)
}

// Run the 'run' function within the context of the test updating the result accordingly.
func (t *Test) Run() error {
	return t.Wrap(t.run)
}

// Wrap will wrap the execution of the given function in the context of the test and update the
// results accordingly.
func (t *Test) Wrap(f func(t *Test) error) error {
	err := f(t)
	if err != nil {
		t.Result.Error(err)
		return err
	}
	t.Result.Success()
	t.PrintSuccess()
	return nil
}

// Exec will execute the test command.
func (t *Test) Exec() (err error) {
	s, e, err := Exec(t.Command)
	if err != nil {
		t.Result.Error(err)
		return err
	}
	t.stdOut = s
	t.stdErr = e
	t.Result.Success()
	return nil
}

// StdOut from executing the test command.
func (t *Test) StdOut() *bytes.Buffer {
	return t.stdOut
}

// StdErr from executing the test command.
func (t *Test) StdErr() *bytes.Buffer {
	return t.stdErr
}

// Exec the command, exit on error
func Exec(command string) (stdOut, stdErr *bytes.Buffer, err error) {
	c := cleanCommand(command)
	cmd := exec.Command("tanzu", c...)

	var stdOutBytes, stdErrBytes []byte
	var errStdout, errStderr error
	stdOutIn, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	stdErrIn, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	fmt.Printf("$ %s \n", strings.Join(cmd.Args, " "))
	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		stdOutBytes, errStdout = copyAndCapture(os.Stdout, stdOutIn)
		wg.Done()
	}()
	stdErrBytes, errStderr = copyAndCapture(os.Stderr, stdErrIn)
	wg.Wait()
	err = cmd.Wait()
	if err != nil {
		fmt.Println(stdOut.String())
		fmt.Println(stdErr.String())
	}
	if errStdout != nil {
		return nil, nil, fmt.Errorf("failed to capture stdout: %w", errStdout)
	}
	if errStderr != nil {
		return nil, nil, fmt.Errorf("failed to capture stderr: %w", errStderr)
	}
	return bytes.NewBuffer(stdOutBytes), bytes.NewBuffer(stdErrBytes), err
}

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}

// cleanCommand will remove the CLIName from the command if exists as first argument.
func cleanCommand(command string) []string {
	c := strings.Split(command, " ")
	if c[0] == "tanzu" {
		c = c[1:]
	}
	return c
}

// ExecContainsString executes the command and checks if the output contains the given string.
func (t *Test) ExecContainsString(contains string) error {
	err := ExecContainsString(t.Command, contains)
	if err != nil {
		t.Result.Error(err)
		return err
	}
	t.Result.Success()
	return nil
}

// ExecContainsAnyString executes the command and checks if the output contains any of the given set of strings.
func (t *Test) ExecContainsAnyString(contains ...string) error {
	err := ExecContainsAnyString(t.Command, contains)
	if err != nil {
		t.Result.Error(err)
		return err
	}
	t.Result.Success()
	return nil
}

// ExecContainsString checks that the given command output contains the string.
func ExecContainsString(command, contains string) error {
	stdOut, _, err := Exec(command)
	if err != nil {
		return err
	}
	return ContainsString(stdOut, contains)
}

// ExecContainsAnyString checks that the given command output contains any of the given set of strings.
func ExecContainsAnyString(command string, contains []string) error {
	stdOut, _, err := Exec(command)
	if err != nil {
		return err
	}
	return ContainsAnyString(stdOut, contains)
}

// ExecContainsErrorString executes the command and checks if the output contains the given string.
func (t *Test) ExecContainsErrorString(contains string) error {
	err := ExecContainsErrorString(t.Command, contains)
	if err != nil {
		t.Result.Error(err)
		return err
	}
	t.Result.Success()
	return nil
}

// ExecContainsErrorString checks that the given command stdErr output contains the string
func ExecContainsErrorString(command, contains string) error {
	_, stdErr, err := Exec(command)
	if err != nil {
		return err
	}
	return ContainsString(stdErr, contains)
}

// ContainsString checks that the given buffer contains the string.
func ContainsString(stdOut *bytes.Buffer, contains string) error {
	so := stdOut.String()
	if !strings.Contains(so, contains) {
		return fmt.Errorf("stdOut %q did not contain %q", so, contains)
	}
	return nil
}

// ContainsAnyString checks that the given buffer contains any of the given set of strings.
func ContainsAnyString(stdOut *bytes.Buffer, contains []string) error {
	var containsAny bool
	so := stdOut.String()

	for _, str := range contains {
		containsAny = containsAny || strings.Contains(so, str)
	}

	if !containsAny {
		return fmt.Errorf("stdOut %q did not contain of the following %q", so, contains)
	}
	return nil
}
