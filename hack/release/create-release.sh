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

# we only allow this to run from GitHub CI/Action
if [[ "${TCE_CI_BUILD}" != "true" ]]; then
    echo "This is only meant to be run within GitHub Actions CI"
    exit 1
fi

# setup
git config user.name github-actions
git config user.email github-actions@github.com

# which branch did this hash/tag come from
# get commit hash for this tag, then find which branch the hash is on
#
# we need to do this in two stages since we could create a tag on main and then
# create a release branch and tag immediately on that release branch. the new tag would appear in
# both main and this new branch because the commit is the same

# first test the release branch because it gets priority
IS_RELEASE_BRANCH=true
WHICH_BRANCH=$(git branch -a --contains "${ACTUAL_COMMIT_SHA}" | grep remotes | grep -v -e detached -e HEAD | grep -E "\brelease-[0-9]+\.[0-9]+\b"  | cut -d "/" -f3)
echo "branch: ${WHICH_BRANCH}"
if [[ "${WHICH_BRANCH}" == "" ]]; then
    # now try main since the release branch doesnt exist
    IS_RELEASE_BRANCH=false
    WHICH_BRANCH=$(git branch -a --contains "${ACTUAL_COMMIT_SHA}" | grep remotes | grep -v -e detached -e HEAD | grep -E "\bmain\b"  | cut -d "/" -f3)
    echo "branch: ${WHICH_BRANCH}"
    if [[ "${WHICH_BRANCH}" == "" ]]; then
        echo "Unable to find the branch associated with this hash."
        exit 1
    fi
fi

# make sure we are running on a clean state before checking out
git reset --hard
git fetch
git checkout "${WHICH_BRANCH}"
git pull origin "${WHICH_BRANCH}"

# perform the updates to all the necessary tags and update the release notes
# on the draft PR by running the helper util (ie BUILD_VERSION, etc)
pushd "./hack/release/release" || exit 1

PREVIOUS_VERSION=$(cat ../../PREVIOUS_VERSION)
if [[ "${IS_RELEASE_BRANCH}" == true ]]; then
    go run ./release.go -previous "${PREVIOUS_VERSION}" -tag "${BUILD_VERSION}" -relbranch
else
    go run ./release.go -previous "${PREVIOUS_VERSION}" -tag "${BUILD_VERSION}"
fi

popd || exit 1

NEW_BUILD_VERSION="${BUILD_VERSION}"
if [[ -f "./hack/NEW_BUILD_VERSION" ]]; then
    NEW_BUILD_VERSION=$(cat ./hack/NEW_BUILD_VERSION)
fi

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
    git checkout -b "automation-${NEW_DEV_BUILD_VERSION}" "${WHICH_BRANCH}"
else
    git checkout "automation-${NEW_DEV_BUILD_VERSION}"
    git rebase -Xtheirs "origin/${WHICH_BRANCH}"
fi

git stash pop

# do the work
git add ./hack/DEV_VERSION
if [[ "${BUILD_VERSION}" != *"-"* ]]; then
    echo "${BUILD_VERSION}" | tee ./hack/PREVIOUS_VERSION
    git add hack/PREVIOUS_VERSION
fi
git commit -s -m "auto-generated - update dev version"
git push origin "automation-${NEW_DEV_BUILD_VERSION}"
gh pr create --title "auto-generated - update dev version for release" --body "auto-generated - update dev version for release" --base "${WHICH_BRANCH}"
gh pr merge "automation-${NEW_DEV_BUILD_VERSION}" --delete-branch --squash --admin

# tag the new dev release
git tag -m "${NEW_DEV_BUILD_VERSION}" "${NEW_DEV_BUILD_VERSION}"
git push origin "${NEW_DEV_BUILD_VERSION}"

# create the signing trigger file
echo "${BUILD_VERSION}" | tee -a ./cayman_trigger.txt
