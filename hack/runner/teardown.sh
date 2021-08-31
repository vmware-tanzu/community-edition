#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

if [[ -z "${GITHUB_TOKEN}" ]]; then
    echo "GITHUB_TOKEN is not set"
    exit 1
fi
if [[ -z "${INSTANCE_ID}" ]]; then
    echo "INSTANCE_ID is not set"
    exit 1
fi

RUNNER_ID=$(curl -H "Accept: application/vnd.github.v3+json" -H "authorization: Bearer ${GITHUB_TOKEN}" https://api.github.com/repos/vmware-tanzu/community-edition/actions/runners | jq -r ".runners[] | select(.name==\"${INSTANCE_ID}\") | .id")
curl -X DELETE -H "Accept: application/vnd.github.v3+json" -H "authorization: Bearer ${GITHUB_TOKEN}" "https://api.github.com/repos/vmware-tanzu/community-edition/actions/runners/${RUNNER_ID}"
aws ec2 terminate-instances --instance-ids "${INSTANCE_ID}" --region us-west-2 > /dev/null 2>&1
aws ec2 wait instance-terminated --instance-ids "${INSTANCE_ID}" --region us-west-2
