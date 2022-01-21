#!/usr/bin/env bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace


TCE_REPO_PATH="$(git rev-parse --show-toplevel)"
E2E_PATH="${TCE_REPO_PATH}/test/e2e/testdata/velero"

BUCKET_PREFIX_random=${RANDOM}
BUCKET_PREFIX="velero-bucket-prefix-${BUCKET_PREFIX_random}"
export BUCKET_PREFIX

cat > "${E2E_PATH}"/velero.env <<EOF
BUCKET_PREFIX="${BUCKET_PREFIX}*"
EOF

# Creating a temp file to substitute environment variable
cat "${E2E_PATH}"/velero_values.yaml > "${E2E_PATH}"/velero_values_template.yaml
envsubst < "${E2E_PATH}"/velero_values_template.yaml > "${E2E_PATH}"/velero_values.yaml


