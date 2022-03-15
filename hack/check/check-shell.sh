#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

# Change directories to the parent directory of the one in which this
# script is located.
cd "$(dirname "${BASH_SOURCE[0]}")/../.."

usage() {
  cat <<EOF
usage: ${0} [FLAGS]
  Lints the project's shell scripts using Docker.

FLAGS
  -l    use your local shellcheck (warning - errors may differ from docker!)
  -h    prints this help screen
  -x    verify license header
EOF
}

verify_license() {
  printf "\nChecking License in shell scripts and Makefiles ...\n"
  local required_keywords=("VMware Tanzu Community Edition contributors" "SPDX-License-Identifier: Apache-2.0")
  local file_patterns_to_check=("*.sh" "Makefile")

  local result
  result=$(mktemp /tmp/tce-licence-check.XXXXXX)
  for ext in "${file_patterns_to_check[@]}"; do
    find . -path ./vendor -prune -o -name "$ext" -type f -print0 |
      while IFS= read -r -d '' path; do
        for rword in "${required_keywords[@]}"; do
          if ! grep -q "$rword" "$path"; then
            echo "   $path" >> "$result"
          fi
        done
      done
  done

  if [ -s "$result" ]; then
    echo "No required license header found in:"
    sort < "$result" | uniq
    echo "License check failed!"
    echo "Please add the license in each listed file and verify using 'make shellcheck'"
    rm "$result"
    return 1
  else
    echo "License check passed!"
    rm "$result"
    return 0
  fi
}

while getopts ':lhx' opt; do
  case "${opt}" in
  l)
    DO_DOCKER=0
    ;;
  x)
    verify_license || exit 1; exit 0
    ;;
  h)
    usage 1>&2; exit 1
    ;;
  \?)
    { echo "invalid option: -${OPTARG}"; usage; } 1>&2; exit 1
    ;;
  :)
    echo "option -${OPTARG} requires an argument" 1>&2; exit 1
    ;;
  esac
done
shift $((OPTIND-1))

if [ "${DO_DOCKER-}" ]; then
  shellcheck --version
  find . -path ./vendor -prune -o -name "*.*sh" -type f -print0 | xargs -0 shellcheck
else
  docker run --rm -t -v "$(pwd)":/build:ro gcr.io/cluster-api-provider-vsphere/extra/shellcheck --version
  docker run --rm -t -v "$(pwd)":/build:ro gcr.io/cluster-api-provider-vsphere/extra/shellcheck
fi

verify_license
