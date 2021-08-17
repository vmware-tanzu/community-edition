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
# JUMPER_SSH_KNOWN_HOSTS_ENTRY - entry to put in the SSH client machine's (from where script is run) known_hosts file

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

"${MY_DIR}"/install-sshuttle.sh

export JUMPER_SSH_HOST_NAME=vmc-jumper-${CLUSTER_NAME}
export JUMPER_SSH_PRIVATE_KEY_LOCATION=${HOME}/.ssh/jumper_private_key

ssh_config_file_template="${MY_DIR}"/ssh-config-template

ssh_config_file=~/.ssh/config
ssh_known_hosts_file=~/.ssh/known_hosts

mkdir -p "$(dirname ${ssh_config_file})"
touch ${ssh_config_file}

mkdir -p "$(dirname ${ssh_known_hosts_file})"
touch ${ssh_known_hosts_file}

envsubst < "${ssh_config_file_template}" >> ${ssh_config_file}

rm -rfv "${JUMPER_SSH_PRIVATE_KEY_LOCATION}"
mkdir -p "$(dirname "${JUMPER_SSH_PRIVATE_KEY_LOCATION}")"
touch "${JUMPER_SSH_PRIVATE_KEY_LOCATION}"

printenv 'JUMPER_SSH_PRIVATE_KEY' > "${JUMPER_SSH_PRIVATE_KEY_LOCATION}"
chmod 400 "${JUMPER_SSH_PRIVATE_KEY_LOCATION}"

printenv 'JUMPER_SSH_KNOWN_HOSTS_ENTRY' >> ${ssh_known_hosts_file}

sshuttle --daemon -vvvvvvvv --remote "${JUMPER_SSH_HOST_NAME}" "${VSPHERE_SERVER}"/32 "${VSPHERE_CONTROL_PLANE_ENDPOINT}"/32
