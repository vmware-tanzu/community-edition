---
title: Documentation
weight: 3
---

## Documentation

This should include a brief overview of the software components contained in the package, a description of configuration parameters, and general usage information. This documentation is not intended to replace, or be as extensive as the official documentation for the software.

The package documentation should highlight dependencies or considerations on other packages, software, Kubernetes distributions, or underlying infrastructure (e.g. AWS, GCP, Docker, vSphere, etc).

## Sample README

```text
# PACKAGE_NAME Package

This package provides << awesome functionality >> using [PACKAGE_NAME](https://INFO_NEEDED).

## Supported Providers

The following tables shows the providers this package can work with.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
| ✅  |  ✅  | ✅  | ❌  |

## Components

* PACKAGE_NAME version: `1.2.3`

## Configuration

The following configuration values can be set to customize the PACKAGE_NAME package installation.

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `foo` | Required | foo description |
| `bar` | Optional | bar description |

## Usage Example

The follow is a basic guide for getting started with PACKAGE_NAME.
```
