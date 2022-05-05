#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

# Change directories to the parent directory of the one in which this
# script is located.
cd "$(dirname "${BASH_SOURCE[0]}")/../.."

# mdlint rules with common errors and possible fixes can be found here:
# https://github.com/DavidAnson/markdownlint/blob/main/doc/Rules.md
docker run --rm -v "$(pwd)":/build \
  gcr.io/cluster-api-provider-vsphere/extra/mdlint:0.23.2 -- /md/lint \
  -i docs/site/themes/template/static/fonts \
  -i docs/site/public/fonts \
  -i docs/site/content/plugins \
  -i docs/site/content/docs/edge/assets \
  -i docs/site/content/contributors \
  -i LICENSE .
