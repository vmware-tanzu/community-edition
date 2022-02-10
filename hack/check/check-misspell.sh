#!/usr/bin/env bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0


set -o errexit
set -o nounset
set -o pipefail


CHECK_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )"; pwd )"
HACK_DIR="$( cd "${CHECK_DIR}/.."; pwd )"

# Install tools we need if it is not present
if [[ ! -f "${HACK_DIR}/tools/bin/misspell" ]]; then
  mkdir -p "${HACK_DIR}/tools/bin"
  curl -L https://git.io/misspell | BINDIR="${HACK_DIR}/tools/bin" bash
fi

# Spell checking
# misspell check Project - https://github.com/client9/misspell
misspellignore_files="${CHECK_DIR}/.misspellignore"
ignore_files=$(cat "${misspellignore_files}")
git ls-files | grep -v "${ignore_files}" | xargs "${HACK_DIR}/tools/bin/misspell" | grep "misspelling" && echo "Please fix the listed misspell errors and verify using 'make misspell'" && exit 1 || echo "misspell check passed!"
