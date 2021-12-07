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
VERSION="0.7.2"

SUDO_CMD="sudo"
if [[ "${TCE_CI_BUILD}" == "true" ]]; then
  SUDO_CMD=""
fi

CMD="shellcheck"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    curl -LO https://github.com/koalaman/shellcheck/releases/download/v${VERSION}/shellcheck-v${VERSION}.linux.x86_64.tar.xz
    tar -xvf shellcheck-v${VERSION}.linux.x86_64.tar.xz
    pushd "./shellcheck-v${VERSION}" || exit 1
    chmod +x ${CMD}
    ${SUDO_CMD} install ./${CMD} /usr/local/bin
    popd || exit 1
    rm -rf ./shellcheck-v${VERSION}
    rm shellcheck-v${VERSION}.linux.x86_64.tar.xz
    ;;
  Darwin)
    case "${BUILD_ARCH}" in
      x86_64)
        curl -LO https://github.com/koalaman/shellcheck/releases/download/v${VERSION}/shellcheck-v${VERSION}.darwin.x86_64.tar.xz
        tar -xvf shellcheck-v${VERSION}.darwin.x86_64.tar.xz
        pushd "./shellcheck-v${VERSION}" || exit 1
        chmod +x ${CMD}
        ${SUDO_CMD} install ./${CMD} /usr/local/bin
        popd || exit 1
        rm -rf ./shellcheck-v${VERSION}
        rm shellcheck-v${VERSION}.linux.x86_64.tar.xz
        ;;
      arm64)
        brew install shellcheck
        ;;
    esac
    ;;
esac
fi
