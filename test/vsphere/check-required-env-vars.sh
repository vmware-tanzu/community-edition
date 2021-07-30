#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e

declare -a required_env_vars=("VSPHERE_CONTROL_PLANE_ENDPOINT"
"VSPHERE_SERVER"
"VSPHERE_SSH_AUTHORIZED_KEY"
"VSPHERE_USERNAME"
"VSPHERE_PASSWORD"
"VSPHERE_DATACENTER"
"VSPHERE_DATASTORE"
"VSPHERE_FOLDER"
"VSPHERE_NETWORK"
"VSPHERE_RESOURCE_POOL"
"JUMPER_SSH_HOST_IP"
"JUMPER_SSH_USERNAME"
"JUMPER_SSH_PRIVATE_KEY")

for env_var in "${required_env_vars[@]}"
do
    if [ -z "$(printenv "${env_var}")" ]; then
        echo "Environment variable ${env_var} is empty! It's a required environment variable, please set it"
        exit 1
    fi
done
