#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

# defaults
TANZU_FRAMEWORK_REPO_HASH="${TANZU_FRAMEWORK_REPO_HASH:-""}"
# by default, this value is passed in from the TCE makefile, but leaving it empty also works for building
# all of tanzu-framework's environments
ENVS="${ENVS:-""}"

if [[ -z "${TCE_RELEASE_DIR}" ]]; then
    echo "TCE_RELEASE_DIR is not set"
    exit 1
fi

# Change directories to a clean build space
ROOT_REPO_DIR="${TCE_RELEASE_DIR}"
rm -fr "${ROOT_REPO_DIR}"
mkdir -p "${ROOT_REPO_DIR}"
cd "${ROOT_REPO_DIR}" || exit 1

if [[ -z "${TCE_BUILD_VERSION}" ]]; then
    echo "TCE_BUILD_VERSION is not set"
    exit 1
fi

rm -rf "${ROOT_REPO_DIR}/tanzu-framework"
set +x
if [[ -n "${TANZU_FRAMEWORK_REPO_HASH}" ]]; then
    TANZU_FRAMEWORK_REPO_BRANCH="main"
fi
git clone --depth 1 --branch "${TANZU_FRAMEWORK_REPO_BRANCH}" "https://git:${GITHUB_TOKEN}@github.com/vmware-tanzu/tanzu-framework.git" "tanzu-framework"
set -x
pushd "${ROOT_REPO_DIR}/tanzu-framework" || exit 1
git reset --hard
if [[ -n "${TANZU_FRAMEWORK_REPO_HASH}" ]]; then
    echo "checking out specific hash: ${TANZU_FRAMEWORK_REPO_HASH}"
    git fetch --depth 1 origin "${TANZU_FRAMEWORK_REPO_HASH}"
    git checkout "${TANZU_FRAMEWORK_REPO_HASH}"
fi
BUILD_SHA="$(git describe --match="$(git rev-parse --short HEAD)" --always)"
sed -i.bak -e "s/ --dirty//g" ./Makefile && rm ./Makefile.bak
sed -i.bak -e "s|--artifacts artifacts/\${OS}/\${ARCH}/cli|--artifacts artifacts|g" ./Makefile && rm ./Makefile.bak
sed -i.bak -e "s|--artifacts artifacts-admin/\${GOHOSTOS}/\${GOHOSTARCH}/cli|--artifacts artifacts-admin|g" ./Makefile && rm ./Makefile.bak
sed -i.bak -e "s|--local \$(ARTIFACTS_DIR)/\$(GOHOSTOS)/\$(GOHOSTARCH)/cli|--local \$(ARTIFACTS_DIR)|g" ./Makefile && rm ./Makefile.bak
sed -i.bak -e "s|--local \$(ARTIFACTS_DIR)-admin/\$(GOHOSTOS)/\$(GOHOSTARCH)/cli|--local \$(ARTIFACTS_DIR)-admin|g" ./Makefile && rm ./Makefile.bak
sed -i.bak -e "s/\$(shell git describe --tags --abbrev=0 2>\$(NUL))/${FRAMEWORK_BUILD_VERSION}/g" ./Makefile && rm ./Makefile.bak

# do not delete this... removing this fails to get pinniped to functiona correctly
go mod download
go mod tidy

# allow unstable (non-GA) version plugins
if [[ "${TCE_BUILD_VERSION}" == *"-"* ]]; then
make controller-gen
make set-unstable-versions
fi

# generate the correct tkg-bom (which references the tkr-bom)
# make configure-bom
# build all "tanzu-framework" CLI plugins
# (e.g. management-cluster, cluster, etc)
TANZI_CLI_NO_INIT=true TANZU_CORE_BUCKET="tce-tanzu-cli-framework" TKG_DEFAULT_IMAGE_REPOSITORY=${TKG_DEFAULT_IMAGE_REPOSITORY} BUILD_SHA=${BUILD_SHA} BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} make ENVS="${ENVS}" clean-catalog-cache clean-cli-plugins configure-bom build-cli

# by default, tanzu-framework only builds admins plugins for the current platform. we need darwin also.
BUILD_SHA=${BUILD_SHA} BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} GOHOSTOS=linux GOHOSTARCH=amd64 make ENVS="${ENVS}" build-plugin-admin
BUILD_SHA=${BUILD_SHA} BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} GOHOSTOS=darwin GOHOSTARCH=amd64 make ENVS="${ENVS}" build-plugin-admin
BUILD_SHA=${BUILD_SHA} BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} GOHOSTOS=darwin GOHOSTARCH=arm64 make ENVS="${ENVS}" build-plugin-admin
BUILD_SHA=${BUILD_SHA} BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} GOHOSTOS=windows GOHOSTARCH=amd64 make ENVS="${ENVS}" build-plugin-admin
popd || exit 1
