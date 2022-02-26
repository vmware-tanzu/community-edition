#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

TCE_REPO_PATH="$(git rev-parse --show-toplevel)"
source "${TCE_REPO_PATH}/test/smoke/packages/utils/smoke-tests-utils.sh"

# Checking package is installed or not
tanzu package installed list | grep "fluent-bit.community.tanzu.vmware.com" || {
    version=$(tanzu package available list fluent-bit.community.tanzu.vmware.com | tail -n 1 | awk '{print $2}')
    tanzu package install fluent-bit --package-name fluent-bit.community.tanzu.vmware.com --version "${version}"
}

pod_name="$(kubectl get pods -n fluent-bit | tail -n 1 | awk '{print $1}')"
kubectl logs "${pod_name}" -n fluent-bit | grep "Fluent Bit v1.7.5" || {
    packageCleanup fluent-bit
    failureMessage fluent-bit
}

packageCleanup fluent-bit
successMessage fluent-bit
