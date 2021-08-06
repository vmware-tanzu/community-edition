#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

# defaults
TANZU_FRAMEWORK_REPO_HASH="${TANZU_FRAMEWORK_REPO_HASH:-""}"

# Change directories to a clean build space
ROOT_REPO_DIR="/tmp/tce-release"
rm -fr "${ROOT_REPO_DIR}"
mkdir -p "${ROOT_REPO_DIR}"
cd "${ROOT_REPO_DIR}" || exit 1

if [[ -z "${TCE_BUILD_VERSION}" ]]; then
    echo "TCE_BUILD_VERSION is not set"
    exit 1
fi

rm -rf "${ROOT_REPO_DIR}/tanzu-framework"
# TODO remove after this issue has been fixed
# https://github.com/vmware-tanzu/tanzu-framework/issues/144
mv -f "${HOME}/.config/tanzu" "${HOME}/.config/tanzu-$(date +"%Y-%m-%d_%H:%M")"
set +x
if [[ -n "${TANZU_FRAMEWORK_REPO_HASH}" ]]; then
    TANZU_FRAMEWORK_REPO_BRANCH="main"
fi
git clone --depth 1 --branch "${TANZU_FRAMEWORK_REPO_BRANCH}" "https://git:${GH_ACCESS_TOKEN}@github.com/vmware-tanzu/tanzu-framework.git" "tanzu-framework"
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
sed -i.bak -e "s/\$(shell git describe --tags --abbrev=0 2>\$(NUL))/${FRAMEWORK_BUILD_VERSION}/g" ./Makefile && rm ./Makefile.bak

go mod download
go mod tidy

# allow unstable (non-GA) version plugins
if [[ "${TCE_BUILD_VERSION}" == *"-"* ]]; then
make controller-gen
make set-unstable-versions
fi
# generate the correct tkg-bom (which references the tkr-bom)
# make configure-bom
# build and install all "tanzu-framework" CLI plugins
# (e.g. management-cluster, cluster, etc)
# BUILD_SHA=${BUILD_SHA} BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} make build-install-cli-all ENVS=linux-amd64

make clean-catalog-cache
make clean-cli-plugins
BUILD_SHA=${BUILD_SHA} BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} make build-plugin-admin ENVS=linux-amd64

./hack/embed-pinniped-binary.sh go linux amd64
GOPRIVATE="github.com/vmware-tanzu/tanzu-framework" go run ./cmd/cli/plugin-admin/builder/main.go cli compile --version ${FRAMEWORK_BUILD_VERSION} -ldflags -s -w -X 'github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli.BuildDate=2021-08-06' -X 'github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli.BuildSHA=5a42a2b' -X 'github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli.BuildVersion=v1.4.0-pre-alpha-2' -X 'main.BuildEdition=tce' -X 'github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/buildinfo.IsOfficialBuild=""' --tags "" --corepath "cmd/cli/tanzu" --artifacts artifacts/linux/amd64/cli --target  linux_amd64 --match foo

make install-cli

# by default, tanzu-framework only builds admins plugins for the current platform. we need darwin also.
# BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} GOHOSTOS=linux GOHOSTARCH=amd64 make build-plugin-admin
# BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} GOHOSTOS=darwin GOHOSTARCH=amd64 make build-plugin-admin
# BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} GOHOSTOS=windows GOHOSTARCH=amd64 make build-plugin-admin
popd || exit 1
