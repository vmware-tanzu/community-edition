#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script tests TCE Standalone cluster in AWS.
# It builds TCE, spins up a standalone cluster in AWS,
# installs the default packages, tests the e2e functionality of Gatekeeper package
# and cleans the environment.
# Note: This is WIP and supports only Linux(Debian) and MacOS
# Following environment variables need to be exported before running the script
# AWS_ACCOUNT_ID
# AWS_ACCESS_KEY_ID
# AWS_SECRET_ACCESS_KEY
# AWS_B64ENCODED_CREDENTIALS
# AWS_SSH_KEY_NAME
# Region is set to us-east-2
# The best way to run this is by calling `make tce-aws-standalone-cluster-e2e-test`
# from the root of the TCE repository.

TCE_REPO_PATH="$(git rev-parse --show-toplevel)"
# shellcheck source=test/util/utils.sh
source "${TCE_REPO_PATH}/test/util/utils.sh"
# shellcheck source=test/util/aws-nuke-tear-down.sh
source "${TCE_REPO_PATH}/test/util/aws-nuke-tear-down.sh"
"${TCE_REPO_PATH}/test/install-jq.sh"
"${TCE_REPO_PATH}/test/install-dependencies.sh" || { error "Dependency installation failed!"; exit 1; }
"${TCE_REPO_PATH}/test/build-tce.sh" || { error "TCE installation failed!"; exit 1; }

# Set standalone cluster name
export CLUSTER_NAME="test${RANDOM}"
echo "Setting CLUSTER_NAME to ${CLUSTER_NAME}..."

# Cleanup function
function delete_cluster {
    echo "$@"
    tanzu standalone-cluster delete ${CLUSTER_NAME} -y || { kubeconfig_cleanup ${CLUSTER_NAME}; aws-nuke-tear-down "STANDALONE CLUSTER DELETION FAILED! Deleting the cluster using AWS-NUKE..."; }
}

function create_standalone_cluster {
    echo "Bootstrapping TCE standalone cluster on AWS..."
    tanzu standalone-cluster create "${CLUSTER_NAME}" -f "${TCE_REPO_PATH}"/test/aws/cluster-config.yaml || { error "STANDALONE CLUSTER CREATION FAILED!"; delete_kind_cluster; kubeconfig_cleanup ${CLUSTER_NAME}; aws-nuke-tear-down "Deleting standalone cluster" "${CLUSTER_NAME}"; exit 1; }
    kubectl config use-context "${CLUSTER_NAME}"-admin@"${CLUSTER_NAME}" || { error "CONTEXT SWITCH TO STANDALONE CLUSTER FAILED!"; delete_cluster "Deleting standalone cluster"; exit 1; }
    kubectl wait --for=condition=ready pod --all --all-namespaces --timeout=300s || { error "TIMED OUT WAITING FOR ALL PODS TO BE UP!"; delete_cluster "Deleting standalone cluster"; exit 1; }
}

# Create standalone cluster
create_standalone_cluster

# Install packages
echo "Installing packages on TCE..."
"${TCE_REPO_PATH}"/test/add-tce-package-repo.sh || { error "PACKAGE REPOSITORY INSTALLATION FAILED!"; delete_cluster "Deleting standalone cluster"; exit 1; }
tanzu package available list || { error "UNEXPECTED FAILURE OCCURRED!"; delete_cluster "Deleting standalone cluster"; exit 1; }

# Run e2e-test
echo "Starting Gatekeeper test..."
"${TCE_REPO_PATH}"/test/gatekeeper/e2e-test.sh || { error "TEST FAILED!"; delete_cluster "Deleting standalone cluster"; exit 1; }

# Clean up
delete_cluster "Cleaning up..."
