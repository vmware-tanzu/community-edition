#! /bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

# we only allow this to run from CI
if [[ "${TCE_CI_BUILD}" != "true" ]]; then
    echo "This is only meant to be run within CI"
    exit 1
fi

BOLD='\033[1m'
CLEAR='\033[0m'

function info {
  BLUE='\033[0;34m'
  local message=${1}

  echo -e "${BLUE}${BOLD}${message}${CLEAR}"
}

function warning {
  YELLOW='\033[0;33m'
  local message=${1}

  echo -e "${YELLOW}${BOLD}${message}${CLEAR}"
}

function error {
  RED='\033[0;31m'
  local message=${1}

  echo -e "${RED}${BOLD}${message}${CLEAR}"
}

function success {
  GREEN='\033[0;32m'
  local message=${1}

  echo -e "${GREEN}${BOLD}${message}${CLEAR}"
}

function check_dependencies {
  # this script requires git & the gh CLI tool
  warning "--- Checking dependencies ---"

  if ! gh --version > /dev/null 2>&1; then
    error "gh binary test failed. GitHub CLI tool required"
    exit 1
  fi

  if ! gh auth status > /dev/null 2>&1; then
    error "gh auth check failed. Ensure you've logged into to the gh CLI using \`gh auth login\` or have GITHUB_TOKEN set in your environment"
    exit 1
  fi

  if [[ -z "${BUILD_VERSION:-}" ]]; then
    error "BUILD_VERSION environment variable is not set"
    exit 1
  fi

  success "--- Dependencies verified --- \n"
}

check_dependencies

warning "Checking if ${BUILD_VERSION} release exists"

if ! gh release view --repo https://github.com/vmware-tanzu/community-edition "${BUILD_VERSION}" > /dev/null 2>&1; then
    error "${BUILD_VERSION} release does not exist. Ensure ${BUILD_VERSION} TCE release has been created"
    exit 1
fi

info "Triggering MacOS and Linux build tests"

# MacOS and Linux Build Test
gh workflow run --repo https://github.com/vmware-tanzu/community-edition/ \
  --ref main \
  release-gate-macos-linux-builds.yaml \
  --raw-field release_version="${BUILD_VERSION}"

info "Triggering Windows build tests"

# Windows Build Test
gh workflow run --repo https://github.com/vmware-tanzu/community-edition/ \
  --ref main \
  release-gate-windows-build.yaml \
  --raw-field release_version="${BUILD_VERSION}"
