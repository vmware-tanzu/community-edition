#!/usr/bin/env bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script can be used to check your local development environment
# for necessary dependencies used in TCE

i=0

if [[ -z "$(command -v go)" ]]; then
    echo "Missing go"
    ((i=i+1))
fi

# if [[ -z "$(command -v imgpkg)" ]]; then
#     echo "Missing imgpkg"
#     ((i=i+1))
# fi

# if [[ -z "$(command -v kbld)" ]]; then
#     echo "Missing kbld"
#     ((i=i+1))
# fi

# if [[ -z "$(command -v ytt)" ]]; then
#     echo "Missing ytt"
#     ((i=i+1))
# fi

# if [[ -z "$(command -v shellcheck)" ]]; then
#     echo "Missing shellcheck"
#     ((i=i+1))
# fi

if [[ -z "$(command -v docker)" ]]; then
    echo "Missing docker"
    ((i=i+1))
fi

# if [[ -z "$(command -v kubectl)" ]]; then
#     echo "Missing kubectl"
#     ((i=i+1))
# fi

if [[ $i -gt 0 ]]; then
    echo "Total missing: $i"
    echo "Please install dependencies to continue"
    exit 1
fi

echo "No missing dependencies!"
exit 0

