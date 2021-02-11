// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package extension

import (
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"time"
	"os"
	"flag"
	"strconv"

	"github.com/adrg/xdg"
	yaml "github.com/ghodss/yaml"
	klog "k8s.io/klog/v2"

	cfg "github.com/vmware-tanzu/tce/pkg/common/config"
	gcp "github.com/vmware-tanzu/tce/pkg/common/gcp"
	github "github.com/vmware-tanzu/tce/pkg/common/github"
	kapp "github.com/vmware-tanzu/tce/pkg/common/kapp"
	types "github.com/vmware-tanzu/tce/pkg/common/types"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Manager
var mgr *Manager

// common vars
var outputFormat string

// App CRD input
var (
	inputAppCrd = &kapp.AppCrdInput{
		Paths: make(map[string]string),
	}
)

// GetDebugLevel default is 2 (aka DefaultLogLevel)
func GetDebugLevel(s string) string {
	_, err := strconv.Atoi(s)
	if err != nil {
		return DefaultLogLevel
	}
	return s
}

// NewManager generates a Manager object
func NewManager() (*Manager, error) {

	// logging...
	klog.InitFlags(nil)

	level := "0"
	if v := os.Getenv("TCE_EXTENSION_DEBUG"); v != "" {
		level = GetDebugLevel(v)
	}
	flag.Set("v", level)
	flag.Parse()

	// read config
	configFile := filepath.Join(xdg.DataHome, "tanzu-repository", cfg.DefaultConfigFile)
	byFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		klog.Errorf("ReadFile failed. Err: %v", err)
		return nil, err
	}

	config, err := cfg.InitConfig(byFile)
	if err != nil {
		klog.Errorf("NewConfig failed. Err: %v", err)
		return nil, err
	}

	mgr := &Manager{
		cfg: config,
	}

	bucket, err := gcp.NewBucket(byFile)
	if err != nil {
		klog.Errorf("NewBucket failed. Err: %v", err)
		return nil, err
	}
	github, err := github.NewManager(byFile, mgr)
	if err != nil {
		klog.Errorf("NewManager failed. Err: %v", err)
		return nil, err
	}
	kapp, err := kapp.NewKapp(byFile)
	if err != nil {
		klog.Errorf("NewYtt failed. Err: %v", err)
		return nil, err
	}

	mgr.gh = github
	mgr.b = bucket
	mgr.kapp = kapp

	return mgr, nil
}

// RawMetadata grabs the raw metadata
func (m *Manager) RawMetadata() ([]byte, error) {

	klog.V(2).Infof("Calling RawMetadata...")

	var err error
	byMeta, err := m.b.FetchMetadata()
	if err != nil {
		klog.Errorf("FetchMetadata failed. Err: %v", err)
		return nil, err
	}

	// just test if meta is ok
	m.metadata = &types.Metadata{
		ExtensionLookup: make(map[string]*types.Extension),
	}
	err = yaml.Unmarshal(byMeta, &m.metadata)
	if err != nil {
		klog.Errorf("Unmarshal failed. Err: ", err)
		return nil, err
	}
	// just test if meta is ok

	klog.V(2).Info("RawMetadata succeeded")

	return byMeta, nil
}

// InitMetadata grabs all the metadata relating to extensions
func (m *Manager) InitMetadata() (*types.Metadata, error) {
	// check for cache
	if m.metadata != nil {
		klog.V(2).Infof("Returned cached metadata")
		klog.V(4).Infof("%v", m.metadata)
		return m.metadata, nil
	}

	klog.V(2).Infof("Get metadata first time...")

	var err error
	byMeta, err := m.b.FetchMetadata()
	if err != nil {
		klog.Errorf("FetchMetadata failed. Err: %v", err)
		return nil, err
	}

	m.metadata = &types.Metadata{
		ExtensionLookup: make(map[string]*types.Extension),
	}
	err = yaml.Unmarshal(byMeta, m.metadata)
	if err != nil {
		klog.Errorf("Unmarshal failed. Err: ", err)
		return nil, err
	}
	for i := 0; i < len(m.metadata.Extensions); i++ {
		extension := m.metadata.Extensions[i]
		klog.V(6).Infof("extension.Name = %s", extension.Name)
		m.metadata.ExtensionLookup[extension.Name] = &extension
	}

	klog.V(2).Info("InitMetadata succeeded")
	klog.V(4).Infof("%v", m.metadata)

	return m.metadata, nil
}

// InitRelease grabs all the release relating to extensions
func (m *Manager) InitRelease() (*types.Release, error) {
	// check for cache
	if m.release != nil {
		klog.V(2).Infof("Returned cached release")
		klog.V(4).Infof("%v", m.release)
		return m.release, nil
	}

	klog.V(2).Infof("Get release first time...")

	var err error
	byRel, err := m.b.FetchRelease()
	if err != nil {
		klog.Errorf("FetchRelease failed. Err: %v", err)
		return nil, err
	}

	m.release = &types.Release{
		VersionLookup: make(map[string]*types.Version),
	}
	err = yaml.Unmarshal(byRel, m.release)
	if err != nil {
		klog.Errorf("Unmarshal failed. Err: ", err)
		return nil, err
	}
	for i := 0; i < len(m.release.Versions); i++ {
		version := m.release.Versions[i]
		klog.V(6).Infof("version.Version = %s", version.Version)
		m.release.VersionLookup[version.Version] = &version
	}

	klog.V(2).Info("InitMetadata succeeded")
	klog.V(4).Infof("%v", m.release)

	return m.release, nil
}
