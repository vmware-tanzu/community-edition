# antrea Package

This package provides networking and network security solution for containers using [antrea](https://antrea.io/).

## Components

## Configuration

The following configuration values can be set to customize the antrea installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `infraProvider` | Required | The cloud provider in use. One of: `aws`, `azure`, `vsphere`, `docker`. |

### antrea Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `antrea.config.serviceCIDR` | Optional | The service IPv4 CIDR to use. Default: `10.96.0.0/12` |
| `antrea.config.serviceCIDRv6` | Optional | The service IPv6 CIDR to use. Default: nil |
| `antrea.config.trafficEncapMode` | Optional | The traffic encapsulation mode. Default: `encap` |
| `antrea.config.noSNAT` | Optional | Boolean flag to enable/disable SNAT. Default: `false` |
| `antrea.config.defaultMTU` | Optional | MTU to use. Default: `null` (Antrea will autodetect) |
| `antrea.config.tlsCipherSuites` | Optional | List of allowed cipher suites  |
| `antrea.config.disableUdpTunnelOffload` | Optional | Disable UDP tunnel offload feature on default NIC. Default: `false` |
| `antrea.config.featureGates.AntreaProxy` | Optional | Boolean flag to enable/disable antrea proxy. Default: `false` |
| `antrea.config.featureGates.EndpointSlice` | Optional | Boolean flag to enable/disable EndpointSlice support in AntreaProxy. Default: `false` |
| `antrea.config.featureGates.AntreaPolicy` | Optional | Boolean flag to enable/disable antrea policy. Default: `true` |
| `antrea.config.featureGates.AntreaTraceFlow` | Optional | Boolean flag to enable/disable antrea traceflow. Default: `false` |
| `antrea.config.featureGates.Egress` | Optional | Boolean flag to enable/disable SNAT IPs of Pod egress traffic. Default: `false` |
| `antrea.config.featureGates.NodePortLocal` | Optional | Boolean flag to enable/disable NodePortLocal feature to make the pods reachable externally through NodePort |
| `antrea.config.featureGates.FlowExporter`| Optional | Boolean flag to enable/disable flow exporter. Default: `false` |
| `antrea.config.featureGates.NetworkPolicyStats` | Optional | Boolean flag to enable/disable network policy stats. Default: `false` |

## Usage Example

To learn more about how to use antrea refer to [antrea documentation](https://antrea.io/docs/v1.2.3/)
