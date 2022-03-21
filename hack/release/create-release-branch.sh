#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

BUILD_VERSION="${BUILD_VERSION:-""}"
TCE_CI_BUILD="${TCE_CI_BUILD:-""}"

# required input
if [[ -z "${BUILD_VERSION}" ]]; then
    echo "BUILD_VERSION is not set"
    exit 1
fi
BUILD_VERSION_WITHOUT_V=${BUILD_VERSION//"v"}

# we only allow this to run from GitHub CI/Action
if [[ "${TCE_CI_BUILD}" != "true" ]]; then
    echo "This is only meant to be run within GitHub Actions CI"
    exit 1
fi

# setup
git config user.name github-actions
git config user.email github-actions@github.com

# make sure we are running on a clean state before checking out
git fetch
git checkout main
git pull origin main
git reset --hard

# let's do the release branch first
# create the release branch before tagging
OLD_VERSION_PROPER=$(echo "${BUILD_VERSION}" | cut -d "-" -f1)
MAJOR_MINOR_VERSION=$(echo "${BUILD_VERSION_WITHOUT_V}" | cut -d "." -f1,2)
git checkout -b "release-${MAJOR_MINOR_VERSION}"
git push origin "release-${MAJOR_MINOR_VERSION}"

# tag the new release (which subsequently creates a new build)
git tag -m "${BUILD_VERSION}" "${BUILD_VERSION}"
git push origin "${BUILD_VERSION}"


# now we need to fix main
# reset again
git fetch
git checkout main
git pull origin main
git reset --hard

# we are going to simulate like we just cut the GA release. this does the following:
# 1. resets the DEV_VERSION to dev.1
# 2. automatically bumps and creates ./hack/NEW_BUILD_VERSION to the new minor version
pushd "./hack/release/release" || exit 1

PREVIOUS_VERSION=$(cat ../../PREVIOUS_VERSION)
go run ./release.go -previous "${PREVIOUS_VERSION}" -tag "${OLD_VERSION_PROPER}" -skip

popd || exit 1

if [[ ! -f "./hack/NEW_BUILD_VERSION" ]]; then
    echo "NEW_BUILD_VERSION does not exist! Unable to bump the minor version!"
    exit 1
fi
NEW_BUILD_VERSION=$(cat ./hack/NEW_BUILD_VERSION)

# use NEW_BUILD_VERSION to determine VERSION_PROPER this handles the major/minor version changes
VERSION_PROPER=$(echo "${NEW_BUILD_VERSION}" | cut -d "-" -f1)
DEV_VERSION=$(cat ./hack/DEV_VERSION)
echo "VERSION_PROPER: ${VERSION_PROPER}"
echo "DEV_VERSION: ${DEV_VERSION}"

# this is the new tag that needs to be create
NEW_DEV_BUILD_VERSION="${VERSION_PROPER}-${DEV_VERSION}"
echo "NEW_DEV_BUILD_VERSION: ${NEW_DEV_BUILD_VERSION}"


# now that we are ready... perform the commit
# login
set +x
echo "${GITHUB_TOKEN}" | gh auth login --with-token
set -x

git stash

# create the branch from main or the release branch
DOES_NEW_BRANCH_EXIST=$(git branch -a | grep remotes | grep "automation-${NEW_DEV_BUILD_VERSION}")
echo "does branch exist: ${DOES_NEW_BRANCH_EXIST}"
if [[ "${DOES_NEW_BRANCH_EXIST}" == "" ]]; then
    git checkout -b "automation-${NEW_DEV_BUILD_VERSION}" main
else
    git checkout "automation-${NEW_DEV_BUILD_VERSION}"
    git rebase -Xtheirs "origin/main"
fi

git stash pop

# do the work
git add ./hack/DEV_VERSION
echo "${OLD_VERSION_PROPER}" | tee ./hack/PREVIOUS_VERSION
git add hack/PREVIOUS_VERSION
git commit -s -m "auto-generated - update dev version after release branch"

# any other file that was changed... let's revert it
git reset --hard

git push origin "automation-${NEW_DEV_BUILD_VERSION}"
gh pr create --title "auto-generated - update dev version after release branch" --body "auto-generated - update dev version after release branch" --base main
gh pr merge "automation-${NEW_DEV_BUILD_VERSION}" --delete-branch --squash --admin

# tag the new dev release
git tag -m "${NEW_DEV_BUILD_VERSION}" "${NEW_DEV_BUILD_VERSION}"
git push origin "${NEW_DEV_BUILD_VERSION}"
