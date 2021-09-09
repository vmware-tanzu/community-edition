#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o nounset
set -o pipefail
set -o xtrace

if [[ -z "${TCE_RELEASE_DIR}" ]]; then
    echo "TCE_RELEASE_DIR is not set"
    exit 1
fi

# Check and change directories to the pre-built tanzu binaries
ROOT_REPO_DIR="${TCE_RELEASE_DIR}"

if [ ! -d "${ROOT_REPO_DIR}" ] || [ ! -d "${ROOT_REPO_DIR}/tanzu-framework" ]; then
    echo "Error! No Tanzu Framework build directory found"
    echo "Please run \`make build-cli\` first to clone the repository"
fi

cd "${ROOT_REPO_DIR}/tanzu-framework" || exit 1

BUILD_SHA="$(git describe --match="$(git rev-parse --short HEAD)" --always)"


 #Only do an install if the environments to build contain the current host OS.
 #The tanzu-framework `build-install-cli-all` target always uses the current host OS, and if that's not being built it will fail.
GOHOSTOS=$(go env GOHOSTOS)
if [[ "$ENVS" == *"${GOHOSTOS}"* ]]; then
    BUILD_SHA=${BUILD_SHA} BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} make ENVS="${ENVS}" install-cli-plugins install-cli
fi
