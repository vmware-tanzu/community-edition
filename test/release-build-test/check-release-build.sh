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
  for binary in "${TCE_INSTALLATION_DIR}"/bin/*; do
    spctl -vv --type install --asses "${binary}"
  done
fi

"${TCE_INSTALLATION_DIR}"/install.sh

tanzu version

tanzu cluster version

tanzu conformance version

tanzu diagnostics version

tanzu kubernetes-release version

tanzu management-cluster version

tanzu package version

tanzu standalone-cluster version

tanzu pinniped-auth version

tanzu builder version

tanzu login version
