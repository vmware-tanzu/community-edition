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

if [[ -z "${TCE_SCRATCH_DIR}" ]]; then
    echo "TCE_SCRATCH_DIR is not set"
    exit 1
fi

# Change directories to a clean build space
rm -fr "${TCE_SCRATCH_DIR}"
mkdir -p "${TCE_SCRATCH_DIR}"

# recreate the TF repo
pushd "${TCE_SCRATCH_DIR}" || exit 1

if [[ -z "${TCE_BUILD_VERSION}" ]]; then
    echo "TCE_BUILD_VERSION is not set"
    exit 1
fi

rm -rf "${TCE_SCRATCH_DIR}/tanzu-framework"
set +x
if [[ -n "${TANZU_FRAMEWORK_REPO_HASH}" ]]; then
    TANZU_FRAMEWORK_REPO_BRANCH="main"
fi
git clone --depth 1 --branch "${TANZU_FRAMEWORK_REPO_BRANCH}" "${TANZU_FRAMEWORK_REPO}" "tanzu-framework"
set -x

popd || exit 1

# now build TF
pushd "${TCE_SCRATCH_DIR}/tanzu-framework" || exit 1
git reset --hard
if [[ -n "${TANZU_FRAMEWORK_REPO_HASH}" ]]; then
    echo "checking out specific hash: ${TANZU_FRAMEWORK_REPO_HASH}"
    git fetch --depth 1 origin "${TANZU_FRAMEWORK_REPO_HASH}"
    git checkout "${TANZU_FRAMEWORK_REPO_HASH}"
fi

BUILD_SHA="$(git describe --match="$(git rev-parse --short HEAD)" --always)"
sed -i.bak -e "s| --dirty||g" ./Makefile && rm ./Makefile.bak

# do not delete this... removing this fails to get pinniped to function correctly
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
TANZI_CLI_NO_INIT=true TANZU_CORE_BUCKET="tce-tanzu-cli-framework" \
TKG_DEFAULT_IMAGE_REPOSITORY=${TKG_DEFAULT_IMAGE_REPOSITORY} \
TKG_DEFAULT_COMPATIBILITY_IMAGE_PATH=${TKG_DEFAULT_COMPATIBILITY_IMAGE_PATH} \
BUILD_SHA=${BUILD_SHA} BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} make ENVS="${ENVS}" clean-catalog-cache clean-cli-plugins configure-bom

for platform in ${ENVS}
do
    OS=$(cut -d"-" -f1 <<< "${platform}")
    ARCH=$(cut -d"-" -f2 <<< "${platform}")
    scriptextension="" && [[ $platform = *windows* ]] && scriptextension=".exe"

    # build everything
    BUILD_SHA=${BUILD_SHA} BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} make "build-cli-local-${platform}"
    BUILD_SHA=${BUILD_SHA} BUILD_VERSION=${FRAMEWORK_BUILD_VERSION} make "build-plugin-admin-local-${platform}"

    # use new discovery structure
    # Workaround!!!
    # For TF 0.17.0 or higher
    # make publish-admin-plugins-all-local ENVS="${platform}" ADMIN_PLUGINS="builder codegen" TANZU_PLUGIN_PUBLISH_PATH="${TCE_SCRATCH_DIR}/tanzu-framework/build/${platform}-default"
    # make "publish-plugins-local-${platform}" PLUGINS="cluster feature kubernetes-release login management-cluster package pinniped-auth secret" TANZU_PLUGIN_PUBLISH_PATH="${TCE_SCRATCH_DIR}/tanzu-framework/build" DISCOVERY_NAME="default"
    # For 0.11.1
    mkdir -p "./build/${OS}-${ARCH}-${DISCOVERY_NAME}"

    pushd "./build/${OS}-${ARCH}-${DISCOVERY_NAME}" || exit 1
        go run github.com/vmware-tanzu/tanzu-framework/cmd/cli/plugin-admin/builder publish --type local \
        --plugins "builder codegen" \
        --version "${FRAMEWORK_BUILD_VERSION}" --os-arch "${OS}-${ARCH}" \
        --local-output-discovery-dir "discovery/admin" \
        --local-output-distribution-dir "." \
        --input-artifact-dir ../../artifacts-admin
        go run github.com/vmware-tanzu/tanzu-framework/cmd/cli/plugin-admin/builder publish --type local \
        --plugins "cluster kubernetes-release login management-cluster package pinniped-auth secret" \
        --version "${FRAMEWORK_BUILD_VERSION}" --os-arch "${OS}-${ARCH}" \
        --local-output-discovery-dir "discovery/${DISCOVERY_NAME}" \
        --local-output-distribution-dir "." \
        --input-artifact-dir ../../artifacts
        # admin plugins will not install from admin folder
        mv -f ./discovery/admin/* "./discovery/${DISCOVERY_NAME}/"
        rm -rf ./discovery/admin
        # end
        mkdir -p ./distribution
        mv "${OS}" ./distribution
    popd || exit 1

    # copy the tanzu cli
    cp -f "./artifacts/${OS}/${ARCH}/cli/core/${FRAMEWORK_BUILD_VERSION}/tanzu-core-${OS}_${ARCH}${scriptextension}" "${TCE_SCRATCH_DIR}/tanzu-framework/build/${platform}-${DISCOVERY_NAME}"
done

popd || exit 1
