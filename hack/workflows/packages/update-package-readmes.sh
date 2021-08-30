#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

# Get the dependencies
make ensure-deps

# Login to GitHub
set +x
echo "${GH_PACKAGING_ACCESS_TOKEN}" | gh auth login --with-token
set -x

BRANCH_NAME=update-package-readmes-$(date +%s)

git stash

# create branch off main
git checkout -b "${BRANCH_NAME}"

git stash pop

git add .

git commit -s -am "auto-generated - update package documentation"

git push origin "${BRANCH_NAME}"

gh pr create --title "auto-generated - update package documentation" --body "auto-generated - update package documentation"

gh pr merge "${BRANCH_NAME}" --admin --squash

echo "Y" | gh auth logout --hostname github.com
