#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script clones TCE repo and builds the latest release

source test/install-dependencies.sh

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
echo "Building TCE release..."
make release || { error "TCE BUILD FAILED!"; exit 1; }
echo "Installing TCE release"
if [[ $BUILD_OS == "Linux" ]]; then
    pushd build/tce-linux-amd64*/ || exit 1
elif [[ $BUILD_OS == "Darwin" ]]; then
    pushd build/tce-darwin-amd64*/ || exit 1
fi
./install.sh || { error "TCE INSTALLATION FAILED!"; exit 1; }
popd || exit 1
echo "TCE version..."
tanzu standalone-cluster version || { error "Unexpected failure during TCE installation"; exit 1; }
