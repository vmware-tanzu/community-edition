# Prometheus Addon

A time series database for your metrics.

## Components

- A prometheus server

## Configuration

The following configuration values can be set to customize the prometheus installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Optional | The namespace in which to deploy prometheus. |
| `replicas` | Optional | The number of prometheus replicas. |

## Usage Example

The prometheus server will scrape metrics from pods that emit metrics
on an endpoint and has the following annotations:

```yaml
prometheus.io/scrape: 'true'
prometheus.io/path: '/metrics'
prometheus.io/port: '8080'
```
