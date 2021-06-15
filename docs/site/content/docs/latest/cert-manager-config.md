# Configuring the cert-manager Package

This package provides certificate management functionality using [cert-manager](https://cert-manager.io/docs/).

## Components

* cert-manager

## Installation
Run the following command to install the cert-manager package, for more information, see [Packages Introduction](packages-intro.md):

```shell
tanzu package install cert-manager.tce.vmware.com
```

## Configuration

The following configuration values can be set to customize the cert-manager installation.

### Global

| Value | Required/Optional | Description |
|:-------|:-------------------|:-------------|
| `namespace` | Optional | The namespace in which to deploy cert-manager. |

<!--## Usage Example

This walkthrough guides you through using cert-manager...-->
