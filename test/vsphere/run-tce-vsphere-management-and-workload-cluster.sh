#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e
set -x
set -o pipefail

# Note: This script supports only Linux(Debian/Ubuntu) and MacOS
# Following environment variables are expected to be exported before running the script.
# The script will fail if any of them is missing
# VSPHERE_MANAGEMENT_CLUSTER_ENDPOINT - virtual and static IP for the management cluster's control plane nodes
# VSPHERE_WORKLOAD_CLUSTER_ENDPOINT - virtual and static IP for the management cluster's control plane nodes
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

declare -a required_env_vars=("VSPHERE_MANAGEMENT_CLUSTER_ENDPOINT"
"VSPHERE_WORKLOAD_CLUSTER_ENDPOINT"
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

# shellcheck source=test/util/utils.sh
source "${TCE_REPO_PATH}/test/util/utils.sh"

# shellcheck source=test/vsphere/cleanup-utils.sh
source "${TCE_REPO_PATH}/test/vsphere/cleanup-utils.sh"

"${TCE_REPO_PATH}/test/install-dependencies.sh" || { error "Dependency installation failed!"; exit 1; }
"${TCE_REPO_PATH}/test/build-tce.sh" || { error "TCE installation failed!"; exit 1; }

random_id="${RANDOM}"

export MANAGEMENT_CLUSTER_NAME="test-management-cluster-${random_id}"
export WORKLOAD_CLUSTER_NAME="test-workload-cluster-${random_id}"

function cleanup_management_cluster {
    delete_kind_cluster
    kubeconfig_cleanup ${MANAGEMENT_CLUSTER_NAME}
    echo "Using govc to cleanup ${MANAGEMENT_CLUSTER_NAME} management cluster resources"
    govc_cleanup ${MANAGEMENT_CLUSTER_NAME} || error "MANAGEMENT CLUSTER CLEANUP USING GOVC FAILED! Please manually delete any ${MANAGEMENT_CLUSTER_NAME} management cluster resources using vCenter Web UI"
}

function cleanup_workload_cluster {
    kubeconfig_cleanup ${WORKLOAD_CLUSTER_NAME}
    error "Using govc to cleanup ${WORKLOAD_CLUSTER_NAME} workload cluster resources"
    govc_cleanup ${WORKLOAD_CLUSTER_NAME} || error "WORKLOAD CLUSTER CLEANUP USING GOVC FAILED! Please manually delete any ${WORKLOAD_CLUSTER_NAME} workload cluster resources using vCenter Web UI"
}

function cleanup_management_and_workload_cluster {
    cleanup_management_cluster
    cleanup_workload_cluster
}

function create_management_cluster {
    management_cluster_config_file="${TCE_REPO_PATH}/test/vsphere/cluster-config.yaml"

    export VSPHERE_CONTROL_PLANE_ENDPOINT=${VSPHERE_MANAGEMENT_CLUSTER_ENDPOINT}
    export CLUSTER_NAME=${MANAGEMENT_CLUSTER_NAME}

    time tanzu management-cluster create ${MANAGEMENT_CLUSTER_NAME} --file "${management_cluster_config_file}" -v 10 || {
        error "MANAGEMENT CLUSTER CREATION FAILED!"
        unset VSPHERE_CONTROL_PLANE_ENDPOINT
        unset CLUSTER_NAME
        return 1
    }

    unset VSPHERE_CONTROL_PLANE_ENDPOINT
    unset CLUSTER_NAME
}

function check_management_cluster_creation {
    tanzu management-cluster get | grep "${MANAGEMENT_CLUSTER_NAME}" | grep running || {
        error "MANAGEMENT CLUSTER CREATION CHECK FAILED!"
        return 1
    }

    tanzu management-cluster kubeconfig get ${MANAGEMENT_CLUSTER_NAME} --admin || {
        error "ERROR GETTING MANAGEMENT CLUSTER KUBECONFIG!"
        return 1
    }

    "${TCE_REPO_PATH}/test/check-tce-cluster-creation.sh" ${MANAGEMENT_CLUSTER_NAME}-admin@${MANAGEMENT_CLUSTER_NAME} || {
        error "MANAGEMENT CLUSTER CREATION CHECK FAILED!"
        return 1
    }
}

function delete_management_cluster {
    echo "Deleting management cluster"
    time tanzu management-cluster delete ${MANAGEMENT_CLUSTER_NAME} -y || {
        error "MANAGEMENT CLUSTER DELETION FAILED!"
        return 1
    }
}

function create_workload_cluster {
    workload_cluster_config_file="${TCE_REPO_PATH}/test/vsphere/cluster-config.yaml"

    export VSPHERE_CONTROL_PLANE_ENDPOINT=${VSPHERE_WORKLOAD_CLUSTER_ENDPOINT}
    export CLUSTER_NAME=${WORKLOAD_CLUSTER_NAME}

    time tanzu cluster create ${WORKLOAD_CLUSTER_NAME} --file "${workload_cluster_config_file}" -v 10 || {
        error "WORKLOAD CLUSTER CREATION FAILED!"
        unset VSPHERE_CONTROL_PLANE_ENDPOINT
        unset CLUSTER_NAME
        return 1
    }

    unset VSPHERE_CONTROL_PLANE_ENDPOINT
    unset CLUSTER_NAME
}

function check_workload_cluster_creation {
    tanzu cluster list | grep "${WORKLOAD_CLUSTER_NAME}" | grep running || {
        error "WORKLOAD CLUSTER CREATION CHECK FAILED!"
        return 1
    }

    tanzu cluster kubeconfig get ${WORKLOAD_CLUSTER_NAME} --admin || {
        error "ERROR GETTING WORKLOAD CLUSTER KUBECONFIG!"
        return 1
    }

    "${TCE_REPO_PATH}/test/check-tce-cluster-creation.sh" ${WORKLOAD_CLUSTER_NAME}-admin@${WORKLOAD_CLUSTER_NAME} || {
        error "WORKLOAD CLUSTER CREATION CHECK FAILED!"
        return 1
    }
}

function add_package_repo {
    echo "Installing package repository on TCE..."
    "${TCE_REPO_PATH}/test/add-tce-package-repo.sh" || {
        error "PACKAGE REPOSITORY INSTALLATION FAILED!";
        return 1;
    }
}

function list_packages {
    echo "Listing available packages..."
    tanzu package available list || {
        error "LISTING PACKAGES FAILED";
        return 1;
    }
}

function test_gate_keeper_package {
    echo "Starting Gatekeeper test..."
    "${TCE_REPO_PATH}/test/gatekeeper/e2e-test.sh" || {
        error "GATEKEEPER PACKAGE TEST FAILED!";
        return 1;
    }
}

function delete_workload_cluster {
    echo "Deleting workload cluster"
    time tanzu cluster delete ${WORKLOAD_CLUSTER_NAME} -y || {
        error "WORKLOAD CLUSTER DELETION FAILED!"
        return 1
    }
}

function wait_for_workload_cluster_deletion {
    wait_iterations=120

    for (( i = 1 ; i <= wait_iterations ; i++))
    do
        echo "Waiting for workload cluster to get deleted..."
        num_of_clusters=$(tanzu cluster list -o json | jq 'length')
        if [[ "$num_of_clusters" == "0" ]]; then
            echo "Workload cluster ${WORKLOAD_CLUSTER_NAME} successfully deleted"
            break
        fi
        if [[ "${i}" == "${wait_iterations}" ]]; then
            echo "Timed out waiting for workload cluster ${WORKLOAD_CLUSTER_NAME} to get deleted"
            return 1
        fi
        sleep 5
    done
}

create_management_cluster || {
    delete_kind_cluster
    cleanup_management_cluster
    exit 1
}

check_management_cluster_creation || {
    cleanup_management_cluster
    exit 1
}

create_workload_cluster || {
    cleanup_management_and_workload_cluster
    exit 1
}

check_workload_cluster_creation || {
    cleanup_management_and_workload_cluster
    exit 1
}

add_package_repo || {
    cleanup_management_and_workload_cluster
    exit 1
}

list_packages || {
    cleanup_management_and_workload_cluster
    exit 1
}

test_gate_keeper_package || {
    cleanup_management_and_workload_cluster
    exit 1
}

echo "Cleaning up"

delete_workload_cluster || {
    cleanup_management_and_workload_cluster
    exit 1
}

wait_for_workload_cluster_deletion || {
    cleanup_management_and_workload_cluster
    exit 1
}

# since tanzu cluster delete does not delete workload cluster kubeconfig entry
kubeconfig_cleanup ${WORKLOAD_CLUSTER_NAME}

delete_management_cluster || {
    cleanup_management_cluster
    exit 1
}
