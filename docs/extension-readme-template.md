# EXT_NAME Extension

This extension provides disaster recovery capabilities using [velero](https://velero.io/). At the moment, it leverages [minio](https://github.com/minio/minio) for object storage.

## Components

## Configuration

The following configuration values can be set to customize the EXT_NAME installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Optional | The namespace in which to deploy EXT_NAME. |

### EXT_NAME Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `provider` | Required | The cloud provider in use. One of: `aws`, `azure`, `vsphere`, `docker`. |

## Usage Example

The follow is a basic guide for getting started with EXT_NAME.
