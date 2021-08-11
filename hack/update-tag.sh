#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

FAKE_RELEASE="${FAKE_RELEASE:-""}"
BUILD_VERSION="${BUILD_VERSION:-""}"

# required input
if [[ -z "${BUILD_VERSION}" ]]; then
    echo "BUILD_VERSION is not set"
    exit 1
fi

# we only allow this to run from GitHub CI/Action
WHOAMI=$(whoami)
if [[ "${WHOAMI}" != "runner" ]]; then
    echo "This is only meant to be run within GitHub Actions CI"
    exit 1
fi

# setup
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

WHICH_BRANCH=$(git branch -a --contains "${WHICH_HASH}" | grep remotes | grep -v -e detached -e HEAD | grep -E "\bmain\b|\brelease-[0-9]+\.[0-9]+\b"  | cut -d "/" -f3)
echo "branch: ${WHICH_BRANCH}"
if [[ "${WHICH_BRANCH}" == "" ]]; then
    echo "Unable to find the branch associated with this hash."
    exit 1
fi

# make sure we are running on a clean state before checking out
git reset --hard
git fetch
git checkout "${WHICH_BRANCH}"
git pull origin "${WHICH_BRANCH}"

# perform the updates to all the necessary tags by running the helper util
# ie DEV_BUILD_VERSION.yaml, FAKE_BUILD_VERSION.yaml, etc
pushd "./hack/tags" || exit 1
if [[ "${BUILD_VERSION}" != *"-"* ]]; then
    go run ./tags.go -tag "${BUILD_VERSION}" -release
else
    go run ./tags.go -tag "${BUILD_VERSION}"
fi
popd || exit 1

NEW_BUILD_VERSION=""
if [[ -f "./hack/NEW_BUILD_VERSION" ]]; then
    NEW_BUILD_VERSION=$(cat ./hack/NEW_BUILD_VERSION)
elif [[ "${NEW_BUILD_VERSION}" == "" ]]; then
    NEW_BUILD_VERSION="${BUILD_VERSION}"
fi

# now that we are ready... perform the commit
# use NEW_BUILD_VERSION to determine VERSION_PROPER this handles the major/minor version changes
VERSION_PROPER=$(echo "${NEW_BUILD_VERSION}" | cut -d "-" -f1)
echo "VERSION_PROPER: ${VERSION_PROPER}"

# is this a fake release to test the process?
if [[ "${FAKE_RELEASE}" != "" ]]; then

DEV_VERSION=$(awk '{print $2}' < ./hack/FAKE_BUILD_VERSION.yaml)
NEW_FAKE_BUILD_VERSION="${VERSION_PROPER}-${DEV_VERSION}"
echo "NEW_FAKE_BUILD_VERSION: ${NEW_FAKE_BUILD_VERSION}"

# commit dev file
git add hack/FAKE_BUILD_VERSION.yaml
git commit -m "auto-generated - update fake version"
git push origin "${WHICH_BRANCH}"

# skip the tagging the dev release... commit the file is a good enough simulation

else

DEV_VERSION=$(awk '{print $2}' < ./hack/DEV_BUILD_VERSION.yaml)
NEW_DEV_BUILD_VERSION="${VERSION_PROPER}-${DEV_VERSION}"
echo "NEW_DEV_BUILD_VERSION: ${NEW_DEV_BUILD_VERSION}"

# commit dev file
git add hack/DEV_BUILD_VERSION.yaml
git commit -m "auto-generated - update dev version"
git push origin "${WHICH_BRANCH}"

# tag the new dev release
git tag -m "${NEW_DEV_BUILD_VERSION}" "${NEW_DEV_BUILD_VERSION}"
git push origin "${NEW_DEV_BUILD_VERSION}"

fi
