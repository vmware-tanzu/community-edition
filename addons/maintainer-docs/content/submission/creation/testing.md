---
title: Testing
weight: 3
---

A package should be tested. The 2 most typical types of tests are Unit and End to End.

The `test` directory should include a separate README.md file that documents how to setup and execute those tests. If needed, include a `Makeflie` as well.

## Unit Tests

Unit tests should test basic functionality and validity of the package. An example of this is to test that the ytt templates and overlays process correctly with default values. Instructions how to write unit tests are outside the scope of this guide. The Harbor package provides a good [example](https://github.com/vmware-tanzu/community-edition/tree/main/addons/packages/harbor/2.3.3/test/unittest) of how to test overlays.

1. Create an `expected.yaml` file that contains the result of the upstream content overlaid by your templates and default values.
2. Write a unit test that accumulates all of the yaml files in the config directory.:
   1. upstream/**/*.yaml
   2. overlay/*.yaml
   3. *.yaml
   4. *.star
3. Run the files through the `ytt.RenderYTTTemplate` function in the `addons/packages/test/pkg/ytt` package
4. Verify that no errors were raised
5. Verify the output of the test matches the expected yaml file

## End to End Tests

End-to-end tests are much more thorough and rigorous tests. In these scenarios, the package is actually deployed to a test cluster. Runtime tests and validation can then be performed. End to End tests can be helpful in ensuring that a package is able to install an run successfully on different versions of Kubernetes and with different providers. For example, if a package might be configured differently based on deployment to a cluster on AWS versus Azure. End to End tests for both AWS and Azure would be useful to show proper functionality.

The [Harbor](https://github.com/vmware-tanzu/community-edition/tree/main/addons/packages/harbor/2.3.3/test/e2e) and [ExternalDNS](https://github.com/vmware-tanzu/community-edition/tree/main/addons/packages/external-dns/0.10.0/test/e2e) packages both provide good examples of end-to-end tests.
