# vSphere PV CSI

> This package provides cloud storage interface using paravirtual vsphere-csi - usually used in guestclusters.

For more information, see the [GitHub page](https://github.com/kubernetes-sigs/vsphere-csi-driver) of vSphere CSI.

## Configuration

The following configuration values can be set to customize the vsphere PV CSI installation.

### Global

None

### vSphere CPI Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `vspherePVCSI.namespace` | Required | The namespace of the Kubernetes cluster in cluster ID. Default value is `kube-system`. |
| `vspherePVCSI.supervisor_master_endpoint_hostname` | Required | <TODO add description>. Default value is `null`. |
| `vspherePVCSI.supervisor_master_port` | Required | <TODO: add description>. Default value is `null`. |
| `vspherePVCSI.tanzukubernetescluster_uid` | Required | <TODO add description>. Default value is `null`. |
| `vspherePVCSI.tanzukubernetescluster_name` | Required | <TODO add description>. Default value is `null`. |

## Usage Example

To learn more about how to use vSphere PV CSI refer to [vSphere CSI website](https://github.com/kubernetes-sigs/vsphere-csi-driver)
