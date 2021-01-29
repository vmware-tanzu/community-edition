// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package github

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	klog "k8s.io/klog/v2"

	download "github.com/vmware-tanzu/tce/pkg/common/download"
	types "github.com/vmware-tanzu/tce/pkg/common/types"
)

// NewManager generates a Manager object
func NewManager(byConfig []byte, mgr types.IManager) (*Manager, error) {

	cfg, err := InitGitHubConfig(byConfig)
	if err != nil {
		klog.Errorf("InitGitHubConfig failed. Err: %v", err)
		return nil, err
	}

	dl, err := download.NewManager(byConfig)
	if err != nil {
		klog.Errorf("download.NewManager failed. Err: %v", err)
		return nil, err
	}

	m := &Manager{
		cfg:                    cfg,
		dl:                     dl,
		extensionDirectoryRoot: filepath.Join(xdg.DataHome, "tanzu-repository", cfg.ExtensionDirectory),
		extensionDirectory:     filepath.Join(xdg.DataHome, "tanzu-repository", cfg.ExtensionDirectory, cfg.originalBranchTag),
		extMgr:                 mgr,
	}

	klog.V(4).Infof("extensionDirectory = %s", m.extensionDirectory)

	return m, nil
}

// GetInternalRepo true if internal/private repo
func (m *Manager) GetInternalRepo() bool {
	return m.cfg.Token != ""
}

// DownloadExtension download extension on github
func (m *Manager) DownloadExtension(extensionName string) error {

	klog.V(6).Infof("DownloadExtension(%s)", extensionName)

	metadata, err := m.extMgr.InitMetadata()
	if err != nil {
		klog.Errorf("InitMetadata failed. Err: %v", err)
		return err
	}

	extension, err := metadata.GetExtension(extensionName)
	if err != nil {
		klog.Errorf("Extension %s not found", extensionName)
		return err
	}
	klog.V(2).Infof("Extension found %s", extension.Name)

	if m.GetInternalRepo() {
		klog.V(2).Infof("Using internal GitHub repo")
		extensionFolder := filepath.Join(m.cfg.ExtensionDirectory, extension.Name)
		extensionDirectory := filepath.Join(m.extensionDirectory, extension.Name)

		return m.dl.DownloadGitHubFilesToDir(extensionFolder, extension.GetFileList(), extensionDirectory, m.cfg.Token)
	}

	klog.V(2).Infof("Using external GitHub repo")
	extensionFolder := m.cfg.GitHubURI + "/" + extension.Name
	extensionDirectory := filepath.Join(m.extensionDirectory, extension.Name)

	return m.dl.DownloadFilesToDir(extensionFolder, extension.GetFileList(), extensionDirectory)
}

// PrintExtension print extension on github
func (m *Manager) PrintExtension(extensionName string) error {

	klog.V(6).Infof("PrintExtension(%s)", extensionName)

	metadata, err := m.extMgr.InitMetadata()
	if err != nil {
		klog.Errorf("InitMetadata failed. Err: %v", err)
		return err
	}

	extension, err := metadata.GetExtension(extensionName)
	if err != nil {
		klog.Errorf("Extension %s not found", extensionName)
		return err
	}
	klog.V(2).Infof("Extension found %s", extension.Name)

	if m.GetInternalRepo() {
		klog.V(2).Infof("Using internal GitHub repo")
		extensionFolder := filepath.Join(m.cfg.ExtensionDirectory, extension.Name)
		extensionDirectory := filepath.Join(m.extensionDirectory, extension.Name)

		return m.dl.PrintGitHubFiles(extensionFolder, extension.GetFileList(), extensionDirectory, m.cfg.Token)
	}

	klog.V(2).Infof("Using external GitHub repo")
	extensionFolder := m.cfg.GitHubURI + "/" + extension.Name
	extensionDirectory := filepath.Join(m.extensionDirectory, extension.Name)

	return m.dl.PrintFiles(extensionFolder, extension.GetFileList(), extensionDirectory)
}

// StageFiles stages files for use
func (m *Manager) StageFiles(workingDir string, extensionName string) error {

	extensionFileCheck := filepath.Join(m.extensionDirectory, extensionName, types.DefaultAppCrdFilename)

	if _, err := os.Stat(extensionFileCheck); os.IsNotExist(err) {
		klog.V(4).Infof("Local file does not exist. Downloading from repo.")

		err = m.DownloadExtension(extensionName)
		if err != nil {
			klog.Errorf("DownloadFilesToDir failed. Err: %v", err)
			return err
		}
	}

	//remove old working/plugin dir
	err := os.RemoveAll(workingDir)
	if err != nil {
		klog.Errorf("RemoveAll failed. Err: %v", err)
		return err
	}

	// copy to working directory
	srcDir := filepath.Dir(extensionFileCheck)

	klog.V(4).Infof("srcDir = %s", srcDir)
	klog.V(4).Infof("dstDir = %s", workingDir)

	err = os.MkdirAll(workingDir, 0755)
	if err != nil {
		klog.Errorf("MkdirAll failed. Err: %v", err)
		return err
	}

	return types.RecursiveCopy(srcDir, workingDir)
}

// Reset all directories to "factory"
func (m *Manager) Reset(workingDir string) error {

	klog.V(2).Infof("workingDir = %s", workingDir)
	klog.V(2).Infof("extensionDirectory = %s", m.extensionDirectoryRoot)

	//remove working dir
	err := os.RemoveAll(workingDir)
	if err != nil {
		klog.Errorf("RemoveAll(workingDir) failed. Err: %v", err)
		return err
	}

	//remove extension dir
	err = os.RemoveAll(m.extensionDirectoryRoot)
	if err != nil {
		klog.Errorf("RemoveAll(extension) failed. Err: %v", err)
		return err
	}

	klog.V(4).Infof("Reset succeessful")
	return nil
}
