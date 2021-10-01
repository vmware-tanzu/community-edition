#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
TCE_REPO_PATH=${MY_DIR}/../../../../..

TCE_VERSION="v0.9.1"

echo "Installing TCE ${TCE_VERSION}"

BUILD_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
TCE_RELEASE_TAR_BALL="tce-${BUILD_OS}-amd64-${TCE_VERSION}.tar.gz"
TCE_RELEASE_DIR="tce-${BUILD_OS}-amd64-${TCE_VERSION}"
INSTALLATION_DIR="${MY_DIR}/tce-installation"

"${TCE_REPO_PATH}"/hack/get-tce-release.sh ${TCE_VERSION} "${BUILD_OS}"-amd64

mkdir -p "${INSTALLATION_DIR}"
tar xzvf "${TCE_RELEASE_TAR_BALL}" --directory="${INSTALLATION_DIR}"

"${INSTALLATION_DIR}"/"${TCE_RELEASE_DIR}"/install.sh || { error "Unexpected failure during TCE installation"; exit 1; }

echo "TCE version: "
tanzu standalone-cluster version || { error "Unexpected failure during TCE installation"; exit 1; }

TANZU_DIAGNOSTICS_PLUGIN_DIR=${MY_DIR}/..
TANZU_DIAGNOSTICS_BIN=${MY_DIR}/tanzu-diagnostics-e2e-bin

echo "Entering ${TANZU_DIAGNOSTICS_PLUGIN_DIR} directory to build tanzu diagnostics plugin"
pushd "${TANZU_DIAGNOSTICS_PLUGIN_DIR}"

go build -o "${TANZU_DIAGNOSTICS_BIN}" -v

echo "Finished building tanzu diagnostics plugin. Leaving ${TANZU_DIAGNOSTICS_PLUGIN_DIR}"
popd

CLUSTER_NAME_SUFFIX=${RANDOM}
CLUSTER_NAME="e2e-diagnostics-${CLUSTER_NAME_SUFFIX}"
CLUSTER_KUBE_CONTEXT="kind-${CLUSTER_NAME}"
NEW_CLUSTER_KUBE_CONTEXT="${CLUSTER_NAME}-admin@${CLUSTER_NAME}"
OUTPUT_DIR=$(mktemp -d)

echo "Creating a kind cluster for the E2E test"

kind create cluster --name ${CLUSTER_NAME} || {
    echo "Error creating kind cluster!"
    exit 1
}

# The context rename is required for workload cluster diagnostics data collection to work
# as it expects the context name to be in a particular format based on cluster name
# and --workload-cluster-context flag is not supported for now.
kubectl config rename-context ${CLUSTER_KUBE_CONTEXT} ${NEW_CLUSTER_KUBE_CONTEXT} || {
    echo "Error renaming kube context!"
    exit 1
}

echo "Running tanzu diagnostics collect command"

"${TANZU_DIAGNOSTICS_BIN}" collect --bootstrap-cluster-name ${CLUSTER_NAME} \
    --management-cluster-kubeconfig "${HOME}/.kube/config" \
    --management-cluster-context ${NEW_CLUSTER_KUBE_CONTEXT} \
    --management-cluster-name ${CLUSTER_NAME} \
    --workload-cluster-infra docker \
    --workload-cluster-name ${CLUSTER_NAME} \
    --output-dir "${OUTPUT_DIR}" || {
        echo "Error running tanzu diagnostics collect command!"
        exit 1
    }

echo "Checking if the diagnostics tar balls for the different clusters have been created"

EXPECTED_BOOTSTRAP_CLUSTER_DIAGNOSTICS="${OUTPUT_DIR}/bootstrap.${CLUSTER_NAME}.diagnostics.tar.gz"
EXPECTED_MANAGEMENT_CLUSTER_DIAGNOSTICS="${OUTPUT_DIR}/management-cluster.${CLUSTER_NAME}.diagnostics.tar.gz"
EXPECTED_WORKLOAD_CLUSTER_DIAGNOSTICS="${OUTPUT_DIR}/workload-cluster.${CLUSTER_NAME}.diagnostics.tar.gz"

errors=0

if [ ! -f "$EXPECTED_BOOTSTRAP_CLUSTER_DIAGNOSTICS" ]; then
    echo "$EXPECTED_BOOTSTRAP_CLUSTER_DIAGNOSTICS does not exist. Expected bootstrap cluster diagnostics tar ball to be present"
    ((errors=errors+1))
fi

if [ ! -f "$EXPECTED_MANAGEMENT_CLUSTER_DIAGNOSTICS" ]; then
    echo "$EXPECTED_MANAGEMENT_CLUSTER_DIAGNOSTICS does not exist. Expected management cluster diagnostics tar ball to be present"
    ((errors=errors+1))
fi

if [ ! -f "$EXPECTED_WORKLOAD_CLUSTER_DIAGNOSTICS" ]; then
    echo "$EXPECTED_WORKLOAD_CLUSTER_DIAGNOSTICS does not exist. Expected workload cluster diagnostics tar ball to be present"
    ((errors=errors+1))
fi

if [[ ${errors} -gt 0 ]]; then
    echo "Total E2E errors in tanzu diagnostics plugin: ${errors}"
fi

echo "Cleaning up"

kind delete cluster --name ${CLUSTER_NAME} || {
    echo "Failed deleting kind cluster. Please delete it manually"
}

# This deletion of context name is required as kind does not delete it
# because we renamed the context to something different from kind's context
# name convention
kubectl config delete-context ${NEW_CLUSTER_KUBE_CONTEXT} || {
    echo "Failed deleting kind cluster kube context. Please delete it manually"
}

rm -rfv "${OUTPUT_DIR}"

if [[ ${errors} -gt 0 ]]; then
    exit 1
fi
