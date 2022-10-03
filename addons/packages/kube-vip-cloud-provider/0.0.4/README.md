# kube-vip-cloud-provider

This package provides NSX Advanced Load Balancer.

## Components

* kube-vip-cloud-provider works with kube-vip to provide L4 load balancing for applications deployed in a kubernetes cluster for north-south traffic.
* kube-vip-cloud-provider itself is responsible for allocating IP address for load balancer type service. While kube-vip is responsible for advertising IP address for service.

## Installation of kube-vip-cloud-provider

kube-vip-cloud-provider(KVCP) supports different cloud infrastructure and we use vsphere as example in the following setup.
Please refer to the following documents if you are using a different cloud infrastructure.

[Upstream reference page](https://github.com/kube-vip/kube-vip-cloud-provider#installing-the-kube-vip-cloud-provider)

### kube-vip-cloud-provider Configuration

Either `loadbalancerCIDRs` or `loadbalancerIPRanges` is required to be set. Not required to set at the same time

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `loadbalancerCIDRs`    | Optional          | a list of comma separated cidrs will be used to allocate IP for external load balancer. Example 192.168.0.200/29,192.168.1.200/29                |
| `loadbalancerIPRanges` | Optinoal          | a list of comma separated cidrs will be used to allocate IP for external load balancer. Example 192.168.0.10-192.168.0.11,192.168.0.10-192.168.0.13       |
