// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// CopyFile src to dst
func CopyFile(source, destination string) error {
	input, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destination, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

// RecursiveCopy src to dst
func RecursiveCopy(source, destination string) error {
	var err error = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		var relPath string = strings.Replace(path, source, "", 1)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			if _, err := os.Stat(destination); os.IsNotExist(err) {
				return os.MkdirAll(filepath.Join(destination, relPath), 0755)
			}
			return nil
		} else {
			newDir := filepath.Dir(relPath)
			if newDir != "/" {
				err := os.MkdirAll(filepath.Join(destination, newDir), 0755)
				if err != nil {
					return err
				}
			}

			var data, err1 = ioutil.ReadFile(filepath.Join(source, relPath))
			if err1 != nil {
				return err1
			}
			return ioutil.WriteFile(filepath.Join(destination, relPath), data, 0777)
		}
	})
	return err
}

// GetFileList as []string
func (e *Extension) GetFileList() []string {
	var files []string
	for _, file := range e.Files {
		files = append(files, file.Name)
	}
	return files
}

// GetExtension by name
func (m *Metadata) GetExtension(name string) (*Extension, error) {
	extension := m.ExtensionLookup[name]
	if extension == nil {
		return nil, ErrExtensionNotFound
	}
	return extension, nil
}

// GetVersion by name
func (r *Release) GetVersion(name string) (*Version, error) {
	version := r.VersionLookup[name]
	if version == nil {
		return nil, ErrVersionNotFound
	}
	return version, nil
}
