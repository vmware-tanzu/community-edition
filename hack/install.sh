#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

BUILD_OS=$(uname 2>/dev/null || echo Unknown)

case "${BUILD_OS}" in
  Linux)
    XDG_DATA_HOME="${HOME}/.local/share"
    ;;
  Darwin)
    XDG_DATA_HOME="${HOME}/Library/Application Support"
    ;;
  *)
    echo "${BUILD_OS} is unsupported"
    exit 1
    ;;
esac

echo "${XDG_DATA_HOME}"
mv -f "${HOME}/.tanzu" "${HOME}/.tanzu-$(date +"%Y-%m-%d_%H:%M")"
rm -rf "${XDG_DATA_HOME}/tanzu-cli"
mkdir -p "${XDG_DATA_HOME}/tanzu-cli"

if [[ "$BUILD_OS" == "Darwin" ]] ;  then
  xattr -d com.apple.quarantine "${MY_DIR}/bin/tanzu"
  xattr -d com.apple.quarantine "${MY_DIR}/bin/tanzu-plugin-alpha"
  xattr -d com.apple.quarantine "${MY_DIR}/bin/tanzu-plugin-cluster"
  xattr -d com.apple.quarantine "${MY_DIR}/bin/tanzu-plugin-management-cluster"
  xattr -d com.apple.quarantine "${MY_DIR}/bin/tanzu-plugin-kubernetes-release"
  xattr -d com.apple.quarantine "${MY_DIR}/bin/tanzu-plugin-login"
  xattr -d com.apple.quarantine "${MY_DIR}/bin/tanzu-plugin-pinniped-auth"
  xattr -d com.apple.quarantine "${MY_DIR}/bin/tanzu-plugin-builder"
  xattr -d com.apple.quarantine "${MY_DIR}/bin/tanzu-plugin-test"
  xattr -d com.apple.quarantine "${MY_DIR}/bin/tanzu-plugin-extension"
fi

# check if the tanzu CLI already exists and remove it to avoid conflicts
TANZU_BIN_PATH=$(which tanzu)
if [[ -n "${TANZU_BIN_PATH}" ]]; then
  # best effort, so just ignore errors
  rm -f "${TANZU_BIN_PATH}" > /dev/null
fi

# check if ~/bin is in PATH if so use that and don't sudo
# fall back to /usr/local/bin with sudo
TANZU_BIN_PATH="/usr/local/bin"
if [[ ":${PATH}:" == *":$HOME/bin:"* && -d "${HOME}/bin" ]]; then
  TANZU_BIN_PATH="${HOME}/bin"
  echo Installing tanzu cli to "${TANZU_BIN_PATH}"
  install "${MY_DIR}/bin/tanzu" "${TANZU_BIN_PATH}"
else
  echo Installing tanzu cli to "${TANZU_BIN_PATH}"
  sudo install "${MY_DIR}/bin/tanzu" "${TANZU_BIN_PATH}"
fi

install "${MY_DIR}/bin/tanzu-plugin-alpha" "${XDG_DATA_HOME}/tanzu-cli"
install "${MY_DIR}/bin/tanzu-plugin-cluster" "${XDG_DATA_HOME}/tanzu-cli"
install "${MY_DIR}/bin/tanzu-plugin-management-cluster" "${XDG_DATA_HOME}/tanzu-cli"
install "${MY_DIR}/bin/tanzu-plugin-kubernetes-release" "${XDG_DATA_HOME}/tanzu-cli"
install "${MY_DIR}/bin/tanzu-plugin-login" "${XDG_DATA_HOME}/tanzu-cli"
install "${MY_DIR}/bin/tanzu-plugin-pinniped-auth" "${XDG_DATA_HOME}/tanzu-cli"
install "${MY_DIR}/bin/tanzu-plugin-builder" "${XDG_DATA_HOME}/tanzu-cli"
install "${MY_DIR}/bin/tanzu-plugin-test" "${XDG_DATA_HOME}/tanzu-cli"
install "${MY_DIR}/bin/tanzu-plugin-package" "${XDG_DATA_HOME}/tanzu-cli"

# repo config
rm -rf "${XDG_DATA_HOME}/tanzu-repository"
mkdir -p "${XDG_DATA_HOME}/tanzu-repository"
mkdir -p "${XDG_DATA_HOME}/tanzu-repository/metadata"
mkdir -p "${XDG_DATA_HOME}/tanzu-repository/extensions"

cp -f "${MY_DIR}/config.yaml" "${XDG_DATA_HOME}/tanzu-repository"
cp -rf "${MY_DIR}/metadata/." "${XDG_DATA_HOME}/tanzu-repository/metadata"
cp -rf "${MY_DIR}/extensions/." "${XDG_DATA_HOME}/tanzu-repository/extensions"

# explicit init of tanzu cli and add tce repo
TANZU_CLI_NO_INIT=true tanzu init
tanzu plugin repo add --name tce --gcp-bucket-name tce-cli-plugins --gcp-root-path artifacts
