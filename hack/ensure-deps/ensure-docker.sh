#!/bin/bash

# Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script installs dependencies needed for running build-tce.sh and deploy-tce.sh

BUILD_OS=$(uname -s)
export BUILD_OS
# Make sure docker is installed
echo "Checking for Docker..."
if [[ -z "$(command -v docker)" ]]; then
    echo "Installing Docker..."
    if [[ "$BUILD_OS" == "Linux" ]]; then
        curl -fsSL https://get.docker.com -o get-docker.sh
        sh ./get-docker.sh
    elif [[ "$(BUILD_OS)" == "Darwin" ]]; then
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
        brew cask install docker
    fi
else
    echo "Found Docker!"
fi

echo "Verifying Docker..."
if ! sudo docker run --rm hello-world > /dev/null; then
    error "Unable to verify docker functionality, make sure docker is installed correctly"
    exit 1
else
    echo "Verified Docker functionality successfully!"
fi
