#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights
# Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

readonly TCE_REPO_PATH="$(git rev-parse --show-toplevel)"
# shellcheck source=test/smoke/packages/utils/smoke-tests-utils.sh
source "${TCE_REPO_PATH}/test/smoke/packages/utils/smoke-tests-utils.sh"

readonly MY_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

main() {
        apply_dependencies

        apply_workload || failureMessage cartographer-catalog
        successMessage cartographer-catalog
}

apply_dependencies() {
        apply_dependency cert-manager
        apply_dependency cartographer
        apply_dependency fluxcd-source-controller

        apply_dependency kpack \
                --values-file="${MY_DIR}"/kpack.yaml
        apply_dependency kpack-dependencies \
                --values-file="${MY_DIR}"/kpack-dependencies.yaml
        apply_dependency cartographer-catalog \
                --values-file="${MY_DIR}"/cartographer-catalog.yaml
        apply_dependency knative-serving \
                --values-file="${MY_DIR}"/knative-serving.yaml
}

apply_dependency() {
        local name=$1
        shift
        local extra_arg=$*
        local package_name=$name.community.tanzu.vmware.com
        local version

        tanzu package installed list | grep "$package_name" || {
                version=$(
                        tanzu package available list "$package_name" |
                                tail -n 1 |
                                awk '{print $2}'
                )

                # shellcheck disable=SC2086
                tanzu package install "$name" \
                        --package-name "$package_name" \
                        --version "${version}" $extra_arg
        }
}

apply_workload() {
        kapp deploy --yes \
                -a cartographer-test \
                -f "${MY_DIR}"/cartographer-catalog-test.yaml
}

main
