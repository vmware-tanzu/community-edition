#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
TANZU_DIAGNOSTICS_BIN=${MY_DIR}/tanzu-diagnostics-e2e-bin

CLUSTER_NAME_SUFFIX=${RANDOM}
CLUSTER_NAME="e2e-diagnostics-${CLUSTER_NAME_SUFFIX}"
CLUSTER_KUBE_CONTEXT="kind-${CLUSTER_NAME}"
OUTPUT_DIR=$(mktemp -d)

echo "Creating a kind cluster for the E2E test"

kind create cluster --name ${CLUSTER_NAME} || {
    echo "Error creating kind cluster!"
    exit 1
}

echo "Running tanzu diagnostics collect command"

if [[ ! -f "${TANZU_DIAGNOSTICS_BIN}" ]]; then
    # we are running on a typical install of TCE
    tanzu diagnostics collect --bootstrap-cluster-name ${CLUSTER_NAME} \
        --management-cluster-kubeconfig "${HOME}/.kube/config" \
        --management-cluster-context ${CLUSTER_KUBE_CONTEXT} \
        --management-cluster-name ${CLUSTER_NAME} \
        --workload-cluster-infra docker \
        --workload-cluster-kubeconfig "${HOME}/.kube/config" \
        --workload-cluster-context ${CLUSTER_KUBE_CONTEXT} \
        --workload-cluster-name ${CLUSTER_NAME} \
        --output-dir "${OUTPUT_DIR}" || {
            echo "Error running tanzu diagnostics collect command!"
            exit 1
        }
else
    "${TANZU_DIAGNOSTICS_BIN}" collect --bootstrap-cluster-name ${CLUSTER_NAME} \
        --management-cluster-kubeconfig "${HOME}/.kube/config" \
        --management-cluster-context ${CLUSTER_KUBE_CONTEXT} \
        --management-cluster-name ${CLUSTER_NAME} \
        --workload-cluster-infra docker \
        --workload-cluster-kubeconfig "${HOME}/.kube/config" \
        --workload-cluster-context ${CLUSTER_KUBE_CONTEXT} \
        --workload-cluster-name ${CLUSTER_NAME} \
        --output-dir "${OUTPUT_DIR}" || {
            echo "Error running tanzu diagnostics collect command!"
            exit 1
        }
fi

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
kubectl config delete-context ${CLUSTER_KUBE_CONTEXT} || {
    echo "Failed deleting kind cluster kube context. Please delete it manually"
}
kind delete cluster --name ${CLUSTER_NAME} || {
    echo "Failed deleting kind cluster. Please delete it manually"
}


rm -rfv "${OUTPUT_DIR}"

if [[ ${errors} -gt 0 ]]; then
    exit 1
fi
