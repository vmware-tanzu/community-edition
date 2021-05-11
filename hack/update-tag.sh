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

git config user.name github-actions
git config user.email github-actions@github.com

if [[ "${FAKE_RELEASE}" != "" ]]; then

DEV_VERSION=$(awk '{print $2}' < ./hack/FAKE_BUILD_VERSION.yaml)
NEW_BUILD_VERSION="${BUILD_VERSION}-${DEV_VERSION}"
echo "NEW_BUILD_VERSION: ${NEW_BUILD_VERSION}"

git add hack/FAKE_BUILD_VERSION.yaml
git commit -m "auto-generated - update fake version"
git push origin main
# skip the tagging... commit the file is a good enough simulation

else

DEV_VERSION=$(awk '{print $2}' < ./hack/DEV_BUILD_VERSION.yaml)
NEW_BUILD_VERSION="${BUILD_VERSION}-${DEV_VERSION}"
echo "NEW_BUILD_VERSION: ${NEW_BUILD_VERSION}"

git add hack/DEV_BUILD_VERSION.yaml
git commit -m "auto-generated - update dev version"
git push origin main
git tag -m "${NEW_BUILD_VERSION}" "${NEW_BUILD_VERSION}"
git push origin "${NEW_BUILD_VERSION}"

fi
