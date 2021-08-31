#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

if [[ -f "/home/ubuntu/.config/tanzu/tkg/configs/mytest_ClusterConfig" ]]; then
    echo "ClusterConfig still present"
    exit 1
fi
if [[ -f "/home/ubuntu/.config/tanzu/clusterconfigs/mytest.yaml" ]]; then
    echo "BootstrapConfig still present"
    exit 1
fi
# IS_CLUSTER_DELETED=$(kubectl config view | grep mytest-admin@mytest)
# if [[ "${IS_CLUSTER_DELETED}" == "" ]]; then
#     echo "Cluster has not been deleted"
#     exit 1
# fi

echo "Cluster deleted successfully!"
