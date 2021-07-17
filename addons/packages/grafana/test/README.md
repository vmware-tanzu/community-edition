# Grafana package e2e tests

This file depicts the basic information about end2end tests of the Grafana package.

## Grafana

Grafana is open source visualization and analytics software. It allows you to query, visualize, alert on, and explore your metrics no matter where they are stored. In plain English, it provides you with tools to turn your time-series database (TSDB) data into beautiful graphs and visualizations.

## Steps involved
  
 In order to test Grafana , the below manual steps are to be performed.The same steps are used for automation by using go Ginkgo framework.

- Step-1 ```tanzu package install grafana.tce.vmware.com```
- Step-2 Make sure that the deployment is ready to serve before going to run tests.
   ```kubectl get deployment -l app.kubernetes.io/name=grafana -n grafana-addon -o jsonpath={..status.conditions[?(@.type=="Available")].status}```
- Step-3 once everything is up and running , the actual tests will do two things.
  - port forward : ```kubectl port-forward   deployment.apps/grafana -n grafana-addon   56016:3000```
  - check the health of grafana api: ```crul -I http://127.0.0.1:56016/api/health``` .But this step is done in the tests using core golang http pacakge.
  - If the status is 200, that means grafana is up and running.
- Step-4 : after tests are run , suite will delete grafana package ```tanzu package delete grafana.tce.vmware.com```

## Prerequisites

- Before running the suite , the cluster should be up and running.Grafana test suite will never create or destroy any cluster.
- Before running the suite, All depending packages should be installed.

## How to run the tests

- To run individual suits , ```cd addons/packages/grafana/test``` from the tce root path and run ```ginkgo -v -r```
