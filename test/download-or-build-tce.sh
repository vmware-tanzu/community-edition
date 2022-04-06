#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script clones TCE repo and builds the latest release

set -x
set -e

TCE_REPO_PATH="$(git rev-parse --show-toplevel)"

# shellcheck source=test/util/utils.sh
source "${TCE_REPO_PATH}/test/util/utils.sh"

# defaults
DAILY_BUILD="${DAILY_BUILD:-""}"

# required input
if [[ -z "${BUILD_VERSION}" ]]; then
    echo "BUILD_VERSION is not set"
    exit 1
fi

# Check for nightly build
if [[ "${DAILY_BUILD}" != "" ]]; then
    echo "Attempt to download the daily build for ${DAILY_BUILD}."
    DOWNLOAD_GCP_URI="https://storage.googleapis.com/tce-cli-plugins-staging/build-daily/${DAILY_BUILD}/tce-${GOHOSTOS}-${GOHOSTARCH}-${BUILD_VERSION}.tar.gz"

    rm -f tce-daily-build.tar.gz
    curl -L "${DOWNLOAD_GCP_URI}" -o tce-daily-build.tar.gz

    set +e
    IS_XML=$(file tce-daily-build.tar.gz | grep XML)
    set -e

    if [ "${IS_XML}" != "" ]; then
        echo "Fatal error.... downloaded file is a HTML error code. Usually a 404."
        exit 1
    fi
else
    echo "Building TCE from source..."
    make release ENVS="${GOHOSTOS}-${GOHOSTARCH}" || { error "TCE BUILD FAILED!"; exit 1; }
fi

TCE_FOLDER=$(find ./release -type d -name "tce-*" -print0 | tr -d '\0')
pushd "${TCE_FOLDER}" || exit 1

    ./uninstall.sh || { error "TCE CLEANUP (UNINSTALLATION) FAILED!"; exit 1; }
    ALLOW_INSTALL_AS_ROOT=true ./install.sh || { error "TCE INSTALLATION FAILED!"; exit 1; }

popd || exit 1

echo "TCE version..."
tanzu management-cluster version || { error "Unexpected failure during TCE installation"; exit 1; }
