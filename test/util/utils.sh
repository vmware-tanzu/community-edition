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

function collect_standalone_cluster_diagnostics {
    cluster_infra=$1
    cluster_name=$2

    if [[ -z "${cluster_infra}" ]]; then
        error "Cluster infra (vsphere, azure, aws, docker) not passed to collect_standalone_cluster_diagnostics function. Ignoring error. Usage example: collect_standalone_cluster_diagnostics vsphere test-cluster-1234"
    fi

    if [[ -z "${cluster_name}" ]]; then
        error "Cluster name not passed to collect_standalone_cluster_diagnostics function. Ignoring error. Usage example: collect_standalone_cluster_diagnostics vsphere test-cluster-1234"
    fi

    echo "Collecting ${cluster_name} standalone cluster diagnostics data"

    tanzu diagnostics collect --workload-cluster-infra "${cluster_infra}" \
        --workload-cluster-name "${cluster_name}" || {
        error "There was an error collecting tanzu diagnostics data. Ignoring the error"
    }
}

function collect_management_cluster_diagnostics {
    cluster_name=$1

    if [[ -z "${cluster_name}" ]]; then
        error "Cluster name not passed to collect_management_cluster_diagnostics function. Ignoring error. Usage example: collect_management_cluster_diagnostics management-cluster-1234"
    fi

    echo "Collecting ${cluster_name} management cluster diagnostics data"

    tanzu diagnostics collect --management-cluster-name "${cluster_name}" || {
        error "There was an error collecting tanzu diagnostics data. Ignoring the error"
    }
}

function collect_management_and_workload_cluster_diagnostics {
    cluster_infra=$1
    management_cluster_name=$2
    workload_cluster_name=$3

    if [[ -z "${cluster_infra}" ]]; then
        error "Workload cluster infra (vsphere, azure, aws, docker) not passed to collect_management_and_workload_cluster_diagnostics function. Ignoring error. Usage example: collect_management_and_workload_cluster_diagnostics vsphere management-cluster-1234 workload-cluster-1234"
    fi

    if [[ -z "${management_cluster_name}" ]]; then
        error "Management cluster name not passed to collect_management_and_workload_cluster_diagnostics function. Ignoring error. Usage example: collect_management_and_workload_cluster_diagnostics vsphere management-cluster-1234 workload-cluster-1234"
    fi

    if [[ -z "${workload_cluster_name}" ]]; then
        error "Management cluster name not passed to collect_management_and_workload_cluster_diagnostics function. Ignoring error. Usage example: collect_management_and_workload_cluster_diagnostics vsphere management-cluster-1234 workload-cluster-1234"
    fi

    echo "Collecting ${management_cluster_name} management cluster and ${workload_cluster_name} workload cluster diagnostics data"

    tanzu diagnostics collect --bootstrap-cluster-skip \
        --management-cluster-name "${management_cluster_name}" \
        --workload-cluster-infra "${cluster_infra}" \
        --workload-cluster-name "${workload_cluster_name}" || {
        error "There was an error collecting tanzu diagnostics data. Ignoring the error"
    }    
}
