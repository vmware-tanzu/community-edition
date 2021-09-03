#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

# Change directories to the parent directory of the one in which this
# script is located.
cd "$(dirname "${BASH_SOURCE[0]}")/.."

usage() {
  cat <<EOF
usage: ${0} [FLAGS]
  Lints the project's shell scripts using Docker.

FLAGS
  -l    use your local shellcheck (warning - errors may differ from docker!)
  -h    prints this help screen
EOF
}

while getopts ':lh' opt; do
  case "${opt}" in
  l)
    DO_DOCKER=0
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
