#!/usr/bin/env bash

set -euo pipefail

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
ME="$( basename "${BASH_SOURCE[0]}" )"
for version_dir in "$MY_DIR"/../*; do
  if [[ "$(basename "$version_dir")" =~ [0-9].[0-9].[0-9] ]]; then
    version_readme="${version_dir}/README.md"
    pushd "$version_dir" >/dev/null
      go run "${MY_DIR}/gen_readme.go" "$ME" "${MY_DIR}/readme.tmpl" "${version_dir}/bundle/config/values.yaml" >"$version_readme"
    popd >/dev/null
  fi
done
