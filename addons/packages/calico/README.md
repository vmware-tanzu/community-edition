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
| `imageInfo.imageRepository` | Optional | The image repository to fetch the images from. |
| `imageInfo.imagePullPolicy` | Optinal | The image pull policy to use. |
| `imageInfo.images.calicoCniImage.imagePath` | Optional | The path of calico cni image in the repository. |
| `imageInfo.images.calicoCniImage.imageTag` | Optional | The tag of calico cni image. |
| `imageInfo.images.calicoKubecontrollerImage.imagePath` | Optional | The path of calico kube controller image in the repository. |
| `imageInfo.images.calicoKubecontrollerImage.imageTag` | Optional | The tag of calico kube controller image. |
| `imageInfo.images.calicoNodeImage.imagePath` | Optional | The path of calico node image in the repository. |
| `imageInfo.images.calicoNodeImage.imageTag` | Optinal | The tag of calico node image. |
| `imageInfo.images.calicoPodDaemonImage.imagePath ` | Optional | The path of calico pod daemon image in the repository. |
| `imageInfo.images.calicoPodDaemonImage.imageTag` | Optional | The tag of calico pod daemon image. |

### calico Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `calico.config.clusterCIDR` | Optional | The pod network CIDR. Default value is `192.168.0.0/16`. |
| `calico.config.vethMTU` | Optional | MTU size. Default is `1440`. |

## Usage Example

To learn more about how to use calico refer to [calico documentation](https://docs.projectcalico.org/about/about-calico)
