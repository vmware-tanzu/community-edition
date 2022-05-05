#!/bin/bash

# Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o nounset
set -o pipefail
set -o xtrace

version="${1:?TCE version argument empty. Example usage: ./hack/choco/update-choco-package.sh v0.10.0}"
: "${GITHUB_TOKEN:?GITHUB_TOKEN is not set}"

# we only allow this to run from GitHub CI/Action
if [[ "${TCE_CI_BUILD}" != "true" ]]; then
    echo "This is only meant to be run within GitHub Actions CI"
    exit 1
fi

TCE_REPO_RELEASES_URL="https://github.com/vmware-tanzu/community-edition/releases"
TCE_WINDOWS_ZIP_FILE="tce-windows-amd64-${version}.zip"
TCE_CHECKSUMS_FILE="tce-checksums.txt"

echo "Checking if the necessary files exist for the TCE ${version} release"
wget --spider -q \
   "${TCE_REPO_RELEASES_URL}/download/${version}/${TCE_WINDOWS_ZIP_FILE}" || {
       echo "${TCE_WINDOWS_ZIP_FILE} is not accessible in TCE ${version} release"
       exit 1
   }

wget "${TCE_REPO_RELEASES_URL}/download/${version}/${TCE_CHECKSUMS_FILE}" || {
   echo "${TCE_CHECKSUMS_FILE} is not accessible in TCE ${version} release"
   exit 1
}

# download checksum file
windows_amd64_shasum=$(grep "${TCE_WINDOWS_ZIP_FILE}" "${TCE_CHECKSUMS_FILE}" | cut -d ' ' -f1)
rm -f "${TCE_CHECKSUMS_FILE}"

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


# Replacing old version with the latest stable released version.
pushd "./hack/choco" || exit 1

    sed -i.bak -E "s/\\\$releaseVersion =.*/\$releaseVersion = \'${version}\'/g" tools/chocolateyinstall.ps1 && rm tools/chocolateyinstall.ps1.bak

    version="${version:1}"
    sed -i.bak -E "s/<version>.*<\/version>/<version>${version}<\/version>/g" tanzu-community-edition.nuspec && rm tanzu-community-edition.nuspec.bak

    sed -i.bak -E "s/\\\$checksum64 =.*/\$checksum64 = \'${windows_amd64_shasum}\'/g" tools/chocolateyinstall.ps1 && rm tools/chocolateyinstall.ps1.bak

popd || exit 1


PR_BRANCH="automation-choco-${version}"

# now that we are ready... perform the commit
# login
set +x
echo "${GITHUB_TOKEN}" | gh auth login --with-token
set -x

git stash

# create the branch from main or the release branch
DOES_NEW_BRANCH_EXIST=$(git branch -a | grep remotes | grep "${PR_BRANCH}")
echo "does branch exist: ${DOES_NEW_BRANCH_EXIST}"
if [[ "${DOES_NEW_BRANCH_EXIST}" == "" ]]; then
    git checkout -b "${PR_BRANCH}" "${WHICH_BRANCH}"
else
    git checkout "${PR_BRANCH}"
    git rebase -Xtheirs "origin/${WHICH_BRANCH}"
fi

git stash pop

# do the work
git add hack/choco/tools/chocolateyinstall.ps1
git add hack/choco/tanzu-community-edition.nuspec
git commit -s -m "auto-generated - update tce choco install scripts for version ${version}"
git push origin "${PR_BRANCH}"
gh pr create --title "auto-generated - update tce choco install scripts for version ${version}" --body "auto-generated - update tce choco install scripts for version ${version}" --base "${WHICH_BRANCH}"
gh pr merge "${PR_BRANCH}" --delete-branch --squash --admin
