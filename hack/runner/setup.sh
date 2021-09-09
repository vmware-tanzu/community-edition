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

succeeded=false
numOfTimesToRetry=3
numOfTimesToPoll=30
mustHaveStatusBefore=10

i=0
while [ $i -ne $numOfTimesToRetry ]
do
    set +x
    RUNNER_TOKEN=$(curl -X POST -H "Accept: application/vnd.github.v3+json" -H "authorization: Bearer ${GITHUB_TOKEN}" https://api.github.com/repos/vmware-tanzu/community-edition/actions/runners/registration-token | jq .token -r)
    INSTANCE_ID=$(aws ec2 run-instances --image-id "${AMI_ID}" --count 1 --instance-type t2.2xlarge --key-name default --security-group-ids sg-03315c6a6ce53b6b7 --region us-west-2 --subnet-id subnet-f0837888 --tag-specifications "ResourceType=instance,Tags=[{Key=token,Value=${RUNNER_TOKEN}}]" | jq .Instances[].InstanceId -r)
    set -x
    aws ec2 wait instance-running --instance-ids "${INSTANCE_ID}" --region us-west-2
    echo "${INSTANCE_ID}" | tee instanceid.txt

    # dont poll yet... give a small head start
    echo "Sleeping for 30s"
    sleep 30

    # retry loop and get signal on what to do...
    j=0
    while [ $j -ne $numOfTimesToPoll ]
    do
        RUNNER_STATUS=$(curl -H "Accept: application/vnd.github.v3+json" -H "authorization: Bearer ${GITHUB_TOKEN}" https://api.github.com/repos/vmware-tanzu/community-edition/actions/runners | jq -r ".runners[] | select(.name==\"${INSTANCE_ID}\") | .status")

        if [[ "${RUNNER_STATUS}" == "" && $j -gt $mustHaveStatusBefore ]]; then
            echo "The node should have already returned some status... retry"
            break
        fi
        if [[ "${RUNNER_STATUS}" == "online" ]]; then
            echo "Succeeded!"
            succeeded=true
            break
        fi

        j=$((j+1))
        echo "Attempt poll $j... sleeping"
        sleep 10
    done

    if [[ $succeeded == "true" ]]; then
        echo "Succeeded! breaking the retry loop!"
        break
    fi

    # delete instance and retry
    i=$((i+1))
    echo "Retrying setup $i..."
    aws ec2 terminate-instances --instance-ids "${INSTANCE_ID}" --region us-west-2 > /dev/null 2>&1
done

if [ $succeeded == "true" ]; then
    echo "setup completed successfully"
else
    echo "Timed out trying to initialize runner"
    echo "timedout" | tee timedout.txt
fi
