#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# rsync is required on host to run this script

BUILD_DIR=$(mktemp -d /tmp/XXXXXX)
echo "Working in ${KIND_DIR}"

if [[ "$1" == "" ]]; then
  echo "specify a kubernetes version (tag)"
  exit 1
fi

K8S_TAG=$1
REPO_NAME="projects.registry.vmware.com/tce/kind/node"

# enter kind directory
pushd "${BUILD_DIR}" || exit
git clone --depth 1 --branch v0.11.1 https://github.com/kubernetes-sigs/kind
git clone --depth 1 --branch "${K8S_TAG}" https://github.com/kubernetes/kubernetes
pushd kind || exit

# build the kind binary
make build

# create the kind node image
bin/kind build node-image \
  --image "projects.registry.vmware.com/tce/kind/node:${K8S_TAG}" \
  "${BUILD_DIR}/kubernetes"

echo "Image created"
echo "To publish this image, run:"
echo "docker push ${REPO_NAME}:${K8S_TAG}"
