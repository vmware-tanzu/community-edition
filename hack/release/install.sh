#!/bin/env bash

# Copyright 2021-2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace
set +x

debug="false"
if [[ $# -eq 1 ]] && [[ "$1" == "-d" ]]; then
    debug="true"
fi

echo_debug () {
    if [[ "${debug}" == "true" ]]; then
        echo "${1}"
    fi
}

error_exit () {
    echo "ERROR: ${1}"
    exit 1
}

echo "===================================="
echo " Installing Tanzu Community Edition"
echo "===================================="
echo

ALLOW_INSTALL_AS_ROOT="${ALLOW_INSTALL_AS_ROOT:-""}"
if [[ "$EUID" -eq 0 && "${ALLOW_INSTALL_AS_ROOT}" != "true" ]]; then
  error_exit "Do not run this script as root"
fi

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

BUILD_OS=$(uname 2>/dev/null || echo Unknown)

case "${BUILD_OS}" in
  Linux)
    XDG_DATA_HOME="${HOME}/.local/share"
    XDG_CONFIG_HOME="${HOME}/.config"
    ;;
  Darwin)
    XDG_DATA_HOME="${HOME}/Library/Application Support"
    XDG_CONFIG_HOME="${HOME}/.config"
    ;;
  *)
    echo "${BUILD_OS} is unsupported"
    exit 1
    ;;
esac

echo_debug "Data home:   ${XDG_DATA_HOME}"
echo_debug "Config home: ${XDG_CONFIG_HOME}"
echo_debug ""

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
  echo Installing tanzu cli to "${TANZU_BIN_PATH}/tanzu"
  install "${MY_DIR}/tanzu" "${TANZU_BIN_PATH}"
else
  echo Installing tanzu cli to "${TANZU_BIN_PATH}/tanzu"
  sudo install "${MY_DIR}/tanzu" "${TANZU_BIN_PATH}" || handle_sudo_failure
fi
echo

# copy the uninstall script to tanzu-cli directory
mkdir -p "${XDG_DATA_HOME}/tce"
install "${MY_DIR}/uninstall.sh" "${XDG_DATA_HOME}/tce"

# if plugin cache pre-exists, remove it so new plugins are detected
TANZU_PLUGIN_CACHE="${HOME}/.cache/tanzu/catalog.yaml"
if [[ -n "${TANZU_PLUGIN_CACHE}" ]]; then
  echo_debug "Removing old plugin cache from ${TANZU_PLUGIN_CACHE}"
  rm -f "${TANZU_PLUGIN_CACHE}" > /dev/null
fi

# install all plugins present in the bundle
platformdir=$(find "${MY_DIR}" -maxdepth 1 -type d -name "*default*" -exec basename {} \;)

# Workaround!!!
# For TF 0.17.0 or higher
# tanzu plugin install all --local "${MY_DIR}/${platformdir}"
# For 0.11.2
# setup
mkdir -p "${XDG_CONFIG_HOME}/tanzu-plugins"
cp -r "${MY_DIR}/${platformdir}/." "${XDG_CONFIG_HOME}/tanzu-plugins"

# install plugins
tanzu plugin install builder
tanzu plugin install codegen
tanzu plugin install cluster
tanzu plugin install kubernetes-release
tanzu plugin install login
tanzu plugin install management-cluster
tanzu plugin install package
tanzu plugin install pinniped-auth
tanzu plugin install secret
tanzu plugin install conformance
tanzu plugin install diagnostics
tanzu plugin install unmanaged-cluster

# explicit init of tanzu cli and add tce repo
# For TF 0.17.0 or higher
# tanzu init
TCE_REPO="$(tanzu plugin repo list | grep tce)"
if [[ -z "${TCE_REPO}"  ]]; then
  tanzu plugin repo add --name tce --gcp-bucket-name tce-tanzu-cli-plugins --gcp-root-path artifacts
fi
TCE_REPO="$(tanzu plugin repo list | grep core-admin)"
if [[ -z "${TCE_REPO}"  ]]; then
  tanzu plugin repo add --name core-admin --gcp-bucket-name tce-tanzu-cli-framework-admin --gcp-root-path artifacts-admin
fi

echo
echo "Installation complete!"
