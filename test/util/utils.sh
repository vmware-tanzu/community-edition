#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

function error {
    printf '\E[31m'; echo "$@"; printf '\E[0m'
}

function delete_kind_cluster {
	echo "Deleting local kind bootstrap cluster(s) running in Docker container(s)"
    docker ps --all --format "{{ .Names }}" | grep tkg-kind | xargs docker rm --force
}