#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

if [[ -z "${BUILD_VERSION}" ]]; then
    echo "BUILD_VERSION is not set"
    exit 1
fi

WHOAMI=$(whoami)
if [[ "${WHOAMI}" != "runner" ]]; then
    echo "This is only meant to be run within GitHub Actions CI"
    exit 1
fi

# which branch did this tag come from
# get commit hash for this tag, then find which branch the hash is on
WHICH_HASH=$(git rev-parse "tags/${BUILD_VERSION}")
echo "hash: ${WHICH_HASH}"
if [[ "${WHICH_HASH}" == "" ]]; then
    echo "Unable to find the hash associated with this tag."
    exit 1
fi

WHICH_BRANCH=$(git branch -a --contains "${WHICH_HASH}" | grep -v -e detached -e HEAD | grep remotes | cut -d "/" -f3)
echo "branch: ${WHICH_BRANCH}"
if [[ "${WHICH_BRANCH}" == "" ]]; then
    echo "Unable to find the branch associated with this hash."
    exit 1
fi

git reset --hard
git checkout "${WHICH_BRANCH}"
