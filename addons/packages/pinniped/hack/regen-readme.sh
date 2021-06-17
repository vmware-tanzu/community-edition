#!/usr/bin/env bash

set -euo pipefail

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
README="${MY_DIR}/../README.md"
pushd "$MY_DIR" >/dev/null
  go run gen_readme.go "${MY_DIR}/readme.tmpl" "${MY_DIR}/../bundle/config/values.yaml" >"$README"
popd >/dev/null
