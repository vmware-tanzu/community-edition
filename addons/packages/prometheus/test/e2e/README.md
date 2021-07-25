# Prometheus package e2e tests

This file depicts the basic information about end2end tests of the Prometheus package.

## Prometheus

A time series database for your metrics.

## Steps involved
  
 In order to test Prometheus , the below manual steps are to be performed.The same steps are used for automation by using go Ginkgo framework.

- Step-1 ```tanzu package install prometheus.tce.vmware.com```
- Step-2 Make sure that the deployment is ready to serve before going to run tests.
   ```kubectl get deployment deployment.apps/prometheus-server -n prometheus-addon -o jsonpath={..status.conditions[?(@.type=="Available")].status}```
- Step-3 once everything is up and running , the actual tests will do two things.
  - port forward : ```kubectl port-forward deployment.apps/prometheus -n prometheus-addon   56018:9090```
  - check the health of prometheus api: ```crul -I http://127.0.0.1:56018/-/healthy``` .But this step is done in the tests using core golang http pacakge.
  - If the status is 200, that means prometheus is up and running.
- Step-4 : after tests are run , suite will delete prometheus package ```tanzu package delete prometheus.tce.vmware.com```

## Prerequisites

- Before running the suite , the cluster should be up and running.Prometheus test suite will never create or destroy any cluster.
- Before running the suite, All depending packages should be installed.

## How to run the tests

- To run individual suits , ```cd addons/packages/prometheus/test/e2e``` from the tce root path and run ```ginkgo -v -r```