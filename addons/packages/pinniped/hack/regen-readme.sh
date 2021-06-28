#!/usr/bin/env bash

set -euo pipefail

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
ME="$( basename "${BASH_SOURCE[0]}" )"
README="${MY_DIR}/../README.md"
pushd "$MY_DIR" >/dev/null
  go run gen_readme.go "$ME" "${MY_DIR}/readme.tmpl" "${MY_DIR}/../bundle/config/values.yaml" >"$README"
popd >/dev/null
