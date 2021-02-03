// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"bufio"
	"os"
	"strings"

	yaml "github.com/ghodss/yaml"
	klog "k8s.io/klog/v2"

	"github.com/vmware-tanzu-private/core/pkg/v1/cli"

	"github.com/vmware-tanzu/tce/pkg/common/types"
)

// GetToken - gets the github token
func (c *Config) GetToken() string {
	return c.GithubToken
}

// UpdateToken - update token
func (c *Config) UpdateToken(token string) error {
	return c.updateField(KeyToken, token)
}

// GetRelease - gets the actively used release/version
func (c *Config) GetRelease() (string, error) {
	if c.ReleaseVersion == "" {
		version := cli.BuildVersion
		err := c.SetRelease(version)
		if err != nil {
			return "", err
		}
		return version, nil
	}

	return c.ReleaseVersion, nil
}

// SetRelease - sets the actively used release/version
func (c *Config) SetRelease(version string) error {
	return c.updateField(KeyUpdate, version)
}

func (c *Config) updateField(key string, value string) error {

	klog.V(2).Infof("key = %s", key)
	klog.V(6).Infof("value = %s", value)

	tmpFile := c.configFile + ".tmp"
	klog.V(4).Infof("configFile = %s", c.configFile)
	klog.V(4).Infof("tmpFile = %s", tmpFile)

	err := types.CopyFile(c.configFile, tmpFile)
	if err != nil {
		klog.Errorf("CopyFile failed. Err: %v", err)
		return err
	}
	defer os.RemoveAll(tmpFile)

	// read file
	fileRead, err := os.OpenFile(c.configFile, os.O_RDONLY, 0755)
	if err != nil {
		klog.Errorf("Open Config for read failed. Err: %v", err)
		return err
	}

	dataReader := bufio.NewReader(fileRead)
	if dataReader == nil {
		klog.Errorf("Datareader creation failed")
		return ErrDatareaderFailed
	}

	//write file
	fileWrite, err := os.OpenFile(tmpFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		klog.Errorf("Open Config for write failed. Err: %v", err)
		return err
	}

	datawriter := bufio.NewWriter(fileWrite)
	if datawriter == nil {
		klog.Errorf("Datawriter creation failed")
		return ErrDatawriterFailed
	}

	// replace
	if strings.Contains(key, KeyUpdate) {
		c.ReleaseVersion = value
	} else if strings.Contains(key, KeyToken) {
		c.GithubToken = value
	}
	klog.V(4).Infof("Config{} = %s", c)

	byRaw, err := yaml.Marshal(c)
	if err != nil {
		klog.V(2).Infof("yaml.Marshal error. Err: %v", err)
		return err
	}
	klog.V(6).Infof("byRaw = %v", byRaw)

	_, err = datawriter.Write(byRaw)
	if err != nil {
		klog.V(2).Infof("datawriter.Write error. Err: %v", err)
		return err
	}
	datawriter.Flush()

	// close everything
	err = fileRead.Close()
	if err != nil {
		klog.Errorf("fileRead.Close failed. Err: %v", err)
		return err
	}

	err = fileWrite.Close()
	if err != nil {
		klog.Errorf("fileWrite.Close failed. Err: %v", err)
		return err
	}

	// switch files
	err = types.CopyFile(tmpFile, c.configFile)
	if err != nil {
		klog.Errorf("CopyFile failed. Err: %v", err)
		return err
	}

	klog.V(2).Infof("updateField succeeded!")
	return nil
}
