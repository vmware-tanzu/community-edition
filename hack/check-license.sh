#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

HACK_DIR="$(dirname "${BASH_SOURCE[0]}")"
REPO_ROOT_DIR="$(dirname "${HACK_DIR}"/)"

echo "Checking License in files ..."
required_keywords=("VMware Tanzu Community Edition contributors" "SPDX-License-Identifier" "Apache-2.0")
extensions_to_check=("sh" "go")

check_output=$(mktemp /tmp/tce-licence-check.XXXXXX)
for ext in "${extensions_to_check[@]}"; do
  find "$REPO_ROOT_DIR" -name "*.$ext" -a \! -path "*vendor/*" -a \! -path "./.*" -a \! -path "./addons/*" -print0 |
    while IFS= read -r -d '' path; do
      for rword in "${required_keywords[@]}"; do
        if ! grep -q "$rword" "$path"; then
          echo "   $path" >> "$check_output"
        fi
      done
    done
done

if [ -s "$check_output" ]; then
  echo "No license header found in:"
  sort < "$check_output" | uniq
  echo "License check failed!"
  echo "Please add the license in each listed file and verify using 'make checklicense'"
  rm "$check_output"
  exit 1
else
  echo "License check passed!"
fi
rm "$check_output"
