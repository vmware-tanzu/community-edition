#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e
set -x

# Note: This script supports only Linux(Debian/Ubuntu) and MacOS
# Following environment variables are expected to be exported before running the script
# VSPHERE_CONTROL_PLANE_ENDPOINT - virtual and static IP for the cluster's control plane nodes
# VSPHERE_SERVER - private IP of the vcenter server
# JUMPER_SSH_HOST_IP - public IP address to access the Jumper host for SSH
# JUMPER_SSH_USERNAME - username to access the Jumper host for SSH
# JUMPER_SSH_PRIVATE_KEY - private key to access to access the Jumper host for SSH

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

"${MY_DIR}"/run-proxy-to-vcenter-server-and-control-plane.sh

trap '{ "${MY_DIR}"/stop-proxy-to-vcenter-server-and-control-plane.sh; }' EXIT
