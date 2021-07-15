// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type PackageConfiguration struct {
	Packages []AddonsConfig `yaml:"packages"`
}

type AddonsConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

func readPackageConfig(filename string) (*PackageConfiguration, error) {
	cfg, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &PackageConfiguration{}
	err = yaml.Unmarshal(cfg, &c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}
