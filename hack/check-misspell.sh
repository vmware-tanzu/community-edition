#!/usr/bin/env bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0


set -o errexit
set -o nounset
set -o pipefail


MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Install tools we need if it is not present
if [[ ! -f "${MY_DIR}/tools/bin/misspell" ]]; then
  curl -L https://git.io/misspell | bash
  mkdir -p "${MY_DIR}/tools/bin"
  mv ./bin/misspell "${MY_DIR}/tools/bin/misspell"
fi

# Spell checking
# misspell check Project - https://github.com/client9/misspell
misspellignore_files="${MY_DIR}/.misspellignore"
ignore_files=$(cat "${misspellignore_files}")
git ls-files | grep -v "${ignore_files}" | xargs "${MY_DIR}/tools/bin/misspell" | grep "misspelling" && echo "Please fix the errors mentioned and try running make misspell" && exit 1 || echo "Misspell check Pass"
