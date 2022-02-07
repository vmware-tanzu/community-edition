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

# make sure we are running on a clean state before checking out
git reset --hard
git fetch
git checkout main
git pull origin main

# login
set +x
echo "${GITHUB_TOKEN}" | gh auth login --with-token
set -x

# post and create rss.xml
pushd "hack/dailybuild" || exit 1
PREVIOUS_DAILY_HASH=$(cat ./PREVIOUS_DAILY_HASH)

echo "Generating release notes..."
set +x
GITHUB_TOKEN="${GITHUB_TOKEN}" release-notes \
  --org vmware-tanzu --repo tce --branch main \
  --start-sha "${PREVIOUS_DAILY_HASH}" --end-sha "${ACTUAL_COMMIT_SHA}" \
  --required-author "" --go-template go-template:./daily.template --output daily-notes.txt
set -x

echo "${ACTUAL_COMMIT_SHA}" | tee ./PREVIOUS_DAILY_HASH

make run

popd || exit 1

# commit rss.xml to github
git stash

DATE=$(date +%F)
DOES_NEW_BRANCH_EXIST=$(git branch -a | grep remotes | grep "${DATE}")
echo "does branch exist: ${DOES_NEW_BRANCH_EXIST}"
if [[ "${DOES_NEW_BRANCH_EXIST}" == "" ]]; then
    git checkout -b "dailybuild-${DATE}" main
else
    git checkout "dailybuild-${DATE}"
    git rebase -Xtheirs origin/main
fi

git stash pop

git add README.md
git add hack/dailybuild/PREVIOUS_DAILY_HASH
git commit -s -m "auto-generated - update daily build in README"
git push origin "dailybuild-${DATE}"
gh pr create --title "auto-generated - update daily build in README" --body "auto-generated - update daily build in README"
gh pr merge "dailybuild-${DATE}" --delete-branch --squash --admin
