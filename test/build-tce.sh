#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script clones TCE repo and builds the latest release.
# Make sure GitHub Access Token is exported.
# This script clones TCE repo and builds the latest release
set -x
set -e 

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# shellcheck source=test/util/utils.sh
source "${MY_DIR}"/util/utils.sh

BUILD_OS=$(uname -s)
export BUILD_OS

if [[ -z "$GITHUB_TOKEN" ]]; then
    echo "Access to GitHub private repo requires a token."
    echo "Please create a token (Settings > Developer Settings > Personal Access Tokens)"

    read -r -p "Please enter your GitHub token: " GITHUB_TOKEN
    echo
    export GITHUB_TOKEN=$GITHUB_TOKEN
fi

git config --global url."https://git:$GITHUB_TOKEN@github.com".insteadOf "https://github.com"

# Build TCE
echo "Building TCE release..."
make release || { error "TCE BUILD FAILED!"; exit 1; }
echo "Installing TCE release"
if [[ $BUILD_OS == "Linux" ]]; then
    pushd build/tce-linux-amd64*/ || exit 1
elif [[ $BUILD_OS == "Darwin" ]]; then
    if [[ "$BUILD_ARCH" == "x86_64" ]]; then
        pushd build/tce-darwin-amd64*/ || exit 1
    else
        pushd build/tce-darwin-arm64*/ || exit 1
    fi
    
fi
./install.sh || { error "TCE INSTALLATION FAILED!"; exit 1; }
popd || exit 1
echo "TCE version..."
tanzu standalone-cluster version || { error "Unexpected failure during TCE installation"; exit 1; }
