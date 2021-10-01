# E2E Test

## How to run E2E test

First go to diagnostics plugin directory

```bash
# if you are in TCE repo root, then
$ cd cli/cmd/plugin/diagnostics
```

Then simply run

```bash
make e2e-test
```

Or you can directly invoke the E2E test script `e2e-test.sh`, present in this directory, from anywhere

## Internals

For testing `tanzu diagnostics collect` command, we use a `kind` cluster to quickly spin up a cluster and check if `tanzu diagnostics collect` can collect logs from the cluster. We use a single `kind` cluster and treat it as a bootstrap cluster, management cluster and a workload cluster. This is done just to keep the test simple and light weight and fast

What we test as part of the E2E test currently - Does running the `tanzu diagnostics collect` command collect the diagnostics data from all the different kinds of clusters and create tar balls in the appropriate directory? We do this check at a very high level - check if the tar balls exist

What we can test in the future

- Check some of the contents of the tar ball
  - Check if some files are present and have file size > 0
  - Check if some logs are present in some log files by running some test apps in namespaces that `tanzu diagnostics collect` collects data from, for example `capi-system`
- Check if the different flags of the `tanzu diagnostics collect` are working as expected
