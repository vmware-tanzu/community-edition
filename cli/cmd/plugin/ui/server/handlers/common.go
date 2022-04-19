// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	yaml "gopkg.in/yaml.v3"

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
