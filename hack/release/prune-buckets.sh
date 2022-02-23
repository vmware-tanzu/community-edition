#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

TCE_CI_BUILD="${TCE_CI_BUILD:-""}"
TCE_SCRATCH_DIR="${TCE_SCRATCH_DIR:-""}"

# required input
if [[ -z "${TCE_SCRATCH_DIR}" ]]; then
    echo "TCE_SCRATCH_DIR is not set"
    exit 1
fi

# make no-op folders in order to prevent plugin updates by deleting
# all binaries and preserving the folder structure using an .empty file
# this needs to be done for both tce and tanzu framework

# do this on TCE
pushd "./artifacts" || exit 1

find ./ -type f | grep -v "yaml" | xargs rm 
find ./ -type d | grep "test" | xargs rm -rf
for i in $(find ./ -type d | grep "v"); do echo "empty" >> "${i}/.empty"; done

popd || exit 1

# do this on tanzu framework
pushd "${TCE_SCRATCH_DIR}/tanzu-framework/artifacts" || exit 1

find ./ -type f | grep -v "yaml" | xargs rm 
find ./ -type d | grep "test" | xargs rm -rf
for i in $(find ./ -type d | grep "v"); do echo "empty" >> "${i}/.empty"; done

popd || exit 1

pushd "${TCE_SCRATCH_DIR}/tanzu-framework/artifacts-admin/" || exit 1

find ./ -type f | grep -v "yaml" | xargs rm 
find ./ -type d | grep "test" | xargs rm -rf
for i in $(find ./ -type d | grep "v"); do echo "empty" >> "${i}/.empty"; done

popd || exit 1
