#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

# Change directories to a clean build space
ROOT_REPO_DIR="/tmp/tce-release"
rm -fr "${ROOT_REPO_DIR}"
mkdir -p "${ROOT_REPO_DIR}"
cd "${ROOT_REPO_DIR}" || exit 1

if [[ -z "${BUILD_VERSION}" ]]; then
    echo "BUILD_VERSION is not set"
    exit 1
fi

TANZU_CORE_REPO_BRANCH=${TANZU_CORE_REPO_BRANCH:-$BUILD_VERSION}
if [[ "${TANZU_CORE_REPO_BRANCH}" != "${BUILD_VERSION}" ]]; then
    echo "**************** WARNING - TANZU_CORE_REPO_BRANCH = ${TANZU_CORE_REPO_BRANCH} ****************"
fi

rm -rf "${ROOT_REPO_DIR}/core"
mv -f "${HOME}/.tanzu" "${HOME}/.tanzu-$(date +"%Y-%m-%d_%H:%M")"
set +x
git clone --depth 1 --branch "${TANZU_CORE_REPO_BRANCH}" "https://git:${GH_ACCESS_TOKEN}@github.com/vmware-tanzu-private/core.git" "core"
set -x
pushd "${ROOT_REPO_DIR}/core" || exit 1
git reset --hard
sed -i.bak -e "s/ --dirty//g" ./Makefile && rm ./Makefile.bak
sed -i.bak -e "s/\$(shell git describe --tags --abbrev=0 2>\$(NUL))/${BUILD_VERSION}/g" ./Makefile && rm ./Makefile.bak

# generate the correct tkg-bom (which references the tkr-bom)
make configure-bom
# build and install all "core" CLI plugins
# (e.g. management-cluster, cluster, etc)
make build-install-cli-all
# by default, core only builds admins plugins for the current platform. we need darwin also.
GOHOSTOS=darwin GOHOSTARCH=amd64 make build-plugin-admin
popd || exit 1
