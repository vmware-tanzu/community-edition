#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e
set -x

BUILD_OS=$(uname -s)
export BUILD_OS

if [[ -z "$(command -v jq)" ]]; then
    if [[ "$BUILD_OS" == "Linux" ]]; then
        curl -L https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 -o jq
        chmod +x jq
        sudo mv jq /usr/local/bin/
    elif [[ "$BUILD_OS" == "Darwin" ]]; then
        brew install jq
    fi
fi
