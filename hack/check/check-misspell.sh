#!/usr/bin/env bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o nounset
set -o pipefail

CHECK_DIR="$(dirname "${BASH_SOURCE[0]}")"
MISSPELL_DIR="$(mktemp -d)"

# We always want a fresh install in case someone has replaced our binary locally
# or even if we just want the latest version. On a fresh build system, this would
# have been the case anyways.
rm -rf "${MISSPELL_DIR}"
mkdir -p "${MISSPELL_DIR}"
curl -L https://git.io/misspell | BINDIR="${MISSPELL_DIR}" bash

# Spell checking - misspell check Project - https://github.com/client9/misspell
misspellignore_files="${CHECK_DIR}/.misspellignore"
ignore_files=$(cat "${misspellignore_files}")

# check spelling
FILES_TO_CHECK=$(git ls-files | grep -v "${ignore_files}")
echo " "
echo " "
echo "Files to check spelling:"
echo "${FILES_TO_CHECK}"
echo " "
echo " "
echo "Errors:"
git ls-files | grep -v "${ignore_files}" | xargs "${MISSPELL_DIR}/misspell"

# delete the directory to return environment to original condition
rm -rf "${MISSPELL_DIR}"
