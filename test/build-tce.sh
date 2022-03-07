#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script clones TCE repo and builds the latest release

set -x
set -e

TCE_REPO_PATH="$(git rev-parse --show-toplevel)"
# shellcheck source=test/util/utils.sh
source "${TCE_REPO_PATH}/test/util/utils.sh"

BUILD_OS=$(uname -s)
export BUILD_OS

# Build TCE
echo "Building TCE release..."
make release || { error "TCE BUILD FAILED!"; exit 1; }
echo "Installing TCE release"
if [[ $BUILD_OS == "Linux" ]]; then
    pushd release/tce-linux-amd64*/ || exit 1
elif [[ $BUILD_OS == "Darwin" ]]; then
    if [[ "$BUILD_ARCH" == "x86_64" ]]; then
        pushd release/tce-darwin-amd64*/ || exit 1
    else
        pushd release/tce-darwin-arm64*/ || exit 1
    fi
fi
./uninstall.sh || { error "TCE CLEANUP (UNINSTALLATION) FAILED!"; exit 1; }
./install.sh || { error "TCE INSTALLATION FAILED!"; exit 1; }
popd || exit 1
echo "TCE version..."
tanzu management-cluster version || { error "Unexpected failure during TCE installation"; exit 1; }
