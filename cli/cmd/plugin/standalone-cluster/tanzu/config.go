// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tanzu

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	TKRLocation string = "TkrLocation"
	Provider    string = "Provider"
	CNI         string = "Cni"
	PodCIDR     string = "PodCidr"
	ServiceCIDR string = "ServiceCidr"
)

var defaultConfigValues = map[string]string{
	TKRLocation: "projects.registry.vmware.com/tkg/tkr-bom:v1.21.2_vmware.1-tkg.1",
	Provider:    "kind",
	CNI:         "antrea",
	PodCIDR:     "10.244.0.0/16",
	ServiceCIDR: "10.96.0.0/16",
}

// PortMap is the mapping between a host port and a container port.
type PortMap struct {
	// HostPort is the port on the host machine.
	HostPort int `yaml:"HostPort"`
	// ContainerPort is the port on the container to map to.
	ContainerPort int `yaml:"ContainerPort"`
}

// LocalClusterConfig contains all the configuration settings for creating a
// local Tanzu cluster.
type LocalClusterConfig struct {
	// ClusterName is the name of the cluster.
	ClusterName string `yaml:"ClusterName"`
	// Provider is the local infrastructure provider to use (e.g. kind).
	Provider string `yaml:"Provider"`
	// ProviderConfiguration offers optional provider-specific configuration.
	// The exact keys and values accepted are determined by the provider.
	ProviderConfiguration map[string]interface{} `yaml:"ProviderConfiguration"`
	// CNI is the networking CNI to use in the cluster. Default is antrea.
	CNI string `yaml:"Cni"`
	// CNIConfiguration offers optional cni-plugin specific configuration.
	// The exact keys and values accepted are determined by the CNI choice.
	CNIConfiguration map[string]interface{} `yaml:"CniConfiguration"`
	// PodCidr is the Pod CIDR range to assign pod IP addresses.
	PodCidr string `yaml:"PodCidr"`
	// ServiceCidr is the Service CIDR range to assign service IP addresses.
	ServiceCidr string `yaml:"ServiceCidr"`
	// TkrLocation is the path to the Tanzu Kubernetes Release (TKR) data.
	TkrLocation string `yaml:"TkrLocation"`
	// PortsToForward contains a mapping of host to container ports that should
	// be exposed.
	PortsToForward []PortMap `yaml:"PortsToForward"`
}

// KubeConfigPath gets the full path to the KubeConfig for this local cluster.
func (lcc *LocalClusterConfig) KubeConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), configDir, tanzuConfigDir, lcc.ClusterName+".yaml")
}

// InitializeConfiguration determines the configuration to use for cluster creation.
//
// There are three places where configuration comes from:
// - default settings
// - configuration file
// - environment variables
// - command line arguments
//
// The effective configuration is determined by combining these sources, in ascending
// order of preference listed. So env variables override values in the config file,
// and explicit CLI arguments override config file and env variable values.
func InitializeConfiguration(commandArgs map[string]string) (*LocalClusterConfig, error) {
	config := &LocalClusterConfig{}

	// First, populate values based on a supplied config file
	if commandArgs["clusterconfigfile"] != "" {
		configData, err := os.ReadFile(commandArgs["clusterconfigfile"])
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(configData, config)
		if err != nil {
			return nil, err
		}
	}

	// Loop through and look up each field
	element := reflect.ValueOf(config).Elem()
	for i := 0; i < element.NumField(); i++ {
		field := element.Type().Field(i)
		if field.Type.Kind() != reflect.String {
			// Not supporting more complex data types yet, will need to see if and
			// how to do this.
			continue
		}

		// Use the yaml name if provided so it matches what we serialize to file
		fieldName := field.Tag.Get("yaml")
		if fieldName == "" {
			fieldName = field.Name
		}

		// Check if an explicit value was passed in
		if value, ok := commandArgs[strings.ToLower(fieldName)]; ok {
			element.FieldByName(field.Name).SetString(value)
		} else if value := os.Getenv(fieldNameToEnvName(fieldName)); value != "" {
			// See if there is an environment variable set for this field
			element.FieldByName(field.Name).SetString(value)
		}

		// Only set to the default value if it hasn't been set already
		if element.FieldByName(field.Name).String() == "" {
			if value, ok := defaultConfigValues[fieldName]; ok {
				element.FieldByName(field.Name).SetString(value)
			}
		}
	}

	// Make sure cluster name was either set on the command line or in the config
	// file.
	if config.ClusterName == "" {
		return nil, fmt.Errorf("cluster name must be provided")
	}

	return config, nil
}

// fieldNameToEnvName converts the config values yaml name to its expected env
// variable name.
func fieldNameToEnvName(field string) string {
	namedArray := []string{"TANZU"}
	re := regexp.MustCompile(`[A-Z][^A-Z]*`)
	allWords := re.FindAllString(field, -1)
	for _, word := range allWords {
		namedArray = append(namedArray, strings.ToUpper(word))
	}
	return strings.Join(namedArray, "_")
}
