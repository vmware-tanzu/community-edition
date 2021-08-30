# Gatekeeper package e2e tests

This file depicts the basic information about end2end tests of the Gatekeeper package.

## Gatekeeper

This package provides custom admission control using
[gatekeeper](https://github.com/open-policy-agent/gatekeeper). Under the hood,
gatekeeper uses [Open Policy Agent](https://www.openpolicyagent.org) to enforce
policy when requests hit the Kubernetes API server.

## Steps involved

 In order to test Gatekeeper , the below manual steps are to be performed.The same steps are used for automation by using go Ginkgo framework.

- Step-1 ```tanzu package install gatekeeper.tce.vmware.com```
- Step-2  ```kubectl apply -f  ${TCE-REPO-DIR}addons/packages/gatekeeper/test/fixtures/constraint-template.yaml/constraint-template.yaml```
- Step-3 Check required CRDs after applying the Step-2 file. ```kubectl get crds | grep -i k8srequiredlabels```
- Step-4   ```kubectl apply -f ${TCE-REPO-DIR}addons/packages/gatekeeper/test/fixtures/constraint.yaml```
- Step-5  ```kubectl create ns test```  User cannot create a namespace without owner (The actual use case)
- Step-6  `` kubectl apply -f ${TCE-REPO-DIR}addons/packages/gatekeeper/test/e2e/fixtures/test-namespace.yaml```
  - After applying the Step-6 file , namespace should be created as the namespace will be created with the owner name.

## Prerequisites

- Before running the suite , the cluster should be up and running.Gatekeeper test suite will never create or destroy any cluster.
- Before running the suite, All depending packages should be installed.

## How to run the tests

- To run individual suits , ```cd addons/packages/gatekeeper/test``` from the tce root path and run ```ginkgo -v -r```
