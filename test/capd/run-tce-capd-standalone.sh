#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e
set -x

MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

guest_cluster_name="guest-cluster-${RANDOM}"

CLUSTER_PLAN=dev CLUSTER_NAME="$guest_cluster_name" tanzu standalone-cluster create ${guest_cluster_name} -i docker -v 10

"${MY_DIR}"/check-tce-cluster-creation.sh ${guest_cluster_name}-admin@${guest_cluster_name}
