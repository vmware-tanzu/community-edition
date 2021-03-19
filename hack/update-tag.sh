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

DEV_VERSION=$(awk '{print $2}' < ./hack/DEV_BUILD_VERSION.yaml)
NEW_BUILD_VERSION="${BUILD_VERSION}-${DEV_VERSION}"

git config user.name github-actions
git config user.email github-actions@github.com
git add hack/DEV_BUILD_VERSION.yaml
git commit -m "auto-generated - update dev version"
git push origin HEAD:main --force
git tag -m "${NEW_BUILD_VERSION}" "${NEW_BUILD_VERSION}"
git push --tags
