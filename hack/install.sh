#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

if [[ "$EUID" -eq 0 ]]; then
  echo "Do not run this script as root"
  exit
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

# check if the tanzu CLI already exists and remove it to avoid conflicts
TANZU_BIN_PATH=$(command -v tanzu)
if [[ -n "${TANZU_BIN_PATH}" ]]; then
  # best effort, so just ignore errors
  sudo rm -f "${TANZU_BIN_PATH}" > /dev/null
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

# install all plugins present in the bundle
mkdir -p "${XDG_DATA_HOME}/tanzu-cli"
for plugin in "${MY_DIR}"/bin/tanzu-plugin*; do
  install "${plugin}" "${XDG_DATA_HOME}/tanzu-cli"
done

# copy the uninstall script to tanzu-cli directory
mkdir -p "${XDG_DATA_HOME}/tce"
install "${MY_DIR}/uninstall.sh" "${XDG_DATA_HOME}/tce"

# explicit init of tanzu cli and add tce repo
TANZU_CLI_NO_INIT=true tanzu init
TCE_REPO="$(tanzu plugin repo list | grep tce)"
if [[ -z "${TCE_REPO}"  ]]; then
  tanzu plugin repo add --name tce --gcp-bucket-name tce-cli-plugins --gcp-root-path artifacts
fi

echo "Installation complete!"
