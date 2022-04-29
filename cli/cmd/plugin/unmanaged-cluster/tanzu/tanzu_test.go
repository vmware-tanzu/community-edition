// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tanzu

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/config"
	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/log"
)

func TestNew(t *testing.T) {
	unmanagedCluster := New(logger.NewLogger(true, 0))
	if unmanagedCluster == nil {
		t.Error("New should generate an unmanaged cluster object")
	}
}

func TestValidateConfiguration(t *testing.T) {
	scConfig := config.UnmanagedClusterConfig{}

	err := validateConfiguration(&scConfig)
	if err == nil {
		t.Error("validate configuration should fail if no cluster name provided")
	}

	scConfig.ClusterName = "test"
	err = validateConfiguration(&scConfig)
	if err != nil {
		t.Error("validate configuration should succeed if a cluster name is provided")
	}
}

func TestGetUnmanagedBomPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	bomPath, err := getUnmanagedBomPath()
	if err != nil {
		t.Error("bom path should always succeed")
	}
	if bomPath == "" {
		t.Error("bom path should always return a valid bomPath")
	}
	if !strings.Contains(bomPath, homeDir) {
		t.Errorf("bom path should always contain the homedir, was %s", bomPath)
	}
	if !strings.Contains(bomPath, ".config") {
		t.Errorf("bom path should always contain .config in the the path, was %s", bomPath)
	}
}

func TestGetUnmanagedCompatibilityPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	compatPath, err := getUnmanagedCompatibilityPath()
	if err != nil {
		t.Error("compatibility path should always succeed")
	}
	if compatPath == "" {
		t.Error("compatibility path should always return a valid compatPath")
	}
	if !strings.Contains(compatPath, homeDir) {
		t.Errorf("compatibility path should always contain the homedir, was %s", compatPath)
	}
	if !strings.Contains(compatPath, ".config") {
		t.Errorf("compatibility path should always contain .config in the the path, was %s", compatPath)
	}
}

func TestBuildFilesystemSafeBomName(t *testing.T) {
	testStr := buildFilesystemSafeBomName("hello/THESE:are$NUMBERS*0123456789")

	if !strings.Contains(testStr, "hello_THESE_areNUMBERS0123456789") {
		t.Errorf("safe bom name should always contain .config in the the path, was %s", testStr)
	}
}

func TestResolveClusterDir(t *testing.T) {
	dir, err := resolveClusterDir("")
	if err == nil {
		t.Error("resolve cluster dir should fail if an empty cluster name is provided")
	}
	if dir != "" {
		t.Errorf("resolve cluster dir should fail if an empty cluster name is provided, was %s", dir)
	}

	dir, err = resolveClusterDir("test")
	if err == nil {
		t.Error("resolve cluster dir should fail if an empty cluster name is provided")
	}
	if dir != "" {
		t.Errorf("resolve cluster dir should fail if an empty cluster name is provided, was %s", dir)
	}

	// Simulates cluster directory, need to clean up later
	_, err = createClusterDirectory("tmp-tanzu-resolve-cluster-dir")
	if err != nil {
		t.Fatal(err)
	}

	dir, err = resolveClusterDir("tmp-tanzu-resolve-cluster-dir")
	if err != nil {
		t.Error("resolve cluster dir should succeeded if directory exist")
	}
	if dir == "" {
		t.Errorf("resolve cluster dir should succeeded if directory exist, was %s", dir)
	}
	if !strings.Contains(dir, "tmp-tanzu-resolve-cluster-dir") {
		t.Errorf("resolve cluster dir should contain cluster name, was %s", dir)
	}

	// clean up
	err = os.RemoveAll(dir)
	if err != nil {
		t.Fatal(err)
	}
}

func TestResolveClusterConfig(t *testing.T) {
	cfg, err := resolveClusterConfig("")
	if err == nil {
		t.Error("resolve cluster config should fail if an empty cluster name is provided")
	}
	if cfg != "" {
		t.Errorf("resolve cluster config should fail if an empty cluster name is provided, was %s", cfg)
	}

	cfg, err = resolveClusterConfig("test")
	if err == nil {
		t.Error("resolve cluster config should fail if an empty cluster name is provided")
	}
	if cfg != "" {
		t.Errorf("resolve cluster config should fail if an empty cluster name is provided, was %s", cfg)
	}

	// Simulates cluster directory, need to clean up later
	dir, err := createClusterDirectory("tmp-tanzu-resolve-cluster-config")
	if err != nil {
		t.Fatal(err)
	}
	testCfg := filepath.Join(dir, configFileName)
	emptyFile, err := os.Create(testCfg)
	if err != nil {
		t.Fatal(err)
	}
	emptyFile.Close()

	cfg, err = resolveClusterConfig("tmp-tanzu-resolve-cluster-config")
	if err != nil {
		t.Error("resolve cluster config should exist")
	}
	if cfg == "" {
		t.Errorf("resolve cluster config should exist, was %s", cfg)
	}
	if !strings.Contains(cfg, "tmp-tanzu-resolve-cluster-config") {
		t.Errorf("resolve cluster config should contain cluster name, was %s", cfg)
	}
	if !strings.Contains(cfg, configFileName) {
		t.Errorf("resolve cluster config should contain %s, was %s", configFileName, cfg)
	}

	// clean up
	err = os.RemoveAll(dir)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateClusterDirectory(t *testing.T) {
	dir, err := createClusterDirectory("")
	if err == nil {
		t.Error("create cluster dir should fail if an empty cluster name is provided")
	}
	if dir != "" {
		t.Errorf("create cluster dir should fail if an empty cluster name is provided, was %s", dir)
	}

	dir, err = createClusterDirectory("tmp-tanzu-create-cluster-dir")
	if err != nil {
		t.Error("resolve cluster dir should succeeded if directory exist")
	}
	if dir == "" {
		t.Errorf("resolve cluster dir should succeeded if directory exist, was %s", dir)
	}
	if !strings.Contains(dir, "tmp-tanzu-create-cluster-dir") {
		t.Errorf("resolve cluster dir should contain cluster name, was %s", dir)
	}

	err = os.RemoveAll(dir)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetCompatibilityFile(t *testing.T) {
	file, err := getCompatibilityFile()
	if err != nil {
		t.Error("get compat file should succeed")
	}
	if file == "" {
		t.Errorf("get compat file should always non-exmpty, was %s", file)
	}
}

func TestGetTkrCompatibility(t *testing.T) {
	compat, err := getTkrCompatibility()
	if err != nil {
		t.Error("get tkr compatibility should always succeed")
	}
	if compat == nil {
		t.Errorf("get tkr compatibility should not be nil")
	}
}

func TestIsTkrCompatible(t *testing.T) {
	compat, err := getTkrCompatibility()
	if err != nil {
		t.Error("get tkr compatibility should always succeed")
	}
	if compat == nil {
		t.Errorf("get tkr compatibility should not be nil")
	}
	if compat != nil && !strings.EqualFold(compat.Version, "v1") {
		t.Errorf("get tkr compatibility should currently always be v1, was %s", compat.Version)
	}

	latest, err := getLatestCompatibleTkr(compat)
	if err != nil {
		t.Error("get tkr compatibility should always succeed")
	}
	if latest == "" {
		t.Errorf("get tkr compatibility should not be empty")
	}

	// reference a supported version
	ret := isTkrCompatible(compat, latest)
	if !ret {
		t.Error("get compat file should be supported")
	}
	// unsupported
	ret = isTkrCompatible(compat, "projects.registry.vmware.com/tce/tkr:v0.0.0")
	if ret {
		t.Error("get compat file should NOT be supported")
	}
}

func TestGetTkrBom(t *testing.T) {
	compat, err := getTkrCompatibility()
	if err != nil {
		t.Error("get tkr compatibility should always succeed")
	}
	if compat == nil {
		t.Errorf("get tkr compatibility should not be nil")
	}
	if compat != nil && !strings.EqualFold(compat.Version, "v1") {
		t.Errorf("get tkr compatibility should currently always be v1, was %s", compat.Version)
	}

	latest, err := getLatestCompatibleTkr(compat)
	if err != nil {
		t.Error("get tkr compatibility should always succeed")
	}
	if latest == "" {
		t.Errorf("get tkr compatibility should not be empty")
	}

	// invalid
	tkrBomFile, err := getTkrBom("junk")
	if err == nil {
		t.Error("get tkr bom should error out when given a junk value")
	}
	if tkrBomFile != "" {
		t.Errorf("get tkr bom should be empty, was %s", tkrBomFile)
	}

	// first time
	tkrBomFile, err = getTkrBom(latest)
	if err != nil {
		t.Error("get tkr bom should always succeed")
	}
	if tkrBomFile == "" {
		t.Errorf("get tkr bom should not be empty")
	}

	// second time should invoke local tkr
	tkrBomFile, err = getTkrBom(latest)
	if err != nil {
		t.Error("get tkr bom should always succeed")
	}
	if tkrBomFile == "" {
		t.Errorf("get tkr bom should not be empty")
	}

	// clean up
	err = os.RemoveAll(tkrBomFile)
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseTkrBom(t *testing.T) {
	compat, err := getTkrCompatibility()
	if err != nil {
		t.Error("get tkr compatibility should always succeed")
	}
	if compat == nil {
		t.Errorf("get tkr compatibility should not be nil")
	}
	if compat != nil && !strings.EqualFold(compat.Version, "v1") {
		t.Errorf("get tkr compatibility should currently always be v1, was %s", compat.Version)
	}

	latest, err := getLatestCompatibleTkr(compat)
	if err != nil {
		t.Error("get tkr compatibility should always succeed")
	}
	if latest == "" {
		t.Errorf("get tkr compatibility should not be empty")
	}

	tkrBomFile, err := getTkrBom(latest)
	if err != nil {
		t.Error("get tkr bom should always succeed")
	}
	if tkrBomFile == "" {
		t.Errorf("get tkr bom should not be empty")
	}

	// deserialize
	tkrBom, err := parseTKRBom("junk")
	if err == nil {
		t.Error("get tkr bom should always fail with junk value")
	}
	if tkrBom != nil {
		t.Errorf("get tkr bom should be empty with junk file")
	}
	tkrBom, err = parseTKRBom(tkrBomFile)
	if err != nil {
		t.Error("get tkr bom should succeed given a valid file")
	}
	if tkrBom == nil {
		t.Errorf("get tkr bom should not be empty")
	}

	// clean up
	err = os.RemoveAll(tkrBomFile)
	if err != nil {
		t.Fatal(err)
	}
}
