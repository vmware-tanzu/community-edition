#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script tests TCE Management cluster on AWS.
# It builds TCE, spins up a management cluster in AWS, 
# creates a workload cluster, installs the default packages, 
# tests the e2e functionality of Gatekeeper package and cleans the environment.
# Note: This is WIP and supports only Linux(Debian) and MacOS
# Following environment variables need to be exported before running the script
# AWS_ACCOUNT_ID
# AWS_ACCESS_KEY_ID
# AWS_SECRET_ACCESS_KEY
# AWS_B64ENCODED_CREDENTIALS
# AWS_SSH_KEY_NAME
# Region is set to us-east-2
# The best way to run this is by calling `make tce-aws-managed-cluster-e2e-test`
# from the root of the TCE repository.

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
TCE_REPO_PATH="${MY_DIR}"/../..
# shellcheck source=test/util/utils.sh
source "${TCE_REPO_PATH}"/test/util/utils.sh
# shellcheck source=test/util/aws-nuke-tear-down.sh
source "${TCE_REPO_PATH}"/test/util/aws-nuke-tear-down.sh
"${TCE_REPO_PATH}"/test/build-tce.sh || { error "TCE installation failed!"; exit 1; }
"${TCE_REPO_PATH}"/test/install-jq.sh
"${TCE_REPO_PATH}"/test/install-dependencies.sh || { error "Dependency installation failed!"; exit 1; }

function delete_management_cluster {
    echo "$@"
    export AWS_REGION="us-east-2"
    tanzu management-cluster delete "${MGMT_CLUSTER_NAME}" -y || { aws-nuke-tear-down "MANAGEMENT CLUSTER DELETION FAILED! Deleting the cluster using AWS-NUKE..." "${MGMT_CLUSTER_NAME}"; }
}

function nuke_management_and_workload_clusters {
    aws-nuke-tear-down "Deleting the MANAGEMENT CLUSTER using AWS-NUKE..." "${MGMT_CLUSTER_NAME}"
    aws-nuke-tear-down "Deleting the WORKLOAD CLUSTER using AWS-NUKE..." "${WLD_CLUSTER_NAME}";
}

function delete_workload_cluster {
    echo "$@"
    tanzu cluster delete "${WLD_CLUSTER_NAME}" --yes || { nuke_management_and_workload_clusters; exit 1; }
    for (( i = 1 ; i <= 120 ; i++))
    do
        echo "Waiting for workload cluster to get deleted..."
        num_of_clusters=$(tanzu cluster list -o json | jq 'length')
        if [[ "$num_of_clusters" != "0" ]]; then
            echo "Workload cluster ${WLD_CLUSTER_NAME} successfully deleted"
            break
        fi
        if [[ "$i" == 120 ]]; then
            echo "Timed out waiting for workload cluster ${WLD_CLUSTER_NAME} to get deleted"
            echo "Using AWS NUKE to delete management and workload clusters"
            nuke_management_and_workload_clusters
            exit 1
        fi
        sleep 5
    done
    echo "Workload cluster ${WLD_CLUSTER_NAME} successfully deleted"
}

function create_management_cluster {
    echo "Bootstrapping TCE management cluster on AWS..."
    # Set management cluster name
    export CLUSTER_NAME_SUFFIX=${RANDOM}
    export MGMT_CLUSTER_NAME="test-mc-${MGMT_CLUSTER_NAME_SUFFIX}"
    echo "Setting MANAGEMENT CLUSTER NAME to ${MGMT_CLUSTER_NAME_SUFFIX}..."
    tanzu management-cluster create "${MGMT_CLUSTER_NAME_SUFFIX}" -f "${TCE_REPO_PATH}"/test/aws/cluster-config.yaml || { error "MANAGEMENT CLUSTER CREATION FAILED!"; delete_kind_cluster; aws-nuke-tear-down "Deleting management cluster" "${MGMT_CLUSTER_NAME}"; exit 1; }
    kubectl config use-context "${MGMT_CLUSTER_NAME_SUFFIX}"-admin@"${MGMT_CLUSTER_NAME_SUFFIX}" || { error "CONTEXT SWITCH TO MANAGEMENT CLUSTER FAILED!"; delete_management_cluster "Deleting management cluster"; exit 1; }
    kubectl wait --for=condition=ready pod --all --all-namespaces --timeout=300s || { error "TIMED OUT WAITING FOR ALL PODS TO BE UP!"; delete_management_cluster "Deleting management cluster"; exit 1; }
    tanzu management-cluster get | grep "${MGMT_CLUSTER_NAME_SUFFIX}" | grep running || { error "MANAGEMENT CLUSTER NOT RUNNING!"; delete_management_cluster "Deleting management cluster"; exit 1; }
}

function create_workload_cluster {
    echo "Creating workload cluster..."
    # Set workload cluster name
    export WLD_CLUSTER_NAME="test-wld-${CLUSTER_NAME_SUFFIX}"
    echo "Setting WORKLOAD CLUSTER NAME to ${WLD_CLUSTER_NAME}..."
    tanzu cluster create "${WLD_CLUSTER_NAME}" -f "${TCE_REPO_PATH}"/test/aws/cluster-config.yaml || { error "WORKLOAD CLUSTER CREATION FAILED!"; nuke_management_and_workload_clusters; exit 1; }
    tanzu cluster kubeconfig get "${WLD_CLUSTER_NAME}" --admin
    kubectl config use-context "${WLD_CLUSTER_NAME}"-admin@"${WLD_CLUSTER_NAME}" || { error "CONTEXT SWITCH TO MANAGEMENT CLUSTER FAILED!"; delete_workload_cluster "Deleting workload cluster"; delete_management_cluster "Deleting management cluster"; exit 1; }
    kubectl wait --for=condition=ready pod --all --all-namespaces --timeout=300s || { error "TIMED OUT WAITING FOR ALL PODS TO BE UP!"; delete_workload_cluster "Deleting workload cluster"; delete_management_cluster "Deleting management cluster"; exit 1; }
}

# Create management and workload clusters
create_management_cluster || exit 1
create_workload_cluster || exit 1

# Install packages
echo "Installing packages on TCE..."
"${TCE_REPO_PATH}"/test/add-tce-package-repo.sh || { error "PACKAGE REPOSITORY INSTALLATION FAILED!"; delete_workload_cluster "Deleting workload cluster"; delete_management_cluster "Deleting management cluster"; exit 1; }
tanzu package available list || { error "UNEXPECTED FAILURE OCCURRED!"; delete_workload_cluster "Deleting workload cluster"; delete_management_cluster "Deleting management cluster"; exit 1; }

# Run e2e test
echo "Starting Gatekeeper test..."
"${TCE_REPO_PATH}"/test/aws/e2e-test.sh || { error "TEST FAILED!"; delete_workload_cluster "Deleting workload cluster"; delete_management_cluster "Deleting management cluster"; exit 1; }

# Clean up
echo "Cleaning up..."
delete_workload_cluster "Deleting workload cluster" || { exit 1; }
delete_management_cluster "Deleting management cluster" || { exit 1; }