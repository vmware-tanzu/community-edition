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
tanzu package installed list | grep "contour.community.tanzu.vmware.com" || {
    version=$(tanzu package available list contour.community.tanzu.vmware.com | tail -n 1 | awk '{print $2}')
    tanzu package install contour --package-name contour.community.tanzu.vmware.com -f "${MY_DIR}"/../../contour/"${version}"/contour-values.yaml --version "${version}"
}

tanzu package installed list | grep "knative-serving.community.tanzu.vmware.com" || {
    version=$(tanzu package available list knative-serving.community.tanzu.vmware.com | tail -n 1 | awk '{print $2}')
    tanzu package install knative-serving --package-name knative-serving.community.tanzu.vmware.com -f "${MY_DIR}"/knative-serving-values.yaml --version "${version}"
}

NAMESPACE_SUFFIX=${RANDOM}
NAMESPACE="knative-serving-${NAMESPACE_SUFFIX}"
kubectl create ns ${NAMESPACE}

kubectl apply -n ${NAMESPACE} --filename "${MY_DIR}"/knative-service.yaml
echo "Waiting for pod to get scale down..."
sleep 10s

kubectl wait --for=delete --all pod -n ${NAMESPACE} --timeout=360s || {
    packageCleanup knative-serving contour
    namespaceCleanup ${NAMESPACE}
    failureMessage knative-serving
}

packageCleanup knative-serving contour
namespaceCleanup ${NAMESPACE}
successMessage knative-serving
