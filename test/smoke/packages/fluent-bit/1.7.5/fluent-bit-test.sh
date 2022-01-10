#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

# Checking package is installed or not

tanzu package installed list | grep "fluent-bit.community.tanzu.vmware.com" || {
    version=$(tanzu package available list fluent-bit.community.tanzu.vmware.com | tail -n 1 | awk '{print $2}')
    tanzu package install fluent-bit --package-name fluent-bit.community.tanzu.vmware.com --version "${version}"
}

pod_name="$(kubectl get pods -n fluent-bit | tail -n 1 | awk '{print $1}')"
kubectl logs "${pod_name}" -n fluent-bit | grep "Fluent Bit v1.7.5" || {
    tanzu package installed delete fluent-bit -y
    printf '\E[31m'; echo "Fluent-bit failed"; printf '\E[0m'
    exit 1
}

tanzu package installed delete fluent-bit -y
printf '\E[32m'; echo "Fluent-bit Passed"; printf '\E[0m'
