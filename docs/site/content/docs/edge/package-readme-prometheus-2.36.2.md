# Prometheus

The package provides the ability to monitor and alert your kubernetes cluster using an open-source system named [Prometheus](https://github.com/prometheus/prometheus).

## Installation

### Installation of package

Install TCE Prometheus package through tanzu command:

```bash
tanzu package install prometheus --package-name prometheus.community.tanzu.vmware.com --version ${PROMETHEUS_PACKAGE_VERSION}
```

> You can get the `${PROMETHEUS_PACKAGE_VERSION}` from running `tanzu package
> available list prometheus.community.tanzu.vmware.com`. Specifying a
> namespace may be required depending on where your package repository was
> installed.

## Options

The following configuration values can be set to customize the Prometheus installation.

### Package configuration values

#### Global

| Value | Required/Optional | Default | Description |
|-------|-------------------|---------|-------------|
| `namespace` | Optional | `prometheus` | The namespace in which to deploy Prometheus  components |

#### Prometheus Configuration

> Note: Ingress for Prometheus server is not available by default, and can be activated using the `ingress.enabled` configuration field.
>
> If you choose to activate the Contour-based Ingress, `Contour` must also be installed on the target cluster. Additionally, enabling the Ingress requires either `Cert Manager` or your own user-provided TLS certificate (`ingress.tlsCertificate.tls.crt` and `ingress.tlsCertificate.tls.key`) to configure TLS settings for the Ingress. For ad-hoc Prometheus UI access without an Ingress, use `kubectl port-forward`.

The following configuration values can be set to customize the Prometheus/Alertmanager installation.

| Value | Required/Optional | Default | Description  |
|:------------------------------------------------------------|:----------------------------------------------------------------------------------------------------------------------|:-----------|:-------------------------|
| `prometheus.deployment.replicas`                             | true                                                                                        | 1   | Number of Prometheus replicas                      |
| `prometheus.deployment.updateStrategy`                  | false                                                      | Recreate      | Type of prometheus upgrade strategy. Supported Values: Recreate, RollingUpdate                    |
| `prometheus.deployment.rollingUpdate.maxUnavailable`    | false                        | null     | Number of maxUnavailable pods when prometheus is upgrading, It only work when updateStrategy is RollingUpdate                       |
| `prometheus.deployment.rollingUpdate.maxSurge`          |  false                      |null     | Number of maxSurge pods when prometheus is upgrading, It only work when updateStrategy is RollingUpdate                        |
| `prometheus.deployment.containers.args`                      | true                                                                                       | `--storage.tsdb.retention.time=15d --config.file=/etc/config/prometheus.yml --storage.tsdb.path=/data --web.console.libraries=/etc/prometheus/console_libraries --web.console.templates=/etc/prometheus/consoles --web.enable-lifecycle`     |   Prometheus container arguments                      |
| `prometheus.deployment.containers.resources`                 | false                                                                   | {}       | Prometheus container resource requests and limits                       |
| `prometheus.deployment.podAnnotations`                       |  false                                                                          | {}       |The Prometheus deployments pod annotations                     |
| `prometheus.deployment.podLabels`                            |  false                                                                              | {}       | The Prometheus deployments pod labels                      |
| `prometheus.deployment.configMapReload.containers.args`      |  true                                                                               | `\--volume-dir=/etc/config  \--webhook-url=<http://127.0.0.1:9090/-/reload>`      |  Configmap-reload container arguments.                       |
| `prometheus.deployment.configMapReload.containers.resources` |false                                                               | {}       | Configmap-reload container resource requests and limits                     |
| `prometheus.service.type`                                    | true        | ClusterIP    |  Type of service to expose Prometheus. Supported Values: ClusterIP              |
| `prometheus.service.port`                                    | true                                                                                              | 80   |  Prometheus service port                     |
| `prometheus.service.targetPort`                              |true                                                                                       | 9090   |   Prometheus service target port                   |
| `prometheus.service.labels`                                  |  false                                                                                           | {}       | Prometheus service labels                     |
| `prometheus.service.annotations`                             | false                                                                                      | {}       | Prometheus service annotations                       |
| `prometheus.pvc.annotations`                                 |  false                                                                                        | {}       | Storage class annotations                        |
| `prometheus.pvc.storageClassName`                            |false           | null    | Storage class to use for persistent volume claim. By default this is null and default provisioner is used                      |
| `prometheus.pvc.accessMode`                                  |true         | ReadWriteOnce    | Define access mode for persistent volume claim. Supported values: ReadWriteOnce, ReadOnlyMany, ReadWriteMany            |
| `prometheus.pvc.storage`                                     |   true          | 150Gi    | 150Gi  Define storage size for persistent volume claim                 |
| `prometheus.config.prometheus_yml`                           |true  | prometheus.yaml | The [global Prometheus configuration](https://www.prometheus.io/docs/prometheus/latest/configuration/configuration/)          |
| `prometheus.config.alerting_rules_yml`                       |false       | alerting_rules.yaml |  The [Prometheus alerting rules](https://www.prometheus.io/docs/prometheus/latest/configuration/alerting_rules/)    |
| `prometheus.config.recording_rules_yml`                      |false     | recording_rules.yaml  |   The [prometheus recording rules](https://www.prometheus.io/docs/prometheus/latest/configuration/recording_rules/) |
| `prometheus.config.alerts_yml`                               | false  | alerts_yml.yaml   |   Additional Prometheus alerting rules can be configured here.     |
| `prometheus.config.rules_yml`                                | false  |rules_yml.yaml  |  Additional Prometheus alerting rules can be configured here.         |
| `alertmanager.deployment.replicas`                           | true                                                                                     | 1   |  Number of alertmanager replicas       |
| `alertmanager.deployment.updateStrategy`                  |  false                                                     | Recreate      |   Type of prometheus upgrade strategy. Supported Values: RollingUpdate, Recreate                  |
| `alertmanager.deployment.rollingUpdate.maxUnavailable`    | false                       | null     | Number of maxUnavailable pods when alertmanager is upgrading, It only work when updateStrategy is RollingUpdate                        |
| `alertmanager.deployment.rollingUpdate.maxSurge`          | false                        | null     | Number of maxSurge pods when alertmanager is upgrading, It only work when updateStrategy is RollingUpdate                       |
| `alertmanager.deployment.containers.resources`               |  false                                      | {}       | Alertmanager container resource requests and limits                     |
| `alertmanager.deployment.podAnnotations`                     |    false                                     | {}       | The Alertmanager deployments pod annotations                      |
| `alertmanager.deployment.podLabels`                          |     false                 | {}       | The Alertmanager deployments pod labels                      |
| `alertmanager.service.type`                                  |       true                                       | ClusterIP    |   Type of service to expose Alertmanager. Supported Values: ClusterIP              |
| `alertmanager.service.port`                                  |   true                                                                                        | 80   |   Alertmanager service port                      |
| `alertmanager.service.targetPort`                            |  true                                                                                   | 9093   | Alertmanager service target port                     |
| `alertmanager.service.labels`                               |   false                                                                                        | {}       | Alertmanager service labels                      |
| `alertmanager.service.annotations`                           | false                                                                                    | {}       | Alertmanager service annotations                       |
| `alertmanager.pvc.annotations`                               |    false                                                                                         | {}       | Storage class annotations                     |
| `alertmanager.pvc.storageClassName`                          | false          | null    | Storage class to use for persistent volume claim. By default this is null and default provisioner is used.                     |
| `alertmanager.pvc.accessMode`                                | false        | ReadWriteOnce    | Define access mode for persistent volume claim. Supported values: ReadWriteOnce, ReadOnlyMany, ReadWriteMany            |
| `alertmanager.pvc.storage`                                   | true                                                                     | 2Gi    |     Define storage size for persistent volume claim                  |
| `alertmanager.config.alertmanager_yml`                       | false   | alertmanager_yml |   The [global yaml configuration for alert manager](https://www.prometheus.io/docs/alerting/latest/configuration/).       |
| `kube_state_metrics.deployment.replicas`                     | Number of kube-state-metrics replicas                                                                                | integer   | 1                       |
| `kube_state_metrics.deployment.containers.resources`         |     false                              | {}       |     kube-state-metrics container resource requests and limits                  |
| `kube_state_metrics.deployment.podAnnotations`               |      false                               | {}       |  The kube-state-metrics deployments pod annotations                     |
| `kube_state_metrics.deployment.podLabels`                    |     false                           | {}       | The kube-state-metrics deployments pod labels                       |
| `kube_state_metrics.service.type`                            |    true                                        | ClusterIP    |   Type of service to expose kube-state-metrics. Supported Values: ClusterIP              |
| `kube_state_metrics.service.port`                            |    true                                                                                  | 80   |    kube-state-metrics service port                    |
| `kube_state_metrics.service.targetPort`                      | true                                                                              | 8080   |      kube-state-metrics service target port                |
| `kube_state_metrics.service.telemetryPort`                   |  true                                                                          | 81   | kube-state-metrics service telemetry port                       |
| `kube_state_metrics.service.telemetryTargetPort`            |  true                                                                  | 8081   | kube-state-metrics service target telemetry  port                     |
| `kube_state_metrics.service.labels`                          |   false                                  | {}       | kube-state-metrics service labels                     |
| ``kube_state_metrics.service.annotations``                     | false                                                                             | {}       |  kube-state-metrics service annotations                      |
| `node_exporter.daemonset.containers.resources`               |   false                               | {}       | node-exporter container resource requests and limits                       |
| `node_exporter.daemonset.hostNetwork`                        |  false                                                                           | false   |  Host networking requested for this pod                     |
| `node_exporter.daemonset.podAnnotations`                     |    false                     | {}       | The node-exporter deployments pod annotations                     |
| `node_exporter.daemonset.podLabels`                          | false                                                                             | {}       | The node-exporter deployments pod labels                     |
| `node_exporter.service.type`                                 |     true                                            | ClusterIP    |   Type of service to expose node-exporter. Supported Values: ClusterIP              |
| `node_exporter.service.port`                                 |   true                                                                                        | 9100   | node-exporter service port                     |
| `node_exporter.service.targetPort`                           |  true                                                                                  | 9100   | node-exporter service target port                     |
| `node_exporter.service.labels`                               |   false                                                                                      | {}       | node-exporter service labels                       |
| `node_exporter.service.annotations`                          |                                                                                    | {}       | node-exporter service annotations                      |
| `pushgateway.deployment.replicas`                            |   true                                                                                   | 1   | Number of pushgateway replicas                         |
| `pushgateway.deployment.containers.resources`                |  false                                                                 | {}       | pushgateway container resource requests and limits                       |
| `pushgateway.deployment.podAnnotations`                      |  false                                                                        | {}       | The pushgateway deployments pod annotations                       |
| `pushgateway.deployment.podLabels`                           |   false                            | {}       |  The pushgateway deployments pod labels                       |
| `pushgateway.service.type`                                  |true                                                    | ClusterIP    | Type of service to expose pushgateway. Supported Values: ClusterIP               |
| `pushgateway.service.port`                                   |    true                                                                                       | 9091   | pushgateway service port                       |
| `pushgateway.service.targetPort`                             |   true                                                                                  | 9091   | pushgateway service target port                      |
| `pushgateway.service.labels`                                | false                                                  | {}       | pushgateway service labels                       |
| `pushgateway.service.annotations`                            |    false                                                                                   | {}       | pushgateway service annotations                      |
| `ingress.enabled`                                            | false                                                               | false   | Activate/deactivate ingress for prometheus and alertmanager                   |
| `ingress.virtual_host_fqdn`                                  |    false                                                               | prometheus.system.tanzu    |  Hostname for accessing prometheus and alertmanager |
| `ingress.prometheus_prefix`                                  |      false                         | /    |      Path prefix for prometheus                   |
| `ingress.alertmanager_prefix`                               |    false                          | /alertmanager/     |    Path prefix for alertmanager       |
| `ingress.prometheusServicePort`                              |   false            | 80   | Prometheus service port to proxy traffic to                       |
| `ingress.alertmanagerServicePort`                            |     false             | 80   | Alertmanager service port to proxy traffic to                       |
| `ingress.tlsCertificate.tls.crt`                             |   false        | Generated cert     |    Optional cert for ingress if you want to use your own TLS cert. A self signed cert is generated by default       |
| `ingress.tlsCertificate.tls.key`                             |   false                                       | Generated cert key    |   Optional cert private key for ingress if you want to use your own TLS cert.     |
| `ingress.tlsCertificate.ca.crt`                              |    false                                                                                        |  CA certificate    |     Optional CA certificate        |

### Application configuration values

The various configuration files can be loaded from `/etc/config/`.
For example, when loading `rule_files`:

```text
rule_files:
- /etc/config/alerting_rules.yml
- /etc/config/recording_rules.yml
- /etc/config/alerts
- /etc/config/rules
```

#### Multi-cloud configuration steps

There are currently no configuration steps necessary for installation of the Prometheus package to any provider.

## What This Package Does

Prometheus scrapes metrics from instrumented jobs, either directly or via an intermediary push gateway for short-lived jobs. It stores all scraped samples locally and runs rules over this data to either aggregate and record new time series from existing data or generate alerts. Grafana or other API consumers can be used to visualize the collected data.

## Components

* Prometheus-server
* Alert-manager
* Kube-state-metrics
* Prometheus-pushgateway
* Node-exporter

### Supported Providers

The following table shows the providers this package can work with.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
| ✅  |  ✅  | ✅  | ✅ |

## Files

Here is an example of the package configuration file [values.yaml](bundle/config/values.yaml).

## Package Limitations

* Multi-replicas in Prometheus-server and multi-instances server are not support.
* Resources like Service Monitor in prometheus-operator are not support.

## Usage Example

The default `prometheus.yml` configuration will deploy a prometheus server
that will scrape metrics from pods that emit metrics
on an endpoint and has the following annotations:

```yaml
prometheus.io/scrape: 'true'
prometheus.io/path: '/metrics'
prometheus.io/port: '8080'
```

For details, see the [Prometheus documentation](https://prometheus.io/docs/prometheus/latest/configuration/configuration/) for more information.

## Troubleshooting

Not applicable.

## Additional Documentation

See the [Prometheus documentation](https://prometheus.io/docs/prometheus/latest/configuration/configuration/) for more information.
