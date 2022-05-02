// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tanzu

import "testing"

// This test checks that the exit codes are their expected integer.
// These ints should not change without going through the lifecycle processes
// since users may rely on them in their ci/cd for fallback and automation purposes

//nolint:gocyclo
func TestExitCodes(t *testing.T) {
	if Success != 0 {
		t.Errorf("expected success exit code to be 0. Actual: %v", Success)
	}

	if InvalidConfig != 1 {
		t.Errorf("expected InvalidConfig exit code to be 1. Actual: %v", InvalidConfig)
	}

	if ErrCreatingClusterDirs != 2 {
		t.Errorf("expected ErrCreatingClusterDirs exict code to be 2. Actual: %v", ErrCreatingClusterDirs)
	}

	if ErrTkrBom != 3 {
		t.Errorf("expected ErrTkrBom exit code to be 3. Actual: %v", ErrTkrBom)
	}

	if ErrRenderingConfig != 4 {
		t.Errorf("expected ErrRenderingConfig exit code to be 4. Actual: %v", ErrRenderingConfig)
	}

	if ErrTkrBomParsing != 5 {
		t.Errorf("expected ErrTkrBomParsing exit code to be 5. Actual: %v", ErrTkrBomParsing)
	}

	if ErrKappBundleResolving != 6 {
		t.Errorf("expected ErrKappBundleResolving exit code to be 6. Actual: %v", ErrKappBundleResolving)
	}

	if ErrCreateCluster != 7 {
		t.Errorf("expected ErrCreateCluster exit code to be 7. Actual: %v", ErrCreateCluster)
	}

	if ErrExistingCluster != 8 {
		t.Errorf("expected ErrExistingCluster exit code to be 8. Actual: %v", ErrExistingCluster)
	}

	if ErrKappInstall != 9 {
		t.Errorf("expected ErrKappInstall exit code to be 9. Actual: %v", ErrKappInstall)
	}

	if ErrCorePackageRepoInstall != 10 {
		t.Errorf("expected ErrCorePackageRepoInstall exit code to be 10. Actual: %v", ErrCorePackageRepoInstall)
	}

	if ErrOtherPackageRepoInstall != 11 {
		t.Errorf("expected ErrOtherPackageRepoInstall exit code to be 11. Actual: %v", ErrOtherPackageRepoInstall)
	}

	if ErrCniInstall != 12 {
		t.Errorf("expected ErrCniInstall exit code to be 12. Actual %v", ErrCniInstall)
	}

	if ErrKubeconfigContextSet != 13 {
		t.Errorf("expected ErrKubeconfigContextSet exit code to be 13. Actual %v", ErrKubeconfigContextSet)
	}

	if ErrInstallPackage != 14 {
		t.Errorf("expected ErrInstallPackage exit code to be 14. Actual %v", ErrInstallPackage)
	}
}
