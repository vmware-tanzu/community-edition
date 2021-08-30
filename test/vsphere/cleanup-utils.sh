#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e

# Note: This script supports only Linux(Debian/Ubuntu) and MacOS
# Following environment variables are expected to be exported before running the script
# VSPHERE_SERVER - private IP of the vcenter server
# VSPHERE_USERNAME - vcenter username
# VSPHERE_PASSWORD - vcenter password

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
    # Install govc if is not already installed
    install_govc

    export GOVC_URL="${VSPHERE_USERNAME}:${VSPHERE_PASSWORD}@${VSPHERE_SERVER}"

    # Delete nodes with the name of the cluster as part of the node / VM name
    govc find -k -type m . -name "${CLUSTER_NAME}*" | \
        xargs govc vm.destroy -k -debug -dump
}
