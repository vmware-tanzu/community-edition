# Azure File CSI Driver Package

This package provides cloud storage interface driver using [azurefile-csi-driver](https://github.com/kubernetes-sigs/azurefile-csi-driver).

## Installation

The Azure File CSI Driver package is installed automatically during cluster creation.

## Options

### Package configuration values

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `nodeSelector` | Optional | NodeSelector configuration applied to all the deployments. Defaults to null. |
| `deployment.updateStrategy` | Optional | The update strategy of deployments to overwrite. Defaults to null. |
| `deployment.rollingUpdate.maxUnavailable` | Optional | The maxUnavailable of rollingUpdate. Applied only if RollingUpdate is used as updateStrategy. Defaults to null. |
| `deployment.rollingUpdate.maxSurge` | Optional | The maxSurge of rollingUpdate. Applied only if RollingUpdate is used as updateStrategy. Defaults to null. |
| `daemonset.updateStrategy` | Optional | The update strategy of daemonsets to overwrite. Defaults to null. |

### Azure File CSI Driver Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `azureFileCSIDriver.namespace` | Required | The namespace of the Kubernetes cluster in cluster ID. Default value is `kube-system`. |
| `azureFileCSIDriver.http_proxy`                     | Optional          | The HTTP proxy to use for network traffic                                                                                                                                                         |
| `azureFileCSIDriver.https_proxy`                    | Optional          | The HTTPS proxy to use for network traffic                                                                                                                                                        |
| `azureFileCSIDriver.no_proxy`                       | Optional          | A comma-separated list of hostnames, IP addresses, or IP ranges in CIDR format that should not use a proxy
| `azureFileCSIDriver.deployment_replicas` | Optional | The number of replicas of csi-azurefile-controller. Default: `3`. |

### Application configuration values

No available options to configure.

#### Multi-cloud configuration steps

No extra configuration steps needed, Azure File CSI Driver is Azure only.

## What This Package Does

The Azure File Container Storage Interface (CSI) Driver provides a CSI interface used by Container Orchestrators to manage the lifecycle of Azure File volumes.

## Components

* Azure File CSI Driver DaemonSet
* Azure File CSI Driver ServiceAccount
* Azure File CSI Driver ClusterRole
* Azure File CSI Driver ClusterRoleBinding
* Azure File CSI Driver Deployment
* Azure File CSI Driver ConfigMap

## Supported Providers

The following tables shows the providers this package can work with.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
|  ❌ |  ✅  | ❌  |  ❌  |

## Files

Here is an example of the package configuration file [values.yaml](bundle/config/values.yaml).

## Package Limitations

Not applicable.

## Usage Example

To learn more about how to use Azure Disk CSI Driver refer to [Azure File CSI Driver website](https://github.com/kubernetes-sigs/azurefile-csi-driver)

## Troubleshooting

Not applicable.

## Additional Documentation

See the [Azure File CSI Driver website](https://github.com/kubernetes-sigs/azurefile-csi-driver) for more information.
