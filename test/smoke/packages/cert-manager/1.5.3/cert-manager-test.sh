#!/usr/bin/env bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

TCE_REPO_PATH="$(git rev-parse --show-toplevel)"
source "${TCE_REPO_PATH}/test/smoke/packages/utils/smoke-tests-utils.sh"

# Checking package is installed or not
tanzu package installed list | grep "cert-manager.community.tanzu.vmware.com" || {
  version=$(tanzu package available list cert-manager.community.tanzu.vmware.com | tail -n 1 | awk '{print $2}')
  tanzu package install cert-manager --package-name cert-manager.community.tanzu.vmware.com --version "${version}"
}

# Providing prerequisite 

NAMESPACE_SUFFIX=${RANDOM}
NAMESPACE="cert-manager-${NAMESPACE_SUFFIX}"
kubectl create ns ${NAMESPACE}

cat <<EOF | kubectl apply --filename -
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: self-signed
  namespace: ${NAMESPACE}
spec:
  selfSigned: {}
EOF


cat <<EOF | kubectl apply --filename -
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: private-ca
  namespace: ${NAMESPACE}
spec:
  isCA: true
  duration: 2160h
  secretName: private-ca
  commonName: private-ca
  subject:
    organizations:
      - cert-manager
  issuerRef:
    name: self-signed
    kind: Issuer
    group: cert-manager.io
EOF

cat <<EOF | kubectl apply --filename -
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: private-ca
  namespace: ${NAMESPACE}
spec:
  ca:
    secretName: private-ca
EOF


cat <<EOF | kubectl apply --filename -
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-com
  namespace: ${NAMESPACE}
spec:
  secretName: example-com-tls
  issuerRef:
    name: private-ca
    kind: Issuer
  commonName: example.com
  dnsNames:
    - example.com
    - www.example.com
EOF

echo "Waiting for certificate/example-com to Ready..."
kubectl wait --for=condition=Ready certificate/example-com -n ${NAMESPACE} --timeout=300s || {
  packageCleanup cert-manager
  namespaceCleanup ${NAMESPACE}
  printf '\E[31m'; echo "failed waiting for certificate status to be ready"; printf '\E[0m';
  failureMessage cert-manager
}

packageCleanup cert-manager
namespaceCleanup ${NAMESPACE}
successMessage cert-manager
