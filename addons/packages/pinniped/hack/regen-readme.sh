#!/usr/bin/env bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
ME="$( basename "${BASH_SOURCE[0]}" )"

# Build gen_readme to a temp directory that will be cleaned up on exit.
TMP="$(mktemp -d)"
cleanup() {
  rm -rf "$TMP"
}
trap "cleanup" EXIT SIGINT
GEN_README="${TMP}/gen_readme"
pushd "$MY_DIR" >/dev/null
  go build -o "$GEN_README" gen_readme.go
popd >/dev/null

# For each version of the package, generate the README.
for version_dir in "$MY_DIR"/../*; do
  if [[ "$(basename "$version_dir")" =~ [0-9]+.[0-9]+.[0-9]+ ]]; then
    version_readme="${version_dir}/README.md"
    pushd "$version_dir" >/dev/null
      "$GEN_README" "$ME" "${MY_DIR}/readme.tmpl" "${version_dir}/bundle/config/values.yaml" >"$version_readme"
    popd >/dev/null
  fi
done
