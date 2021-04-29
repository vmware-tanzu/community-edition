# Node exporter

Prometheus style exporter for hardware and OS metrics

## Components

- The node-exporter daemonset

## Configuration

The following configuration values can be set to customize node-exporter

### Global

| Value                                   | Required/Optional | Description                                                                                                                           |
|-----------------------------------------|-------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| `namespace`                             | Required          | The namespace in which to deploy node-exporter                                                                                        |
| `daemonset.updatestrategy`              | Required          | The update strategy for the daemonset on new images being available                                                                   |
| `hostNetwork`                           | Required          | true or false. Wether to use the host machine network or not                                                                          |


