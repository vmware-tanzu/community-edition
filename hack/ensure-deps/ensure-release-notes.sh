#!/usr/bin/env bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script can be used to check your local development environment
# for necessary dependencies used in TCE

set -o nounset
set -o pipefail
set -o xtrace

TCE_CI_BUILD="${TCE_CI_BUILD:-""}"
BUILD_OS=$(uname 2>/dev/null || echo Unknown)
BUILD_ARCH=$(uname -m 2>/dev/null || echo Unknown)
VERSION="0.10.0-tce.3"

SUDO_CMD="sudo"
if [[ "${TCE_CI_BUILD}" == "true" ]]; then
  SUDO_CMD=""
fi

CMD="release-notes"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    curl -LO https://storage.googleapis.com/tce-tanzu-cli-plugins/build-tools/release-notes/v${VERSION}/release-notes-linux-amd64
    mv release-notes-linux-amd64 ${CMD}
    chmod +x ${CMD}
    ${SUDO_CMD} install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
  Darwin)
    case "${BUILD_ARCH}" in
      x86_64)
        curl -LO https://storage.googleapis.com/tce-tanzu-cli-plugins/build-tools/release-notes/v${VERSION}/release-notes-darwin-amd64
        mv release-notes-darwin-amd64 ${CMD}
        chmod +x ${CMD}
        ${SUDO_CMD} install ./${CMD} /usr/local/bin
        rm ./${CMD}
        ;;
     arm64)
        curl -LO https://storage.googleapis.com/tce-tanzu-cli-plugins/build-tools/release-notes/v${VERSION}/release-notes-darwin-arm64
        mv release-notes-darwin-arm64 ${CMD}
        chmod +x ${CMD}
        ${SUDO_CMD} install ./${CMD} /usr/local/bin
        rm ./${CMD}
        ;;
    esac
    ;;
esac
fi
