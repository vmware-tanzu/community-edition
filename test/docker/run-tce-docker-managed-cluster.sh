#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e
set -x

TCE_REPO_PATH="$(git rev-parse --show-toplevel)"

"${TCE_REPO_PATH}/test/install-dependencies.sh"
"${TCE_REPO_PATH}/test/download-or-build-tce.sh"
"${TCE_REPO_PATH}/test/install-jq.sh"

random_id="${RANDOM}"

export MGMT_CLUSTER_NAME="management-cluster-${random_id}"
export GUEST_CLUSTER_NAME="guest-cluster-${random_id}"

tanzu management-cluster create -i docker --name ${MGMT_CLUSTER_NAME} -v 10 --plan dev --ceip-participation=false

# Check management cluster details
tanzu management-cluster get

# Get kube config of management cluster
tanzu management-cluster kubeconfig get ${MGMT_CLUSTER_NAME} --admin

"${TCE_REPO_PATH}/test/check-tce-cluster-creation.sh" ${MGMT_CLUSTER_NAME}-admin@${MGMT_CLUSTER_NAME}

tanzu cluster create ${GUEST_CLUSTER_NAME} --plan dev

tanzu cluster list

tanzu cluster kubeconfig get ${GUEST_CLUSTER_NAME} --admin

"${TCE_REPO_PATH}/test/check-tce-cluster-creation.sh" ${GUEST_CLUSTER_NAME}-admin@${GUEST_CLUSTER_NAME}

"${TCE_REPO_PATH}/test/add-tce-package-repo.sh"

# wait for packages to be available
sleep 10

tanzu package available list

tanzu package available list fluent-bit.community.tanzu.vmware.com

fluentbit_version=$(tanzu package available list fluent-bit.community.tanzu.vmware.com -o json | jq -r '.[0].version | select(. !=null)')

tanzu package install fluent-bit --package-name fluent-bit.community.tanzu.vmware.com --version "${fluentbit_version}"

tanzu package installed list

kubectl -n fluent-bit get all

kubectl get installedpackage,apps --all-namespaces

tanzu cluster delete ${GUEST_CLUSTER_NAME} --yes

num_of_clusters=$(tanzu cluster list -o json | jq 'length')

while [ "$num_of_clusters" != "0" ]
do
    echo "Waiting for workload cluster to get deleted..."
    sleep 2;
    num_of_clusters=$(tanzu cluster list -o json | jq 'length')
done

echo "Workload cluster ${GUEST_CLUSTER_NAME} successfully deleted"

tanzu management-cluster delete ${MGMT_CLUSTER_NAME} --yes

echo "Management cluster ${MGMT_CLUSTER_NAME} successfully deleted"
