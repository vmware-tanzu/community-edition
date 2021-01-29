// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package download

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	klog "k8s.io/klog/v2"
)

// DownloadGitHubFile - downloads a file
func (m *Manager) DownloadGitHubFile(fromURI string, toDirFile string, token string) error {

	klog.V(2).Infof("fromURI: %s", fromURI)
	klog.V(2).Infof("toDirFile: %s", toDirFile)
	klog.V(6).Infof("token: %s", token)

	if _, err := os.Stat(toDirFile); os.IsNotExist(err) {
		klog.V(4).Infof("Local file does not exist. Downloading from repo.")
		toDir := filepath.Dir(toDirFile)

		err := os.MkdirAll(toDir, 0755)
		if err != nil {
			klog.Errorf("MkdirAll failed. Err: %v", err)
			return err
		}

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)

		opts := &github.RepositoryContentGetOptions{}
		fileGH, _, _, err := client.Repositories.GetContents(ctx, "vmware-tanzu", "tce", fromURI, opts)
		if err != nil {
			klog.Errorf("client.Repositories failed. Err: %v", err)
			return err
		}

		klog.V(2).Infof("Name: %s", *fileGH.Name)
		klog.V(6).Infof("DownloadURL: %s", *fileGH.DownloadURL)

		response, err := http.Get(*fileGH.DownloadURL)
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

	klog.V(2).Info("DownloadGitHubFile succeeded")
	return nil
}

// DownloadGitHubFileToDir - downloads a file
func (m *Manager) DownloadGitHubFileToDir(fromURI string, toDir string, token string) error {

	klog.V(2).Infof("fromURI: %s", fromURI)
	klog.V(2).Infof("toDir: %s", toDir)

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

	return m.DownloadGitHubFile(fromURI, toDirFile, token)
}

// DownloadGitHubFilesToDir - download files
func (m *Manager) DownloadGitHubFilesToDir(fromURIDir string, fromFiles []string, toDir string, token string) error {

	klog.V(4).Infof("fromURIDir: %s", fromURIDir)
	for _, file := range fromFiles {
		klog.V(4).Infof("file: %s", file)
	}
	klog.V(4).Infof("toDir: %s", toDir)

	for _, fromFile := range fromFiles {
		dirfileURI := fromURIDir + "/" + fromFile

		toDirFile := filepath.Join(toDir, fromFile)
		toNewDirFile := filepath.Dir(toDirFile)

		err := m.DownloadGitHubFileToDir(dirfileURI, toNewDirFile, token)
		if err != nil {
			klog.Errorf("Failed to download %s -> %s", dirfileURI, toNewDirFile)
			return err
		}
	}

	return nil
}

// PrintGitHubFile - prints a file
func (m *Manager) PrintGitHubFile(fromURI string, toDirFile string, token string) error {

	klog.V(2).Infof("fromURI: %s", fromURI)
	klog.V(2).Infof("toDirFile: %s", toDirFile)

	if _, err := os.Stat(toDirFile); os.IsNotExist(err) {

		klog.V(6).Infof("File missing. Download file from URI")
		err := m.DownloadGitHubFile(fromURI, toDirFile, token)
		if err != nil {
			klog.Errorf("DownloadFile failed. Err: %v", err)
			return err
		}
	}

	klog.V(6).Infof("Print local copy of file")
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

// PrintGitHubFiles - prints files
func (m *Manager) PrintGitHubFiles(fromURIDir string, fromFiles []string, toDir string, token string) error {

	klog.V(4).Infof("fromURIDir: %s", fromURIDir)
	for _, file := range fromFiles {
		klog.V(4).Infof("file: %s", file)
	}
	klog.V(4).Infof("toDir: %s", toDir)

	for _, fromFile := range fromFiles {
		dirfileURI := fromURIDir + "/" + fromFile

		toDirFile := filepath.Join(toDir, fromFile)

		err := m.PrintGitHubFile(dirfileURI, toDirFile, token)
		if err != nil {
			klog.Errorf("Failed to print %s", dirfileURI)
			return err
		}
	}

	return nil
}
