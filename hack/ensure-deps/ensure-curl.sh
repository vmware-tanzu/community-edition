#!/usr/bin/env bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script can be used to check your local development environment
# for necessary dependencies used in TCE

set -o nounset
set -o pipefail
set -o xtrace

TCE_CI_BUILD="${TCE_CI_BUILD:-""}"

SUDO_CMD="sudo"
if [[ "${TCE_CI_BUILD}" == "true" ]]; then
  SUDO_CMD=""
fi

if [[ -z "$(command -v curl)" ]]; then
  if [[ -f "/etc/redhat-release" ]]; then
    ${SUDO_CMD} yum -y install curl
  elif [[ "$(grep Ubuntu /etc/os-release)" != "" ]]; then
    ${SUDO_CMD} apt-get -y install curl
  else
    echo "**** Please install curl before proceeding *****"
    exit 1
  fi
fi
