#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# WIP script to test TCE in AWS
# Basic idea is to build the latest release of TCE, spin up a standalone cluster
# in AWS, install the default packages, test the packages and clean the environment 
# using aws-nuke
# Note: This is WIP and supports only Linux(Debian) and MacOS
# Following environment variables are expected to be set before running the script
# AWS_ACCOUNT_ID
# AWS_ACCESS_KEY_ID
# AWS_SECRET_ACCESS_KEY
# AWS_B64ENCODED_CREDENTIALS
# AWS_SSH_KEY_NAME
# Region is set to us-east-2

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

# Make sure docker is installed
echo "Docker check!"
if [[ -z "$(which docker)" ]]; then
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
if [[ -z "$(which kubectl)" ]]; then
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

git config --global url."https://git:$GH_ACCESS_TOKEN@github.com".insteadOf "https://github.com"

# Build TCE
cd "$TCE_REPO_PATH" || echo "Unexpected failure, exiting..."; exit 1
git clone https://git:"${GH_ACCESS_TOKEN}"@github.com/vmware-tanzu/tce
echo "Building TCE release"
cd tce || echo "TCE repo not cloned, exiting..."; exit 1
make release
echo "Installing TCE release"
if [[ $BUILD_OS == "Linux" ]]; then
    cd build/tce-linux-amd64*/ || echo "Unexpected artifact structuring, exiting..."; exit 1
elif [[ $BUILD_OS == "Darwin" ]]; then
    cd build/tce-darwin-amd64*/ || echo "Unexpected artifact structuring, exiting..."; exit 1
fi
./install.sh
echo "TCE version"
if [[ -z "$(tanzu standalone-cluster version)" ]]; then
    error "Unexpected failure during TCE installation"
    exit 1
fi

# Set standalone cluster name
echo "Setting GUEST_CLUSTER_NAME to guest"
export GUEST_CLUSTER_NAME="guest"

# Substitute env variables in aws-template
echo "Bootstrapping TCE standalone cluster on AWS"
envsubst < "${TCE_REPO_PATH}"/hack/e2e-tests/aws/standalone-template.yaml > "${TCE_REPO_PATH}"/hack/e2e-tests/aws/standalone.yaml
tanzu standalone-cluster create "${GUEST_CLUSTER_NAME}" -f "${TCE_REPO_PATH}"/hack/e2e-tests/aws/standalone.yaml
rm -rf "${TCE_REPO_PATH}"/hack/e2e-tests/aws/standalone.yaml

kubectl config use-context "${GUEST_CLUSTER_NAME}"-admin@"${GUEST_CLUSTER_NAME}"
tanzu package repository install --default

tanzu package list

#AWS cleanup
envsubst < "${TCE_REPO_PATH}"/hack/e2e-tests/aws/nuke-config-template.yml > "${TCE_REPO_PATH}"/hack/e2e-tests/aws/nuke-config.yml
aws-nuke -c "${TCE_REPO_PATH}"/hack/e2e-tests/aws/nuke-config.yml --access-key-id "$AWS_ACCESS_KEY_ID" --secret-access-key "$AWS_SECRET_ACCESS_KEY" --force --no-dry-run
rm -rf "${TCE_REPO_PATH}"/hack/e2e-tests/aws/nuke-config.yml