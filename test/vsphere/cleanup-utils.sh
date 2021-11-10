#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e

# Note: This script supports only Linux(Debian/Ubuntu) and MacOS
# Following environment variables are expected to be exported before running the script
# VSPHERE_SERVER - private IP of the vcenter server
# VSPHERE_USERNAME - vcenter username
# VSPHERE_PASSWORD - vcenter password

TCE_REPO_PATH="$(git rev-parse --show-toplevel)"

function install_govc {
    installation_error_message="Unable to automatically install govc for this platform. Please install govc."

    if [[ -z "$(command -v govc)" ]]; then
        {
            curl -L -o - \
                "https://github.com/vmware/govmomi/releases/latest/download/govc_$(uname -s)_$(uname -m).tar.gz" | \
                sudo tar -C /usr/local/bin -xvzf - govc
        } || echo "${installation_error_message}"
    fi
}

function govc_cleanup {
    vsphere_cluster_name=$1

    if [[ -z "${vsphere_cluster_name}" ]]; then
        echo "Cluster name not passed to govc_cleanup function. Usage example: govc_cleanup management-cluster-1234"
        exit 1
    fi

    declare -a required_env_vars=("VSPHERE_SERVER"
    "VSPHERE_USERNAME"
    "VSPHERE_PASSWORD")

    "${TCE_REPO_PATH}/test/vsphere/check-required-env-vars.sh" "${required_env_vars[@]}"

    # Install govc if is not already installed
    install_govc

    export GOVC_URL="${VSPHERE_USERNAME}:${VSPHERE_PASSWORD}@${VSPHERE_SERVER}"

    govc find -k -type m . -name "${vsphere_cluster_name}*" | \
      xargs -I{} -r govc vm.destroy -k -debug -dump {}
}
