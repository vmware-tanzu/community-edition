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
VERSION="0.30.0"
GOBIN=$(go env GOBIN)
GOPATH=$(go env GOPATH)
GOBINDIR="${GOBIN:-$GOPATH/bin}"

SUDO_CMD="sudo"
if [[ "${TCE_CI_BUILD}" == "true" ]]; then
  SUDO_CMD=""
fi

CMD="kbld"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    curl -LO https://github.com/vmware-tanzu/carvel-kbld/releases/download/v${VERSION}/kbld-linux-amd64
    chmod +x ${CMD}-linux-amd64
    mv ${CMD}-linux-amd64 ${CMD}
    ${SUDO_CMD} install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
  Darwin)
    go install github.com/k14s/kbld/cmd/${CMD}@v${VERSION}
    ${SUDO_CMD} install "${GOBINDIR}"/${CMD} /usr/local/bin
    ;;
esac
fi
