#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o nounset
set -o pipefail
set -o xtrace

if [[ -z "${TCE_SCRATCH_DIR}" ]]; then
    echo "TCE_SCRATCH_DIR is not set"
    exit 1
fi

# Check and change directories to the pre-built tanzu binaries
if [ ! -d "${TCE_SCRATCH_DIR}" ] || [ ! -d "${TCE_SCRATCH_DIR}/tanzu-framework" ]; then
    echo "Error! No Tanzu Framework build directory found"
    echo "Please run \`make build-cli\` first to clone the repository"
fi

pushd "${TCE_SCRATCH_DIR}/tanzu-framework" || exit 1

BUILD_SHA="$(git describe --match="$(git rev-parse --short HEAD)" --always)"
sed -i.bak -e "s| --dirty||g" ./Makefile && rm ./Makefile.bak

 #Only do an install if the environments to build contain the current host OS.
 #The tanzu-framework `build-install-cli-all` target always uses the current host OS, and if that's not being built it will fail.
GOHOSTOS=$(go env GOHOSTOS)
GOHOSTARCH=$(go env GOHOSTARCH)
if [[ "$ENVS" == *"${GOHOSTOS}-${GOHOSTARCH}"* ]]; then
    BUILD_SHA=${BUILD_SHA} BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} make install-cli
    tanzu plugin install all --local "${TCE_SCRATCH_DIR}/tanzu-framework/build/${GOHOSTOS}-${GOHOSTARCH}-default"
fi

popd || exit 1
