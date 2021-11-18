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

# Providing prerequisite 
NAMESPACE_SUFFIX=${RANDOM}
NAMESPACE="local-path-storage-${NAMESPACE_SUFFIX}"
kubectl create ns ${NAMESPACE}

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
kubectl apply -n ${NAMESPACE} -f "${MY_DIR}"/pvc.yaml
kubectl apply -n ${NAMESPACE} -f "${MY_DIR}"/pod.yaml

echo "Waiting for local-path-pvc to get BOUND..."
timeout=300
while [[ $(kubectl get pvc local-path-pvc -n ${NAMESPACE} -o 'jsonpath={.status.phase}') != "Bound" && ${timeout} -ne 0 ]]; do sleep 1 && timeout=${timeout}-1 ; done

echo "Waiting for pod volume-test to Ready..."
kubectl wait --for=condition=Ready pod volume-test -n ${NAMESPACE}

Pod_status="$(kubectl get pod volume-test -n ${NAMESPACE} -o 'jsonpath={..status.conditions[?(@.type=="Ready")].status}')"
PVC_status="$(kubectl get pvc local-path-pvc -n ${NAMESPACE} -o 'jsonpath={.status.phase}')"

if [ "${PVC_status}" != "Bound" ]
then
    printf '\E[31m'; echo "PVC did not bound to PV. check the storage class name"; printf '\E[0m';
fi

if [ "${Pod_status}" != "True" ]
then
    printf '\E[31m'; echo "Pod for PV is not in Running state"; printf '\E[0m';
fi

kubectl exec volume-test -n ${NAMESPACE} -- sh -c "echo local-path-storage-test > /data/test"


kubectl delete -n ${NAMESPACE} -f "${MY_DIR}"/pod.yaml 
kubectl apply -n ${NAMESPACE} -f "${MY_DIR}"/pod.yaml

echo "Waiting for pod volume-test to Ready..."
kubectl wait --for=condition=Ready pod volume-test -n ${NAMESPACE}

Output="$(kubectl exec volume-test -n ${NAMESPACE} cat /data/test)"
if [ "${Output}" != "local-path-storage-test" ]
then
    ffailureMessage local-path-storage 
else
    successMessage local-path-storage 
fi

kubectl delete pod/volume-test -n ${NAMESPACE}
kubectl delete persistentvolumeclaim/local-path-pvc -n ${NAMESPACE}
packageCleanup local-path-storage 
