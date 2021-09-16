#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e
set -x

# Helper functions
function error {
    printf '\E[31m'; echo "$@"; printf '\E[0m'
}

# This script works only in Linux and MacOS (Darwin) and amd64 (64 bit) platforms

TCE_VERSION="$1"

if [[ -z ${TCE_VERSION} ]]; then
    echo "TCE release version is not passed to the script!"
    echo "Please pass the TCE release version as the first argument to this script, for example like this ./test/fetch-and-install-tce-release.sh v0.7.0-rc.2"
    echo "You can find the available TCE release verions in the TCE GitHub releases page - https://github.com/vmware-tanzu/community-edition/releases/"
    exit 1
fi

if [[ -z ${GH_E2E_ACCESS_TOKEN} ]]; then
    echo "GitHub Access Token is not set!"
    echo "Please set the GitHub Access Token using \$GH_E2E_ACCESS_TOKEN environment variable"
    echo "You can create a GitHub Access Token using the following doc - https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token"
    exit 1
fi

BUILD_OS=$(uname -s | tr '[:upper:]' '[:lower:]')

if [[ "$BUILD_OS" != "linux" ]] && [[ "$BUILD_OS" != "darwin" ]]; then
    error "Installation on $BUILD_OS is not supported."
    exit 1
fi

TCE_RELEASE_TAR_BALL="tce-${BUILD_OS}-amd64-${TCE_VERSION}.tar.gz"
TCE_RELEASE_DIR="tce-${BUILD_OS}-amd64-${TCE_VERSION}"
MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
INSTALLATION_DIR="${MY_DIR}/tce-installation"

mkdir -p "${INSTALLATION_DIR}"

if [[ -z "$(command -v fetch)" ]]; then
    echo "Installing fetch CLI tool as it is not present..."
    curl -L "https://github.com/gruntwork-io/fetch/releases/download/v0.4.2/fetch_${BUILD_OS}_amd64" -o "${INSTALLATION_DIR}"/fetch
    chmod +x "${INSTALLATION_DIR}"/fetch
    sudo install "${INSTALLATION_DIR}"/fetch /usr/local/bin/fetch
fi

fetch --repo "https://github.com/vmware-tanzu/community-edition" \
    --tag "${TCE_VERSION}" \
    --release-asset "${TCE_RELEASE_TAR_BALL}" \
    --github-oauth-token "${GH_E2E_ACCESS_TOKEN}" \
    "${INSTALLATION_DIR}"

tar xzvf "${INSTALLATION_DIR}"/"${TCE_RELEASE_TAR_BALL}" --directory="${INSTALLATION_DIR}"


"${INSTALLATION_DIR}"/"${TCE_RELEASE_DIR}"/install.sh

rm -rf "${INSTALLATION_DIR}"