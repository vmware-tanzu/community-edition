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
if [[ -z "${AMI_ID}" ]]; then
    echo "AMI_ID is not set"
    exit 1
fi

RUNNER_TOKEN=$(curl -X POST -H "Accept: application/vnd.github.v3+json" -H "authorization: Bearer ${GITHUB_TOKEN}" https://api.github.com/repos/vmware-tanzu/community-edition/actions/runners/registration-token | jq .token -r)
INSTANCE_ID=$(aws ec2 run-instances --image-id "${AMI_ID}" --count 1 --instance-type t2.2xlarge --key-name default --security-group-ids sg-03315c6a6ce53b6b7 --region us-west-2 --subnet-id subnet-f0837888 --tag-specifications "ResourceType=instance,Tags=[{Key=token,Value=${RUNNER_TOKEN}}]" | jq .Instances[].InstanceId -r)
aws ec2 wait instance-running --instance-ids "${INSTANCE_ID}" --region us-west-2
echo "${INSTANCE_ID}" | tee instanceid.txt

sleep 30

i=0
while [ $i -ne 30 ]
do
    RUNNER_STATUS=$(curl -H "Accept: application/vnd.github.v3+json" -H "authorization: Bearer ${GITHUB_TOKEN}" https://api.github.com/repos/vmware-tanzu/community-edition/actions/runners | jq -r ".runners[] | select(.name==\"${INSTANCE_ID}\") | .status")

    if [[ "${RUNNER_STATUS}" == "online" ]]; then
        break
    fi

    i=$((i+1))
    echo "Attempt $i... sleeping"
    sleep 10
done

if [ $i == 30 ]; then
    echo "Timed out trying to initialize runner"
    echo "timedout" | tee timedout.txt
else
    echo "setup completed successfully"
fi
