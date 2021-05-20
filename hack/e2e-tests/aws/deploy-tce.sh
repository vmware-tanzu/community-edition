#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# WIP script to test TCE in AWS
# install-dependencies.sh and build-tce.sh should be executed before this script.
# Basic idea is to build the latest release of TCE, spin up a standalone cluster
# in AWS, install the default packages, test the packages and clean the environment 
# using aws-nuke
# Note: This is WIP and supports only Linux(Debian) and MacOS
# Following environment variables are expected to be exported before running the script
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