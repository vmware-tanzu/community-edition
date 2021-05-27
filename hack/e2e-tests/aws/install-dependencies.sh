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
echo "Docker check!"
if [[ -z "$(command -v docker)" ]]; then
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
fi

if ! sudo docker run hello-world > /dev/null; then
    error "Unable to verify docker functionality, make sure docker is installed correctly"
    exit 1
fi

# Make sure kubectl is installed
if [[ -z "$(command -v kubectl)" ]]; then
    if [[ "$BUILD_OS" == "Linux" ]]; then
        curl -LO https://dl.k8s.io/release/v1.20.1/bin/linux/amd64/kubectl
    elif [[ "$BUILD_OS" == "Darwin" ]]; then
        curl -LO https://dl.k8s.io/release/v1.20.1/bin/darwin/amd64/kubectl
    else
        error "$BUILD_OS NOT SUPPORTED!!!"
        exit 1
    fi
    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
fi

# Installing aws-nuke and envsubst
# TODO(rajaskakodkar): Add installation steps for mac
sudo apt-get update > /dev/null
sudo apt-get -y install gettext-base wget > /dev/null 2>&1
# TODO(rajaskakodkar): Optimize the following bash
wget -q https://github.com/rebuy-de/aws-nuke/releases/download/v2.15.0/aws-nuke-v2.15.0-linux-amd64.tar.gz
tar xvzf aws-nuke-v2.15.0-linux-amd64.tar.gz && mv aws-nuke-v2.15.0-linux-amd64 aws-nuke
sudo mv aws-nuke /usr/local/bin/ 