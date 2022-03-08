#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script clones TCE repo and builds the latest release

set -x
set -e

TCE_REPO_PATH="$(git rev-parse --show-toplevel)"
# shellcheck source=test/util/utils.sh
source "${TCE_REPO_PATH}/test/util/utils.sh"

BUILD_OS=$(uname -s)
export BUILD_OS

# required input
if [[ -z "${BUILD_VERSION}" ]]; then
    echo "BUILD_VERSION is not set"
    exit 1
fi

# Check for nightly build
TODAY=$(date +%F)
DOWNLOAD_GCP_URI="https://storage.googleapis.com/tce-cli-plugins-staging/build-daily/${TODAY}/tce-linux-amd64-${BUILD_VERSION}.tar.gz"
if [[ $BUILD_OS == "Darwin" ]]; then
    if [[ "$BUILD_ARCH" == "x86_64" ]]; then
        DOWNLOAD_GCP_URI="https://storage.googleapis.com/tce-cli-plugins-staging/build-daily/${TODAY}/tce-darwin-amd64-${BUILD_VERSION}.tar.gz"
    else
        DOWNLOAD_GCP_URI="https://storage.googleapis.com/tce-cli-plugins-staging/build-daily/${TODAY}/tce-darwin-arm64-${BUILD_VERSION}.tar.gz"
    fi
fi

rm -f tce-daily-build.tar.gz
curl -L "${DOWNLOAD_GCP_URI}" -o tce-daily-build.tar.gz

set +e
IS_XML=$(file tce-daily-build.tar.gz | grep XML)
set -e

if [ "${IS_XML}" == "" ]; then
    # extract the daily
    tar -xvf tce-daily-build.tar.gz --one-top-level=tce-daily-build --strip-components 1

    pushd tce-daily-build || exit 1
else
    # delete the file
    rm -f tce-daily-build.tar.gz

    # No daily... So build TCE
    echo "Building TCE release..."
    make release || { error "TCE BUILD FAILED!"; exit 1; }

    echo "Installing TCE release"
    if [[ $BUILD_OS == "Linux" ]]; then
        pushd release/tce-linux-amd64*/ || exit 1
    elif [[ $BUILD_OS == "Darwin" ]]; then
        if [[ "$BUILD_ARCH" == "x86_64" ]]; then
            pushd release/tce-darwin-amd64*/ || exit 1
        else
            pushd release/tce-darwin-arm64*/ || exit 1
        fi
    fi
fi

./uninstall.sh || { error "TCE CLEANUP (UNINSTALLATION) FAILED!"; exit 1; }
./install.sh || { error "TCE INSTALLATION FAILED!"; exit 1; }
popd || exit 1
echo "TCE version..."
tanzu management-cluster version || { error "Unexpected failure during TCE installation"; exit 1; }
