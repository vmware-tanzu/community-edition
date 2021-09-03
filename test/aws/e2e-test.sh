#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
TCE_REPO_PATH="${MY_DIR}"/../../
# shellcheck source=test/util/utils.sh
source "${TCE_REPO_PATH}"/test/util/utils.sh

function test-gatekeeper {
    echo "Installing Gatekeeper..."
    gatekeeper_version=$(tanzu package available list gatekeeper.community.tanzu.vmware.com -o json | jq -r '.[0].version | select(. != null)')
    tanzu package install gatekeeper --package-name gatekeeper.community.tanzu.vmware.com --version "${gatekeeper_version}" || { error "Gatekeeper installation failed. TEST FAILED."; exit 1; }
    # Added this as it takes time to create namespace for Gatekeeper
    sleep 10s
    echo "Verifying Gatekeeper installation..."
    kubectl wait --for=condition=ready pod --all -n gatekeeper-system --timeout=300s || { error "Timed out waiting for Gatekeeper pods to come up. TEST FAILED."; exit 1; }
    echo "Applying constraint template..."
    kubectl apply -f "${TCE_REPO_PATH}"/test/gatekeeper/constraint-template.yaml || { error "Unexpected error. TEST FAILED."; exit 1; }
    echo "Verifying creation of k8srequiredlabels CRD..."
    kubectl get crds | grep -i k8srequiredlabels || { error "Unexpected error. TEST FAILED."; exit 1; }
    echo "Creating constraint..."
    kubectl apply -f "${TCE_REPO_PATH}"/test/gatekeeper/constraint.yaml || { error "Unexpected error. TEST FAILED."; exit 1; }
    echo "Creating test namespace..."
    # It takes time for Gatekeeper webhook service to come up. Added retires to get around Internal Server Error.
    retries=1
    while [ $retries -le 5 ]
    do
        error_message=$(kubectl create ns test 2>&1)
        if echo "$error_message" | grep 'All namespaces must have an owner label'; then
            echo "Expected failure"
            echo "Creating test namespace with owner label..."
            kubectl apply -f "${TCE_REPO_PATH}"/test/gatekeeper/test-namespace.yaml || { error "TEST FAILED."; exit 1; }
            printf '\E[32m'; echo "TEST PASSED!"; printf '\E[0m'
            return
        else
            echo "Retrying..."
            sleep 60s
        fi
        retries=$(( retries+1 ))
    done
    error "TEST FAILED"; exit 1;
}


echo "Starting Gatekeeper test..."
test-gatekeeper



