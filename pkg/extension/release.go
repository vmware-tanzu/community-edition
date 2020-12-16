// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package extension

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	klog "k8s.io/klog/v2"
)

var printRelease bool
var releaseName string

// ReleaseCmd represents the release command
var ReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Manage releases for extensions",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		mgr, err = NewManager()
		return err
	},
	RunE: release,
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func init() {
	ReleaseCmd.Flags().BoolVarP(&printRelease, "list", "l", false, "List releases")
	ReleaseCmd.Flags().StringVarP(&releaseName, "set", "s", "", "Set release version")
}

func currentRelease(cmd *cobra.Command, args []string) error {

	version, err := mgr.cfg.GetRelease()
	if err == nil {
		fmt.Printf("Version: %s\n", version)
	} else {
		fmt.Printf("Version not found. Err: %v\n", err)
	}

	return nil
}

func listReleases(cmd *cobra.Command, args []string) error {
	release, err := mgr.InitRelease()
	if err != nil {
		klog.Errorf("InitRelease failed. Err: %v", err)
		return err
	}

	// Or print the version
	fmt.Printf("Stable version: %s\n", release.Stable)
	if release.Date != "" {
		fmt.Printf("Stable date: %s\n", release.Date)
	}
	fmt.Printf("\n")

	// List all releases
	for _, version := range release.Versions {
		if version.Date != "" {
			fmt.Printf("Version (%s): %s\n", version.Date, version.Version)
		} else {
			fmt.Printf("Version: %s\n", version.Version)
		}
	}

	return nil
}

func setRelease(cmd *cobra.Command, args []string) error {
	release, err := mgr.InitRelease()
	if err != nil {
		klog.Errorf("InitRelease failed. Err: %v", err)
		return err
	}

	// check to see if version is valid
	if !strings.EqualFold(releaseName, DefaultReleaseLatest) && !strings.EqualFold(releaseName, DefaultReleaseStable) {
		_, err := release.GetVersion(releaseName)
		if err != nil {
			klog.Errorf("Invalid version")
			return err
		}
	}

	if releaseName == DefaultReleaseLatest {
		klog.V(2).Infof("Setting to %s", releaseName)
	} else if releaseName == DefaultReleaseStable {
		klog.V(2).Infof("Setting to %s", releaseName)
		releaseName = release.Stable
	}

	err = mgr.cfg.SetRelease(releaseName)
	if err != nil {
		fmt.Printf("Failed to update to %s\n", releaseName)
		return err
	}

	fmt.Printf("Updated to %s\n", releaseName)
	return nil
}

func release(cmd *cobra.Command, args []string) error {

	// List releases
	if printRelease {
		return listReleases(cmd, args)
	}

	// Set release
	if releaseName != "" {
		return setRelease(cmd, args)
	}

	// Print current release
	return currentRelease(cmd, args)
}
