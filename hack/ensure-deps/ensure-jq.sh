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
# BUILD_ARCH=$(uname -m 2>/dev/null || echo Unknown)
VERSION="1.6"

SUDO_CMD="sudo"
if [[ "${TCE_CI_BUILD}" == "true" ]]; then
  SUDO_CMD=""
fi

CMD="jq"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    curl -L "https://github.com/stedolan/jq/releases/download/jq-${VERSION}/jq-linux64" -o "jq"
    chmod +x jq
    ${SUDO_CMD} mv ${CMD} /usr/local/bin
    ;;
  Darwin)
    curl -L "https://github.com/stedolan/jq/releases/download/jq-${VERSION}/jq-linux64" -o "jq"
    chmod +x jq
    ${SUDO_CMD} mv ${CMD} /usr/local/bin
    ;;
esac
fi
