#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

version="${1:?TCE version argument empty. Example usage: ./test/release-build-test/check-release-build.sh v0.10.0}"
: "${GITHUB_TOKEN:?GITHUB_TOKEN is not set}"

TCE_REPO_URL="https://github.com/vmware-tanzu/community-edition"
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH="amd64"

temp_dir=$(mktemp -d)

TCE_TAR_BALL="${temp_dir}/tce-${OS}-${ARCH}-${version}.tar.gz"
TCE_INSTALLATION_DIR="${temp_dir}/tce-${OS}-${ARCH}-${version}"

gh release download "${version}" --repo ${TCE_REPO_URL} --pattern "tce-${OS}-${ARCH}-${version}.tar.gz" --dir "${temp_dir}"

tar xvzf "${TCE_TAR_BALL}" --directory "${temp_dir}"

if [ "${OS}" == 'darwin' ]; then
  pushd "${TCE_INSTALLATION_DIR}" || exit 1
    # tanzu cli
    spctl -vv --type install --asses "tanzu"

    # tanzu plugins
    pushd "./default-local/distribution/darwin/amd64/cli" || exit 1
      for file in $(find ./ -type f); do
        DARWIN_DIRECTORY=$(dirname ${file})
        DARWIN_FILENAME=$(basename ${file})

        pushd "./${DARWIN_DIRECTORY}" || exit 1
          spctl -vv --type install --asses "${DARWIN_FILENAME}"
        popd || exit 1
      done
    popd || exit 1
  popd || exit 1
fi

"${TCE_INSTALLATION_DIR}/install.sh"

tanzu version

tanzu cluster version

tanzu conformance version

tanzu diagnostics version

tanzu kubernetes-release version

tanzu management-cluster version

tanzu package version

tanzu pinniped-auth version

tanzu builder version

tanzu login version

tanzu unmanaged-cluster version
