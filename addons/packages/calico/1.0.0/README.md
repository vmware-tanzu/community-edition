# calico Package

This package provides networking and network security solution for containers using [calico](https://www.projectcalico.org/).

## Components

## Configuration

The following configuration values can be set to customize the calico installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Optional | The namespace in which to deploy calico. |
| `infraProvider` | Required | The cloud provider in use. One of: `aws`, `azure`, `vsphere`, `docker`. |
| `ipFamily` | Optional | The IP family calico should be configured with. Defaults to `ipv4` One of: `ipv4`, `ipv6` |

### calico Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `calico.config.clusterCIDR` | Optional | The pod network CIDR. Default value is `192.168.0.0/16`. |
| `calico.config.vethMTU` | Optional | MTU size. Default is `1440`. |

## Usage Example

To learn more about how to use calico refer to [calico documentation](https://docs.projectcalico.org/about/about-calico)
