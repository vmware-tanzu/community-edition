#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script installs dependencies needed for running build-tce.sh and deploy-tce.sh

TCE_REPO_PATH=$(pwd)
export TCE_REPO_PATH
BUILD_OS=$(uname -s)
export BUILD_OS

# Helper functions
function error {
    printf '\E[31m'; echo "$@"; printf '\E[0m'
}

# Make sure docker is installed
echo "Checking for Docker..."
if [[ -z "$(command -v docker)" ]]; then
    echo "Installing Docker..."
    if [[ "$BUILD_OS" == "Linux" ]]; then
        sudo apt-get update > /dev/null
        sudo apt-get install -y \
            apt-transport-https \
            ca-certificates \
            curl \
            gnupg \
            lsb-release > /dev/null 2>&1
        curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
        echo \
            "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian \
            $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

        sudo apt update > /dev/null
        sudo apt install -y docker-ce docker-ce-cli containerd.io > /dev/null 2>&1
        sudo service docker start
        sleep 30s

        if [[ $(id -u) -ne 0 ]]; then
            sudo usermod -aG docker "$(whoami)"
        fi
    elif [[ "$(BUILD_OS)" == "Darwin" ]]; then
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
        brew cask install docker
    fi
else
    echo "Found Docker!"
fi

echo "Verifying Docker..."
if ! sudo docker run hello-world > /dev/null; then
    error "Unable to verify docker functionality, make sure docker is installed correctly"
    exit 1
else
    echo "Verified Docker functionality successfully!"
fi

# Make sure kubectl is installed
echo "Checking for kubectl..."
if [[ -z "$(command -v kubectl)" ]]; then
    echo "Installing kubectl..."
    if [[ "$BUILD_OS" == "Linux" ]]; then
        curl -LO https://dl.k8s.io/release/v1.20.1/bin/linux/amd64/kubectl
    elif [[ "$BUILD_OS" == "Darwin" ]]; then
        curl -LO https://dl.k8s.io/release/v1.20.1/bin/darwin/amd64/kubectl
    else
        error "$BUILD_OS NOT SUPPORTED!!!"
        exit 1
    fi
    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
else
    echo "Found kubectl!"
fi

echo "Installing envsubst"
if [[ -z "$(command -v envsubst)" ]]; then
	if [[ "$BUILD_OS" == "Linux" ]]; then
		sudo apt-get update > /dev/null
		sudo apt-get -y install gettext-base wget
	elif [[ "$BUILD_OS" == "Darwin" ]]; then
		echo "Please install gettext"
		exit 1
	fi
fi

# echo "Installing aws-nuke"
# if [[ -z "$(command -v aws-nuke)" ]]; then
#     if [[ "$BUILD_OS" == "Linux" ]]; then
#         wget -q https://github.com/rebuy-de/aws-nuke/releases/download/v2.15.0/aws-nuke-v2.15.0-linux-amd64.tar.gz
# 		tar xvzf aws-nuke-v2.15.0-linux-amd64.tar.gz && mv aws-nuke-v2.15.0-linux-amd64 aws-nuke
# 		sudo mv aws-nuke /usr/local/bin/
#     elif [[ "$BUILD_OS" == "Darwin" ]]; then
#         echo "Please install aws-nuke"
#         exit 1
#     fi
# fi
