#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

TCE_REPO_PATH="$(git rev-parse --show-toplevel)"
source "${TCE_REPO_PATH}/test/smoke/packages/utils/smoke-tests-utils.sh"

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Checking package is installed or not
tanzu package installed list | grep "contour.community.tanzu.vmware.com" || {
    version=$(tanzu package available list contour.community.tanzu.vmware.com | tail -n 1 | awk '{print $2}')
    tanzu package install contour --package-name contour.community.tanzu.vmware.com -f "${MY_DIR}"/contour-values.yaml --version "${version}"
}

# Providing prerequisite 

NAMESPACE_SUFFIX=${RANDOM}
NAMESPACE="contour-${NAMESPACE_SUFFIX}"
kubectl create ns ${NAMESPACE}
kubectl create deployment nginx-example --image nginx --namespace ${NAMESPACE}
kubectl create service clusterip nginx-example --tcp 8080:80 --namespace ${NAMESPACE}

kubectl apply -f - <<EOF
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  name: nginx-example-proxy
  namespace: ${NAMESPACE}
  labels:
    app: ingress
spec:
  virtualhost:
    fqdn: nginx-example.projectcontour.io
  routes:
  - conditions:
    - prefix: /
    services:
    - name: nginx-example
      port: 8080
EOF

sleep 10s
kubectl --namespace projectcontour port-forward svc/envoy 5436:80 &
sleep 5s

curl -s -H "Host: nginx-example.projectcontour.io" http://localhost:5436 | grep title || {
  packageCleanup contour
  namespaceCleanup ${NAMESPACE}
  failureMessage contour
}

packageCleanup contour
namespaceCleanup ${NAMESPACE}
successMessage contour
