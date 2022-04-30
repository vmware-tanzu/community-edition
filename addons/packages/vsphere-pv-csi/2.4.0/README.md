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
| `vspherePVCSI.namespace` | Required | The namespace where to deploy the pv csi k8s resources. Default value is `vmware-system-csi`. |
| `vspherePVCSI.supervisor_master_endpoint_hostname` | Required | The DNS hostname to reference the supervisor cluster. Default value is `supervisor.default.svc`. |
| `vspherePVCSI.supervisor_master_port` | Required | The IP port number through which to communicate with the supervisor cluster. Default value is `6443`. |
| `vspherePVCSI.cluster_uid` | Required | The unique id of the guest cluster. Default value is `null`. |
| `vspherePVCSI.cluster_name` | Required | The name of the guest cluster. Default value is `null`. |

## Usage Example

To learn more about how to use vSphere PV CSI refer to [vSphere CSI website](https://github.com/kubernetes-sigs/vsphere-csi-driver)
