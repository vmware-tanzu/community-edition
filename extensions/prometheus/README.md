# Prometheus Extension

A time series database for your metrics.

## Components

## Configuration

The following configuration values can be set to customize the prometheus installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Optional | The namespace in which to deploy prometheus. |

### Prometheus Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `provider` | Required | The cloud provider in use. One of: `aws`, `azure`, `vsphere`, `docker`. |

## Usage Example

The follow is a basic guide for getting started with Prometheus.

