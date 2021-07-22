# fluent-bit package e2e tests

This file depicts the basic information about end2end tests of the fluent-bit package.

## fluent-bit

This package collect any data like metrics and logs from different sources, enrich them with filters and send them to multiple destinations using [fluent-bit](https://github.com/fluent/fluent-bit).

## Steps involved
  
 In order to test fluent-bit , the below manual steps are to be performed.The same steps are used for automation by using go Ginkgo framework.

- Step-1 First checking the installation is done without error or not by storing the output the command ``` tanzu package install fluent-bit.tce.vmware.com ``` and checking for the desired result in the output.
- Step-2 Make sure that the number of ports are 2 before going to run tests.
  ```Kubectl get daemonset.apps/fluent-bit -n fluent-bit -o jsonpath={..status.desiredNumberScheduled}```
- Step-3 Make sure that the number of availabe ports is 2 before going to run tests.  
  ```Kubectl get daemonset.apps/fluent-bit -n fluent-bit -o jsonpath={..status.numberAvailable}`}```
- Step-3 once everything is up and running , the actual tests will do two things.
  - port forward : ```Kubectl port-forward daemonset/fluent-bit -n fuent-bit 56017:2020```
  - check the health of fluent-bit api: ```curl -I http://127.0.0.1:56017/api/v1/health``` .But this step is done in the tests using core golang http pacakge.
  - If the status is 200, that means fluent-bit is up and running.
- Step-4 : after tests are run , suite will delete fluent-bit package ```tanzu package delete fluent-bit.tce.vmware.com```

## Prerequisites

- Before running the suite , the cluster should be up and running.fluent-bit test suite will never create or destroy any cluster.
- Before running the suite, All depending packages should be installed.

## How to run the tests

- To run individual suits , ```cd addons/packages/fluent-bit/test/e2e``` from the tce root path and run ```ginkgo -v -r```
