# Grafana Addon

Grafana is open source visualization and analytics software. It allows you to query, visualize, alert on, and explore your metrics no matter where they are stored. In plain English, it provides you with tools to turn your time-series database (TSDB) data into beautiful graphs and visualizations.

## Components

- Grafana server.

## Configuration

The following configuration values can be set to customize the grafana installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Required | The namespace in which to deploy prometheus. |
| `grafana.deployment.replicas` | Required | The number of grafana replicas. |
| `grafana.deployment.image` | Required | The grafana image to deploy. |
| `grafana.config.grafana_ini` | Optional | The (grafana configuration.)[https://github.com/grafana/grafana/blob/master/conf/defaults.ini] |


## Usage Example

- Set up data sources for your metrics
you can add one or more data sources to Grafana to start. See the Grafana (documentation)[https://grafana.com/docs/grafana/latest/datasources/add-a-data-source/] for detailed description of how to add a data source.
- Create Dashboards
There are many prebuilt Grafana dashboard templates available for various data sources. You can check out the templates (here)[https://grafana.com/grafana/dashboards]
- Enable Ingress on Grafana as per your requirement.


