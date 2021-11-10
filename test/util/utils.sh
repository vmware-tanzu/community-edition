#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

function error {
    printf '\E[31m'; echo "$@"; printf '\E[0m'
}

function delete_kind_cluster {
	echo "Deleting local kind bootstrap cluster(s) running in Docker container(s)"
    docker ps --all --format "{{ .Names }}" | grep tkg-kind | xargs -r docker rm --force
}

function kubeconfig_cleanup {
    cluster_name=$1

    if [[ -z "${cluster_name}" ]]; then
        echo "Cluster name not passed to kubeconfig_cleanup function. Usage example: kubeconfig_cleanup management-cluster-1234"
        return 1
    fi

    echo "Removing cluster context, cluster user and cluster entry from kubeconfig for ${cluster_name} cluster"

    # ignore errors assuming the error happens only when the values are already deleted
    kubectl config delete-context "${cluster_name}-admin@${cluster_name}" || true
    kubectl config delete-user "${cluster_name}-admin" || true
    kubectl config delete-cluster "${cluster_name}" || true
}
