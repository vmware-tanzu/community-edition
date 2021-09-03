#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

kubectl config use-context mytest-admin@mytest

IS_KAPP_CRD_INSTALLED=$(kubectl get crds -A | grep apps.kappctrl.k14s.io)
if [[ "${IS_KAPP_CRD_INSTALLED}" == "" ]]; then
    echo "Kapp controller CRD is not installed"
    exit 1
fi
IS_KAPP_RUNNING=$(kubectl get pods -A | grep -E 'kapp-controller.+Running')
if [[ "${IS_KAPP_RUNNING}" == "" ]]; then
    echo "Kapp controller is not running"
    exit 1
fi

echo "Cluster created successfully!"
