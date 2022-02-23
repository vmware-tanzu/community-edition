#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

# Change directories to the parent directory of the one in which this
# script is located.
ROOT_REPO_DIR="$(dirname "${BASH_SOURCE[0]}")/../.."
cd "${ROOT_REPO_DIR}" || exit 1

if [[ -z "${BUILD_VERSION}" ]]; then
    echo "BUILD_VERSION is not set"
    exit 1
fi
if [[ -z "${FRAMEWORK_BUILD_VERSION}" ]]; then
    echo "FRAMEWORK_BUILD_VERSION is not set"
    exit 1
fi
if [[ -z "${TCE_SCRATCH_DIR}" ]]; then
    echo "TCE_SCRATCH_DIR is not set"
    exit 1
fi

# Required directories
RELEASE_ROOT_DIR="${ROOT_REPO_DIR}/release"
ROOT_TCE_BUILD_DIR="${ROOT_REPO_DIR}/build"
ROOT_FRAMEWORK_BUILD_DIR="${TCE_SCRATCH_DIR}/tanzu-framework/build"

FRAMEWORK_BUILD_VERSION="${FRAMEWORK_BUILD_VERSION:-latest}"
TCE_BUILD_VERSION="${BUILD_VERSION:-latest}"

rm -rf "${RELEASE_ROOT_DIR}"
mkdir -p "${RELEASE_ROOT_DIR}"

# change settings
chmod +x "${ROOT_REPO_DIR}/hack/release/install.sh"
chmod +x "${ROOT_REPO_DIR}/hack/release/uninstall.sh"

for env in ${ENVS}; do
    OS=$(cut -d"-" -f1 <<< "${env}")
    ARCH=$(cut -d"-" -f2 <<< "${env}")
    scriptextension=".sh" && [[ "${OS}" == "windows" ]] && scriptextension=".bat"
    executableextension="" && [[ "${OS}" == "windows" ]] && executableextension=".exe"

    PACKAGE_DIR="${RELEASE_ROOT_DIR}/tce-${env}-${BUILD_VERSION}"
    rm -rf "${PACKAGE_DIR}"
    mkdir -p "${PACKAGE_DIR}"

    cp -r "${ROOT_TCE_BUILD_DIR}/${env}-default" "${PACKAGE_DIR}"
    cp -r "${ROOT_FRAMEWORK_BUILD_DIR}/${env}-default" "${PACKAGE_DIR}"
    cp -f "${ROOT_FRAMEWORK_BUILD_DIR}/${env}-default/tanzu-core-${OS}_${ARCH}${executableextension}" "${PACKAGE_DIR}/tanzu${executableextension}"
    cp -f "${ROOT_REPO_DIR}/hack/release/install${scriptextension}" "${PACKAGE_DIR}"
    cp -f "${ROOT_REPO_DIR}/hack/release/uninstall${scriptextension}" "${PACKAGE_DIR}"

    chown -R "$(id -u -n)":"$(id -g -n)" "${PACKAGE_DIR}"

    # packaging
    rm -f "tce-${env}-*.tar.gz"
    pushd "./release" || exit 1
        if [[ $env = *windows* ]]
        then
            zip -r "tce-${env}-${TCE_BUILD_VERSION}.zip" "tce-${env}-${BUILD_VERSION}"
        else
            tar -czvf "tce-${env}-${TCE_BUILD_VERSION}.tar.gz" "tce-${env}-${BUILD_VERSION}"
        fi
    popd || exit 1
done
