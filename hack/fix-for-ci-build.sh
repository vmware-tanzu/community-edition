#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

# make user CI directory
if [[ "${GITLAB_CI_BUILD}" == "true" ]]; then
apt update
apt install zip unzip

rm -rf /tmp/tce-release
rm -rf /tmp/mylocal/tanzu-cli
mkdir -p /tmp/tce-release
mkdir -p /tmp/mylocal/tanzu-cli
fi

# override https
git config --global url."https://git:${GH_ACCESS_TOKEN}@github.com".insteadOf "https://github.com"

# docker container has no user account
sed -i.bak -e "s/\"\$(id -g -n \"\$USER\")\"/\$(id -g)/g" ./hack/package-release.sh && rm ./hack/package-release.sh.bak
sed -i.bak -e "s/\"\$USER\"/\$(id -u)/g" ./hack/package-release.sh && rm ./hack/package-release.sh.bak
