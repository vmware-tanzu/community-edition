#!/bin/bash

# Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# fix the Makefile for the plugins before building
pushd "${MY_DIR}" || exit 1

sed -i.bak -e "s|--artifacts[ ]\+../../\${ARTIFACTS_DIR}/\${OS}/\${ARCH}/cli|--artifacts ../../\${ARTIFACTS_DIR}|g" ./Makefile && rm ./Makefile.bak
sed -i.bak -e "s|--artifacts[ ]\+../../\${ARTIFACTS_DIR}/\${GOHOSTOS}/\${ARCH}/cli|--artifacts ../../\${ARTIFACTS_DIR}|g" ./Makefile && rm ./Makefile.bak

popd || exit 1
