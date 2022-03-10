#!/usr/bin/env bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

# Clean up func

function packageCleanup() {
   arr=("$@")
   for i in "${arr[@]}";
      do
          tanzu package installed delete "$i" -y
      done
}

function namespaceCleanup() {
    kubectl delete ns "$1"
}

function successMessage() {
    printf '\E[32m'; echo "$1 passed"; printf '\E[0m'
}

function failureMessage() {
    printf '\E[31m'; echo "$1 failed"; printf '\E[0m'
    exit 1
}
