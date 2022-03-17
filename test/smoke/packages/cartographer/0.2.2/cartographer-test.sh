#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

TCE_REPO_PATH="$(git rev-parse --show-toplevel)"
# shellcheck source=test/smoke/packages/utils/smoke-tests-utils.sh
source "${TCE_REPO_PATH}/test/smoke/packages/utils/smoke-tests-utils.sh"

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Checking package is installed or not
tanzu package installed list | grep "cert-manager.community.tanzu.vmware.com" || {
    version=$(tanzu package available list cert-manager.community.tanzu.vmware.com | tail -n 1 | awk '{print $2}')
    tanzu package install cert-manager --package-name cert-manager.community.tanzu.vmware.com --version "${version}"
}

tanzu package installed list | grep "cartographer.community.tanzu.vmware.com" || {
    version=$(tanzu package available list cartographer.community.tanzu.vmware.com | tail -n 1 | awk '{print $2}')
    tanzu package install cartographer --package-name cartographer.community.tanzu.vmware.com --version "${version}"
}


NAMESPACE_SUFFIX=${RANDOM}
NAMESPACE="cartographer-${NAMESPACE_SUFFIX}"
kubectl create ns ${NAMESPACE}

kubectl apply -n ${NAMESPACE} --filename "${MY_DIR}"/testdata.yaml

for sleep_duration in {1..10}; do
    echo "sleeping ${sleep_duration}s to wait for configmap"
    sleep "$sleep_duration"

    kubectl get configmap workload-test-basic && {
        packageCleanup cartographer cert-manager
        namespaceCleanup ${NAMESPACE}
        successMessage cartographer
        exit 0
    }
done

packageCleanup cartographer cert-manager
namespaceCleanup ${NAMESPACE}
failureMessage cartographer
