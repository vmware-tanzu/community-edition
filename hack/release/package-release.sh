#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

ENVS="${ENVS:-"linux-amd64 windows-amd64 darwin-amd64 darwin-arm64"}"
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
if [[ -z "${TCE_RELEASE_DIR}" ]]; then
    echo "TCE_RELEASE_DIR is not set"
    exit 1
fi

# Required directories
BUILD_ROOT_DIR="${ROOT_REPO_DIR}/build"
FRAMEWORK_BUILD_VERSION="${FRAMEWORK_BUILD_VERSION:-latest}"
TCE_BUILD_VERSION="${BUILD_VERSION:-latest}"

DEP_BUILD_DIR="${TCE_RELEASE_DIR}"
ROOT_FRAMEWORK_DIR="${DEP_BUILD_DIR}/tanzu-framework"

rm -rf "${BUILD_ROOT_DIR}"
mkdir -p "${BUILD_ROOT_DIR}"

# Common directories
ROOT_FRAMEWORK_ARTFACTS_DIR="${ROOT_FRAMEWORK_DIR}/artifacts"
ROOT_FRAMEWORK_ARTFACTS_ADMIN_DIR="${ROOT_FRAMEWORK_DIR}/artifacts-admin"
ROOT_TCE_ARTIFACTS_DIR="${ROOT_REPO_DIR}/artifacts"

# change settings
chmod +x "${ROOT_REPO_DIR}/hack/install.sh"

for env in ${ENVS}; do
    binaryname=${env//-/_}
    extension="" && [[ $binaryname = *windows* ]] && extension=".exe"
    PACKAGE_DIR="${BUILD_ROOT_DIR}/tce-${env}-${BUILD_VERSION}"
    rm -rf "${PACKAGE_DIR}"
    mkdir -p "${PACKAGE_DIR}/bin"
    # Tanzu bits
    cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/core/${FRAMEWORK_BUILD_VERSION}/tanzu-core-${binaryname}${extension}" "${PACKAGE_DIR}/bin/tanzu${extension}"
    cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/cluster/${FRAMEWORK_BUILD_VERSION}/tanzu-cluster-${binaryname}${extension}" "${PACKAGE_DIR}/bin/tanzu-plugin-cluster${extension}"
    cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/kubernetes-release/${FRAMEWORK_BUILD_VERSION}/tanzu-kubernetes-release-${binaryname}${extension}" "${PACKAGE_DIR}/bin/tanzu-plugin-kubernetes-release${extension}"
    cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/login/${FRAMEWORK_BUILD_VERSION}/tanzu-login-${binaryname}${extension}" "${PACKAGE_DIR}/bin/tanzu-plugin-login${extension}"
    cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/package/${FRAMEWORK_BUILD_VERSION}/tanzu-package-${binaryname}${extension}" "${PACKAGE_DIR}/bin/tanzu-plugin-package${extension}"
    cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/pinniped-auth/${FRAMEWORK_BUILD_VERSION}/tanzu-pinniped-auth-${binaryname}${extension}" "${PACKAGE_DIR}/bin/tanzu-plugin-pinniped-auth${extension}"
    cp -f "${ROOT_FRAMEWORK_ARTFACTS_DIR}/management-cluster/${FRAMEWORK_BUILD_VERSION}/tanzu-management-cluster-${binaryname}${extension}" "${PACKAGE_DIR}/bin/tanzu-plugin-management-cluster${extension}"
    cp -f "${ROOT_FRAMEWORK_ARTFACTS_ADMIN_DIR}/builder/${FRAMEWORK_BUILD_VERSION}/tanzu-builder-${binaryname}${extension}" "${PACKAGE_DIR}/bin/tanzu-plugin-builder${extension}"

    # TCE bits (New folder structure using tanzu-framework main)
    cp -f "${ROOT_TCE_ARTIFACTS_DIR}/standalone-cluster/${TCE_BUILD_VERSION}/tanzu-standalone-cluster-${binaryname}${extension}" "${PACKAGE_DIR}/bin/tanzu-plugin-standalone-cluster${extension}"
    cp -f "${ROOT_TCE_ARTIFACTS_DIR}/conformance/${TCE_BUILD_VERSION}/tanzu-conformance-${binaryname}${extension}" "${PACKAGE_DIR}/bin/tanzu-plugin-conformance${extension}"
    cp -f "${ROOT_TCE_ARTIFACTS_DIR}/diagnostics/${TCE_BUILD_VERSION}/tanzu-diagnostics-${binaryname}${extension}" "${PACKAGE_DIR}/bin/tanzu-plugin-diagnostics${extension}"

    cp -f "${ROOT_REPO_DIR}/hack/install.sh" "${PACKAGE_DIR}"
    cp -f "${ROOT_REPO_DIR}/hack/uninstall.sh" "${PACKAGE_DIR}"
    chown -R "$(id -u -n)":"$(id -g -n)" "${PACKAGE_DIR}"

    # packaging
    rm -f tce-${env}-*.tar.gz
    pushd "${BUILD_ROOT_DIR}" || exit 1
        if [[ $env = *windows* ]]
        then
            zip -r "tce-${env}-${TCE_BUILD_VERSION}.zip" "tce-${env}-${BUILD_VERSION}"
        else
            tar -czvf "tce-${env}-${TCE_BUILD_VERSION}.tar.gz" "tce-${env}-${BUILD_VERSION}"
        fi
    popd || exit 1
done
