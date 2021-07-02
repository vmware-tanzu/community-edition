# antrea Package

This package provides networking and network security solution for containers using [antrea](https://antrea.io/).

## Components

## Configuration

The following configuration values can be set to customize the antrea installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Optional | The namespace in which to deploy antrea. |
| `infraProvider` | Required | The cloud provider in use. One of: `aws`, `azure`, `vsphere`, `docker`. |

### antrea Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `antrea.config.serviceCIDR` | Optional | The service CIDR to use. Default: `10.96.0.0/12` |
| `antrea.config.trafficEncapMode` | Optional | The traffic encapsulation mode. Default: `encap` |
| `antrea.config.noSNAT` | Optional | Boolean flag to enable/disable SNAT. Default: `false`. |
| `antrea.config.defaultMTU` | Optional | MTU to use. Default: `null` (Antrea will autodetect). |
| `antrea.config.featureGates.AntreaProxy` | Optional | Boolean flag to enable/disable antrea proxy. Default: `false`. |
| `antrea.config.featureGates.AntreaPolicy` | Optional | Boolean flag to enable/disable antrea policy. Default: `true`. |
| `antrea.config.featureGates.AntreaTraceFlow` | Optional | Boolean flag to enable/disable antrea traceflow. Default: `false`. |
| `antrea.config.featureGates.FlowExporter`| Optional | Boolean flag to enable/disable flow exporter. Default: `false`. |
| `antrea.config.featureGates.NetworkPolicyStats` | Optional | Boolean flag to enable/disable network policy stats. Default: `false`. |

## Usage Example

To learn more about how to use antrea refer to [antrea documentation](https://antrea.io/docs/v0.11.3/)
