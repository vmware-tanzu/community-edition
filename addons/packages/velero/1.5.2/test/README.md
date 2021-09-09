# Velero tests

## End-to-End tests

End-to-End tests for Velero are located in [vmware-tanzu/velero](https://github.com/vmware-tanzu/velero/tree/main/test/e2e).
The `make e2e-test` target defined in this directory runs the `Basic` end-to-end test from velero test suite.

### Prerequisites

To run the velero end-to-end tests you need:

* A TCE cluster and the cluster needs to be the current-context.

* `velero.community.tanzu.vmware.com` package must be installed on the cluster.

### Test Configuration

Set the `GITHUB_TOKEN` environment variable which is needed to install velero cli while running the tests.
Set the `GOPATH` environment variable if it is different from `~/go`.
Set the `PROVIDER` environment variable if the cloud provider is not minio.

### Run Tests

Run the tests from the test directory:

```bash
cd addons/packages/velero/1.5.2/test/
make e2e-test
```
