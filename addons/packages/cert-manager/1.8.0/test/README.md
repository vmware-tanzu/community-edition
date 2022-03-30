# Tests

## End-to-End Tests

The end-to-end tests assume that a cluster is running and that the correct version of the package under test is available in a package repository that is already installed to the cluster. The test will install the package and then create an Issuer and a Self Signed certificate following the steps to [verify an installation](https://cert-manager.io/v1.6-docs/installation/verify/) from the cert-manager documentation.

To execute the end-to-end tests, run

```shell
make e2e-test
```

## Unit Tests

To execute the unit tests, run

```shell
make test
```
