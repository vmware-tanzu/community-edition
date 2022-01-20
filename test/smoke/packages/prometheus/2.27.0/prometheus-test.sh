#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

TCE_REPO_PATH="$(git rev-parse --show-toplevel)"
source "${TCE_REPO_PATH}/test/smoke/packages/utils/smoke-tests-utils.sh"


# Checking package is installed or not
tanzu package installed list | grep "local-path-storage.community.tanzu.vmware.com" || {
  version=$(tanzu package available list local-path-storage.community.tanzu.vmware.com | tail -n 1 | awk '{print $2}')
  tanzu package install local-path-storage --package-name local-path-storage.community.tanzu.vmware.com --version "${version}"
}


tanzu package installed list | grep "prometheus.community.tanzu.vmware.com" || {
  version=$(tanzu package available list prometheus.community.tanzu.vmware.com | tail -n 1 | awk '{print $2}')
  tanzu package install prometheus --package-name prometheus.community.tanzu.vmware.com --version "${version}"
}

prometheus_Pod="$(kubectl get pod -n prometheus | grep "prometheus-kube-state-metrics" | awk '{print $1}')"

kubectl port-forward "${prometheus_Pod}"  -n prometheus  8080:8080 &
sleep 5s

curl localhost:8080/metrics | less | grep "kube_deployment_created gauge" || {
    packageCleanup prometheus local-path-storage 
    failure_msg prometheus
}

curl localhost:8080/healthz | grep "OK" || {
    packageCleanup prometheus local-path-storage 
    failureMessage prometheus
}

packageCleanup prometheus local-path-storage 
successMessage prometheus
