# Velero e2e tests

## Description

The e2e tests validate end-to-end behavior of Velero backup and restore operations. For full documentation, see: [velero/README.md e2e tests at v1.9.2 · vmware-tanzu/velero](https://github.com/vmware-tanzu/velero/blob/v1.9.2/test/e2e/README.md).

The `make e2e-test` target defined in this directory runs the [Velero e2e tests hosted on the Velero GitHub repository](https://github.com/vmware-tanzu/velero/tree/v1.9.2/test/e2e).

For testing this TCE package, the e2e tests will optout of installing Velero, since we want to test the already installed TCE Velero package. All of the configurations will reflect this condition.

## Prerequisites

Tests can be run in a TCE cluster hosted in any of the cloud providers supported by the Velero package.

1. A running TCE cluster
1. The necessary storage drivers/provisioners installed.
1. `kubectl` installed locally.
1. The `velero.community.tanzu.vmware.com` v1.9.2 package must be installed on the cluster.
1. The local context needs to be set to the TCE cluster.

⚠️ Please be aware of the limitations for the v1.9.2 of the tests: [velero/test/e2e at v1.9.2 · vmware-tanzu/velero](https://github.com/vmware-tanzu/velero/tree/v1.9.2/test/e2e#limitations).

Note: the tests do not delete any backup/restore resource created during the running of the tests. Feel free to delete them manually, but otherwise that should not interfere with running subsequent tests.

## Test Configuration

1. `GOPATH` - environment variable if it is different from `~/go`. Default is "~/go".
1. `CLOUD_PROVIDER`-  because we are not asking the Velero e2e tests to install Velero, the `CLOUD_PROVIDER` variable is only meant to indicate if the tests are being run in an environment that supports taking volume snapshots, like AWS, Azure or vSphere, or not. If the tests are running on Docker, set this variable to "kind". Otherwise, set it to anything else. Default is "notkind".
1. `REGISTRY_CREDENTIAL_FILE` - only needed for tests that trigger the creation of a workload, usually so there can be a snapshot taken. Required in some cases. See format below:

```sh
{
    "auths": {
        "https://index.docker.io/v1/": {
            "Username": <dockerusername>,
            "Secret": <dockerpwd>
        }
    },
    "credsStore": "desktop"
}
```

## Run Tests

Run the tests from the test directory: `addons/packages/velero/1.9.2/test/` .

🚨 Note: for any test to run successfully, there has to be a BackupStorageLocation that is the default and and that has the status phase of `available` . To verify:

`kubectl get backupstoragelocations.velero.io -n velero` .

⚠️ In the absence of a "focus" flag, all Velero e2e tests will be run, which is not recommended for TCE.

It is recommended that these two tests be run:

* `GINKGO_FOCUS='Basic'` - This is the default. It tests the simple backup and restore of 2 namespaces :

```bash
GITHUB_TOKEN=$GITHUB_TOKEN make e2e-test
```

* `GINKGO_FOCUS='Snapshot'` - This test will create a workload

```bash
GINKGO_FOCUS='Snapshot' GITHUB_TOKEN=$GITHUB_TOKEN REGISTRY_CREDENTIAL_FILE=/Users/carlisiac/.docker/config.json make e2e-test
```
