---
title: Testing
---

A package should be tested. The 2 most typical types of tests are Unit and End to End.

## Unit Tests

Unit tests should test basic functionality and validity of the package. An example of this is to test that the ytt templates and overlays process correctly with default values.

## End to End Tests

End-to-end tests are much more thorough and rigorous tests. In these scenarios, the package is actually deployed to a test cluster. Runtime tests and validation can then be performed. End to End tests can be helpful in ensuring that a package is able to install an run successfully on different versions of Kubernetes and with different providers. For example, if a package might be configured differently based on deployment to a cluster on AWS versus Azure. End to End tests for both AWS and Azure would be useful to show proper functionality.
