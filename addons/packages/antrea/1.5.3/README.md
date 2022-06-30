# antrea Package

This package provides networking and network security solution for containers using [antrea](https://antrea.io/).

## Components

## Configuration

The following configuration values can be set to customize the antrea installation.

### Global

| Value           | Required/Optional | Description                                                             |
|-----------------|-------------------|-------------------------------------------------------------------------|
| `infraProvider` | Required          | The cloud provider in use. One of: `aws`, `azure`, `vsphere`, `docker`. |

### antrea Configuration

| Value                                            | Required/Optional | Description                                                                                                   |
|--------------------------------------------------|-------------------|---------------------------------------------------------------------------------------------------------------|
| `antrea.config.egress.exceptCIDRs`               | Optional          | The CIDR ranges to which outbound Pod traffic will not be SNAT'd by Egresses                                  |
| `antrea.config.nodePortLocal.enabled`            | Optional          | Enable NodePortLocal feature. Default: true                                                                   |
| `antrea.config.nodePortLocal.portRange`          | Optional          | Provide the port range used by NodePortLocal                                                                  |
| `antrea.config.antreaProxy.proxyAll`             | Optional          | ProxyAll tells antrea-agent to proxy all Service traffic. Default: false                                      |
| `antrea.config.antreaProxy.nodePortAddresses`    | Optional          | Specifies the host IPv4/IPv6 addresses for NodePort                                                           |
| `antrea.config.antreaProxy.skipServices`         | Optional          | List of Services which should be ignored by AntreaProxy                                                       |
| `antrea.config.antreaProxy.proxyLoadBalancerIPs` | Optional          | Load-balance traffic destined to the External IPs of LoadBalancer services. Default: false                    |
| `antrea.config.flowExporter.address`             | Optional          | Provide the IPFIX collector address as a string. Default: `flow-aggregator.flow-aggregator.svc:4739:tls`      |
| `antrea.config.flowExporter.pollInterval`        | Optional          | Provide flow poll interval as a duration string. Default: `5s`                                                |
| `antrea.config.flowExporter.activeTimeout`       | Optional          | Provide the active flow export timeout. Default: `30s`                                                        |
| `antrea.config.flowExporter.idleTimeout`         | Optional          | Provide the idle flow export timeout. Default: `15s`                                                          |
| `antrea.config.kubeAPIServerOverride`            | Optional          | Provide the address of Kubernetes apiserver. Default: nil                                                     |
| `antrea.config.transportInterface`               | Optional          | The name of the interface on Node which is used for tunneling or routing the traffic. Default: empty          |
| `antrea.config.transportInterfaceCIDRs`          | Optional          | The network CIDRs of the interface on Node which is used for tunneling or routing the traffic. Default: empty |
| `antrea.config.multicastInterfaces`              | Optional          | The names of the interfaces on Nodes that are used to forward multicast traffic. Default: nil                 |
| `antrea.config.trafficEncryptionMode`            | Optional          | Determines how tunnel traffic is encrypted. Default: none                                                     |
| `antrea.config.wireGuard.port`                   | Optional          | The port for WireGuard to receive traffic. Default: 51820                                                     |
| `antrea.config.enableUsageReporting`             | Optional          | Enable usage reporting (telemetry) to VMware. Default: false                                                  |
| `antrea.config.serviceCIDR`                      | Optional          | The service IPv4 CIDR to use. Default: `10.96.0.0/12`                                                         |
| `antrea.config.serviceCIDRv6`                    | Optional          | The service IPv6 CIDR to use. Default: nil                                                                    |
| `antrea.config.trafficEncapMode`                 | Optional          | The traffic encapsulation mode. Default: `encap`                                                              |
| `antrea.config.noSNAT`                           | Optional          | Boolean flag to enable/disable SNAT. Default: `false`                                                         |
| `antrea.config.disableUdpTunnelOffload`          | Optional          | Disable UDP tunnel offload feature on default NIC. Default: `false`                                           |
| `antrea.config.defaultMTU`                       | Optional          | MTU to use. Default: `null` (Antrea will autodetect)                                                          |
| `antrea.config.tlsCipherSuites`                  | Optional          | List of allowed cipher suites                                                                                 |
| `antrea.config.featureGates.AntreaProxy`         | Optional          | Boolean flag to enable/disable antrea proxy. Default: `true`                                                  |
| `antrea.config.featureGates.EndpointSlice`       | Optional          | Boolean flag to enable/disable EndpointSlice support in AntreaProxy. Default: `false`                         |
| `antrea.config.featureGates.AntreaTraceFlow`     | Optional          | Boolean flag to enable/disable antrea traceflow. Default: `false`                                             |
| `antrea.config.featureGates.NodePortLocal`       | Optional          | Boolean flag to enable/disable antrea proxy. Default: `false`                                                 |
| `antrea.config.featureGates.AntreaPolicy`        | Optional          | Boolean flag to enable/disable antrea policy. Default: `true`                                                 |
| `antrea.config.featureGates.FlowExporter`        | Optional          | Boolean flag to enable/disable flow exporter. Default: `false`                                                |
| `antrea.config.featureGates.NetworkPolicyStats`  | Optional          | Boolean flag to enable/disable network policy stats. Default: `false`                                         |
| `antrea.config.featureGates.Egress`              | Optional          | Boolean flag to enable/disable SNAT IPs of Pod egress traffic. Default: `false`                               |
| `antrea.config.featureGates.AntreaIPAM`          | Optional          | Boolean flag to enable/disable NodePortLocal feature to make the pods reachable externally through NodePort   |
| `antrea.config.featureGates.ServiceExternalIP`   | Optional          | Boolean flag to enable/disable NodePortLocal feature to make the pods reachable externally through NodePort   |
| `antrea.config.featureGates.Multicast`           | Optional          | Boolean flag to enable/disable NodePortLocal feature to make the pods reachable externally through NodePort   |

## Usage Example

The follow is a basic guide for getting started with antrea.
