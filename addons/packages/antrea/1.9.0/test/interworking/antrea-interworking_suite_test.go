// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package interworking_test

import (
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

type AntreaInterworkingConfig struct {
	AntreaNSX struct {
		Enable bool `yaml:"enable"`
	} `yaml:"antrea_nsx"`
	AntreaInterworking struct {
		Config struct {
			NsxCert       string   `yaml:"nsxCert"`
			NSXKey        string   `yaml:"nsxKey"`
			NSXUser       string   `yaml:"nsxUser"`
			NSXPassword   string   `yaml:"nsxPassword"`
			ClusterName   string   `yaml:"clusterName"`
			NSXManagers   []string `yaml:"NSXManagers"`
			vpcPath       string   `yaml:"vpcPath"`
			MpAdapterConf struct {
				NSXClientTimeout     int  `yaml:"NSXClientTimeout"`
				InventoryBatchSize   int  `yaml:"InventoryBatchSize"`
				InventoryBatchPeriod int  `yaml:"InventoryBatchPeriod"`
				EnableDebugServer    bool `yaml:"EnableDebugServer"`
				APIServerPort        int  `yaml:"APIServerPort"`
				DebugServerPort      int  `yaml:"DebugServerPort"`
				NSXRPCDebug          bool `yaml:"NSXRPCDebug"`
				ConditionTimeout     int  `yaml:"ConditionTimeout"`
			} `yaml:"mp_adapter_conf"`
			CCPAdapterConf struct {
				EnableDebugServer               bool    `yaml:"EnableDebugServer"`
				APIServerPort                   int     `yaml:"APIServerPort"`
				DebugServerPort                 int     `yaml:"DebugServerPort"`
				NSXRPCDebug                     bool    `yaml:"NSXRPCDebug"`
				RealizeTimeoutSeconds           int     `yaml:"RealizeTimeoutSeconds"`
				RealizeErrorSyncIntervalSeconds int     `yaml:"RealizeErrorSyncIntervalSeconds"`
				ReconcilerWorkerCount           int     `yaml:"ReconcilerWorkerCount"`
				ReconcilerQPS                   float32 `yaml:"ReconcilerQPS"`
				ReconcilerBurst                 int     `yaml:"ReconcilerBurst"`
				ReconcilerResyncSeconds         int     `yaml:"ReconcilerResyncSeconds"`
			} `yaml:"ccp_adapter_conf"`
		} `yaml:"config"`
	} `yaml:"antrea_interworking"`
}

func TestAntreaInterworking(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Antrea-interworking Addons Templates Suite")
}
