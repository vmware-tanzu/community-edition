#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o nounset
set -o pipefail
set -o xtrace

FORCE_UPDATE_PLUGIN="${FORCE_UPDATE_PLUGIN:-"false"}"
GOHOSTOS=$(go env GOHOSTOS)
GOHOSTARCH=$(go env GOHOSTARCH)

if [[ -z "${TCE_SCRATCH_DIR}" ]]; then
    echo "TCE_SCRATCH_DIR is not set"
    exit 1
fi
if [[ -z "${TCE_BUILD_VERSION}" ]]; then
    echo "TCE_BUILD_VERSION is not set"
    exit 1
fi

# check if the tanzu CLI already exists and remove it to avoid conflicts
TANZU_BIN_PATH=$(command -v tanzu)
if [[ -n "${TANZU_BIN_PATH}" ]]; then
    echo "Tanzu CLI is already installed"

   if [ "${FORCE_UPDATE_PLUGIN}" == "true" ] && [ -d "${TCE_SCRATCH_DIR}/tanzu-framework" ]; then
      # For TF 0.17.0 or higher
      # tanzu plugin install all --local "${TCE_SCRATCH_DIR}/tanzu-framework/build/${GOHOSTOS}-${GOHOSTARCH}-default"
      # For 0.11.2
      pushd "${TCE_SCRATCH_DIR}/tanzu-framework" || exit 1
        mkdir -p "${XDG_CONFIG_HOME}/tanzu-plugins"
        find "./build/${GOHOSTOS}-${GOHOSTARCH}-${DISCOVERY_NAME}/." -maxdepth 1 -type d | grep -E -v "/.$" | xargs -I _ cp -rf _ "${XDG_CONFIG_HOME}/tanzu-plugins"
      popd || exit 1

      tanzu plugin install builder
      tanzu plugin install codegen
      tanzu plugin install cluster
      tanzu plugin install kubernetes-release
      tanzu plugin install login
      tanzu plugin install management-cluster
      tanzu plugin install package
      tanzu plugin install pinniped-auth
      tanzu plugin install secret
    fi

    exit 0
fi

# Check and change directories to the pre-built tanzu binaries
if [ ! -d "${TCE_SCRATCH_DIR}" ] || [ ! -d "${TCE_SCRATCH_DIR}/tanzu-framework" ]; then
    echo "Error! No Tanzu Framework build directory found"
    echo "Please run \`make build-cli\` first to clone the repository"
fi

pushd "${TCE_SCRATCH_DIR}/tanzu-framework" || exit 1

BUILD_SHA="$(git describe --match="$(git rev-parse --short HEAD)" --always)"
sed -i.bak -e "s| --dirty||g" ./Makefile && rm ./Makefile.bak

# allow unstable (non-GA) version plugins
if [[ "${TCE_BUILD_VERSION}" == *"-"* ]]; then
make controller-gen
make set-unstable-versions
fi

 #Only do an install if the environments to build contain the current host OS.
 #The tanzu-framework `build-install-cli-all` target always uses the current host OS, and if that's not being built it will fail.
if [[ "$ENVS" == *"${GOHOSTOS}-${GOHOSTARCH}"* ]]; then
    BUILD_SHA=${BUILD_SHA} BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} TCE_BUILD_VERSION=${TCE_BUILD_VERSION}  make install-cli

    if [ "${FORCE_UPDATE_PLUGIN}" == "true" ] ||
      [ ! -d "${XDG_CONFIG_HOME}/tanzu-plugins" ] && [ -d "${TCE_SCRATCH_DIR}/tanzu-framework/build" ]; then
        # For TF 0.17.0 or higher
        # tanzu plugin install all --local "${TCE_SCRATCH_DIR}/tanzu-framework/build/${GOHOSTOS}-${GOHOSTARCH}-default"
        # For 0.11.2
        mkdir -p "${XDG_CONFIG_HOME}/tanzu-plugins"
        find "./build/${GOHOSTOS}-${GOHOSTARCH}-${DISCOVERY_NAME}/." -maxdepth 1 -type d | grep -E -v "/.$" | xargs -I _ cp -rf _ "${XDG_CONFIG_HOME}/tanzu-plugins"

        tanzu plugin install builder
        tanzu plugin install codegen
        tanzu plugin install cluster
        tanzu plugin install kubernetes-release
        tanzu plugin install login
        tanzu plugin install management-cluster
        tanzu plugin install package
        tanzu plugin install pinniped-auth
        tanzu plugin install secret
    fi
fi

popd || exit 1
