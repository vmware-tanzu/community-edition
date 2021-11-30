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

SUDO_CMD="sudo"
if [[ "${TCE_CI_BUILD}" == "true" ]]; then
  SUDO_CMD=""
fi

CMD="aws"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
    unzip awscliv2.zip
    ${SUDO_CMD} ./aws/install -i
    rm -rf ./aws
    rm awscliv2.zip
    ;;
  Darwin)
    curl "https://awscli.amazonaws.com/AWSCLIV2.pkg" -o "awscliv2.pkg"
    ${SUDO_CMD} installer -pkg awscliv2.pkg -target /
    rm -rf ./aws
    rm awscliv2.pkg
    ;;
esac
fi
