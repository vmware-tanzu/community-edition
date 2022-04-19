// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"time"

	yaml "gopkg.in/yaml.v3"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/client"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgconfigupdater"
)

const (
	trueStr  = "true"
	falseStr = "false"
)

// transformConfigToString provides a generic way to convert a cluster config
// into a string.
func transformConfigToString(config interface{}) (out string, err error) {
	var configMap map[string]string
	var configByte []byte

	// turn the configuration object into a map
	configMap, err = tkgconfigupdater.CreateConfigMap(config)
	if err == nil {
		// turn the map into a byte array
		configByte, err = yaml.Marshal(&configMap)
	}
	if err == nil {
		return string(configByte), nil
	}
	return "", err
}

// createManagementCluster can be called by all provider endpoints. The InitRegionOptions
// need to be initialized with the configuration options for the given provider. This
// is a common create call for any provider. It assumes log streaming has been set up
// prior to invocation using something like `go app.StartSendingLogsToUI()`.
func createManagementCluster(tkgClient *client.TkgClient, initOptions *client.InitRegionOptions) {
	err := tkgClient.InitRegion(initOptions)
	if err != nil {
		log.Error(err, "unable to set up management cluster, ")
	} else {
		log.Infof("\nManagement cluster created!\n\n")
		log.Info("\nYou can now create your first workload cluster by running the following:\n\n")
		log.Info("  tanzu cluster create [name] -f [file]\n\n")
		// wait for the logs to be dispatched to UI before exit
		time.Sleep(sleepTimeForLogsPropogation)
	}
}
