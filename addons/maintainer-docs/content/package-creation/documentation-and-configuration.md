---
title: Documentation and Configuration
---

## Documentation

This should include a brief overview of the software components contained in the package, a description of configuration parameters, and general usage information. This documentation is not intended to replace, or be as extensive as the official documentation for the software.

The package documentation should highlight dependencies or considerations on other packages, software, Kubernetes distributions, or underlying infrastructure (e.g. AWS, GCP, Docker, vSphere, etc).

## Configuration Parameters

The package maintainer is responsible for creating a [schema](https://carvel.dev/ytt/docs/latest/how-to-write-schema/) that details the available configuration parameters. This schema can be used to generate package documentation and a default values.yaml file for use when installing the package.

An example of a schema...
