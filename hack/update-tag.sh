#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

FAKE_RELEASE="${FAKE_RELEASE:-""}"

if [[ -z "${BUILD_VERSION}" ]]; then
    echo "BUILD_VERSION is not set"
    exit 1
fi

WHOAMI=$(whoami)
if [[ "${WHOAMI}" != "runner" ]]; then
    echo "This is only meant to be run within GitHub Actions CI"
    exit 1
fi

CONFIG_VERSION=$(echo "${BUILD_VERSION}" | cut -d "-" -f1)
echo "CONFIG_VERSION: ${CONFIG_VERSION}"

git config user.name github-actions
git config user.email github-actions@github.com

# which branch did this tag come from
# get commit hash for this tag, then find which branch the hash is on
WHICH_HASH=$(git rev-parse "tags/${BUILD_VERSION}")
echo "hash: ${WHICH_HASH}"
if [[ "${WHICH_HASH}" == "" ]]; then
    echo "Unable to find the hash associated with this tag."
    exit 1
fi

WHICH_BRANCH=$(git branch -a --contains "${WHICH_HASH}" | grep remotes | grep -v -e detached -e HEAD | grep -E ""\bmain\b|\brelease-[0-9]+\.[0-9]+\b"  | cut -d "/" -f3)
echo "branch: ${WHICH_BRANCH}"
if [[ "${WHICH_BRANCH}" == "" ]]; then
    echo "Unable to find the branch associated with this hash."
    exit 1
fi

# handle the case when a PR is merged before the commit/tag can complete
git stash
git fetch
git pull origin "${WHICH_BRANCH}"
git stash pop

if [[ "${FAKE_RELEASE}" != "" ]]; then

DEV_VERSION=$(awk '{print $2}' < ./hack/FAKE_BUILD_VERSION.yaml)
NEW_BUILD_VERSION="${CONFIG_VERSION}-${DEV_VERSION}"
echo "NEW_BUILD_VERSION: ${NEW_BUILD_VERSION}"

git add hack/FAKE_BUILD_VERSION.yaml
git commit -m "auto-generated - update fake version"
git push origin "${WHICH_BRANCH}"
# skip the tagging... commit the file is a good enough simulation

else

DEV_VERSION=$(awk '{print $2}' < ./hack/DEV_BUILD_VERSION.yaml)
NEW_BUILD_VERSION="${CONFIG_VERSION}-${DEV_VERSION}"
echo "NEW_BUILD_VERSION: ${NEW_BUILD_VERSION}"

git add hack/DEV_BUILD_VERSION.yaml
git commit -m "auto-generated - update dev version"
git push origin "${WHICH_BRANCH}"
git tag -m "${NEW_BUILD_VERSION}" "${NEW_BUILD_VERSION}"
git push origin "${NEW_BUILD_VERSION}"

fi
