#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e
set -x

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
TCE_REPO_PATH="${MY_DIR}"/..

# shellcheck source=test/util/utils.sh
source "${TCE_REPO_PATH}"/test/util/utils.sh

kube_context=$1

if [ -z "$kube_context" ]; then
    error "Error: Kube context name not provided. Please provide kube context name"
    exit 1
fi

kubectl config use-context "$kube_context" || {
    error "CONTEXT SWITCH TO CLUSTER FAILED!"
    exit 1
}

kubectl wait --for=condition=ready pod --all --all-namespaces --timeout=600s || {
    error "TIMED OUT WAITING FOR ALL PODS TO BE UP!"
    exit 1
}

kubectl cluster-info || {
    error "ERROR GETTING CLUSTER INFO!"
    exit 1
}

kubectl get nodes || {
    error "ERROR GETTING CLUSTER NODES!"
    exit 1
}

kubectl get pods -A || {
    error "ERROR GETTING ALL PODS IN CLUSTER!"
    exit 1
}
