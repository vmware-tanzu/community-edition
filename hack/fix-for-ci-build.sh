#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

# Handle differences in MacOS sed
SEDARGS="-i"
if [[ "$(uname -s)" == "Darwin" ]]; then
    SEDARGS="-i '' -e"
fi

# override https
git config --global url."https://git:${GH_ACCESS_TOKEN}@github.com".insteadOf "https://github.com"

# skip the token because we are using gitlab mirrors
sed "$SEDARGS" "s/https:\/\/git:\${GH_ACCESS_TOKEN}@github.com\/vmware-tanzu-private\/tkg-providers.git/git@gitlab.eng.vmware.com:TKG\/tkg-cli-providers/g" ./hack/build-tanzu.sh
sed "$SEDARGS" "s/https:\/\/git:\${GH_ACCESS_TOKEN}@github.com\/vmware-tanzu-private\/tkg-cli.git/git@gitlab.eng.vmware.com:core-build\/mirrors_github_vmware-tanzu-private_tkg-cli.git/g" ./hack/build-tanzu.sh
sed "$SEDARGS" "s/https:\/\/git:\${GH_ACCESS_TOKEN}@github.com\/vmware-tanzu-private\/core.git/git@gitlab.eng.vmware.com:core-build\/mirrors_github_vmware-tanzu-private_core.git/g" ./hack/build-tanzu.sh
sed "$SEDARGS" "s/https:\/\/git:\${GH_ACCESS_TOKEN}@github.com\/vmware-tanzu-private\/tanzu-cli-tkg-plugins.git/git@gitlab.eng.vmware.com:core-build\/mirrors_github_vmware-tanzu-private_tanzu-cli-tkg-plugins.git/g" ./hack/build-tanzu.sh

# docker container has no user account
sed "$SEDARGS" "s/\"\$(id -g -n \"\$USER\")\"/\$(id -g)/g" ./hack/package-release.sh
sed "$SEDARGS" "s/\"\$USER\"/\$(id -u)/g" ./hack/package-release.sh

# TCE overrides for gitlab
go mod edit --replace github.com/vmware-tanzu-private/tkg-providers=/tmp/tce-release/tkg-providers
go mod edit --replace github.com/vmware-tanzu-private/tkg-cli=/tmp/tce-release/tkg-cli
go mod edit --replace github.com/vmware-tanzu-private/core=/tmp/tce-release/core
