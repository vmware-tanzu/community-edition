// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	kappdiff "github.com/k14s/kapp/pkg/kapp/diff"
	kappresources "github.com/k14s/kapp/pkg/kapp/resources"
	"github.com/stretchr/testify/require"
)

const (
	fixtureDir               = "fixtures"
	fixtureInputDir          = "values"
	fixtureExpectedOutputDir = "expected"
)

func TestTemplate(t *testing.T) {
	for _, pkg := range findPkgs(t) {
		pkg := pkg
		t.Run(pkg.version, func(t *testing.T) {
			t.Parallel()

			for _, testCase := range findTestCases(t, pkg) {
				testCase := testCase
				t.Run(testCase.name, func(t *testing.T) {
					t.Parallel()

					runTestCase(t, testCase)
				})
			}
		})
	}
}

type pkg struct {
	version string
	dir     string
}

func findPkgs(t *testing.T) []*pkg {
	t.Helper()

	var pkgs []*pkg

	pkgsRoot := ".."
	pkgsRootEntries, err := os.ReadDir(pkgsRoot)
	require.NoError(t, err)
	for _, pkgsRootEntry := range pkgsRootEntries {
		if !pkgsRootEntry.Type().IsDir() {
			continue
		}

		dir := filepath.Join(pkgsRoot, pkgsRootEntry.Name())
		if fileExists(t, filepath.Join(dir, "bundle")) {
			pkgs = append(pkgs, &pkg{
				version: pkgsRootEntry.Name(),
				dir:     filepath.Join(wd(t), dir),
			})
		} else {
			t.Logf("findPkgs: ignoring non-package directory %q", dir)
		}
	}

	sort.Slice(pkgs, func(i, j int) bool { return pkgs[i].version < pkgs[j].version })

	return pkgs
}

type testCase struct {
	name string

	pkg *pkg

	inFile, outFile string
}

func findTestCases(t *testing.T, pkg *pkg) []*testCase {
	t.Helper()

	var testCases []*testCase

	pkgTestDir := filepath.Join(pkg.dir, fixtureDir)
	pkgTestInDir := filepath.Join(pkgTestDir, fixtureInputDir)
	pkgTestOutDir := filepath.Join(pkgTestDir, fixtureExpectedOutputDir)

	require.Truef(t, fileExists(t, pkgTestDir), "expected package test directory %q to exist", pkgTestDir)
	require.Truef(t, fileExists(t, pkgTestDir), "expected package test input directory %q to exist", pkgTestInDir)
	require.Truef(t, fileExists(t, pkgTestDir), "expected package test output directory %q to exist", pkgTestOutDir)

	pkgTestInDirEntries, err := os.ReadDir(pkgTestInDir)
	require.NoError(t, err)
	for _, pkgTestInDirEntry := range pkgTestInDirEntries {
		if !pkgTestInDirEntry.Type().IsRegular() {
			continue
		}

		pkgTestInDirFile := filepath.Join(pkgTestInDir, pkgTestInDirEntry.Name())
		pkgTestOutDirFile := filepath.Join(pkgTestOutDir, pkgTestInDirEntry.Name())
		require.Truef(
			t,
			fileExists(t, pkgTestOutDirFile),
			"expected out file %q for in file %q",
			pkgTestOutDirFile,
			pkgTestInDirFile,
		)

		testCases = append(testCases, &testCase{
			name:    pkgTestInDirEntry.Name(),
			pkg:     pkg,
			inFile:  pkgTestInDirFile,
			outFile: pkgTestOutDirFile,
		})
	}

	sort.Slice(testCases, func(i, j int) bool { return testCases[i].name < testCases[j].name })

	return testCases
}

func runTestCase(t *testing.T, testCase *testCase) {
	t.Helper()

	wantOut := readFile(t, testCase.outFile)
	yttCommand, gotOut := ytt(t, filepath.Join(testCase.pkg.dir, "bundle", "config"), testCase.inFile)
	if diff := diff(t, wantOut, gotOut); len(diff) > 0 {
		t.Logf("to see diff: kapp tools diff --changes --summary=false --file %s --file2 <(%s)", testCase.outFile, yttCommand)
		t.Logf("to update expected: %s >%s", yttCommand, testCase.outFile)
		t.Errorf("-want,+got;\n%s", diff)
	}
}

func ytt(t *testing.T, inputs ...string) (string, string) {
	t.Helper()

	args := []string{"--ignore-unknown-comments"}
	for _, input := range inputs {
		args = append(args, "-f")
		args = append(args, input)
	}

	cmd := exec.Command("ytt", args...)

	stdout, stderr := bytes.NewBuffer([]byte{}), bytes.NewBuffer([]byte{})
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	t.Log("running", cmd.Args)
	err := cmd.Run()
	require.NoError(t, err, "ytt command %s failed (stderr:\n%s)", cmd.Args, stderr.String())

	return strings.Join(cmd.Args, " "), stdout.String()
}

func diff(t *testing.T, want, got string) string {
	t.Helper()

	wantResources, err := kappresources.NewFileResource(kappresources.NewBytesSource([]byte(want))).Resources()
	require.NoError(t, err)
	gotResources, err := kappresources.NewFileResource(kappresources.NewBytesSource([]byte(got))).Resources()
	require.NoError(t, err)

	changeFactory := kappdiff.NewChangeFactory(nil, nil)
	changes, err := kappdiff.NewChangeSet(wantResources, gotResources, kappdiff.ChangeSetOpts{AgainstLastApplied: false}, changeFactory).Calculate()
	require.NoError(t, err)

	var changeset strings.Builder
	for _, change := range changes {
		if diff := change.ConfigurableTextDiff().Full().MinimalString(); len(diff) > 0 {
			_, err := changeset.WriteString(change.ConfigurableTextDiff().Full().FullString())
			require.NoError(t, err)
		}
	}

	return changeset.String()
}

func wd(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	require.NoError(t, err)

	return dir
}

func fileExists(t *testing.T, path string) bool {
	t.Helper()

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		require.NoError(t, err)
	}

	return true
}

func readFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	return string(data)
}
