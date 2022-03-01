#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
TCE_REPO_PATH="$(git rev-parse --show-toplevel)"

TCE_VERSION="v0.10.0"

echo "Installing TCE ${TCE_VERSION}"

BUILD_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
TCE_RELEASE_TAR_BALL="tce-${BUILD_OS}-amd64-${TCE_VERSION}.tar.gz"
TCE_RELEASE_DIR="tce-${BUILD_OS}-amd64-${TCE_VERSION}"
INSTALLATION_DIR="${MY_DIR}/tce-installation"

"${TCE_REPO_PATH}"/hack/get-tce-release.sh ${TCE_VERSION} "${BUILD_OS}"-amd64

mkdir -p "${INSTALLATION_DIR}"
tar xzvf "${TCE_RELEASE_TAR_BALL}" --directory="${INSTALLATION_DIR}"

"${INSTALLATION_DIR}"/"${TCE_RELEASE_DIR}"/install.sh || { error "Unexpected failure during TCE installation"; exit 1; }

echo "TCE version: "
tanzu unmanaged-cluster version || { error "Unexpected failure during TCE installation"; exit 1; }

TANZU_DIAGNOSTICS_PLUGIN_DIR=${MY_DIR}/..
TANZU_DIAGNOSTICS_BIN=${MY_DIR}/tanzu-diagnostics-e2e-bin

echo "Entering ${TANZU_DIAGNOSTICS_PLUGIN_DIR} directory to build tanzu diagnostics plugin"
pushd "${TANZU_DIAGNOSTICS_PLUGIN_DIR}"

go build -o "${TANZU_DIAGNOSTICS_BIN}" -v

echo "Finished building tanzu diagnostics plugin. Leaving ${TANZU_DIAGNOSTICS_PLUGIN_DIR}"
popd
