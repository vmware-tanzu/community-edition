# Azure Disk CSI Driver

This package provides cloud storage interface driver using [azuredisk-csi-driver](https://github.com/kubernetes-sigs/azuredisk-csi-driver).

## Supported Providers

The following tables shows the providers this package can work with.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
|  ❌ |  ✅  | ❌  |  ❌  |

## Configuration

The following configuration values can be set to customize the azuredisk-csi-driver installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `nodeSelector` | Optional | NodeSelector configuration applied to all the deployments. Defaults to null. |
| `deployment.updateStrategy` | Optional | The update strategy of deployments to overwrite. Defaults to null. |
| `deployment.rollingUpdate.maxUnavailable` | Optional | The maxUnavailable of rollingUpdate. Applied only if RollingUpdate is used as updateStrategy. Defaults to null. |
| `deployment.rollingUpdate.maxSurge` | Optional | The maxSurge of rollingUpdate. Applied only if RollingUpdate is used as updateStrategy. Defaults to null. |
| `daemonset.updateStrategy` | Optional | The update strategy of daemonsets to overwrite. Defaults to null. |

### azuredisk-csi-driver Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `azureDiskCSIDriver.namespace` | Required | The namespace of the Kubernetes cluster in cluster ID. Default value is `kube-system`. |
| `azureDiskCSIDriver.http_proxy`                     | Optional          | The HTTP proxy to use for network traffic                                                                                                                                                         |
| `azureDiskCSIDriver.https_proxy`                    | Optional          | The HTTPS proxy to use for network traffic                                                                                                                                                        |
| `azureDiskCSIDriver.no_proxy`                       | Optional          | A comma-separated list of hostnames, IP addresses, or IP ranges in CIDR format that should not use a proxy
| `azureDiskCSIDriver.deployment_replicas` | Optional | The number of replicas of csi-azuredisk-controller and csi-snapshot-controller deployment. Default: `3`. |

## What This Package Does

The Azure Disk Container Storage Interface (CSI) Driver provides a CSI interface used by Container Orchestrators to manage the lifecycle of Azure Disk volumes.

## Usage Example

To learn more about how to use Azure Disk CSI Driver refer to [Azure Disk CSI Driver website](https://github.com/kubernetes-sigs/azuredisk-csi-driver)
