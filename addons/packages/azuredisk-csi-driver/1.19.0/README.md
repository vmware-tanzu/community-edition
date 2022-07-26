# Azure Disk CSI Driver Package

This package provides cloud storage interface driver using [azuredisk-csi-driver](https://github.com/kubernetes-sigs/azuredisk-csi-driver).

## Installation

The Azure Disk CSI Driver package is installed automatically during cluster creation.

## Options

### Package configuration values

#### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `nodeSelector` | Optional | NodeSelector configuration applied to all the deployments. Defaults to null. |
| `deployment.updateStrategy` | Optional | The update strategy of deployments to overwrite. Defaults to null. |
| `deployment.rollingUpdate.maxUnavailable` | Optional | The maxUnavailable of rollingUpdate. Applied only if RollingUpdate is used as updateStrategy. Defaults to null. |
| `deployment.rollingUpdate.maxSurge` | Optional | The maxSurge of rollingUpdate. Applied only if RollingUpdate is used as updateStrategy. Defaults to null. |
| `daemonset.updateStrategy` | Optional | The update strategy of daemonsets to overwrite. Defaults to null. |

#### Azure Disk CSI Driver Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `azureDiskCSIDriver.namespace` | Required | The namespace of the Kubernetes cluster in cluster ID. Default value is `kube-system`. |
| `azureDiskCSIDriver.http_proxy`                     | Optional          | The HTTP proxy to use for network traffic                                                                                                                                                         |
| `azureDiskCSIDriver.https_proxy`                    | Optional          | The HTTPS proxy to use for network traffic                                                                                                                                                        |
| `azureDiskCSIDriver.no_proxy`                       | Optional          | A comma-separated list of hostnames, IP addresses, or IP ranges in CIDR format that should not use a proxy
| `azureDiskCSIDriver.deployment_replicas` | Optional | The number of replicas of csi-azuredisk-controller and csi-snapshot-controller deployment. Default: `3`. |

### Application configuration values

No available options to configure.

#### Multi-cloud configuration steps

No extra configuration steps needed, Azure Disk CSI Driver is Azure only.

## What This Package Does

The Azure Disk Container Storage Interface (CSI) Driver provides a CSI interface used by Container Orchestrators to manage the lifecycle of Azure Disk volumes.

## Components

* Azure Disk CSI Driver DaemonSet
* Azure Disk CSI Driver ServiceAccount
* Azure Disk CSI Driver ClusterRole
* Azure Disk CSI Driver ClusterRoleBinding
* Azure Disk CSI Driver Deployment
* Azure Disk CSI Driver ConfigMap
* CSI Snapshot Custom Resources
* CSI Snapshot ServiceAccount
* CSI Snapshot ClusterRole
* CSI Snapshot ClusterRoleBinding
* CSI Snapshot Deployment

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

To learn more about how to use Azure Disk CSI Driver refer to [Azure Disk CSI Driver website](https://github.com/kubernetes-sigs/azuredisk-csi-driver)

## Troubleshooting

Not applicable.

## Additional Documentation

See the [Azure Disk CSI Driver website](https://github.com/kubernetes-sigs/azuredisk-csi-driver) for more information.
