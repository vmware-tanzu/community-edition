#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# WIP script to test TCE in AWS
# Basic idea is to build the latest release of TCE, spin up a standalone cluster
# in AWS, install the default packages, test the packages and clean the environment 
# Note: This is WIP and supports only Linux(Debian) and MacOS
# Following environment variables are expected to be exported before running the script
# AWS_ACCOUNT_ID
# AWS_ACCESS_KEY_ID
# AWS_SECRET_ACCESS_KEY
# AWS_B64ENCODED_CREDENTIALS
# AWS_SSH_KEY_NAME
# Region is set to us-east-2

source test/build-tce.sh

# Set standalone cluster name
export CLUSTER_NAME="test${RANDOM}"
echo "Setting CLUSTER_NAME to ${CLUSTER_NAME}..."

# Cleanup function
function deletecluster {
    echo "$@"
    tanzu standalone-cluster delete ${CLUSTER_NAME} -y || { aws-nuke-tear-down "STANDALONE CLUSTER DELETION FAILED! Deleting the cluster using AWS-NUKE..."; }
}

# Back-up function to cleanup incase deletecluster fails
function aws-nuke-tear-down {
	echo "$@"
	envsubst < "${TCE_REPO_PATH}"/test/aws/nuke-config-template.yml > "${TCE_REPO_PATH}"/test/aws/nuke-config.yml
	aws-nuke -c "${TCE_REPO_PATH}"/test/aws/nuke-config.yml --access-key-id "$AWS_ACCESS_KEY_ID" --secret-access-key "$AWS_SECRET_ACCESS_KEY" --force --no-dry-run || { error "STANDALONE CLUSTER DELETION FAILED!!!"; rm -rf "${TCE_REPO_PATH}"/test/aws/nuke-config.yml; exit 1; }
	rm -rf "${TCE_REPO_PATH}"/test/aws/nuke-config.yml
	echo "STANDALONE CLUSTER DELETED using aws-nuke!"
}

# E2E test for Gatekeeper
function test-gatekeeper {
    echo "Installing Gatekeeper..."
    tanzu package install gatekeeper.tce.vmware.com || { error "Gatekeeper installation failed. TEST FAILED."; deletecluster "Deleting standalone cluster"; exit 1; }
    # Added this as it takes time to create namespace for Gatekeeper
    sleep 10s
    echo "Verifying Gatekeeper installation..."
    kubectl wait --for=condition=ready pod --all -n gatekeeper-system --timeout=300s || { error "Timed out waiting for Gatekeeper pods to come up. TEST FAILED."; deletecluster "Deleting standalone cluster"; exit 1; }
    echo "Applying constraint template..."
    kubectl apply -f "${TCE_REPO_PATH}"/test/gatekeeper/constraint-template.yaml || { error "Unexpected error. TEST FAILED."; deletecluster "Deleting standalone cluster"; exit 1; }
    echo "Verifying creation of k8srequiredlabels CRD..."
    kubectl get crds | grep -i k8srequiredlabels || { error "Unexpected error. TEST FAILED."; deletecluster "Deleting standalone cluster"; exit 1; }
    echo "Creating constraint..."
    kubectl apply -f "${TCE_REPO_PATH}"/test/gatekeeper/constraint.yaml || { error "Unexpected error. TEST FAILED."; deletecluster "Deleting standalone cluster"; exit 1; }
    echo "Creating test namespace..."
    # It takes time for Gatekeeper webhook service to come up. Added retires to get around Internal Server Error.
    retries=1
    while [ $retries -le 5 ]
    do
        error_message=$(kubectl create ns test 2>&1)
        if echo "$error_message" | grep 'All namespaces must have an owner label'; then
            echo "Expected failure"
            echo "Creating test namespace with owner label..."
            kubectl apply -f "${TCE_REPO_PATH}"/test/gatekeeper/test-namespace.yaml || { error "TEST FAILED."; deletecluster "Deleting standalone cluster"; exit 1; }
            printf '\E[32m'; echo "TEST PASSED!"; printf '\E[0m'
            return
        else
            echo "Retrying..."
            sleep 60s
        fi
        retries=$(( retries+1 ))
    done
    error "TEST FAILED"; deletecluster "Deleting standalone cluster"; exit 1;
}

# Substitute env variables in aws-template
echo "Bootstrapping TCE standalone cluster on AWS..."
envsubst < "${TCE_REPO_PATH}"/test/aws/standalone-template.yaml > "${TCE_REPO_PATH}"/test/aws/standalone.yaml
tanzu standalone-cluster create "${CLUSTER_NAME}" -f "${TCE_REPO_PATH}"/test/aws/standalone.yaml || { error "STANDALONE CLUSTER CREATION FAILED!"; rm -rf "${TCE_REPO_PATH}"/test/aws/standalone.yaml; deletecluster "Deleting standalone cluster"; exit 1; }
rm -rf "${TCE_REPO_PATH}"/test/aws/standalone.yaml

kubectl config use-context "${CLUSTER_NAME}"-admin@"${CLUSTER_NAME}" || { error "CONTEXT SWITCH TO STANDALONE CLUSTER FAILED!"; deletecluster "Deleting standalone cluster"; exit 1; }

kubectl wait --for=condition=ready pod --all --all-namespaces --timeout=300s || { error "TIMED OUT WAITING FOR ALL PODS TO BE UP!"; deletecluster "Deleting standalone cluster"; exit 1; }

echo "Installing packages on TCE..."
tanzu package repository add tce-repo --namespace default --url projects.registry.vmware.com/tce/main@stable || { error "PACKAGE REPOSITORY INSTALLATION FAILED!"; deletecluster "Deleting standalone cluster"; exit 1; }

# Added this as the above command takes time to install packages
sleep 60s

tanzu package available list || { error "UNEXPECTED FAILURE OCCURRED!"; deletecluster "Deleting standalone cluster"; exit 1; }

echo "Starting Gatekeeper test..."
test-gatekeeper

# Clean up
deletecluster "Cleaning up..."