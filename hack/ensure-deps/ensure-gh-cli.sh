#!/usr/bin/env bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script can be used to check your local development environment
# for necessary dependencies used in TCE

set -o nounset
set -o pipefail
set -o xtrace

BUILD_OS=$(uname 2>/dev/null || echo Unknown)
BUILD_ARCH=$(uname -m 2>/dev/null || echo Unknown)
GH_CLI_VERSION="2.2.0"

CMD="gh"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    curl -LO https://github.com/cli/cli/releases/download/v${GH_CLI_VERSION}/gh_${GH_CLI_VERSION}_linux_amd64.tar.gz
    tar -xvf gh_${GH_CLI_VERSION}_linux_amd64.tar.gz
    pushd "./gh_${GH_CLI_VERSION}_linux_amd64/bin" || exit 1
    chmod +x ${CMD}
    sudo install ./${CMD} /usr/local/bin
    popd || exit 1
    rm -rf ./gh_${GH_CLI_VERSION}_linux_amd64
    rm gh_${GH_CLI_VERSION}_linux_amd64.tar.gz
    ;;
  Darwin)
    case "${BUILD_ARCH}" in
      x86_64)
        curl -LO https://github.com/cli/cli/releases/download/v${GH_CLI_VERSION}/gh_${GH_CLI_VERSION}_macOS_amd64.tar.gz
        tar -xvf gh_${GH_CLI_VERSION}_macOS_amd64.tar.gz
        pushd "./gh_${GH_CLI_VERSION}_macOS_amd64/bin" || exit 1
        chmod +x ${CMD}
        sudo install ./${CMD} /usr/local/bin
        popd || exit 1
        rm -rf ./gh_${GH_CLI_VERSION}_macOS_amd64
        rm gh_${GH_CLI_VERSION}_macOS_amd64.tar.gz
        ;;
      arm64)
        brew install gh
        ;;
    esac
    ;;
esac
fi
