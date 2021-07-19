# Deploy Tanzu Kubernetes Clusters to vSphere

When you deploy Tanzu Kubernetes clusters to vSphere, you must specify options in the cluster configuration file to connect to vCenter Server and identify the vSphere resources that the cluster will use. You can also specify standard sizes for the control plane and worker node VMs, or configure the CPU, memory, and disk sizes for control plane and worker nodes explicitly. If you use custom image templates, you can identify which template to use to create node VMs.

For the full list of options that you must specify when deploying Tanzu Kubernetes clusters to vSphere, see the [Tanzu CLI Configuration File Variable Reference](../tanzu-config-reference.md).

## Tanzu Kubernetes Cluster Template

The template below includes all of the options that are relevant to deploying Tanzu Kubernetes clusters on vSphere. You can copy this template and update it to deploy Tanzu Kubernetes clusters to vSphere.

Mandatory options are uncommented. Optional settings are commented out. Default values are included where applicable. 

With the exception of the options described in the sections below the template, the way in which you configure the variables for Tanzu Kubernetes clusters that are specific to vSphere is identical for both management clusters and workload clusters. For information about how to configure the variables, see [Create a Management Cluster Configuration File from a Template](../mgmt-clusters/create-config-file.md) and [Management Cluster Configuration for vSphere](../mgmt-clusters/config-vsphere.md). Options that are specific to workload clusters that are common to all infrastructure providers are described in [Deploy Tanzu Kubernetes Clusters](deploy.md).

```
#! ---------------------------------------------------------------------
#! Basic cluster creation configuration
#! ---------------------------------------------------------------------

# CLUSTER_NAME:
CLUSTER_PLAN: dev
NAMESPACE: default
CNI: antrea

#! ---------------------------------------------------------------------
#! Node configuration
#! ---------------------------------------------------------------------

# SIZE: 
# CONTROLPLANE_SIZE:  
# WORKER_SIZE:

# VSPHERE_NUM_CPUS: 2
# VSPHERE_DISK_GIB: 40
# VSPHERE_MEM_MIB: 4096

# VSPHERE_CONTROL_PLANE_NUM_CPUS: 2
# VSPHERE_CONTROL_PLANE_DISK_GIB: 40
# VSPHERE_CONTROL_PLANE_MEM_MIB: 8192
# VSPHERE_WORKER_NUM_CPUS: 2
# VSPHERE_WORKER_DISK_GIB: 40
# VSPHERE_WORKER_MEM_MIB: 4096

# CONTROL_PLANE_MACHINE_COUNT: 1
# WORKER_MACHINE_COUNT: 1
# WORKER_MACHINE_COUNT_0:
# WORKER_MACHINE_COUNT_1:
# WORKER_MACHINE_COUNT_2:

#! ---------------------------------------------------------------------
#! vSphere configuration
#! ---------------------------------------------------------------------

VSPHERE_NETWORK: VM Network
# VSPHERE_TEMPLATE:
VSPHERE_SSH_AUTHORIZED_KEY:
VSPHERE_USERNAME:
VSPHERE_PASSWORD:
VSPHERE_SERVER:
VSPHERE_DATACENTER:
VSPHERE_RESOURCE_POOL:
VSPHERE_DATASTORE:
# VSPHERE_STORAGE_POLICY_ID
VSPHERE_FOLDER:
VSPHERE_TLS_THUMBPRINT:
VSPHERE_INSECURE: false
VSPHERE_CONTROL_PLANE_ENDPOINT:

#! ---------------------------------------------------------------------
#! NSX-T specific configuration for enabling NSX-T routable pods
#! ---------------------------------------------------------------------

# NSXT_POD_ROUTING_ENABLED: false
# NSXT_ROUTER_PATH: ""
# NSXT_USERNAME: ""
# NSXT_PASSWORD: ""
# NSXT_MANAGER_HOST: ""
# NSXT_ALLOW_UNVERIFIED_SSL: false
# NSXT_REMOTE_AUTH: false
# NSXT_VMC_ACCESS_TOKEN: ""
# NSXT_VMC_AUTH_HOST: ""
# NSXT_CLIENT_CERT_KEY_DATA: ""
# NSXT_CLIENT_CERT_DATA: ""
# NSXT_ROOT_CA_DATA: ""
# NSXT_SECRET_NAME: "cloud-provider-vsphere-nsxt-credentials"
# NSXT_SECRET_NAMESPACE: "kube-system"

#! ---------------------------------------------------------------------
#! Machine Health Check configuration
#! ---------------------------------------------------------------------

ENABLE_MHC: true
MHC_UNKNOWN_STATUS_TIMEOUT: 5m
MHC_FALSE_STATUS_TIMEOUT: 12m

#! ---------------------------------------------------------------------
#! Common configuration
#! ---------------------------------------------------------------------

# TKG_CUSTOM_IMAGE_REPOSITORY: ""
# TKG_CUSTOM_IMAGE_REPOSITORY_CA_CERTIFICATE: ""

# TKG_HTTP_PROXY: ""
# TKG_HTTPS_PROXY: ""
# TKG_NO_PROXY: ""

ENABLE_AUDIT_LOGGING: true

ENABLE_DEFAULT_STORAGE_CLASS: true

CLUSTER_CIDR: 100.96.0.0/11
SERVICE_CIDR: 100.64.0.0/13

# OS_NAME: ""
# OS_VERSION: ""
# OS_ARCH: ""

#! ---------------------------------------------------------------------
#! Autoscaler configuration
#! ---------------------------------------------------------------------

ENABLE_AUTOSCALER: false
# AUTOSCALER_MAX_NODES_TOTAL: "0"
# AUTOSCALER_SCALE_DOWN_DELAY_AFTER_ADD: "10m"
# AUTOSCALER_SCALE_DOWN_DELAY_AFTER_DELETE: "10s"
# AUTOSCALER_SCALE_DOWN_DELAY_AFTER_FAILURE: "3m"
# AUTOSCALER_SCALE_DOWN_UNNEEDED_TIME: "10m"
# AUTOSCALER_MAX_NODE_PROVISION_TIME: "15m"
# AUTOSCALER_MIN_SIZE_0:
# AUTOSCALER_MAX_SIZE_0:
# AUTOSCALER_MIN_SIZE_1:
# AUTOSCALER_MAX_SIZE_1:
# AUTOSCALER_MIN_SIZE_2:
# AUTOSCALER_MAX_SIZE_2:

#! ---------------------------------------------------------------------
#! Antrea CNI configuration
#! ---------------------------------------------------------------------
# ANTREA_NO_SNAT: false
# ANTREA_TRAFFIC_ENCAP_MODE: "encap"
# ANTREA_PROXY: false
# ANTREA_POLICY: true 
# ANTREA_TRACEFLOW: false
```

## <a id="custom-ova"></a> Deploy a Cluster with a Custom OVA Image

If you are using a single custom OVA image for each version of Kubernetes to deploy clusters on one operating system, follow [Deploy Tanzu Kubernetes Clusters with Different Kubernetes Versions](k8s-versions.md). In that procedure, you import the OVA into vSphere and then specify it for `tanzu cluster create` with the `--tkr` option.

If you are using multiple custom OVA images for the same Kubernetes version, then the `--tkr` value is ambiguous. This happens when the OVAs for the same Kubernetes version:

* Have different operating systems, for example created by `make build-node-ova-vsphere-ubuntu-1804`, `make build-node-ova-vsphere-photon-3`, and `make build-node-ova-vsphere-rhel-7`.
* Have the same name but reside in different vCenter folders.

To resolve this ambiguity, set the `VSPHERE_TEMPLATE` option to the desired OVA image before you run `tanzu cluster create`.

If the OVA template image name is unique, set `VSPHERE_TEMPLATE` to just the image name.

If multiple images share the same name, set `VSPHERE_TEMPLATE` to the full inventory path of the image in vCenter. This path follows the form `/MY-DC/vm/MY-FOLDER-PATH/MY-IMAGE`, where:

  - `MY_DC` is the datacenter containing the OVA template image
  - `MY_FOLDER_PATH` is the path to the image from the datacenter, as shown in the vCenter **VMs and Templates** view
  - `MY_IMAGE` is the image name

For example:

```
 VSPHERE_TEMPLATE: "/TKG_DC/vm/TKG_IMAGES/ubuntu-1804-kube-v1.18.8-vmware.1"
```

You can determine the image's full vCenter inventory path manually, or use the `govc` CLI:

  1. Install `govc`, for example with `brew install govc`
  1. Set environment variables for `govc` to access your vCenter:
      - `export GOVC_USERNAME=VCENTER-USERNAME`
      - `export GOVC_PASSWORD=VCENTER-PASSWORD`
      - `export GOVC_URL=VCENTER-URL`
      - `export GOVC_INSECURE=1`
  1. Run `govc find / -type m` and find the image name in the output, which lists objects by their complete inventory paths.

For more information about custom OVA images, see [Building Machine Images](../build-images/index.md).

## Configure DHCP Reservations for the Control Plane Nodes

After you deploy a cluster to vSphere, each control plane node requires a static IP address. This includes both management and Tanzu Kubernetes clusters. These static IP addresses are required in addition to the static IP address that you assigned to Kube-VIP when you deploy a managment cluster.

To make the IP addresses that your DHCP server assigned to the control plane nodes static, you can configure a DHCP reservation for each control plane node in the cluster. For instructions on how to configure DHCP reservations, see your DHCP server documentation.

## What to Do Next

Advanced options that are applicable to all infrastructure providers are described in the following topics:

- [Deploy Tanzu Kubernetes Clusters with Different Kubernetes Versions](k8s-versions.md)
- [Customize Tanzu Kubernetes Cluster Networking](networking.md)
- [Create Persistent Volumes with Storage Classes](storage.md)
- [Configure Tanzu Kubernetes Cluster Plans](config-plans.md)

After you have deployed your cluster, see [Managing Cluster Lifecycles](../cluster-lifecycle/index.md).