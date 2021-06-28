#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e
set -x

if [[ -z "$(command -v jq)" ]]; then
    curl -L https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 -o jq
    chmod +x jq
    sudo mv jq /usr/local/bin/
fi
