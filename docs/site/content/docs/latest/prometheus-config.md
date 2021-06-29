# Configuring the Prometheus Package

The [Prometheus](https://prometheus.io/) package provides a monitoring system and time series database.

## Components

- A Prometheus server and corresponding [Alertmanager](https://github.com/prometheus/alertmanager)

## Installation
Run the following command to install the Prometheus package, for more information, see [Packages Introduction](packages-intro.md).

```shell
tanzu package install prometheus.tce.vmware.com
```
## Configuration

The following global configuration values can be set to customize the Prometheus/Alertmanager installation.


| Parameter                                                  | Description                                                                                                          | Type      | Default                 |
|------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------|-----------|-------------------------|
| namespace                                                  | Namespace where Prometheus will be deployed                                                                          | string    | tanzu-system-monitoring |
| prometheus.deployment.replicas                             | Number of Prometheus replicas                                                                                        | integer   | 1                       |
| prometheus.deployment.containers.args                      | Prometheus container arguments                                                                                       | list      |                         |
| prometheus.deployment.containers.resources                 | Prometheus container resource requests and limits                                                                    | map       | {}                      |
| prometheus.deployment.podAnnotations                       | The Prometheus deployments pod annotations                                                                           | map       | {}                      |
| prometheus.deployment.podLabels                            | The Prometheus deployments pod labels                                                                                | map       | {}                      |
| prometheus.deployment.configMapReload.containers.args      | Configmap-reload container arguments.                                                                                | list      |                         |
| prometheus.deployment.configMapReload.containers.resources | Configmap-reload container resource requests and limits                                                              | map       | {}                      |
| prometheus.service.type                                    | Type of service to expose Prometheus. Supported Values: ClusterIP                                                    | string    | ClusterIP               |
| prometheus.service.port                                    | Prometheus service port                                                                                              | integer   | 80                      |
| prometheus.service.targetPort                              | Prometheus service target port                                                                                       | integer   | 9090                    |
| prometheus.service.labels                                  | Prometheus service labels                                                                                            | map       | {}                      |
| prometheus.service.annotations                             | Prometheus service annotations                                                                                       | map       | {}                      |
| prometheus.pvc.annotations                                 | Storage class annotations                                                                                            | map       | {}                      |
| prometheus.pvc.storageClassName                            | Storage class to use for persistent volume claim. By default this is null and default provisioner is used            | string    | null                    |
| prometheus.pvc.accessMode                                  | Define access mode for persistent volume claim. Supported values: ReadWriteOnce, ReadOnlyMany, ReadWriteMany         | string    | ReadWriteOnce           |
| prometheus.pvc.storage                                     | Define storage size for persistent volume claim                                                                      | string    | 150Gi                   |
| prometheus.config.prometheus_yml                           | The [global prometheus configuration](https://www.prometheus.io/docs/prometheus/latest/configuration/configuration/) | yaml file | prometheus.yaml         |
| prometheus.config.alerting_rules_yml                       | The [prometheus alerting rules](https://www.prometheus.io/docs/prometheus/latest/configuration/alerting_rules/)      | yaml file | alerting_rules.yaml     |
| prometheus.config.recording_rules_yml                      | The [prometheus recording rules](https://www.prometheus.io/docs/prometheus/latest/configuration/recording_rules/)    | yaml file | recording_rules.yaml    |
| alertmanager.deployment.replicas                           | Number of alertmanager replicas                                                                                      | integer   | 1                       |
| alertmanager.deployment.containers.resources               | Alertmanager container resource requests and limits                                                                  | map       | {}                      |
| alertmanager.deployment.podAnnotations                     | The Alertmanager deployments pod annotations                                                                         | map       | {}                      |
| alertmanager.deployment.podLabels                          | The Alertmanager deployments pod labels                                                                              | map       | {}                      |
| alertmanager.service.type                                  | Type of service to expose Alertmanager. Supported Values: ClusterIP                                                  | string    | ClusterIP               |
| alertmanager.service.port                                  | Alertmanager service port                                                                                            | integer   | 80                      |
| alertmanager.service.targetPort                            | Alertmanager service target port                                                                                     | integer   | 9093                    |
| alertmanager.service.labels                                | Alertmanager service labels                                                                                          | map       | {}                      |
| alertmanager.service.annotations                           | Alertmanager service annotations                                                                                     | map       | {}                      |
| alertmanager.pvc.annotations                               | Storage class annotations                                                                                            | map       | {}                      |
| alertmanager.pvc.storageClassName                          | Storage class to use for persistent volume claim. By default this is null and default provisioner is used.           | string    | null                    |
| alertmanager.pvc.accessMode                                | Define access mode for persistent volume claim. Supported values: ReadWriteOnce, ReadOnlyMany, ReadWriteMany         | string    | ReadWriteOnce           |
| alertmanager.pvc.storage                                   | Define storage size for persistent volume claim                                                                      | string    | 2Gi                     |
| alertmanager.config.alertmanager_yml                       | The [global yaml configuration for alert manager](https://www.prometheus.io/docs/alerting/latest/configuration/).    | yaml file | alertmanager_yml        |
| kube_state_metrics.deployment.replicas                     | Number of kube-state-metrics replicas                                                                                | integer   | 1                       |
| kube_state_metrics.deployment.containers.resources         | kube-state-metrics container resource requests and limits                                                            | map       | {}                      |
| kube_state_metrics.deployment.podAnnotations               | The kube-state-metrics deployments pod annotations                                                                   | map       | {}                      |
| kube_state_metrics.deployment.podLabels                    | The kube-state-metrics deployments pod labels                                                                        | map       | {}                      |
| kube_state_metrics.service.type                            | Type of service to expose kube-state-metrics. Supported Values: ClusterIP                                            | string    | ClusterIP               |
| kube_state_metrics.service.port                            | kube-state-metrics service port                                                                                      | integer   | 80                      |
| kube_state_metrics.service.targetPort                      | kube-state-metrics service target port                                                                               | integer   | 8080                    |
| kube_state_metrics.service.telemetryPort                   | kube-state-metrics service telemetry port                                                                            | integer   | 81                      |
| kube_state_metrics.service.telemetryTargetPort             | kube-state-metrics service target telemetry  port                                                                    | integer   | 8081                    |
| kube_state_metrics.service.labels                          | kube-state-metrics service labels                                                                                    | map       | {}                      |
| kube_state_metrics.service.annotations                     | kube-state-metrics service annotations                                                                               | map       | {}                      |
| node_exporter.deployment.replicas                          | Number of node-exporter replicas                                                                                     | integer   | 1                       |
| node_exporter.deployment.containers.resources              | node-exporter container resource requests and limits                                                                 | map       | {}                      |
| node_exporter.deployment.podAnnotations                    | The node-exporter deployments pod annotations                                                                        | map       | {}                      |
| node_exporter.deployment.podLabels                         | The node-exporter deployments pod labels                                                                             | map       | {}                      |
| node_exporter.service.type                                 | Type of service to expose node-exporter. Supported Values: ClusterIP                                                 | string    | ClusterIP               |
| node_exporter.service.port                                 | node-exporter service port                                                                                           | integer   | 9100                    |
| node_exporter.service.targetPort                           | node-exporter service target port                                                                                    | integer   | 9100                    |
| node_exporter.service.labels                               | node-exporter service labels                                                                                         | map       | {}                      |
| node_exporter.service.annotations                          | node-exporter service annotations                                                                                    | map       | {}                      |
| pushgateway.deployment.replicas                            | Number of pushgateway replicas                                                                                       | integer   | 1                       |
| pushgateway.deployment.containers.resources                | pushgateway container resource requests and limits                                                                   | map       | {}                      |
| pushgateway.deployment.podAnnotations                      | The pushgateway deployments pod annotations                                                                          | map       | {}                      |
| pushgateway.deployment.podLabels                           | The pushgateway deployments pod labels                                                                               | map       | {}                      |
| pushgateway.service.type                                   | Type of service to expose pushgateway. Supported Values: ClusterIP                                                   | string    | ClusterIP               |
| pushgateway.service.port                                   | pushgateway service port                                                                                             | integer   | 9091                    |
| pushgateway.service.targetPort                             | pushgateway service target port                                                                                      | integer   | 9091                    |
| pushgateway.service.labels                                 | pushgateway service labels                                                                                           | map       | {}                      |
| pushgateway.service.annotations                            | pushgateway service annotations                                                                                      | map       | {}                      |
| cadvisor.deployment.replicas                               | Number of cadvisor replicas                                                                                          | integer   | 1                       |
| cadvisor.deployment.containers.resources                   | cadvisor container resource requests and limits                                                                      | map       | {}                      |
| cadvisor.deployment.podAnnotations                         | The cadvisor deployments pod annotations                                                                             | map       | {}                      |
| cadvisor.deployment.podLabels                              | The cadvisor deployments pod labels                                                                                  | map       | {}                      |
| ingress.enabled                                            | Enable/disable ingress for prometheus and alertmanager                                                               | boolean   | false                   |
| ingress.virtual_host_fqdn                                  | Hostname for accessing promethues and alertmanager                                                                   | string    | prometheus.system.tanzu |
| ingress.prometheus_prefix                                  | Path prefix for prometheus                                                                                           | string    | /                       |
| ingress.alertmanager_prefix                                | Path prefix for alertmanager                                                                                         | string    | /alertmanager/          |
| ingress.tlsCertificate.tls.crt                             | Optional cert for ingress if you want to use your own TLS cert. A self signed cert is generated by default           | string    | Generated cert          |
| ingress.tlsCertificate.tls.key                             | Optional cert private key for ingress if you want to use your own TLS cert.                                          | string    | Generated cert key      |


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
