#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e
set -x

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

"${MY_DIR}"/install-jq.sh

echo "Adding TCE package repository..."

REPO_NAME="tce-main-latest"
REPO_URL="projects.registry.vmware.com/tce/main:stable"
REPO_NAMESPACE="default"

# TODO: Use stable version of the tce/main repo once https://github.com/vmware-tanzu/community-edition/issues/1250 is fixed
tanzu package repository add ${REPO_NAME} --namespace ${REPO_NAMESPACE} --url ${REPO_URL}

# Wait for reconciliation to happen within ~ 80 x 5 = 400 seconds . 80 iterations, 5 seconds sleep time.
# Check status every ~5 seconds interval
for (( i = 1 ; i <= 80 ; i++))
do
    repo_status=$(tanzu package repository get ${REPO_NAME} -o json | jq -r '.[0].status | select (. != null)')
    if [[ ${repo_status} == "Reconcile succeeded" ]]; then
        echo "TCE package repository added!"
        exit 0
    fi
    sleep 5
done

echo "Error: Timed out while TCE package repository was reconciling!"
exit 1
