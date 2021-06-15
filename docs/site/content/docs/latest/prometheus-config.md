# Configuring the Prometheus Package

The [Prometheus](https://prometheus.io/) package provides a monitoring system and time series database.

## Components

- A Prometheus server and corresponding [Alertmanager](https://github.com/prometheus/alertmanager)

## Configuration

The following configuration values can be set to customize the Prometheus/Alertmanager installation.

### Global

| Value | Required/Optional | Description |
|:-------|:-------------------|:-------------|
| `namespace` | Required | The namespace in which to deploy Prometheus|
| `prometheus.deployment.replicas` | Required | The number of Prometheus replicas |
| `prometheus.config.prometheus_yml` | Optional | The [global Prometheus configuration](https://www.prometheus.io/docs/prometheus/latest/configuration/configuration/) |
| `prometheus.config.alerting_rules_yml` | Optional | The [Prometheus alerting rules](https://www.prometheus.io/docs/prometheus/latest/configuration/alerting_rules/) |
| `prometheus.config.recording_rules_yml` | Optional | The [Prometheus recording rules](https://www.prometheus.io/docs/prometheus/latest/configuration/recording_rules/) |
| `prometheus.config.alerting_yml` | Optional | Additional [Prometheus alerts can be configured here](https://www.prometheus.io/docs/prometheus/latest/configuration/alerting_rules/) |
| `prometheus.config.rules_yml` | Optional | Additional [Prometheus rules](https://www.prometheus.io/docs/prometheus/latest/configuration/recording_rules/) |
| `alertmanager.deployment.replicas` | Required | The number of Alertmanager replicas |
| `alertmanager.config.alertmanager_yml` | Required | The [global YAML configuration for Alertmanager](https://www.prometheus.io/docs/alerting/latest/configuration/) |

### Config files

The configuration files can be loaded from `/etc/config/`.
For example, when loading `rule_files`:

```text
rule_files:
- /etc/config/alerting_rules.yml
- /etc/config/recording_rules.yml
- /etc/config/alerts
- /etc/config/rules
```

### Alertmanager service

The Alertmanager can be targeted by the deployed Prometheus through it's service.
For example:

```text
targets:
- alertmanager.prometheus-addon.svc:9093
```

`alertmanager` is the default service name and `prometheus-addon` is the namespace of the alertmanager deployment.

## Usage Example

The default `prometheus.yml` configuration will deploy a Prometheus server that will scrape metrics from pods that emit metrics on an endpoint and has the following annotations:

```yaml
prometheus.io/scrape: 'true'
prometheus.io/path: '/metrics'
prometheus.io/port: '8080'
```
