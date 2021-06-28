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

if [[ -z "$(command -v go)" ]]; then
    echo "Missing go"
    ((i=i+1))
fi

if [[ -z "$(command -v docker)" ]]; then
    echo "Missing docker"
    ((i=i+1))
fi

if [[ -z "$(command -v curl)" ]]; then
    echo "Missing curl. Trying wget..."
if [[ -z "$(command -v wget)" ]]; then
    echo "Missing curl and wget"
    ((i=i+1))
fi
fi

if [[ $i -gt 0 ]]; then
    echo "Total missing: $i"
    echo "Please install dependencies to continue"
    exit 1
fi

DOWNLOAD_COMMAND="curl -LO"
if [[ -z "$(command -v curl)" ]] && [[ ! -z "$(command -v wget)" ]]; then
    DOWNLOAD_COMMAND="wget"
fi

CMD="imgpkg"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    $(${DOWNLOAD_COMMAND} https://github.com/vmware-tanzu/carvel-imgpkg/releases/download/v0.12.0/imgpkg-linux-amd64)
    chmod +x ${CMD}-linux-amd64
    mv ${CMD}-linux-amd64 ${CMD}
    sudo install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
  Darwin)
    $(${DOWNLOAD_COMMAND} https://github.com/vmware-tanzu/carvel-imgpkg/releases/download/v0.12.0/imgpkg-darwin-amd64)
    chmod +x ${CMD}-darwin-amd64
    mv ${CMD}-darwin-amd64 ${CMD}
    sudo install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
  *)
    echo "${BUILD_OS} is unsupported"
    exit 1
    ;;
esac
fi

CMD="kbld"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    $(${DOWNLOAD_COMMAND} https://github.com/vmware-tanzu/carvel-kbld/releases/download/v0.30.0/kbld-linux-amd64)
    chmod +x ${CMD}-linux-amd64
    mv ${CMD}-linux-amd64 ${CMD}
    sudo install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
  Darwin)
    $(${DOWNLOAD_COMMAND} https://github.com/vmware-tanzu/carvel-kbld/releases/download/v0.30.0/kbld-darwin-amd64)
    chmod +x ${CMD}-darwin-amd64
    mv ${CMD}-darwin-amd64 ${CMD}
    sudo install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
  *)
    echo "${BUILD_OS} is unsupported"
    exit 1
    ;;
esac
fi

CMD="ytt"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    $(${DOWNLOAD_COMMAND} https://github.com/vmware-tanzu/carvel-ytt/releases/download/v0.34.0/ytt-linux-amd64)
    chmod +x ${CMD}-linux-amd64
    mv ${CMD}-linux-amd64 ${CMD}
    sudo install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
  Darwin)
    $(${DOWNLOAD_COMMAND} https://github.com/vmware-tanzu/carvel-ytt/releases/download/v0.34.0/ytt-darwin-amd64)
    chmod +x ${CMD}-darwin-amd64
    mv ${CMD}-darwin-amd64 ${CMD}
    sudo install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
  *)
    echo "${BUILD_OS} is unsupported"
    exit 1
    ;;
esac
fi

CMD="shellcheck"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    $(${DOWNLOAD_COMMAND} https://github.com/koalaman/shellcheck/releases/download/v0.7.2/shellcheck-v0.7.2.linux.x86_64.tar.xz)
    tar -xvf shellcheck-v0.7.2.linux.x86_64.tar.xz
    pushd "./shellcheck-v0.7.2" || exit 1
    chmod +x ${CMD}
    sudo install ./${CMD} /usr/local/bin
    popd || exit 1
    rm -rf ./shellcheck-v0.7.2
    rm shellcheck-v0.7.2.linux.x86_64.tar.xz
    ;;
  Darwin)
    $(${DOWNLOAD_COMMAND} https://github.com/koalaman/shellcheck/releases/download/v0.7.2/shellcheck-v0.7.2.darwin.x86_64.tar.xz)
    tar -xvf shellcheck-v0.7.2.darwin.x86_64.tar.xz
    pushd "./shellcheck-v0.7.2" || exit 1
    chmod +x ${CMD}
    sudo install ./${CMD} /usr/local/bin
    popd || exit 1
    rm -rf ./shellcheck-v0.7.2
    rm shellcheck-v0.7.2.linux.x86_64.tar.xz
    ;;
  *)
    echo "${BUILD_OS} is unsupported"
    exit 1
    ;;
esac
fi

CMD="kubectl"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    chmod +x ${CMD}-linux-amd64
    mv ${CMD}-linux-amd64 ${CMD}
    sudo install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
  Darwin)
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/darwin/amd64/kubectl"
    chmod +x ${CMD}-darwin-amd64
    mv ${CMD}-darwin-amd64 ${CMD}
    sudo install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
  *)
    echo "${BUILD_OS} is unsupported"
    exit 1
    ;;
esac
fi

echo "No missing dependencies!"
