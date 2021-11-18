#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

TCE_REPO_PATH="$(git rev-parse --show-toplevel)"
# shellcheck source=test/util/utils.sh
source "${TCE_REPO_PATH}/test/smoke/packages/utils/smoke-tests-utils.sh"

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Checking package is installed or not
tanzu package installed list | grep "local-path-storage.community.tanzu.vmware.com" || {
  version=$(tanzu package available list local-path-storage.community.tanzu.vmware.com | tail -n 1 | awk '{print $2}')
  tanzu package install local-path-storage --package-name local-path-storage.community.tanzu.vmware.com --version "${version}"
}

tanzu package installed list | grep "contour.community.tanzu.vmware.com" || {
    version=$(tanzu package available list contour.community.tanzu.vmware.com | tail -n 1 | awk '{print $2}')
    tanzu package install contour --package-name contour.community.tanzu.vmware.com -f "${MY_DIR}"/../../contour/"${version}"/contour-values.yaml --version "${version}"
}

tanzu package installed list | grep "cert-manager.community.tanzu.vmware.com" || {
  version=$(tanzu package available list cert-manager.community.tanzu.vmware.com | tail -n 1 | awk '{print $2}')
  tanzu package install cert-manager --package-name cert-manager.community.tanzu.vmware.com --version "${version}"
}

tanzu package installed list | grep "grafana.community.tanzu.vmware.com" || {
    version=$(tanzu package available list grafana.community.tanzu.vmware.com | tail -n 1 | awk '{print $2}')
    tanzu package install grafana --package-name grafana.community.tanzu.vmware.com --version "${version}"
}

grafana_pod="$(kubectl get pod -A | grep "grafana" | awk '{print $2}')" 
kubectl port-forward "${grafana_pod}" -n grafana 56016:3000 &
sleep 5s

curl -I http://127.0.0.1:56016/api/health | grep "200" || {
  packageCleanup local-path-storage contour cert-manager grafana
  failureMessage grafana
}

packageCleanup local-path-storage contour cert-manager grafana
successMessage grafana
