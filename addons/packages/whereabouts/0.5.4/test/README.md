# WHEREABOUTS E2E TESTS

## Prerequisites

* A Tanzu Community Edition cluster and the cluster needs to be the
  current-context. See the [Getting Started
  Guide](https://tanzucommunityedition.io/docs/getting-started/) for
  instuctions on how to create a cluster.
* The packages `multus-cni.community.tanzu.vmware.com`  with version `3.8.0`
  and `whereabouts.community.tanzu.vmware.com` with version `0.5.4` must
  exist on the cluster so they can be installed by the test.

## USAGE

1. Run below command to run the e2e tests.

```bash
make e2e-test
```
