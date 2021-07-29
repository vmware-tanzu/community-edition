#!/usr/bin/env bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script can be used to check your local development environment
# for necessary dependencies used in TCE

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

BUILD_OS=$(uname 2>/dev/null || echo Unknown)

case "${BUILD_OS}" in
  Linux)
    ;;
  Darwin)
    ;;
  *)
    echo "${BUILD_OS} is unsupported"
    exit 1
    ;;
esac

i=0

# these are must have dependencies to just get going
if [[ -z "$(command -v go)" ]]; then
    echo "Missing go"
    ((i=i+1))
fi

if [[ -z "$(command -v docker)" ]]; then
    echo "Missing docker"
    ((i=i+1))
fi

if [[ -z "$(command -v curl)" ]]; then
    echo "Missing curl"
    ((i=i+1))
fi
# these are must have dependencies to just get going

if [[ $i -gt 0 ]]; then
    echo "Total missing: $i"
    echo "Please install these minimal dependencies in order to continue"
    exit 1
fi

CMD="imgpkg"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    curl -LO https://github.com/vmware-tanzu/carvel-imgpkg/releases/download/v0.12.0/imgpkg-linux-amd64
    chmod +x ${CMD}-linux-amd64
    mv ${CMD}-linux-amd64 ${CMD}
    sudo install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
  Darwin)
    curl -LO https://github.com/vmware-tanzu/carvel-imgpkg/releases/download/v0.12.0/imgpkg-darwin-amd64
    chmod +x ${CMD}-darwin-amd64
    mv ${CMD}-darwin-amd64 ${CMD}
    sudo install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
esac
fi

CMD="kbld"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    curl -LO https://github.com/vmware-tanzu/carvel-kbld/releases/download/v0.30.0/kbld-linux-amd64
    chmod +x ${CMD}-linux-amd64
    mv ${CMD}-linux-amd64 ${CMD}
    sudo install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
  Darwin)
    curl -LO https://github.com/vmware-tanzu/carvel-kbld/releases/download/v0.30.0/kbld-darwin-amd64
    chmod +x ${CMD}-darwin-amd64
    mv ${CMD}-darwin-amd64 ${CMD}
    sudo install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
esac
fi

CMD="ytt"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    curl -LO https://github.com/vmware-tanzu/carvel-ytt/releases/download/v0.34.0/ytt-linux-amd64
    chmod +x ${CMD}-linux-amd64
    mv ${CMD}-linux-amd64 ${CMD}
    sudo install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
  Darwin)
    curl -LO https://github.com/vmware-tanzu/carvel-ytt/releases/download/v0.34.0/ytt-darwin-amd64
    chmod +x ${CMD}-darwin-amd64
    mv ${CMD}-darwin-amd64 ${CMD}
    sudo install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
esac
fi

CMD="shellcheck"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    curl -LO https://github.com/koalaman/shellcheck/releases/download/v0.7.2/shellcheck-v0.7.2.linux.x86_64.tar.xz
    tar -xvf shellcheck-v0.7.2.linux.x86_64.tar.xz
    pushd "./shellcheck-v0.7.2" || exit 1
    chmod +x ${CMD}
    sudo install ./${CMD} /usr/local/bin
    popd || exit 1
    rm -rf ./shellcheck-v0.7.2
    rm shellcheck-v0.7.2.linux.x86_64.tar.xz
    ;;
  Darwin)
    curl -LO https://github.com/koalaman/shellcheck/releases/download/v0.7.2/shellcheck-v0.7.2.darwin.x86_64.tar.xz
    tar -xvf shellcheck-v0.7.2.darwin.x86_64.tar.xz
    pushd "./shellcheck-v0.7.2" || exit 1
    chmod +x ${CMD}
    sudo install ./${CMD} /usr/local/bin
    popd || exit 1
    rm -rf ./shellcheck-v0.7.2
    rm shellcheck-v0.7.2.linux.x86_64.tar.xz
    ;;
esac
fi

CMD="gh"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    curl -LO https://github.com/cli/cli/releases/download/v1.13.1/gh_1.13.1_linux_amd64.tar.gz
    tar -xvf gh_1.13.1_linux_amd64.tar.gz
    pushd "./gh_1.13.1_linux_amd64/bin" || exit 1
    chmod +x ${CMD}
    sudo install ./${CMD} /usr/local/bin
    popd || exit 1
    rm -rf ./gh_1.13.1_linux_amd64
    rm gh_1.13.1_linux_amd64.tar.gz
    ;;
  Darwin)
    curl -LO https://github.com/cli/cli/releases/download/v1.13.1/gh_1.13.1_macOS_amd64.tar.gz
    tar -xvf gh_1.13.1_macOS_amd64.tar.gz
    pushd "./gh_1.13.1_macOS_amd64/bin" || exit 1
    chmod +x ${CMD}
    sudo install ./${CMD} /usr/local/bin
    popd || exit 1
    rm -rf ./gh_1.13.1_macOS_amd64
    rm gh_1.13.1_macOS_amd64.tar.gz
    ;;
esac
fi

CMD="kubectl"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    chmod +x ${CMD}
    sudo install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
  Darwin)
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/darwin/amd64/kubectl"
    chmod +x ${CMD}
    sudo install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
esac
fi

echo "No missing dependencies!"
