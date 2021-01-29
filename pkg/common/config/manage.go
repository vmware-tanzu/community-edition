// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/vmware-tanzu/tce/pkg/common/types"
	klog "k8s.io/klog/v2"
)

// GetRaw - gets the actively used release/version
func (c *Config) GetRaw() []byte {
	return c.byRaw
}

// GetToken - gets the github token
func (c *Config) GetToken() string {
	if c.githubToken != "" {
		return c.githubToken
	}

	return c.ReadToken()
}

// ReadToken - from config file
func (c *Config) ReadToken() string {

	file, err := os.OpenFile(c.configFile, os.O_RDONLY, 0755)
	if err != nil {
		klog.Errorf("Open Config for read failed. Err: %v", err)
		return ""
	}
	defer file.Close()

	dataReader := bufio.NewReader(file)
	if dataReader == nil {
		klog.Errorf("Datareader creation failed")
		return ""
	}

	for {
		line, err := dataReader.ReadString('\n')

		if err == io.EOF {
			klog.V(6).Infof("ReadString error. Err: EOF")
			break
		}
		if err != nil {
			klog.Errorf("ReadString error. Err: %v", err)
			break
		}

		if strings.Contains(line, "token:") {
			klog.V(4).Info("found version line")
			myStrings := strings.Split(line, " ")
			if len(myStrings) == 2 {
				c.githubToken = myStrings[1]
				klog.V(4).Infof("version = %s", c.githubToken)
				break
			}
		}
	}

	klog.V(6).Infof("token = %s", c.githubToken)
	return c.githubToken
}

// UpdateToken - update token
func (c *Config) UpdateToken(token string) error {
	return c.updateField(KeyToken, token)
}

// GetRelease - gets the actively used release/version
func (c *Config) GetRelease() (string, error) {

	file, err := os.OpenFile(c.configFile, os.O_RDONLY, 0755)
	if err != nil {
		klog.Errorf("Open Config for read failed. Err: %v", err)
		return "", err
	}
	defer file.Close()

	dataReader := bufio.NewReader(file)
	if dataReader == nil {
		klog.Errorf("Datareader creation failed")
		return "", ErrDatareaderFailed
	}

	var myVersion string
	for {
		line, err := dataReader.ReadString('\n')

		if err == io.EOF {
			klog.V(6).Infof("ReadString error. Err: EOF")
			break
		}
		if err != nil {
			klog.Errorf("ReadString error. Err: %v", err)
			break
		}

		if strings.Contains(line, "version:") {
			klog.V(4).Info("found version line")
			myStrings := strings.Split(line, " ")
			if len(myStrings) == 2 {
				myVersion = myStrings[1]
				klog.V(4).Infof("version = %s", myVersion)
				break
			}
		}
	}

	if len(myVersion) == 0 {
		klog.Error("Invalid version")
		return "", types.ErrVersionNotFound
	}

	klog.V(2).Infof("Current release = %s", myVersion)
	return myVersion, nil
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

	// do the replace
	found := false
	for {
		line, err := dataReader.ReadString('\n')
		klog.V(4).Infof("line: %s", line)

		if err == io.EOF {
			klog.V(6).Infof("ReadString error. Err: EOF")
			break
		}
		if err != nil {
			klog.Errorf("ReadString error. Err: %v", err)
			break
		}

		lineToWrite := line
		if strings.Contains(line, key) {
			klog.V(4).Infof("Replacing line!")
			found = true
			lineToWrite = key + " " + value + "\n"
		}

		if len(value) == 0 {
			klog.V(4).Infof("Clear %s value", key)
			found = true
			continue
		}

		klog.V(4).Infof("lineToWrite = %s", lineToWrite)
		_, err = datawriter.WriteString(lineToWrite)
		if err != nil {
			klog.Errorf("Fail to write line. Err: %v", err)
			return err
		}
		datawriter.Flush()
	}

	if !found {
		klog.V(4).Infof("%s was not found. Writing line!", key)
		lineToWrite := key + " " + value + "\n"
		klog.V(4).Infof("lineToWrite = %s", lineToWrite)
		_, err = datawriter.WriteString(lineToWrite)
		if err != nil {
			klog.Errorf("Fail to write line. Err: %v", err)
			return err
		}
		datawriter.Flush()
	}

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

	err = types.CopyFile(tmpFile, c.configFile)
	if err != nil {
		klog.Errorf("CopyFile failed. Err: %v", err)
		return err
	}

	klog.V(2).Infof("updateField succeeded!")
	return nil
}
