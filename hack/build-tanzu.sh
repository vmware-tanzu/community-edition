#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

# Handle differences in MacOS sed
SEDARGS="-i"
if [[ "$(uname -s)" == "Darwin" ]]; then
    SEDARGS="-i '' -e"
fi

# Change directories to a clean build space
ROOT_REPO_DIR="/tmp/tce-release"
rm -fr "${ROOT_REPO_DIR}"
mkdir -p "${ROOT_REPO_DIR}"
cd "${ROOT_REPO_DIR}" || exit 1

if [[ -z "${BUILD_VERSION}" ]]; then
    echo "BUILD_VERSION is not set"
    exit 1
fi

TKG_PROVIDERS_REPO_BRANCH=${TKG_PROVIDERS_REPO_BRANCH:-$BUILD_VERSION}
if [[ "${TKG_PROVIDERS_REPO_BRANCH}" != "${BUILD_VERSION}" ]]; then
    echo "**************** WARNING - TKG_PROVIDERS_REPO_BRANCH = ${TKG_PROVIDERS_REPO_BRANCH} ****************"
fi
TANZU_CORE_REPO_BRANCH=${TANZU_CORE_REPO_BRANCH:-$BUILD_VERSION}
if [[ "${TANZU_CORE_REPO_BRANCH}" != "${BUILD_VERSION}" ]]; then
    echo "**************** WARNING - TANZU_CORE_REPO_BRANCH = ${TANZU_CORE_REPO_BRANCH} ****************"
fi

rm -rf "${ROOT_REPO_DIR}/tkg-providers"
set +x
git clone --depth 1 --branch "${TKG_PROVIDERS_REPO_BRANCH}" "https://git:${GH_ACCESS_TOKEN}@github.com/vmware-tanzu-private/tkg-providers.git"
set -x
pushd "${ROOT_REPO_DIR}/tkg-providers" || exit 1
git reset --hard
popd || exit 1

rm -rf "${ROOT_REPO_DIR}/core"
mv -f "${HOME}/.tanzu" "${HOME}/.tanzu-$(date +"%Y-%m-%d_%H:%M")"
set +x
git clone --depth 1 --branch "${TANZU_CORE_REPO_BRANCH}" "https://git:${GH_ACCESS_TOKEN}@github.com/vmware-tanzu-private/core.git"
set -x
pushd "${ROOT_REPO_DIR}/core" || exit 1
git reset --hard
# go mod edit --replace github.com/vmware-tanzu-private/tkg-cli=../tkg-cli
go mod edit --replace github.com/vmware-tanzu-private/tkg-providers=../tkg-providers
sed "$SEDARGS" "s/ --dirty//g" ./Makefile
sed "$SEDARGS" "s/\$(shell git describe --tags --abbrev=0 2>\$(NUL))/${TANZU_CORE_REPO_BRANCH}/g" ./Makefile
sed "$SEDARGS" "s/^ENVS.*/ENVS := linux-amd64 darwin-amd64/g" ./Makefile
make build-install-cli-all
popd || exit 1
