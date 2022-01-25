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

### CVE Remediation

Security concerns and CVE's can occur within a package in 2 different manners. The first is in the underlying, packaged software. The other is in the package itself.

If a concern is present in the underlying software, the package maintainer is expected to update the package with patched or updated versions of the effected software when available. The package maintainer is not expected to patch the underlying software.

If the concern is in the package itself, the package maintainer is expected to provide a patch as soon as possible. For example, if an overlay introduces a vulnerability, it is expected that the package maintainer fix the overlay to remedy the concern.

### Exposing Schemas

The package maintainer is responsible for creating a [schema](https://carvel.dev/ytt/docs/latest/how-to-write-schema/) that details the available configuration parameters. This schema can be used to generate package documentation and a default values.yaml file for use when installing the package.

### Testing (and Coverage)

The packages should some form of basic validation or tests. Tests can include unit or end-to-end. Unit tests should test basic functionality and validity of the package. An example of this is to test that the ytt templates and overlays process correctly with default values. End-to-end tests can include can verify basic functionality, like a successful installation, or demonstrate more advanced features of the packages on a running cluster.

### Publishing Container Images

### Publishing Package Images

### Documentation

Providing updated documentation on package usage and configuration

Maintainers should provide documentation on the usage of the package, including details on installation and configuration.

### Supporting Package End-Users

Community support is expected of Package Maintainers. Maintainers shall monitor and respond in appropriate Slack workspaces or other online messaging channels relevant to the usage of the package. Issues/Bugs/Enhancements that are filed against the package shall be addressed and responded to by maintainers.
