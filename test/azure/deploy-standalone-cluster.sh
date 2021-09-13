#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script tests TCE Standalone cluster in Azure.
# It builds TCE, spins up a standalone cluster in Azure,
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

set -e
set -x

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
TCE_REPO_PATH="${MY_DIR}"/../..

"${TCE_REPO_PATH}"/test/azure/check-required-env-vars.sh

# shellcheck source=test/util/utils.sh
source "${TCE_REPO_PATH}"/test/util/utils.sh
"${TCE_REPO_PATH}"/test/install-dependencies.sh || { error "Dependency installation failed!"; exit 1; }
"${TCE_REPO_PATH}"/test/build-tce.sh || { error "TCE installation failed!"; exit 1; }

export CLUSTER_NAME="test${RANDOM}"
echo "Setting CLUSTER_NAME to ${CLUSTER_NAME}..."

export AZURE_RESOURCE_GROUP="${CLUSTER_NAME}-resource-group"
export AZURE_VNET_RESOURCE_GROUP="${AZURE_RESOURCE_GROUP}"
export AZURE_VNET_NAME="${CLUSTER_NAME}-vnet"
export AZURE_CONTROL_PLANE_SUBNET_NAME="${CLUSTER_NAME}-control-plane-subnet"
export AZURE_NODE_SUBNET_NAME="${CLUSTER_NAME}-worker-node-subnet"

function az_docker {
    docker run --user "$(id -u)":"$(id -g)" \
        --volume "${HOME}":/home/az \
        --env HOME=/home/az \
        --rm \
        mcr.microsoft.com/azure-cli az "$@"
}

function azure_cluster_cleanup {
    failure_message="STANDLONE CLUSTER CLEANUP USING azure CLI FAILED! Please manually delete any ${CLUSTER_NAME} standalone cluster resources using Azure Web UI"
    echo "Cleaning up ${CLUSTER_NAME} cluster resources using azure CLI"

    az_docker login --service-principal --username "${AZURE_CLIENT_ID}" --password "${AZURE_CLIENT_SECRET}" \
        --tenant "${AZURE_TENANT_ID}" || {
        error "azure CLI LOGIN FAILED!"
        error "${failure_message}"
        return 1
    }

    az_docker account set --subscription "${AZURE_SUBSCRIPTION_ID}" || {
        error "azure CLI SETTING ACCOUNT SUBSCRIPTION ID FAILED!"
        error "${failure_message}"
        return 1
    }

    az_docker group delete --name "${AZURE_RESOURCE_GROUP}" --yes || {
        error "azure CLI RESOURCE GROUP DELETION FAILED!"
        error "${failure_message}"
        return 1
    }
}

function delete_cluster_or_cleanup {
    echo "Deleting standalone cluster"
    time tanzu standalone-cluster delete ${CLUSTER_NAME} -y || {
        error "STANDALONE CLUSTER DELETION FAILED!";
        azure_cluster_cleanup
        return 1
    }
}

function create_standalone_cluster {
    echo "Bootstrapping TCE standalone cluster on Azure..."
    time tanzu standalone-cluster create "${CLUSTER_NAME}" -f "${TCE_REPO_PATH}"/test/azure/cluster-config.yaml || {
        error "STANDALONE CLUSTER CREATION FAILED!";
        return 1;
    }
}

function wait_for_pods {
    kubectl config use-context "${CLUSTER_NAME}"-admin@"${CLUSTER_NAME}" || {
        error "CONTEXT SWITCH TO STANDALONE CLUSTER FAILED!";
        return 1;
    }
    kubectl wait --for=condition=ready pod --all --all-namespaces --timeout=600s || {
        error "TIMED OUT WAITING FOR ALL PODS TO BE UP!";
        return 1;
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
    "${TCE_REPO_PATH}"/test/aws/e2e-test.sh || {
        error "GATEKEEPER PACKAGE TEST FAILED!";
        return 1;
    }
}

create_standalone_cluster || {
    delete_kind_cluster
    azure_cluster_cleanup
    exit 1
}

wait_for_pods || {
    delete_cluster_or_cleanup
    exit 1
}


add_package_repo || {
    delete_cluster_or_cleanup
    exit 1
}

list_packages || {
    delete_cluster_or_cleanup
    exit 1
}

test_gate_keeper_package || {
    delete_cluster_or_cleanup
    exit 1
}

echo "Cleaning up..."
delete_cluster_or_cleanup
