# MULTUS-CNI E2E TESTS

## Prerequisites

* A Tanzu Community Edition cluster and the cluster needs to be the
  current-context. See the [Getting Started
  Guide](https://tanzucommunityedition.io/docs/getting-started/) for
  instuctions on how to create a cluster.
* The `multus-cni.community.tanzu.vmware.com` Package must exist on the
  cluster so it can be installed by the test.

## USAGE

1. Run below command will run the e2e tests.

```bash
make e2e-test
```
