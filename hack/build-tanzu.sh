#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

# Change directories to the parent directory of the one in which this
# script is located.
ROOT_REPO_DIR="$(dirname "${BASH_SOURCE[0]}")/.."
cd "${ROOT_REPO_DIR}" || return 1

if [[ -z "${CORE_BUILD_VERSION}" ]]; then
    echo "CORE_BUILD_VERSION is not set"
    exit 1
fi

# TODO check script

# TKG_CLI_REPO_BRANCH=v1.3.0-rc.2
TKG_PROVIDERS_REPO_BRANCH=${CORE_BUILD_VERSION}
TANZU_CORE_REPO_BRANCH=${CORE_BUILD_VERSION}
TANZU_TKG_CLI_PLUGINS_REPO_BRANCH=${CORE_BUILD_VERSION}

# rm -rf "${ROOT_REPO_DIR}/tkg-cli"
# git clone --depth 1 --branch ${TKG_CLI_REPO_BRANCH} git@github.com:vmware-tanzu-private/tkg-cli.git
# pushd "${ROOT_REPO_DIR}/tkg-cli" || return 1
# git reset --hard
# popd || return 1

rm -rf "${ROOT_REPO_DIR}/tkg-providers"
git clone --depth 1 --branch "${TKG_PROVIDERS_REPO_BRANCH}" git@github.com:vmware-tanzu-private/tkg-providers.git
pushd "${ROOT_REPO_DIR}/tkg-providers" || return 1
git reset --hard
popd || return 1

rm -rf "${ROOT_REPO_DIR}/core"
mv -f ~/.tanzu ~/.tanzu-old
git clone --depth 1 --branch "${TANZU_CORE_REPO_BRANCH}" git@github.com:vmware-tanzu-private/core.git
pushd "${ROOT_REPO_DIR}/core" || return 1
git reset --hard
# go mod edit --replace github.com/vmware-tanzu-private/tkg-cli=../tkg-cli
go mod edit --replace github.com/vmware-tanzu-private/tkg-providers=../tkg-providers
sed -i "s/ --dirty//g" ./Makefile
sed -i "s/\$(shell git describe --tags --abbrev=0 2>\$(NUL))/${TANZU_CORE_REPO_BRANCH}/g" ./Makefile
make build-install-cli-all
popd || return 1

rm -rf "${ROOT_REPO_DIR}/tanzu-cli-tkg-plugins"
git clone --depth 1 --branch "${TANZU_TKG_CLI_PLUGINS_REPO_BRANCH}" git@github.com:vmware-tanzu-private/tanzu-cli-tkg-plugins.git
pushd "${ROOT_REPO_DIR}/tanzu-cli-tkg-plugins" || return 1
git reset --hard
# go mod edit --replace github.com/vmware-tanzu-private/tkg-cli=../tkg-cli
go mod edit --replace github.com/vmware-tanzu-private/tkg-providers=../tkg-providers
go mod edit --replace github.com/vmware-tanzu-private/core=../core
sed -i "s/ --dirty//g" ./Makefile
sed -i "s/\$(shell git describe --tags --abbrev=0 2>\$(NUL))/${TANZU_TKG_CLI_PLUGINS_REPO_BRANCH}/g" ./Makefile
sed -i "s/tanzu plugin install all --local \$(ARTIFACTS_DIR)/TANZU_CLI_NO_INIT=true \$(GO) run -ldflags \"\$(LD_FLAGS)\" github.com\/vmware-tanzu-private\/core\/cmd\/cli\/tanzu plugin install all --local \$(ARTIFACTS_DIR)/g" ./Makefile
make build
make install-cli-plugins
popd || return 1
