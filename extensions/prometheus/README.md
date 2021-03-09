# Prometheus Addon

A time series database for your metrics.

## Components

- A prometheus server and corrisponding alert manager

## Configuration

The following configuration values can be set to customize the prometheus / alert manager installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Required | The namespace in which to deploy prometheus. |
| `prometheus.deployment.replicas` | Required | The number of prometheus replicas. |
| `prometheus.deployment.prometheus_yml` | Optional | The (global prometheus configuration)[https://www.prometheus.io/docs/prometheus/latest/configuration/configuration/] |
| `prometheus.deployment.alerting_rules_yml` | Optional | The (prometheus alerting rules)[https://www.prometheus.io/docs/prometheus/latest/configuration/alerting_rules/] |
| `prometheus.deployment.recording_rules_yml` | Optional | The (prometheus recording rules)[https://www.prometheus.io/docs/prometheus/latest/configuration/recording_rules/] |
| `prometheus.deployment.alerting_yml` | Optional | Additional (prometheus alerts can be configured here)[https://www.prometheus.io/docs/prometheus/latest/configuration/alerting_rules/] |
| `prometheus.deployment.rules_yml` | Optional | Additional (prometheus rules)[https://www.prometheus.io/docs/prometheus/latest/configuration/recording_rules/] |
| `alertmanager.deployment.replicas` | Required | The number of alertmanager replicas. |
| `alertmanager.deployment.alertmanager_yml` | Required | The (global yaml configuration for alert manager)[https://www.prometheus.io/docs/alerting/latest/configuration/] |

### Config files

The various configuration files can be loaded from `/etc/config/`.
For example, when loading `rule_files`:
```
rule_files:
- /etc/config/alerting_rules.yml
- /etc/config/recording_rules.yml
- /etc/config/alerts
- /etc/config/rules
```

### Alert manager service
The alert manager can be targeted by the deployed prometheus through it's service.
For example:
```
- targets:
    - alertmanager.prometheus-addon.svc:9093
```
`alertmanager` is the default service name and `prometheus-addon` is the namespace of the alertmanager deployment.

## Usage Example

The default `prometheus.yml` configuration will deploy a prometheus server
that will scrape metrics from pods that emit metrics
on an endpoint and has the following annotations:

```yaml
prometheus.io/scrape: 'true'
prometheus.io/path: '/metrics'
prometheus.io/port: '8080'
```
