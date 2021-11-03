#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0



# Checking package is installed or not
tanzu package installed list | grep "contour.community.tanzu.vmware.com"
exitcode_contour=$?
MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
echo "${MY_DIR}"
if [ "${exitcode_contour}" == 1 ] 
then
    version=$(tanzu package available list contour.community.tanzu.vmware.com | tail -n 1 | awk '{print $2}')
    tanzu package install contour --package-name contour.community.tanzu.vmware.com -f "${MY_DIR}"/contour-values.yaml --version "${version}"
fi
sleep 15s

# Providing prerequisite 
kubectl create namespace contour-example-workload
kubectl create deployment nginx-example --image nginx --namespace contour-example-workload
kubectl create service clusterip nginx-example --tcp 80:80 --namespace contour-example-workload

kubectl apply -f - <<EOF
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  name: nginx-example-proxy
  namespace: contour-example-workload
  labels:
    app: ingress
spec:
  virtualhost:
    fqdn: nginx-example.projectcontour.io
  routes:
  - conditions:
    - prefix: /
    services:
    - name: nginx-example
      port: 80
EOF

sleep 15s
kubectl --namespace projectcontour port-forward svc/envoy 5436:80 &
sleep 5s

curl -s -H "Host: nginx-example.projectcontour.io" http://localhost:5436 | grep "<h1>Welcome to nginx!</h1>"
exitcode=$?
kubectl delete ns contour-example-workload
tanzu package installed delete contour -y

if [ "${exitcode}" == 1 ] 
then
    echo "Contour failed"
    exit 1
else
    echo "Countour Passed"
fi
