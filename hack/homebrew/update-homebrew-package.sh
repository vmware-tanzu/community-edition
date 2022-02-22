#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

version="${1:?TCE version argument empty. Example usage: ./hack/homebrew/update-homebrew-package.sh v0.10.0}"
: "${GITHUB_TOKEN:?GITHUB_TOKEN is not set}"

temp_dir=$(mktemp -d)

pushd "${temp_dir}"

TCE_REPO_RELEASES_URL="https://github.com/vmware-tanzu/community-edition/releases"
TCE_DARWIN_TAR_BALL_FILE="tce-darwin-amd64-${version}.tar.gz"
TCE_LINUX_TAR_BALL_FILE="tce-linux-amd64-${version}.tar.gz"
TCE_CHECKSUMS_FILE="tce-checksums.txt"
TCE_HOMEBREW_TAP_REPO="https://github.com/vmware-tanzu/homebrew-tanzu"

echo "Checking if the necessary files exist for the TCE ${version} release"

curl -f -I -L \
    "${TCE_REPO_RELEASES_URL}/download/${version}/${TCE_DARWIN_TAR_BALL_FILE}" > /dev/null || {
        echo "${TCE_DARWIN_TAR_BALL_FILE} is not accessible in TCE ${version} release"
        exit 1
    }

curl -f -I -L \
    "${TCE_REPO_RELEASES_URL}/download/${version}/${TCE_LINUX_TAR_BALL_FILE}" > /dev/null || {
        echo "${TCE_LINUX_TAR_BALL_FILE} is not accessible in TCE ${version} release"
        exit 1
    }

wget "${TCE_REPO_RELEASES_URL}/download/${version}/${TCE_CHECKSUMS_FILE}" || {
    echo "${TCE_CHECKSUMS_FILE} is not accessible in TCE ${version} release"
    exit 1
}

darwin_amd64_shasum=$(grep "${TCE_DARWIN_TAR_BALL_FILE}" ${TCE_CHECKSUMS_FILE} | cut -d ' ' -f 1)

linux_amd64_shasum=$(grep "${TCE_LINUX_TAR_BALL_FILE}" ${TCE_CHECKSUMS_FILE} | cut -d ' ' -f 1)

git clone ${TCE_HOMEBREW_TAP_REPO}

cd homebrew-tanzu

# make sure we are on main branch before checking out
git checkout main

PR_BRANCH="update-tce-to-${version}-${RANDOM}"

# Random number in branch name in case there's already some branch for the version update,
# though there shouldn't be one. There could be one if the other branch's PR tests failed and didn't merge
git checkout -b "${PR_BRANCH}"

# setup
git config user.name github-actions
git config user.email github-actions@github.com

# Replacing old version with the latest stable released version.
# Using -i so that it works on Mac and Linux OS, so that it's useful for local development.
sed -i.bak "s/version \"v.*/version \"${version}\"/" tanzu-community-edition.rb
rm -fv tanzu-community-edition.rb.bak

# First occurrence of sha256 is for MacOS SHA sum
awk "/sha256 \".*/{c+=1}{if(c==1){sub(\"sha256 \\\".*\",\"sha256 \\\"${darwin_amd64_shasum}\\\"\",\$0)};print}" tanzu-community-edition.rb > tanzu-community-edition-updated.rb
mv tanzu-community-edition-updated.rb tanzu-community-edition.rb

# Second occurrence of sha256 is for Linux SHA sum
awk "/sha256 \".*/{c+=1}{if(c==2){sub(\"sha256 \\\".*\",\"sha256 \\\"${linux_amd64_shasum}\\\"\",\$0)};print}" tanzu-community-edition.rb > tanzu-community-edition-updated.rb
mv tanzu-community-edition-updated.rb tanzu-community-edition.rb

git add tanzu-community-edition.rb

git commit -s -m "auto-generated - update tce homebrew formula for version ${version}"

git push origin "${PR_BRANCH}"

gh pr create --repo ${TCE_HOMEBREW_TAP_REPO} --title "auto-generated - update tce homebrew formula for version ${version}" --body "auto-generated - update tce homebrew formula for version ${version}"

gh pr merge --repo ${TCE_HOMEBREW_TAP_REPO} "${PR_BRANCH}" --squash --delete-branch --auto

popd
