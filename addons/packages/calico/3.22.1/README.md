# Calico Package

This package provides networking and network security solution for containers using [calico](https://www.projectcalico.org/).

## Installation

As a primary CNI option, like Antrea, the Calico package is installed automatically during cluster creation.

## Options

The following configuration values can be set to customize the calico installation.

### Package configuration values

#### Global

| Value | Required/Optional | Default | Description |
|-------|-------------------|---------|-------------|
| `infraProvider` | Required | `vsphere` | The infrastructure provider in use. One of: `aws`, `azure`, `vsphere`, `docker`. |
| `ipFamily` | Optional | `null` | The IP family calico should be configured with. Defaults to `ipv4`. One of: `ipv4`, `ipv6`, `ipv4,ipv6` (IPv4-primary dualstack), or `ipv6,ipv4` (IPv6-primary dualstack). |
| `nodeSelector` | Optional | `null` | NodeSelector configuration applied to all the deployments. Defaults to null. |
| `deployment.updateStrategy` | Optional | `null` | The update strategy of deployments to overwrite. Defaults to null. |
| `deployment.rollingUpdate.maxUnavailable` | Optional | `null` | The maxUnavailable of rollingUpdate. Applied only if RollingUpdate is used as updateStrategy. Defaults to null. |
| `deployment.rollingUpdate.maxSurge` | Optional | `null` | The maxSurge of rollingUpdate. Applied only if RollingUpdate is used as updateStrategy. Defaults to null. |
| `daemonset.updateStrategy` | Optional | `null` | The update strategy of daemonsets to overwrite. Defaults to null. |

#### Calico Configuration

| Value | Required/Optional |  Default | Description |
|-------|-------------------|----------|-------------|
| `calico.config.clusterCIDR` | Optional | `null` | The pod network CIDR. Default value is empty because it can be auto detected. |
| `calico.config.vethMTU` | Optional | `"0"` | MTU size. Default is `0`, which means it will be auto detected. |
| `calico.config.skipCNIBinaries` | Optional | `false` | Skip the cni plugin(bandwidth, flannel, host-local, loopback, portmap, tuning) binaries installation from calico to avoid the overwrite when, for some cases, the nodes already contains these cni plugins. Default to `false`|

### Application configuration values

No available options to configure.

#### Multi-cloud configuration steps

Set `infraProvider` to one of `aws`, `azure`, `vsphere`, `docker` to specify the infrastructure provider.

## What This Package Does

Calico is an open source networking and network security solution for containers, virtual machines, and native host-based workloads.

## Components

* Calico Custom Resources
* Calico DaemonSet
* Calico ServiceAccount
* Calico ClusterRole
* Calico ClusterRoleBinding
* Calico Deployment
* Calico ConfigMap
* Calico PodDisruptionBudget

### Supported Providers

The following table shows the providers this package can work with.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
| ✅  |  ✅  | ✅  | ✅ |

## Files

Here is an example of the package configuration file [values.yaml](bundle/config/values.yaml).

## Package Limitations

As a core package, Calico should be installed automatically during cluster creation. No manual installation required, especially on a cluster who already has CNI.

## Usage Example

To learn more about how to use calico refer to [calico documentation](https://docs.projectcalico.org/about/about-calico).

## Troubleshooting

To learn more about how to troubleshoot calico installation and operation refer to [calico documentation](https://docs.projectcalico.org/about/about-calico).

## Additional Documentation

See the [calico documentation](https://docs.projectcalico.org/about/about-calico) for more information.
