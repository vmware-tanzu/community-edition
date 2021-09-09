#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

# make user CI directory
apt update
apt install zip unzip

# override https
git config --global url."https://git:${GITHUB_TOKEN}@github.com".insteadOf "https://github.com"
