# Prometheus

A time series database for your metrics.

## Components

- A Prometheus server and corresponding alert manager

## Configuration

> Note: Ingress for Prometheus server is not available by default, and can be activated using the `ingress.enabled` configuration field.
>
> If you choose to activate the Contour-based Ingress, `Contour` must also be installed on the target cluster. Additionally, enabling the Ingress requires either `Cert Manager` or your own user-provided TLS certificate (`ingress.tlsCertificate.tls.crt` and `ingress.tlsCertificate.tls.key`) to configure TLS settings for the Ingress. For ad-hoc Prometheus UI access without an Ingress, use `kubectl port-forward`.

The following configuration values can be set to customize the Prometheus/Alertmanager installation.

| Parameter                                                  | Description                                                                                                          | Type      | Default                 |
|:------------------------------------------------------------|:----------------------------------------------------------------------------------------------------------------------|:-----------|:-------------------------|
| namespace                                                  | Namespace where Prometheus will be deployed                                                                          | string    | prometheus              |
| `prometheus.deployment.replicas`                             | Number of Prometheus replicas                                                                                        | integer   | 1                       |
| `prometheus.deployment.containers.args`                      | Prometheus container arguments                                                                                       | list      |                         |
| `prometheus.deployment.containers.resources`                 | Prometheus container resource requests and limits                                                                    | map       | {}                      |
| `prometheus.deployment.podAnnotations`                       | The Prometheus deployments pod annotations                                                                           | map       | {}                      |
| `prometheus.deployment.podLabels`                            | The Prometheus deployments pod labels                                                                                | map       | {}                      |
| `prometheus.deployment.configMapReload.containers.args`      | Configmap-reload container arguments.                                                                                | list      |                         |
| `prometheus.deployment.configMapReload.containers.resources` | Configmap-reload container resource requests and limits                                                              | map       | {}                      |
| `prometheus.service.type`                                    | Type of service to expose Prometheus. Supported Values: ClusterIP                                                    | string    | ClusterIP               |
| `prometheus.service.port`                                    | Prometheus service port                                                                                              | integer   | 80                      |
| `prometheus.service.targetPort`                              | Prometheus service target port                                                                                       | integer   | 9090                    |
| `prometheus.service.labels`                                  | Prometheus service labels                                                                                            | map       | {}                      |
| `prometheus.service.annotations`                             | Prometheus service annotations                                                                                       | map       | {}                      |
| `prometheus.pvc.annotations`                                 | Storage class annotations                                                                                            | map       | {}                      |
| `prometheus.pvc.storageClassName`                            | Storage class to use for persistent volume claim. By default this is null and default provisioner is used            | string    | null                    |
| `prometheus.pvc.accessMode`                                  | Define access mode for persistent volume claim. Supported values: ReadWriteOnce, ReadOnlyMany, ReadWriteMany         | string    | ReadWriteOnce           |
| `prometheus.pvc.storage`                                     | Define storage size for persistent volume claim                                                                      | string    | 150Gi                   |
| `prometheus.config.prometheus_yml`                           | The [global Prometheus configuration](https://www.prometheus.io/docs/prometheus/latest/configuration/configuration/) | yaml file | prometheus.yaml         |
| `prometheus.config.alerting_rules_yml`                       | The [Prometheus alerting rules](https://www.prometheus.io/docs/prometheus/latest/configuration/alerting_rules/)      | yaml file | alerting_rules.yaml     |
| `prometheus.config.recording_rules_yml`                      | The [prometheus recording rules](https://www.prometheus.io/docs/prometheus/latest/configuration/recording_rules/)    | yaml file | recording_rules.yaml    |
| `prometheus.config.alerts_yml`                               | Additional Prometheus alerting rules can be configured here.                                                         | yaml file | alerts_yml.yaml         |
| `prometheus.config.rules_yml`                                | Additional Prometheus recording rules can be configured here.                                                        | yaml file | rules_yml.yaml          |
| `alertmanager.deployment.replicas`                           | Number of alertmanager replicas                                                                                      | integer   | 1                       |
| `alertmanager.deployment.containers.resources`               | Alertmanager container resource requests and limits                                                                  | map       | {}                      |
| `alertmanager.deployment.podAnnotations`                     | The Alertmanager deployments pod annotations                                                                         | map       | {}                      |
| `alertmanager.deployment.podLabels`                          | The Alertmanager deployments pod labels                                                                              | map       | {}                      |
| `alertmanager.service.type`                                  | Type of service to expose Alertmanager. Supported Values: ClusterIP                                                  | string    | ClusterIP               |
| `alertmanager.service.port`                                  | Alertmanager service port                                                                                            | integer   | 80                      |
| `alertmanager.service.targetPort`                            | Alertmanager service target port                                                                                     | integer   | 9093                    |
| `alertmanager.service.labels`                               | Alertmanager service labels                                                                                          | map       | {}                      |
| `alertmanager.service.annotations`                           | Alertmanager service annotations                                                                                     | map       | {}                      |
| `alertmanager.pvc.annotations`                               | Storage class annotations                                                                                            | map       | {}                      |
| `alertmanager.pvc.storageClassName`                          | Storage class to use for persistent volume claim. By default this is null and default provisioner is used.           | string    | null                    |
| `alertmanager.pvc.accessMode`                                | Define access mode for persistent volume claim. Supported values: ReadWriteOnce, ReadOnlyMany, ReadWriteMany         | string    | ReadWriteOnce           |
| `alertmanager.pvc.storage`                                   | Define storage size for persistent volume claim                                                                      | string    | 2Gi                     |
| `alertmanager.config.alertmanager_yml`                       | The [global yaml configuration for alert manager](https://www.prometheus.io/docs/alerting/latest/configuration/).    | yaml file | alertmanager_yml        |
| `kube_state_metrics.deployment.replicas`                     | Number of kube-state-metrics replicas                                                                                | integer   | 1                       |
| `kube_state_metrics.deployment.containers.resources`         | kube-state-metrics container resource requests and limits                                                            | map       | {}                      |
| `kube_state_metrics.deployment.podAnnotations`               | The kube-state-metrics deployments pod annotations                                                                   | map       | {}                      |
| `kube_state_metrics.deployment.podLabels`                    | The kube-state-metrics deployments pod labels                                                                        | map       | {}                      |
| `kube_state_metrics.service.type`                            | Type of service to expose kube-state-metrics. Supported Values: ClusterIP                                            | string    | ClusterIP               |
| `kube_state_metrics.service.port`                            | kube-state-metrics service port                                                                                      | integer   | 80                      |
| `kube_state_metrics.service.targetPort`                      | kube-state-metrics service target port                                                                               | integer   | 8080                    |
| `kube_state_metrics.service.telemetryPort`                   | kube-state-metrics service telemetry port                                                                            | integer   | 81                      |
| `kube_state_metrics.service.telemetryTargetPort`            | kube-state-metrics service target telemetry  port                                                                    | integer   | 8081                    |
| `kube_state_metrics.service.labels`                          | kube-state-metrics service labels                                                                                    | map       | {}                      |
| ``kube_state_metrics.service.annotations``                     | kube-state-metrics service annotations                                                                               | map       | {}                      |
| `node_exporter.daemonset.replicas`                           | Number of node-exporter replicas                                                                                     | integer   | 1                       |
| `node_exporter.daemonset.containers.resources`               | node-exporter container resource requests and limits                                                                 | map       | {}                      |
| `node_exporter.daemonset.hostNetwork`                        |  Host networking requested for this pod                                                                              | boolean   | false                   |
| `node_exporter.daemonset.podAnnotations`                     | The node-exporter deployments pod annotations                                                                        | map       | {}                      |
| `node_exporter.daemonset.podLabels`                          | The node-exporter deployments pod labels                                                                             | map       | {}                      |
| `node_exporter.service.type`                                 | Type of service to expose node-exporter. Supported Values: ClusterIP                                                 | string    | ClusterIP               |
| `node_exporter.service.port`                                 | node-exporter service port                                                                                           | integer   | 9100                    |
| `node_exporter.service.targetPort`                           | node-exporter service target port                                                                                    | integer   | 9100                    |
| `node_exporter.service.labels`                               | node-exporter service labels                                                                                         | map       | {}                      |
| `node_exporter.service.annotations`                          | node-exporter service annotations                                                                                    | map       | {}                      |
| `pushgateway.deployment.replicas`                            | Number of pushgateway replicas                                                                                       | integer   | 1                       |
| `pushgateway.deployment.containers.resources`                | pushgateway container resource requests and limits                                                                   | map       | {}                      |
| `pushgateway.deployment.podAnnotations`                      | The pushgateway deployments pod annotations                                                                          | map       | {}                      |
| `pushgateway.deployment.podLabels`                           | The pushgateway deployments pod labels                                                                               | map       | {}                      |
| `pushgateway.service.type`                                  | Type of service to expose pushgateway. Supported Values: ClusterIP                                                   | string    | ClusterIP               |
| `pushgateway.service.port`                                   | pushgateway service port                                                                                             | integer   | 9091                    |
| `pushgateway.service.targetPort`                             | pushgateway service target port                                                                                      | integer   | 9091                    |
| `pushgateway.service.labels`                                 | pushgateway service labels                                                                                           | map       | {}                      |
| `pushgateway.service.annotations`                            | pushgateway service annotations                                                                                      | map       | {}                      |
| `cadvisor.daemonset.replicas`                                | Number of cadvisor replicas                                                                                          | integer   | 1                       |
| `cadvisor.daemonset.containers.resources`                    | cadvisor container resource requests and limits                                                                      | map       | {}                      |
| `cadvisor.daemonset.podAnnotations`                          | The cadvisor deployments pod annotations                                                                             | map       | {}                      |
| `cadvisor.daemonset.podLabels`                               | The cadvisor deployments pod labels                                                                                  | map       | {}                      |
| `ingress.enabled`                                            | Activate/deactivate ingress for prometheus and alertmanager                                                               | boolean   | false                   |
| `ingress.virtual_host_fqdn`                                  | Hostname for accessing prometheus and alertmanager                                                                   | string    | prometheus.system.tanzu |
| `ingress.prometheus_prefix`                                  | Path prefix for prometheus                                                                                           | string    | /                       |
| `ingress.alertmanager_prefix`                               | Path prefix for alertmanager                                                                                         | string    | /alertmanager/          |
| `ingress.prometheusServicePort`                              | Prometheus service port to proxy traffic to                                                                          | integer   | 80                      |
| `ingress.alertmanagerServicePort`                            | Alertmanager service port to proxy traffic to                                                                        | integer   | 80                      |
| `ingress.tlsCertificate.tls.crt`                             | Optional cert for ingress if you want to use your own TLS cert. A self signed cert is generated by default           | string    | Generated cert          |
| `ingress.tlsCertificate.tls.key`                             | Optional cert private key for ingress if you want to use your own TLS cert.                                          | string    | Generated cert key      |
| `ingress.tlsCertificate.ca.crt`                              | Optional CA certificate                                                                                              | string    | CA certificate          |

### Config files

The various configuration files can be loaded from `/etc/config/`.
For example, when loading `rule_files`:

```text
rule_files:
- /etc/config/alerting_rules.yml
- /etc/config/recording_rules.yml
- /etc/config/alerts
- /etc/config/rules
```

### Alert manager service

The alert manager can be targeted by the deployed prometheus through it's service.
For example:

```text
targets:
- alertmanager.prometheus.svc:9093
```

`alertmanager` is the default service name and `prometheus` is the namespace of the alertmanager deployment.

## Usage Example

The default `prometheus.yml` configuration will deploy a prometheus server
that will scrape metrics from pods that emit metrics
on an endpoint and has the following annotations:

```yaml
prometheus.io/scrape: 'true'
prometheus.io/path: '/metrics'
prometheus.io/port: '8080'
```
