#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

ALLOW_INSTALL_AS_ROOT="${ALLOW_INSTALL_AS_ROOT:-""}"
if [[ "$EUID" -eq 0 && "${ALLOW_INSTALL_AS_ROOT}" != "true" ]]; then
  echo "Do not run this script as root"
  exit 1
fi

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

handle_sudo_failure() {
    echo "sudo access required to install to ${TANZU_BIN_PATH}"
    exit 1
}

# check if the tanzu CLI already exists and remove it to avoid conflicts
TANZU_BIN_PATH=$(command -v tanzu)
if [[ -n "${TANZU_BIN_PATH}" ]]; then
  # best effort, so just ignore errors
  rm -f "${TANZU_BIN_PATH}" > /dev/null

  TANZU_BIN_PATH=$(command -v tanzu)
  if [[ -n "${TANZU_BIN_PATH}" ]]; then
    # best effort, so just ignore errors
    echo "Unable to delete Tanzu CLI. Retrying using sudo."
    sudo rm -f "${TANZU_BIN_PATH}" > /dev/null
  fi
fi

# check if ~/bin is in PATH if so use that and don't sudo
# fall back to /usr/local/bin with sudo
TANZU_BIN_PATH="/usr/local/bin"
if [[ ":${PATH}:" == *":$HOME/bin:"* && -d "${HOME}/bin" ]]; then
  TANZU_BIN_PATH="${HOME}/bin"
  echo Installing tanzu cli to "${TANZU_BIN_PATH}"
  install "${MY_DIR}/tanzu" "${TANZU_BIN_PATH}"
else
  echo Installing tanzu cli to "${TANZU_BIN_PATH}"
  sudo install "${MY_DIR}/tanzu" "${TANZU_BIN_PATH}" || handle_sudo_failure
fi

# copy the uninstall script to tanzu-cli directory
mkdir -p "${XDG_DATA_HOME}/tce"
install "${MY_DIR}/uninstall.sh" "${XDG_DATA_HOME}/tce"

# if plugin cache pre-exists, remove it so new plugins are detected
TANZU_PLUGIN_CACHE="${HOME}/.cache/tanzu/catalog.yaml"
if [[ -n "${TANZU_PLUGIN_CACHE}" ]]; then
  echo "Removing old plugin cache from ${TANZU_PLUGIN_CACHE}"
  rm -f "${TANZU_PLUGIN_CACHE}" > /dev/null
fi

# install all plugins present in the bundle
platformdir=$(find "${MY_DIR}" -maxdepth 1 -type d -name "*-default" -exec basename {} \;)
tanzu plugin install all --local "${MY_DIR}/${platformdir}"

# explicit init of tanzu cli and add tce repo
tanzu init
TCE_REPO="$(tanzu plugin repo list | grep tce)"
if [[ -z "${TCE_REPO}"  ]]; then
  tanzu plugin repo add --name tce --gcp-bucket-name tce-tanzu-cli-plugins --gcp-root-path artifacts
fi
TCE_REPO="$(tanzu plugin repo list | grep core-admin)"
if [[ -z "${TCE_REPO}"  ]]; then
  tanzu plugin repo add --name core-admin --gcp-bucket-name tce-tanzu-cli-framework-admin --gcp-root-path artifacts-admin
fi

echo "Installation complete!"
