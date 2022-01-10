#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script tests TCE Management and Workload cluster in Azure.
# It builds TCE, spins up a Management and Workload cluster in Azure,
# installs the default packages,
# and cleans the environment.
# Note: This script supports only Linux(Debian) and MacOS
# Following environment variables need to be exported before running the script
# AZURE_TENANT_ID
# AZURE_SUBSCRIPTION_ID
# AZURE_CLIENT_ID
# AZURE_CLIENT_SECRET
# AZURE_SSH_PUBLIC_KEY_B64
# Azure location is set to australiacentral using AZURE_LOCATION
# The best way to run this is by calling `make azure-management-and-workload-cluster-e2e-test`
# from the root of the TCE repository.

set -e
set -x
set -o pipefail

TCE_REPO_PATH="$(git rev-parse --show-toplevel)"

declare -a required_env_vars=("AZURE_CLIENT_ID"
"AZURE_CLIENT_SECRET"
"AZURE_SSH_PUBLIC_KEY_B64"
"AZURE_SUBSCRIPTION_ID"
"AZURE_TENANT_ID")

"${TCE_REPO_PATH}/test/azure/check-required-env-vars.sh" "${required_env_vars[@]}"

# shellcheck source=test/util/utils.sh
source "${TCE_REPO_PATH}/test/util/utils.sh"
# shellcheck source=test/azure/utils.sh
source "${TCE_REPO_PATH}/test/azure/utils.sh"

"${TCE_REPO_PATH}/test/install-dependencies.sh" || { error "Dependency installation failed!"; exit 1; }
"${TCE_REPO_PATH}/test/build-tce.sh" || { error "TCE installation failed!"; exit 1; }

export CLUSTER_NAME_SUFFIX="${RANDOM}"
export MANAGEMENT_CLUSTER_NAME="test-mc-${CLUSTER_NAME_SUFFIX}"
export WORKLOAD_CLUSTER_NAME="test-wld-${CLUSTER_NAME_SUFFIX}"

echo "Setting MANAGEMENT_CLUSTER_NAME to ${MANAGEMENT_CLUSTER_NAME}"
echo "Setting WORKLOAD_CLUSTER_NAME to ${WORKLOAD_CLUSTER_NAME}"

export VM_IMAGE_PUBLISHER="vmware-inc"
# The value k8s-1dot21dot5-ubuntu-2004 comes from latest TKG BOM file based on OS arch, OS name and OS version
# provided in test/azure/cluster-config.yaml. This value needs to be changed manually whenever there's going to
# be a change in the underlying Tanzu Framework CLI version (management-cluster and cluster plugins) causing new
# TKr BOMs to be used with new Azure VM images which have different image billing plan SKU
export VM_IMAGE_BILLING_PLAN_SKU="k8s-1dot21dot5-ubuntu-2004"
export VM_IMAGE_OFFER="tkg-capi"


function cleanup_management_cluster {
    echo "Using azure CLI to cleanup ${MANAGEMENT_CLUSTER_NAME} management cluster resources"
    export CLUSTER_NAME="${MANAGEMENT_CLUSTER_NAME}"
    set_azure_env_vars
    delete_kind_cluster
    kubeconfig_cleanup ${CLUSTER_NAME}
    azure_cluster_cleanup || error "MANAGEMENT CLUSTER CLEANUP USING azure CLI FAILED! Please manually delete any ${MANAGEMENT_CLUSTER_NAME} management cluster resources using Azure Web UI"
    unset_azure_env_vars
    unset CLUSTER_NAME
}

function cleanup_workload_cluster {
    echo "Using azure CLI to cleanup ${WORKLOAD_CLUSTER_NAME} workload cluster resources"
    export CLUSTER_NAME="${WORKLOAD_CLUSTER_NAME}"
    set_azure_env_vars
    kubeconfig_cleanup ${CLUSTER_NAME}
    azure_cluster_cleanup || error "WORKLOAD CLUSTER CLEANUP USING azure CLI FAILED! Please manually delete any ${WORKLOAD_CLUSTER_NAME} workload cluster resources using Azure Web UI"
    unset_azure_env_vars
    unset CLUSTER_NAME
}

function cleanup_management_and_workload_cluster {
    cleanup_management_cluster
    cleanup_workload_cluster
}

function set_azure_env_vars {
    export AZURE_RESOURCE_GROUP="${CLUSTER_NAME}-resource-group"
    export AZURE_VNET_RESOURCE_GROUP="${AZURE_RESOURCE_GROUP}"
    export AZURE_VNET_NAME="${CLUSTER_NAME}-vnet"
    export AZURE_CONTROL_PLANE_SUBNET_NAME="${CLUSTER_NAME}-control-plane-subnet"
    export AZURE_NODE_SUBNET_NAME="${CLUSTER_NAME}-worker-node-subnet"
}

function unset_azure_env_vars {
    unset AZURE_RESOURCE_GROUP
    unset AZURE_VNET_RESOURCE_GROUP
    unset AZURE_VNET_NAME
    unset AZURE_CONTROL_PLANE_SUBNET_NAME
    unset AZURE_NODE_SUBNET_NAME
}

function create_management_cluster {
    echo "Bootstrapping TCE management cluster on Azure..."
    export CLUSTER_NAME="${MANAGEMENT_CLUSTER_NAME}"
    set_azure_env_vars

    management_cluster_config_file="${TCE_REPO_PATH}"/test/azure/cluster-config.yaml
    time tanzu management-cluster create ${MANAGEMENT_CLUSTER_NAME} --file "${management_cluster_config_file}" -v 10 || {
        error "MANAGEMENT CLUSTER CREATION FAILED!"
        unset_azure_env_vars
        unset CLUSTER_NAME
        return 1
    }

    unset_azure_env_vars
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

    "${TCE_REPO_PATH}"/test/check-tce-cluster-creation.sh ${MANAGEMENT_CLUSTER_NAME}-admin@${MANAGEMENT_CLUSTER_NAME} || {
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
    echo "Creating workload cluster on Azure..."
    export CLUSTER_NAME="${WORKLOAD_CLUSTER_NAME}"
    set_azure_env_vars

    workload_cluster_config_file="${TCE_REPO_PATH}"/test/azure/cluster-config.yaml
    time tanzu cluster create ${WORKLOAD_CLUSTER_NAME} --file "${workload_cluster_config_file}" -v 10 || {
        error "WORKLOAD CLUSTER CREATION FAILED!"
        unset_azure_env_vars
        unset CLUSTER_NAME
        return 1
    }

    unset_azure_env_vars
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

    "${TCE_REPO_PATH}"/test/check-tce-cluster-creation.sh ${WORKLOAD_CLUSTER_NAME}-admin@${WORKLOAD_CLUSTER_NAME} || {
        error "WORKLOAD CLUSTER CREATION CHECK FAILED!"
        return 1
    }
}

function add_package_repo {
    echo "Installing package repository on TCE..."
    "${TCE_REPO_PATH}"/test/add-tce-package-repo.sh || {
        error "PACKAGE REPOSITORY INSTALLATION FAILED!";
        return 1;
    }
}

function list_packages {
    tanzu package available list || {
        error "LISTING PACKAGES FAILED";
        return 1;
    }
}

function test_gate_keeper_package {
    echo "Starting Gatekeeper test..."
    "${TCE_REPO_PATH}"/test/gatekeeper/e2e-test.sh || {
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

    for (( i = 1 ; i <= wait_iterations ; i++ ))
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

accept_vm_image_terms || exit 1

create_management_cluster || {
    collect_management_cluster_diagnostics ${MANAGEMENT_CLUSTER_NAME}
    delete_kind_cluster
    cleanup_management_cluster
    exit 1
}

check_management_cluster_creation || {
    collect_management_cluster_diagnostics ${MANAGEMENT_CLUSTER_NAME}
    cleanup_management_cluster
    exit 1
}

create_workload_cluster || {
    collect_management_and_workload_cluster_diagnostics azure ${MANAGEMENT_CLUSTER_NAME} ${WORKLOAD_CLUSTER_NAME}
    cleanup_management_and_workload_cluster
    exit 1
}

check_workload_cluster_creation || {
    collect_management_and_workload_cluster_diagnostics azure ${MANAGEMENT_CLUSTER_NAME} ${WORKLOAD_CLUSTER_NAME}
    cleanup_management_and_workload_cluster
    exit 1
}

add_package_repo || {
    collect_management_and_workload_cluster_diagnostics azure ${MANAGEMENT_CLUSTER_NAME} ${WORKLOAD_CLUSTER_NAME}
    cleanup_management_and_workload_cluster
    exit 1
}

list_packages || {
    collect_management_and_workload_cluster_diagnostics azure ${MANAGEMENT_CLUSTER_NAME} ${WORKLOAD_CLUSTER_NAME}
    cleanup_management_and_workload_cluster
    exit 1
}

test_gate_keeper_package || {
    collect_management_and_workload_cluster_diagnostics azure ${MANAGEMENT_CLUSTER_NAME} ${WORKLOAD_CLUSTER_NAME}
    cleanup_management_and_workload_cluster
    exit 1
}

echo "Cleaning up"

delete_workload_cluster || {
    collect_management_and_workload_cluster_diagnostics azure ${MANAGEMENT_CLUSTER_NAME} ${WORKLOAD_CLUSTER_NAME}
    cleanup_management_and_workload_cluster
    exit 1
}

wait_for_workload_cluster_deletion || {
    collect_management_and_workload_cluster_diagnostics azure ${MANAGEMENT_CLUSTER_NAME} ${WORKLOAD_CLUSTER_NAME}
    cleanup_management_and_workload_cluster
    exit 1
}

# since tanzu cluster delete does not delete workload cluster kubeconfig entry
kubeconfig_cleanup ${WORKLOAD_CLUSTER_NAME}

delete_management_cluster || {
    collect_management_cluster_diagnostics ${MANAGEMENT_CLUSTER_NAME}
    cleanup_management_cluster
    exit 1
}
