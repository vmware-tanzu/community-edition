#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o nounset
set -o pipefail
set -o xtrace

if [[ -z "${BUILD_VERSION}" ]]; then
    echo "BUILD_VERSION is not set"
    exit 1
fi

# we only allow this to run from GitHub CI/Action
if [[ "${TCE_CI_BUILD}" != "true" ]]; then
    echo "This is only meant to be run within GitHub Actions CI"
    exit 1
fi

TCE_HOMEBREW_TAP_REPO="https://github.com/vmware-tanzu/homebrew-tanzu"

# clone the homebrew repo
rm -rf ./homebrew-tanzu
git clone --depth 1 --branch main "${TCE_HOMEBREW_TAP_REPO}"

pushd "./homebrew-tanzu" || exit 1

    # This will do a check on current version before adding it to main.
    BUILD_VERSION="${BUILD_VERSION}" ./test/check-homebrew-upgrade.sh

popd || exit 1

# clean up
rm -rf ./homebrew-tanzu
