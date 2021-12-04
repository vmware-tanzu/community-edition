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
BUILD_ARCH=$(uname -m 2>/dev/null || echo Unknown)
GOBIN=$(go env GOBIN)
GOPATH=$(go env GOPATH)
GOBINDIR="${GOBIN:-$GOPATH/bin}"
TCE_REPO_PATH="$(git rev-parse --show-toplevel)"

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
    go install github.com/k14s/imgpkg/cmd/imgpkg@v0.12.0
    sudo install "${GOBINDIR}"/${CMD} /usr/local/bin
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
    go install github.com/k14s/kbld/cmd/kbld@v0.30.0
    sudo install "${GOBINDIR}"/${CMD} /usr/local/bin
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
    go install github.com/k14s/ytt/cmd/ytt@v0.34.0
    sudo install "${GOBINDIR}"/${CMD} /usr/local/bin
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
    case "${BUILD_ARCH}" in
      x86_64)
        curl -LO https://github.com/koalaman/shellcheck/releases/download/v0.7.2/shellcheck-v0.7.2.darwin.x86_64.tar.xz
        tar -xvf shellcheck-v0.7.2.darwin.x86_64.tar.xz
        pushd "./shellcheck-v0.7.2" || exit 1
        chmod +x ${CMD}
        sudo install ./${CMD} /usr/local/bin
        popd || exit 1
        rm -rf ./shellcheck-v0.7.2
        rm shellcheck-v0.7.2.linux.x86_64.tar.xz
        ;;
      arm64)
        brew install shellcheck
        ;;
    esac
    ;;
esac
fi

"${TCE_REPO_PATH}/hack/ensure-deps/ensure-gh-cli.sh"

CMD="release-notes"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    curl -LO https://storage.googleapis.com/tce-tanzu-cli-plugins/build-tools/release-notes/v0.10.0-tce.3/release-notes-linux-amd64
    mv release-notes-linux-amd64 ${CMD}
    chmod +x ${CMD}
    sudo install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
  Darwin)
    case "${BUILD_ARCH}" in
      x86_64)
        curl -LO https://storage.googleapis.com/tce-tanzu-cli-plugins/build-tools/release-notes/v0.10.0-tce.3/release-notes-darwin-amd64
        mv release-notes-darwin-amd64 ${CMD}
        chmod +x ${CMD}
        sudo install ./${CMD} /usr/local/bin
        rm ./${CMD}
        ;;
     arm64)
        curl -LO https://storage.googleapis.com/tce-tanzu-cli-plugins/build-tools/release-notes/v0.10.0-tce.3/release-notes-darwin-arm64
        mv release-notes-darwin-arm64 ${CMD}
        chmod +x ${CMD}
        sudo install ./${CMD} /usr/local/bin
        rm ./${CMD}
        ;;
    esac
    ;;
esac
fi

CMD="aws"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
    unzip awscliv2.zip
    sudo ./aws/install -i
    rm -rf ./aws
    rm awscliv2.zip
    ;;
  Darwin)
    curl "https://awscli.amazonaws.com/AWSCLIV2.pkg" -o "awscliv2.pkg"
    sudo installer -pkg awscliv2.pkg -target /
    rm -rf ./aws
    rm awscliv2.pkg
    ;;
esac
fi

CMD="jq"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    curl -L "https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64" -o "jq"
    chmod +x jq
    sudo mv jq /usr/local/bin
    ;;
  Darwin)
    curl -L "https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64" -o "jq"
    chmod +x jq
    sudo mv jq /usr/local/bin
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
    case "${BUILD_ARCH}" in
      x86_64)
        curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/darwin/amd64/kubectl"
        chmod +x ${CMD}
        sudo install ./${CMD} /usr/local/bin
        rm ./${CMD}
        ;;
      arm64)
        curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/darwin/arm64/kubectl"
        chmod +x ${CMD}
        sudo install ./${CMD} /usr/local/bin
        rm ./${CMD}
        ;;
    esac
    ;;
esac
fi

CMD="kind"
if [[ -z "$(command -v ${CMD})" ]]; then
echo "Attempting install of ${CMD}..."
case "${BUILD_OS}" in
  Linux)
    curl -L https://github.com/kubernetes-sigs/kind/releases/download/v0.11.1/kind-linux-amd64 -o kind
    chmod +x ${CMD}
    sudo install ./${CMD} /usr/local/bin
    rm ./${CMD}
    ;;
  Darwin)
    case "${BUILD_ARCH}" in
      x86_64)
        curl -L https://github.com/kubernetes-sigs/kind/releases/download/v0.11.1/kind-darwin-amd64 -o kind
        chmod +x ${CMD}
        sudo install ./${CMD} /usr/local/bin
        rm ./${CMD}
        ;;
      arm64)
        curl -L https://github.com/kubernetes-sigs/kind/releases/download/v0.11.1/kind-darwin-arm64 -o kind
        chmod +x ${CMD}
        sudo install ./${CMD} /usr/local/bin
        rm ./${CMD}
        ;;
    esac
    ;;
esac
fi

echo "No missing dependencies!"
