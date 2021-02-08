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

echo ${XDG_DATA_HOME}
rm -rf "${XDG_DATA_HOME}/tanzu-cli"
mkdir -p "${XDG_DATA_HOME}/tanzu-cli"

# check if ~/bin is in PATH if so use that and don't sudo
# fall back to /usr/local/bin with sudo
case ":${PATH}:" in
  *:${HOME}/bin:*)
    echo Installing tanzu cli to ${HOME}/bin
    install ${MY_DIR}/bin/tanzu ${HOME}/bin
    ;;
  *)
    echo Installing tanzu cli to /usr/local/bin
    sudo install ${MY_DIR}/bin/tanzu /usr/local/bin
    ;;
esac

install ${MY_DIR}/bin/tanzu-plugin-alpha  "${XDG_DATA_HOME}/tanzu-cli"
install ${MY_DIR}/bin/tanzu-plugin-cluster  "${XDG_DATA_HOME}/tanzu-cli"
# install ${MY_DIR}/bin/tanzu-plugin-clustergroup  "${XDG_DATA_HOME}/tanzu-cli"
install ${MY_DIR}/bin/tanzu-plugin-management-cluster  "${XDG_DATA_HOME}/tanzu-cli"
install ${MY_DIR}/bin/tanzu-plugin-kubernetes-release  "${XDG_DATA_HOME}/tanzu-cli"
install ${MY_DIR}/bin/tanzu-plugin-login  "${XDG_DATA_HOME}/tanzu-cli"
install ${MY_DIR}/bin/tanzu-plugin-pinniped-auth  "${XDG_DATA_HOME}/tanzu-cli"
install ${MY_DIR}/bin/tanzu-plugin-builder  "${XDG_DATA_HOME}/tanzu-cli"
install ${MY_DIR}/bin/tanzu-plugin-test  "${XDG_DATA_HOME}/tanzu-cli"
install ${MY_DIR}/bin/tanzu-plugin-extension  "${XDG_DATA_HOME}/tanzu-cli"

if [[ "$BUILD_OS" == "Darwin" ]] ;  then
  xattr -d com.apple.quarantine /usr/local/bin/tanzu
  xattr -d com.apple.quarantine  "${XDG_DATA_HOME}/tanzu-cli/tanzu-plugin-alpha"
  xattr -d com.apple.quarantine  "${XDG_DATA_HOME}/tanzu-cli/tanzu-plugin-cluster"
  # xattr -d com.apple.quarantine  "${XDG_DATA_HOME}/tanzu-cli/tanzu-plugin-clustergroup"
  xattr -d com.apple.quarantine  "${XDG_DATA_HOME}/tanzu-cli/tanzu-plugin-management-cluster"
  xattr -d com.apple.quarantine  "${XDG_DATA_HOME}/tanzu-cli/tanzu-plugin-kubernetes-release"
  xattr -d com.apple.quarantine  "${XDG_DATA_HOME}/tanzu-cli/tanzu-plugin-login"
  xattr -d com.apple.quarantine  "${XDG_DATA_HOME}/tanzu-cli/tanzu-plugin-pinniped-auth"
  xattr -d com.apple.quarantine  "${XDG_DATA_HOME}/tanzu-cli/tanzu-plugin-builder"
  xattr -d com.apple.quarantine  "${XDG_DATA_HOME}/tanzu-cli/tanzu-plugin-test"
  xattr -d com.apple.quarantine  "${XDG_DATA_HOME}/tanzu-cli/tanzu-plugin-extension"
fi

# repo config
rm -rf "${XDG_DATA_HOME}/tanzu-repository"
mkdir -p "${XDG_DATA_HOME}/tanzu-repository"

cp -f ${MY_DIR}/config.yaml "${XDG_DATA_HOME}/tanzu-repository"
