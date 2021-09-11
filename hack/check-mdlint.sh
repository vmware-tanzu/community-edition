#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

# Change directories to the parent directory of the one in which this
# script is located.
cd "$(dirname "${BASH_SOURCE[0]}")/.."

docker run --rm -v "$(pwd)":/build \
  gcr.io/cluster-api-provider-vsphere/extra/mdlint:0.17.0 /md/lint \
  -i docs/site -i LICENSE.md .
