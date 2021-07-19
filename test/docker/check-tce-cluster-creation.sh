#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e # Fail on errors
set -x # See what commands are running

kube_context=$1

if [ -z "$kube_context" ]; then
    echo "Error: Kube context name not provided. Please provide kube context name"
    exit 1
fi

kubectl config use-context "$kube_context"

kubectl cluster-info

kubectl get nodes

kubectl get pods -A
