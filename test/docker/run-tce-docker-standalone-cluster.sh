#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e
set -x

TCE_REPO_PATH="$(git rev-parse --show-toplevel)"

"${TCE_REPO_PATH}/test/install-dependencies.sh"
"${TCE_REPO_PATH}/test/build-tce.sh"

guest_cluster_name="guest-cluster-${RANDOM}"

CLUSTER_PLAN=dev CLUSTER_NAME="$guest_cluster_name" tanzu standalone-cluster create ${guest_cluster_name} -i docker -v 10

"${TCE_REPO_PATH}/test/check-tce-cluster-creation.sh" ${guest_cluster_name}-admin@${guest_cluster_name}
