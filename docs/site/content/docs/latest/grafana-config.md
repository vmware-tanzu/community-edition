# Configuring the Grafana Package

Grafana is open source visualization and analytics software. It allows you to query, visualize, alert on, and explore your metrics no matter where they are stored. It provides you with tools to turn your time-series database (TSDB) data into beautiful graphs and visualizations.

## Components

- Grafana server

## Installation
Run the following command to install the Grafana package, for more information, see [Packages Introduction](packages-intro.md).

```shell
tanzu package install grafana.tce.vmware.com
```
## Configuration

The following global configuration values can be set to customize the Grafana installation.

| Value | Required/Optional | Description |
|:-------|:-------------------|:-------------|
| `namespace` | Required | The namespace in which to deploy prometheus. |
| `grafana.deployment.replicas` | Required | The number of Grafana replicas. |
| `grafana.deployment.image` | Required | The Grafana image to deploy. |
| `grafana.config.grafana_ini` | Optional | The [Grafana configuration](https://github.com/grafana/grafana/blob/master/conf/defaults.ini). |

## Usage Example

1. Add data sources for your metrics, you can add multiple data sources. For more information, see the Add a data source topic in the [Grafana documentation](https://grafana.com/docs/grafana/latest/datasources/add-a-data-source/).
2. Create Dashboards. There are many prebuilt Grafana dashboard templates available for various data sources. For more information, see the Dashboard Overview topic in the [Grafana documentation](https://grafana.com/grafana/dashboards).
3. Enable Ingress on Grafana as per your requirements.
