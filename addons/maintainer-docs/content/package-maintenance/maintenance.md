---
title: Maintenance
weight: 1
---

### Metadata

The package maintainer will keep relevant metadata up to date.

### Software Updates

Creating new packages to track the latest version of the underlying software

### Package Updates

Providing package updates containing bug fixes and/or enhancements to the package itself

Creating new packages to track the latest version of the underlying software

### Exposing Schemas

The package maintainer is responsible for creating a [schema](https://carvel.dev/ytt/docs/latest/how-to-write-schema/) that details the available configuration parameters. This schema can be used to generate package documentation and a default values.yaml file for use when installing the package.

### Testing (and Coverage)

The packages should some form of basic validation or tests. Tests can include unit or end-to-end. Unit tests should test basic functionality and validity of the package. An example of this is to test that the ytt templates and overlays process correctly with default values. End-to-end tests can include can verify basic functionality, like a successful installation, or demonstrate more advanced features of the packages on a running cluster.

### Publishing Container Images

### Publishing Package Images

### Documentation

Providing updated documentation on package usage and configuration

Maintainers should provide documentation on the usage of the package, including details on installation and configuration.
