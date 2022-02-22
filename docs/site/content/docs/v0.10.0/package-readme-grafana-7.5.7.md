# Grafana

Grafana is an open source visualization and analytics software. It allows you to query, visualize, alert on, and explore your metrics no matter where they are stored. It provides you with tools to turn your time-series database (TSDB) data into beautiful graphs and visualizations.

This package contains Grafana 7.5.7 (under the Apache 2.0 license).

## Supported Providers

The following table shows the providers this package can work with.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
| ✅  |  ✅  | ✅  | ✅ |

## Components

- Grafana server.

## Configuration

The following configuration values can be set to customize the Grafana installation.

> Note: Ingress for Grafana server is enabled by default, and can be disabled using the `ingress.enabled` configuration field. For clusters running in Docker, disabling the Ingress is the easiest way to get started, as setting up Contour on a Docker cluster requires additional configuration.
> If you choose to enable the Contour-based Ingress, Contour must also be installed on the target cluster. Additionally, enabling the Ingress requires either cert-manager or your own user-provided TLS certificate (`ingress.tlsCertificate.tls.crt` and `ingress.tlsCertificate.tls.key`) to configure TLS settings for the Ingress. For ad-hoc Grafana UI access without an Ingress, use kubectl port-forward.

| Parameter                                          | Description                                                                                                                       | Type        | Default                                    |
|:----------------------------------------------------|:-----------------------------------------------------------------------------------------------------------------------------------|:-------------|:--------------------------------------------|
| namespace                                          | Namespace where Grafana will be deployed                                                                                          | string      | grafana                              |
| `grafana.deployment.replicas`                        | Number of Grafana replicas                                                                                                        | integer     | 1                                          |
| `grafana.deployment.containers.resources`            | Grafana container resource requests and limits                                                                                    | map         | {}                                         |
| `grafana.deployment.k8sSidecar.containers.resources` | k8s-sidecar container resource requests and limits                                                                                | map         | {}                                         |
| `grafana.deployment.podAnnotations`                  | The Grafana deployments pod annotations                                                                                           | map         | {}                                         |
| `grafana.deployment.podLabels`                       | The Grafana deployments pod labels                                                                                                | map         | {}                                         |
| `grafana.service.type`                               | Type of service to expose Grafana. Supported Values: ClusterIP, NodePort, LoadBalancer. (For vSphere set this to NodePort)        | string      | LoadBalancer                               |
| `grafana.service.port`                               | Grafana service port                                                                                                              | integer     | 80                                         |
| `grafana.service.targetPort`                         | Grafana service target port                                                                                                       | integer     | 9093                                       |
| `grafana.service.labels`                             | Grafana service labels                                                                                                            | map         | {}                                         |
| `grafana.service.annotations`                        | Grafana service annotations                                                                                                       | map         | {}                                         |
| `grafana.config.grafana_ini`                         | The [Grafana configuration](https://github.com/grafana/grafana/blob/master/conf/defaults.ini)                                     | config file | grafana.ini                                |
| `grafana.config.datasource_yaml`                     | Grafana [datasource config](https://grafana.com/docs/grafana/latest/administration/provisioning/#example-data-source-config-file) | string      | prometheus                                 |
| `grafana.config.dashboardProvider_yaml`              | Grafana [dashboard provider config](https://grafana.com/docs/grafana/latest/administration/provisioning/#dashboards)              | yaml file   | provider.yaml                              |
| `grafana.pvc.annotations`                            | Storage class to use for persistent volume claim. By default this is null and default provisioner is used                         | string      | null                                       |
| `grafana.pvc.storageClassName`                       | Storage class to use for persistent volume claim. By default this is null and default provisioner is used                         | string      | null                                       |
| `grafana.pvc.accessMode`                             | Define access mode for persistent volume claim. Supported values: ReadWriteOnce, ReadOnlyMany, ReadWriteMany                      | string      | ReadWriteOnce                              |
| `grafana.pvc.storage`                                | Define storage size for persistent volume claim                                                                                   | string      | 2Gi                                        |
| `grafana.secret.type`                                | Secret type defined for Grafana dashboard                                                                                         | string      | Opaque                                     |
| `grafana.secret.admin_user`                          | username to access Grafana dashboard                                                                                              | string      | YWRtaW4=                                   |
| `grafana.secret.admin_password`                      | password to access Grafana dashboard                                                                                              | string      | admin                                      |
| `ingress.enabled`                                    | Enable/disable ingress for grafana                                                                                                | boolean     | true                                       |
| `ingress.virtual_host_fqdn`                         | Hostname for accessing grafana                                                                                                    | string      | grafana.system.tanzu                       |
| `ingress.prefix`                                    | Path prefix for grafana                                                                                                           | string      | /                                          |
| `ingress.servicePort`                                | Grafana service port to proxy traffic to                                                                                          | integer     | 80                                         |
| `ingress.tlsCertificate.tls.crt`                    | Optional cert for ingress if you want to use your own TLS cert. A self signed cert is generated by default                        | string      | Generated cert                             |
| `ingress.tlsCertificate.tls.key`                     | Optional cert private key for ingress if you want to use your own TLS cert.                                                       | string      | Generated cert private key                 |
| `ingress.tlsCertificate.ca.crt`                      | Optional CA certificate                                                                                                           | string      | CA certificate                             |

## Usage Example

- Set up data sources for your metrics, you can add one or more data sources to Grafana. See the Grafana [documentation](https://grafana.com/docs/grafana/latest/datasources/add-a-data-source/) for detailed description of how to add a data source.
- Create Dashboards
There are many prebuilt Grafana dashboard templates available for various data sources. You can check out the templates [here](https://grafana.com/grafana/dashboards).
- Activate Ingress on Grafana as per your requirement.
