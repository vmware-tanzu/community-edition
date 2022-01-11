#!/usr/bin/env bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script can be used to check your local development environment
# for necessary dependencies used in TCE

set -o nounset
set -o pipefail
set -o xtrace

TCE_CI_BUILD="${TCE_CI_BUILD:-""}"

if [[ "${TCE_CI_BUILD}" != "true" ]]; then

if [[ -z "$(command -v zip)" ]]; then
  echo "**** Please install zip before proceeding *****"
  exit 1
fi
if [[ -z "$(command -v unzip)" ]]; then
    echo "**** Please install unzip before proceeding *****"
    exit 1
fi

fi
