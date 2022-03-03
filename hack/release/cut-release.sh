#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

FAKE_RELEASE="${FAKE_RELEASE:-""}"
BUILD_VERSION="${BUILD_VERSION:-""}"
FRAMEWORK_BUILD_VERSION="${FRAMEWORK_BUILD_VERSION:-""}"
TCE_CI_BUILD="${TCE_CI_BUILD:-""}"
TCE_RELEASE_DIR="${TCE_RELEASE_DIR:-""}"
GITHUB_WORKSPACE="${GITHUB_WORKSPACE:-""}"

# required input
if [[ -z "${BUILD_VERSION}" ]]; then
    echo "BUILD_VERSION is not set"
    exit 1
fi
if [[ -z "${FRAMEWORK_BUILD_VERSION}" ]]; then
    echo "FRAMEWORK_BUILD_VERSION is not set"
    exit 1
fi
if [[ -z "${TCE_RELEASE_DIR}" ]]; then
    echo "TCE_RELEASE_DIR is not set"
    exit 1
fi
if [[ -z "${GITHUB_WORKSPACE}" ]]; then
    echo "Use current working directory..."
    GITHUB_WORKSPACE=$(pwd)
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
WHICH_BRANCH=$(git branch -a --contains "${ACTUAL_COMMIT_SHA}" | grep remotes | grep -v -e detached -e HEAD | grep -E "\brelease-[0-9]+\.[0-9]+\b"  | cut -d "/" -f3)
echo "branch: ${WHICH_BRANCH}"
if [[ "${WHICH_BRANCH}" == "" ]]; then
    # now try main since the release branch doesnt exist
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
# on the draft PR by running the helper util
# ie DEV_BUILD_VERSION.yaml, FAKE_BUILD_VERSION.yaml, etc
pushd "./hack/release/release" || exit 1
PREVIOUS_RELEASE_HASH=$(cat ../../PREVIOUS_RELEASE_HASH)

echo "Generating release notes..."
set +x
GITHUB_TOKEN="${GITHUB_TOKEN}" release-notes \
  --org vmware-tanzu --repo tce --branch "${WHICH_BRANCH}" \
  --start-sha "${PREVIOUS_RELEASE_HASH}" --end-sha "${ACTUAL_COMMIT_SHA}" \
  --required-author "" --go-template go-template:./release.template --output release-notes.txt
set -x

sed -i.bak -e "s/{<VERSION>}/${BUILD_VERSION}/g" ./release-notes.txt && rm ./release-notes.txt.bak

if [[ "${BUILD_VERSION}" != *"-"* ]]; then
    go run ./release.go -tag "${BUILD_VERSION}" -notes ./release-notes.txt -release
else
    go run ./release.go -tag "${BUILD_VERSION}" -notes ./release-notes.txt
fi

rm ./release-notes.txt
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

# login
set +x
echo "${GITHUB_TOKEN}" | gh auth login --with-token
set -x

# is this a fake release to test the process?
if [[ "${FAKE_RELEASE}" != "" ]]; then

DEV_VERSION=$(awk '{print $2}' < ./hack/FAKE_BUILD_VERSION.yaml)
NEW_FAKE_BUILD_VERSION="${VERSION_PROPER}-${DEV_VERSION}"
echo "NEW_FAKE_BUILD_VERSION: ${NEW_FAKE_BUILD_VERSION}"

git stash

DOES_NEW_BRANCH_EXIST=$(git branch -a | grep remotes | grep "${NEW_FAKE_BUILD_VERSION}")
echo "does branch exist: ${DOES_NEW_BRANCH_EXIST}"
if [[ "${DOES_NEW_BRANCH_EXIST}" == "" ]]; then
    git checkout -b "${WHICH_BRANCH}-update-${NEW_FAKE_BUILD_VERSION}" "${WHICH_BRANCH}"
else
    git checkout "${WHICH_BRANCH}-update-${NEW_FAKE_BUILD_VERSION}"
    git rebase -Xtheirs "origin/${WHICH_BRANCH}"
fi

git stash pop

git add hack/FAKE_BUILD_VERSION.yaml
git commit -s -m "auto-generated - update fake version"
git push origin "${WHICH_BRANCH}-update-${NEW_FAKE_BUILD_VERSION}"
gh pr create --title "auto-generated - update fake version" --body "auto-generated - update fake version"
gh pr merge "${WHICH_BRANCH}-update-${NEW_FAKE_BUILD_VERSION}" --delete-branch --squash --admin

# skip the tagging the dev release... commit the file is a good enough simulation

else

DEV_VERSION=$(awk '{print $2}' < ./hack/DEV_BUILD_VERSION.yaml)
NEW_DEV_BUILD_VERSION="${VERSION_PROPER}-${DEV_VERSION}"
echo "NEW_DEV_BUILD_VERSION: ${NEW_DEV_BUILD_VERSION}"

git stash

DOES_NEW_BRANCH_EXIST=$(git branch -a | grep remotes | grep "${NEW_DEV_BUILD_VERSION}")
echo "does branch exist: ${DOES_NEW_BRANCH_EXIST}"
if [[ "${DOES_NEW_BRANCH_EXIST}" == "" ]]; then
    git checkout -b "${WHICH_BRANCH}-update-${NEW_DEV_BUILD_VERSION}" "${WHICH_BRANCH}"
else
    git checkout "${WHICH_BRANCH}-update-${NEW_DEV_BUILD_VERSION}"
    git rebase -Xtheirs "origin/${WHICH_BRANCH}"
fi

git stash pop

git add hack/DEV_BUILD_VERSION.yaml
if [[ "${BUILD_VERSION}" != *"-"* ]]; then
    echo "${ACTUAL_COMMIT_SHA}" | tee ./hack/PREVIOUS_RELEASE_HASH
    git add hack/PREVIOUS_RELEASE_HASH
fi
git commit -s -m "auto-generated - update dev version"
git push origin "${WHICH_BRANCH}-update-${NEW_DEV_BUILD_VERSION}"
gh pr create --title "auto-generated - update dev version" --body "auto-generated - update dev version"
gh pr merge "${WHICH_BRANCH}-update-${NEW_DEV_BUILD_VERSION}" --delete-branch --squash --admin

# tag the new dev release
git tag -m "${NEW_DEV_BUILD_VERSION}" "${NEW_DEV_BUILD_VERSION}"
git push origin "${NEW_DEV_BUILD_VERSION}"

fi

# make no-op folders in order to prevent plugin updates by deleting
# all binaries and preserving the folder structure using an .empty file
# this needs to be done for both tce and tanzu framework

# do this on TCE
pushd "${GITHUB_WORKSPACE}/artifacts" || exit 1

for dir in ./*/
do
  dir=${dir%*/}
  echo "dir: ${dir}"

  # if there isnt a new version of a plugin, delete the plugin folder
  # because the plugin is now deprecated
  if [[ ! -d "${dir}/${BUILD_VERSION}" ]]; then
    echo "skipping ${dir}/${BUILD_VERSION}..."
    rm -rf "./${dir}"
    continue
  fi

  pushd "${dir}/${BUILD_VERSION}" || exit 1
  # delete all binaries
  rm -rf ./test
  rm -f ./*
  # drop a no-op file
  echo "empty" | tee .empty
  popd || exit 1
done

popd || exit 1

# do this on tanzu framework
pushd "${TCE_RELEASE_DIR}/tanzu-framework/artifacts" || exit 1

for dir in ./*/
do
  dir=${dir%*/}
  echo "dir: ${dir}"

  # some tf directories also contain a "latest" subfolder that needs to be deleted
  rm -rf ./latest

  pushd "${dir}/${FRAMEWORK_BUILD_VERSION}" || exit 1
  # delete all binaries
  rm -rf ./test
  rm -f ./*
  # drop a no-op file
  echo "empty" | tee .empty
  popd || exit 1
done

popd || exit 1

pushd "${TCE_RELEASE_DIR}/tanzu-framework/artifacts-admin/" || exit 1

for dir in ./*/
do
  dir=${dir%*/}
  echo "dir: ${dir}"

  # some tf directories also contain a "latest" subfolder that needs to be deleted
  rm -rf ./latest

  pushd "${dir}/${FRAMEWORK_BUILD_VERSION}" || exit 1
  # delete all binaries
  rm -rf ./test
  rm -f ./*
  # drop a no-op file
  echo "empty" | tee .empty
  popd || exit 1
done

popd || exit 1
