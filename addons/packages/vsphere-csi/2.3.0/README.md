# vSphere CSI

> This package provides cloud storage interface using vsphere-csi.

For more information, see the [GitHub page](https://github.com/kubernetes-sigs/vsphere-csi-driver) of vSphere CSI.

## Configuration

The following configuration values can be set to customize the vsphere CSI installation.

### Global

None

### vSphere CSI Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `vsphereCSI.namespace` | Required | The namespace in which vsphere-csi is deployed. Default value is `null`. |
| `vsphereCSI.clusterName` | Required | The name of the Kubernetes cluster in cluster ID. Default value is `null`. |
| `vsphereCSI.server` | Required | The IP address or FQDN of the vCenter endpoint. Default value is `null`. |
| `vsphereCSI.datacenter` | Required | The Datacenter in which VMs are located. Default value is `null`. |
| `vsphereCSI.publicNetwork` | Required | The public network to be used. Default value is `null`. |
| `vsphereCSI.username` | Required | vCenter username in clear text. Default value is `null`. |
| `vsphereCSI.password` | Required | vCenter password in clear text. Default value is `null`. |
| `vsphereCSI.tlsThumbprint` | Optional | The cryptographic thumbprint of the vSphere endpoint's certificate. Default value is `""`. |
| `vsphereCSI.provisionTimeout` | Optional | The timeout period for csi-provisioner container. Default value is `300s`. |
| `vsphereCSI.attachTimeout` | Optional | The timeout period for csi-attacher container. Default value is `300s`. |
| `vsphereCSI.resizerTimeout` | Optional | The timeout period for csi-resizer container. Default: `300s` |
| `vsphereCSI.vSphereVersion` | Optional | The vSphere version used. Default: `false`. |

## Usage Example

To learn more about how to use vSphere CSI refer to [vSphere CSI website](https://github.com/kubernetes-sigs/vsphere-csi-driver)
