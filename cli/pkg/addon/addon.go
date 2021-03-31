// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	"flag"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/adrg/xdg"
	"k8s.io/klog/v2"

	"github.com/vmware-tanzu/tce/cli/pkg/common/kapp"
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
	_ = flag.Set("v", level)
	flag.Parse()

	// read config
	configFile := filepath.Join(xdg.DataHome, "tanzu-repository", DefaultConfigFile)
	byFile, err := os.ReadFile(configFile)
	if err != nil {
		klog.Errorf("ReadFile failed. Err: %v", err)
		return nil, err
	}

	mngr := &Manager{}

	kap, err := kapp.NewKapp(byFile)
	if err != nil {
		klog.Errorf("NewYtt failed. Err: %v", err)
		return nil, err
	}

	mngr.kapp = kap

	return mngr, nil
}
