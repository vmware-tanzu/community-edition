#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script clones TCE repo and builds the latest release

TCE_REPO_PATH=$(pwd)
export TCE_REPO_PATH
BUILD_OS=$(uname -s)
export BUILD_OS

# Helper functions
function error {
    printf '\E[31m'; echo "$@"; printf '\E[0m'
}

# Make sure Github Access Token is exported
if [[ -z "$GH_ACCESS_TOKEN" ]]; then
    echo "Access to GitHub private repo requires a token."
    echo "Please create a token (Settings > Developer Settings > Personal Access Tokens)"

    read -r -p "Please enter your GitHub token: " GH_ACCESS_TOKEN
    echo
    export GH_ACCESS_TOKEN=$GH_ACCESS_TOKEN
fi

git config --global url."https://git:$GH_ACCESS_TOKEN@github.com".insteadOf "https://github.com"

# Build TCE
cd "$TCE_REPO_PATH" || exit 
echo "4"
git clone https://git:"${GH_ACCESS_TOKEN}"@github.com/vmware-tanzu/tce
echo "Building TCE release"
cd tce || exit 
make release
echo "Installing TCE release"
if [[ $BUILD_OS == "Linux" ]]; then
    cd build/tce-linux-amd64*/ || exit 
elif [[ $BUILD_OS == "Darwin" ]]; then
    cd build/tce-darwin-amd64*/ || exit 
fi
./install.sh
echo "TCE version"
if [[ -z "$(tanzu standalone-cluster version)" ]]; then
    error "Unexpected failure during TCE installation"
    exit 1
fi