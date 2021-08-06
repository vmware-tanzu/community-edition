#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e
set -x

installation_error_message="Unable to automatically install sshuttle for this platform. Please install sshuttle."

build_os=$(uname -s)

if [[ -z "$(command -v sshuttle)" ]]; then
    if [[ "$build_os" == "Linux" ]]; then
        distro_id=$(awk -F= '/^ID=/{print $2}' /etc/os-release)

        if [[ "${distro_id}" == "ubuntu" || "${distro_id}" == "debian" ]]; then
            sudo apt-get install sshuttle --yes
        else
            echo "${installation_error_message}"
            echo "Exiting..."
            exit 1
        fi
    elif [[ "$build_os" == "Darwin" ]]; then
        brew install sshuttle
    else
        echo "${installation_error_message}"
        echo "Exiting..."
        exit 1
    fi
fi
