#!/usr/bin/env bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script can be used to check your local development environment
# for necessary dependencies used in TCE

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

BUILD_OS=$(uname 2>/dev/null || echo Unknown)
TCE_REPO_PATH="$(git rev-parse --show-toplevel)"

case "${BUILD_OS}" in
  Linux)
    ;;
  Darwin)
    ;;
  *)
    echo "${BUILD_OS} is unsupported"
    exit 1
    ;;
esac

i=0

# these are must have dependencies to just get going
if [[ -z "$(command -v go)" ]]; then
    echo "Missing go"
    ((i=i+1))
fi

if [[ -z "$(command -v docker)" ]]; then
    echo "Missing docker"
    ((i=i+1))
fi
# these are must have dependencies to just get going

if [[ $i -gt 0 ]]; then
    echo "Total missing: $i"
    echo "Please install these minimal dependencies in order to continue"
    exit 1
fi

"${TCE_REPO_PATH}/hack/ensure-deps/ensure-curl.sh"
"${TCE_REPO_PATH}/hack/ensure-deps/ensure-zip.sh"
"${TCE_REPO_PATH}/hack/ensure-deps/ensure-imgpkg.sh"
"${TCE_REPO_PATH}/hack/ensure-deps/ensure-kbld.sh"
"${TCE_REPO_PATH}/hack/ensure-deps/ensure-ytt.sh"
"${TCE_REPO_PATH}/hack/ensure-deps/ensure-shellcheck.sh"
"${TCE_REPO_PATH}/hack/ensure-deps/ensure-gh-cli.sh"
"${TCE_REPO_PATH}/hack/ensure-deps/ensure-aws-cli.sh"
"${TCE_REPO_PATH}/hack/ensure-deps/ensure-jq.sh"
"${TCE_REPO_PATH}/hack/ensure-deps/ensure-kubectl.sh"
"${TCE_REPO_PATH}/hack/ensure-deps/ensure-kind.sh"

echo "No missing dependencies!"
