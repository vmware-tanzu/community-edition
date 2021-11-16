#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e
set -x

# Note: This script supports only Linux(Debian/Ubuntu) and MacOS
# Following environment variables are expected to be exported before running the script.
# The script will fail if any of them is missing
# VSPHERE_CONTROL_PLANE_ENDPOINT - virtual and static IP for the cluster's control plane nodes
# VSPHERE_SERVER - private IP of the vcenter server
# VSPHERE_SSH_AUTHORIZED_KEY - SSH public key to inject into control plane nodes and worker nodes for SSHing into them later
# VSPHERE_USERNAME - vcenter username
# VSPHERE_PASSWORD - vcenter password
# VSPHERE_DATACENTER - SDDC path
# VSPHERE_DATASTORE - Name of the vSphere datastore to deploy the Tanzu Kubernetes cluster as it appears in the vSphere inventory
# VSPHERE_FOLDER - name of an existing VM folder in which to place Tanzu Kubernetes Grid VMs
# VSPHERE_NETWORK - The network portgroup to assign each VM node
# VSPHERE_RESOURCE_POOL - Name of an existing resource pool in which to place this Tanzu Kubernetes cluster

TCE_REPO_PATH="$(git rev-parse --show-toplevel)"

declare -a required_env_vars=("VSPHERE_CONTROL_PLANE_ENDPOINT"
"VSPHERE_SERVER"
"VSPHERE_SSH_AUTHORIZED_KEY"
"VSPHERE_USERNAME"
"VSPHERE_PASSWORD"
"VSPHERE_DATACENTER"
"VSPHERE_DATASTORE"
"VSPHERE_FOLDER"
"VSPHERE_NETWORK"
"VSPHERE_RESOURCE_POOL")

"${TCE_REPO_PATH}/test/vsphere/check-required-env-vars.sh" "${required_env_vars[@]}"

"${TCE_REPO_PATH}/test/install-dependencies.sh" || { error "Dependency installation failed!"; exit 1; }
"${TCE_REPO_PATH}/test/build-tce.sh" || { error "TCE installation failed!"; exit 1; }

# shellcheck source=test/util/utils.sh
source "${TCE_REPO_PATH}/test/util/utils.sh"

# shellcheck source=test/vsphere/cleanup-utils.sh
source "${TCE_REPO_PATH}/test/vsphere/cleanup-utils.sh"

export CLUSTER_NAME="test-standalone-cluster-${RANDOM}"

cluster_config_file="${TCE_REPO_PATH}/test/vsphere/cluster-config.yaml"

function cleanup {
    kubeconfig_cleanup ${CLUSTER_NAME}
    govc_cleanup ${CLUSTER_NAME}
}

time tanzu standalone-cluster create ${CLUSTER_NAME} --file "${cluster_config_file}" -v 10 || {
    error "STANDALONE CLUSTER CREATION FAILED!"
    delete_kind_cluster
    cleanup
    exit 1
}

"${TCE_REPO_PATH}/test/check-tce-cluster-creation.sh" ${CLUSTER_NAME}-admin@${CLUSTER_NAME} || {
    error "STANDALONE CLUSTER CREATION CHECK FAILED!"
    cleanup
    exit 1
}

echo "Cleaning up"
echo "Deleting standalone cluster"

time tanzu standalone-cluster delete ${CLUSTER_NAME} -y || {
    error "STANDALONE CLUSTER DELETION FAILED!"
    delete_kind_cluster
    cleanup
    exit 1
}
