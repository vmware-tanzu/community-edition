# Grafana

The package provides the ability to explore your metrics with kubernetes cluster using an open-source system named [Grafana](https://github.com/grafana/grafana).

This package contains Grafana 7.5.16 (under the Apache 2.0 license).

## Installation

### Installation of dependencies

By default, grafana depends on Contour.  Install TCE Contour package through tanzu command:

```bash
tanzu package install contour --package-name contour.community.tanzu.vmware.com --version ${CONTOUR_PACKAGE_VERSION}
```

> You can get the `${CONTOUR_PACKAGE_VERSION}` from running `tanzu package
> available list contour.community.tanzu.vmware.com`. Specifying a
> namespace may be required depending on where your package repository was
> installed.

### Installation of package

Install Grafana package through tanzu command:

```bash
tanzu package install grafana --package-name grafana.community.tanzu.vmware.com --version ${GRAFANA_PACKAGE_VERSION}
```

> You can get the `${GRAFANA_PACKAGE_VERSION}` from running `tanzu package
> available list grafana.community.tanzu.vmware.com`. Specifying a
> namespace may be required depending on where your package repository was
> installed.

## Options

The following configuration values can be set to customize the Grafana installation.

### Package configuration values

#### Global

| Value | Required/Optional | Default | Description |
|:-------|:-------------------|:---------|:-------------|
| `namespace` | Optional | `grafana` | The namespace in which to deploy Grafana  components |

#### Grafana Configuration

> Note: Ingress for Grafana server is enabled by default, and can be disabled using the `ingress.enabled` configuration field. For clusters running in Docker, disabling the Ingress is the easiest way to get started, as setting up Contour on a Docker cluster requires additional configuration.
> If you choose to enable the Contour-based Ingress, Contour must also be installed on the target cluster. Additionally, enabling the Ingress requires either cert-manager or your own user-provided TLS certificate (`ingress.tlsCertificate.tls.crt` and `ingress.tlsCertificate.tls.key`) to configure TLS settings for the Ingress. For ad-hoc Grafana UI access without an Ingress, use kubectl port-forward.

| Value                                                | Required/Optional | Default                    | Description                                                                                                                       |
|:------------------------------------------------------|:-------------------|:----------------------------|:-----------------------------------------------------------------------------------------------------------------------------------|
| `grafana.deployment.replicas`                        | true              | 1                          | Number of Grafana replicas                                                                                                        |
| `grafana.deployment.updateStrategy`                  | false             | Recreate                   | Type of Grafana upgrade strategy. Supported Values: RollingUpdate, Recreate                                                       |
| `grafana.deployment.rollingUpdate.maxUnavailable`    | false             | null                       | Number of maxUnavailable pods when grafana is upgrading, It only work when updateStrategy is RollingUpdate                        |
| `grafana.deployment.rollingUpdate.maxSurge`          | false             | null                       | Number of maxSurge pods when grafana is upgrading, It only work when updateStrategy is RollingUpdate                              |
| `grafana.deployment.containers.resources`            | false             | {}                         | Grafana container resource requests and limits                                                                                    |
| `grafana.deployment.k8sSidecar.containers.resources` | false             | {}                         | k8s-sidecar container resource requests and limits                                                                                |
| `grafana.deployment.podAnnotations`                  | false             | {}                         | The Grafana deployments pod annotations                                                                                           |
| `grafana.deployment.podLabels`                       | false             | {}                         | The Grafana deployments pod labels                                                                                                |
| `grafana.service.type`                               | true              | LoadBalancer               | Type of service to expose Grafana. Supported Values: ClusterIP, NodePort, LoadBalancer. (For vSphere set this to NodePort)        |
| `grafana.service.port`                               | true              | 80                         | Grafana service port                                                                                                              |
| `grafana.service.targetPort`                         | true              | 9093                       | Grafana service target port                                                                                                       |
| `grafana.service.labels`                             | false             | {}                         | Grafana service labels                                                                                                            |
| `grafana.service.annotations`                        | false             | {}                         | Grafana service annotations                                                                                                       |
| `grafana.config.grafana_ini`                         | true              | grafana.ini                | The [Grafana configuration](https://github.com/grafana/grafana/blob/master/conf/defaults.ini)                                     |
| `grafana.config.datasource_yaml`                     | true              | prometheus                 | Grafana [datasource config](https://grafana.com/docs/grafana/latest/administration/provisioning/#example-data-source-config-file) |
| `grafana.config.dashboardProvider_yaml`              | true              | provider.yaml              | Grafana [dashboard provider config](https://grafana.com/docs/grafana/latest/administration/provisioning/#dashboards)              |
| `grafana.pvc.annotations`                            | false             | null                       | Storage class annotations                                                                                                         |
| `grafana.pvc.storageClassName`                       | false             | null                       | Storage class to use for persistent volume claim. By default this is null and default provisioner is used                         |
| `grafana.pvc.accessMode`                             | true              | ReadWriteOnce              | Define access mode for persistent volume claim. Supported values: ReadWriteOnce, ReadOnlyMany, ReadWriteMany                      |
| `grafana.pvc.storage`                                | false             | 2Gi                        | Define storage size for persistent volume claim                                                                                   |
| `grafana.secret.type`                                | false             | Opaque                     | Secret type defined for Grafana dashboard                                                                                         |
| `grafana.secret.admin_user`                          | false             | YWRtaW4=                   | username to access Grafana dashboard                                                                                              |
| `grafana.secret.admin_password`                      | false             | admin                      | password to access Grafana dashboard                                                                                              |
| `ingress.enabled`                                    | false             | true                       | Enable/disable ingress for grafana                                                                                                |
| `ingress.virtual_host_fqdn`                          | false             | grafana.system.tanzu       | Hostname for accessing grafana                                                                                                    |
| `ingress.prefix`                                     | false             | /                          | Path prefix for grafana                                                                                                           |
| `ingress.servicePort`                                | false             | 80                         | Grafana service port to proxy traffic to                                                                                          |
| `ingress.tlsCertificate.tls.crt`                     | false             | Generated cert             | Optional cert for ingress if you want to use your own TLS cert. A self signed cert is generated by default                        |
| `ingress.tlsCertificate.tls.key`                     | false             | Generated cert private key | Optional cert private key for ingress if you want to use your own TLS cert.                                                       |
| `ingress.tlsCertificate.ca.crt`                      | false             | CA certificate             | Optional CA certificate    |

### Application configuration values

No available options to configure.

#### Multi-cloud configuration steps

There are currently no configuration steps necessary for installation of the Grafana package to any provider.

## What This Package Does

Grafana is an open source visualization and analytics software. It allows you to query, visualize, alert on, and explore your metrics no matter where they are stored. It provides you with tools to turn your time-series database (TSDB) data into beautiful graphs and visualizations.

## Components

- Grafana server
- k8s-sidecar

## Supported Providers

The following table shows the providers this package can work with.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
| ✅  |  ✅  | ✅  | ✅ |

## Files

Here is an example of the package configuration file [values.yaml](bundle/config/values.yaml).

## Package Limitations

- Multi-replicas in Grafana-server and multi-instances server are not support.

## Usage Example

- Set up data sources for your metrics, you can add one or more data sources to Grafana. See the Grafana [documentation](https://grafana.com/docs/grafana/latest/datasources/add-a-data-source/) for detailed description of how to add a data source.
- Create Dashboards
There are many prebuilt Grafana dashboard templates available for various data sources. You can check out the templates [here](https://grafana.com/grafana/dashboards).
- Activate Ingress on Grafana as per your requirement.

## Troubleshooting

Not applicable.

## Additional Documentation

See the [documentation](https://grafana.com/docs/grafana/v7.5/) for more information.
