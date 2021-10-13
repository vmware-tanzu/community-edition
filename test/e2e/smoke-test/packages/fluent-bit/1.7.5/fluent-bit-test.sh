#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# Checking package is installed or not
tanzu package installed list | grep "fluent-bit.community.tanzu.vmware.com"
exitcode=$?

if [ "${exitcode}" == 1 ] 
then
    version=$(tanzu package available list fluent-bit.community.tanzu.vmware.com | tail -n 1 | awk '{print $2}')
    tanzu package install fluent-bit --package-name fluent-bit.community.tanzu.vmware.com --version "${version}"
    
fi

echo "Wating for Pods to reach ready state ..."
#sleep 30s
pod_name="$(kubectl get pods -n fluent-bit | tail -n 1 | awk '{print $1}')"
kubectl logs "${pod_name}" -n fluent-bit | grep "Fluent Bit v1.7.5"
exitcode=$?

tanzu package installed delete fluent-bit -y

if [ "${exitcode}" == 1 ] 
then
    echo "Fluent-bit failed"
    exit 1
else
    echo "Fluent-bit Passed"
fi
