// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package download

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	klog "k8s.io/klog/v2"
)

// DownloadFile - downloads a file
func (m *Manager) DownloadFile(fromURI string, toDirFile string) error {

	klog.V(2).Infof("fromURI: %s", fromURI)
	klog.V(2).Infof("toDirFile: %s", toDirFile)

	if _, err := os.Stat(toDirFile); os.IsNotExist(err) {
		klog.V(4).Infof("Local file does not exist. Downloading from repo.")
		toDir := filepath.Dir(toDirFile)

		err := os.MkdirAll(toDir, 0755)
		if err != nil {
			klog.Errorf("MkdirAll failed. Err: %v", err)
			return err
		}

		response, err := http.Get(fromURI)
		if err != nil {
			klog.Errorf("Http Get failed. Err: %v", err)
			return err
		}
		defer response.Body.Close()

		file, err := os.Create(toDirFile)
		if err != nil {
			klog.Errorf("Create file failed. Err: %v", err)
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, response.Body)
		if err != nil {
			klog.Errorf("Copy bits failed. Err: %v", err)
			return err
		}
	} else {
		klog.V(4).Infof("Local file already exists!")
	}

	klog.V(2).Info("DownloadFile succeeded")
	return nil
}

// DownloadFileToDir - downloads a file
func (m *Manager) DownloadFileToDir(fromURI string, toDir string) error {

	klog.V(4).Infof("fromURI: %s", fromURI)
	klog.V(4).Infof("toDir: %s", toDir)

	url, err := url.Parse(fromURI)
	if err != nil {
		klog.Errorf("url.Parse failed. Err: %v", err)
		return err
	}

	ss := strings.Split(url.Path, "/")
	filename := ss[len(ss)-1]
	klog.V(4).Infof("Filename: %s", filename)

	toDirFile := filepath.Join(toDir, filename)
	klog.V(4).Infof("toDirFile: %s", toDirFile)

	return m.DownloadFile(fromURI, toDirFile)
}

// DownloadFilesToDir - download files
func (m *Manager) DownloadFilesToDir(fromURIDir string, fromFiles []string, toDir string) error {

	klog.V(6).Infof("fromURIDir: %s", fromURIDir)
	for _, file := range fromFiles {
		klog.V(6).Infof("file: %s", file)
	}
	klog.V(6).Infof("toDir: %s", toDir)

	for _, fromFile := range fromFiles {
		dirfileURI := fromURIDir + "/" + fromFile

		toDirFile := filepath.Join(toDir, fromFile)
		toNewDirFile := filepath.Dir(toDirFile)

		err := m.DownloadFileToDir(dirfileURI, toNewDirFile)
		if err != nil {
			klog.Errorf("Failed to download %s -> %s", dirfileURI, toNewDirFile)
			return err
		}
	}

	klog.V(6).Infof("DownloadFilesToDir succeeded")

	return nil
}

// PrintFile - prints a file
func (m *Manager) PrintFile(fromURI string, toDirFile string) error {

	klog.V(2).Infof("fromURI: %s", fromURI)
	klog.V(2).Infof("toDirFile: %s", toDirFile)

	if _, err := os.Stat(toDirFile); os.IsNotExist(err) {

		klog.V(2).Infof("File missing. Download file from URI")
		err := m.DownloadFile(fromURI, toDirFile)
		if err != nil {
			klog.Errorf("DownloadFile failed. Err: %v", err)
			return err
		}
	}

	klog.V(2).Infof("Print local copy of file")
	byFile, err := ioutil.ReadFile(toDirFile)
	if err != nil {
		klog.Errorf("ReadAll failed. Err: %v", err)
		return err
	}

	fmt.Printf("File: %s\n", fromURI)
	fmt.Printf("-----------------------------------------------\n")
	fmt.Printf("%s\n\n", string(byFile))

	return nil
}

// PrintFiles - prints files
func (m *Manager) PrintFiles(fromURIDir string, fromFiles []string, toDir string) error {

	klog.V(4).Infof("fromURIDir: %s", fromURIDir)
	for _, file := range fromFiles {
		klog.V(4).Infof("file: %s", file)
	}
	klog.V(4).Infof("toDir: %s", toDir)

	for _, fromFile := range fromFiles {
		dirfileURI := fromURIDir + "/" + fromFile

		toDirFile := filepath.Join(toDir, fromFile)

		err := m.PrintFile(dirfileURI, toDirFile)
		if err != nil {
			klog.Errorf("Failed to print %s", dirfileURI)
			return err
		}
	}

	klog.V(4).Infof("PrintFiles succeeded")
	return nil
}
