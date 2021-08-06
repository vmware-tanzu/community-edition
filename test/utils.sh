#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e

function error {
    printf '\E[31m'; echo "$@"; printf '\E[0m'
}
