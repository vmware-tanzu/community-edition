#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e

declare -a required_env_vars=("AZURE_CLIENT_ID"
"AZURE_CLIENT_SECRET"
"AZURE_SSH_PUBLIC_KEY_B64"
"AZURE_SUBSCRIPTION_ID"
"AZURE_TENANT_ID")

for env_var in "${required_env_vars[@]}"
do
    if [ -z "$(printenv "${env_var}")" ]; then
        echo "Environment variable ${env_var} is empty! It's a required environment variable, please set it"
        exit 1
    fi
done
