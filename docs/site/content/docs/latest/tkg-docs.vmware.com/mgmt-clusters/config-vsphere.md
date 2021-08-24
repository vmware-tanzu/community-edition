# Management Cluster Configuration for vSphere

To create a cluster configuration file, you can copy an existing configuration file for a previous deployment to vSphere and update it. Alternatively, you can create a file from scratch by using an empty template.

## Management Cluster Configuration Template

The template below includes all of the options that are relevant to deploying management clusters on vSphere. You can copy this template and use it to deploy management clusters to vSphere.

- For information about how to update the settings that are common to all infrastructure providers, see [Create a Management Cluster Configuration File](create-config-file.md)
- For information about all configuration file variables, see the [Tanzu CLI Configuration File Variable Reference](../tanzu-config-reference.md).
- For examples of how to configure the vSphere settings, see the sections below the template.

Mandatory options are uncommented. Optional settings are commented out. Default values are included where applicable.

```
#! ---------------------------------------------------------------------
#! Basic cluster creation configuration
#! ---------------------------------------------------------------------

CLUSTER_NAME:
CLUSTER_PLAN: dev
INFRASTRUCTURE_PROVIDER: vsphere
ENABLE_CEIP_PARTICIPATION: true
# TMC_REGISTRATION_URL:
ENABLE_AUDIT_LOGGING: true
CLUSTER_CIDR: 100.96.0.0/11
SERVICE_CIDR: 100.64.0.0/13

#! ---------------------------------------------------------------------
#! Image repository configuration
#! ---------------------------------------------------------------------

# TKG_CUSTOM_IMAGE_REPOSITORY: ""
# TKG_CUSTOM_IMAGE_REPOSITORY_CA_CERTIFICATE: ""

#! ---------------------------------------------------------------------
#! Proxy configuration
#! ---------------------------------------------------------------------

# TKG_HTTP_PROXY: ""
# TKG_HTTPS_PROXY: ""
# TKG_NO_PROXY: ""

#! ---------------------------------------------------------------------
#! vSphere configuration
#! ---------------------------------------------------------------------

VSPHERE_SERVER:
VSPHERE_USERNAME:
VSPHERE_PASSWORD:
VSPHERE_DATACENTER:
VSPHERE_RESOURCE_POOL:
VSPHERE_DATASTORE:
VSPHERE_FOLDER:
VSPHERE_NETWORK: VM Network
VSPHERE_CONTROL_PLANE_ENDPOINT:
VIP_NETWORK_INTERFACE: "eth0"
# VSPHERE_TEMPLATE:
VSPHERE_SSH_AUTHORIZED_KEY:
# VSPHERE_STORAGE_POLICY_ID: ""
VSPHERE_TLS_THUMBPRINT:
VSPHERE_INSECURE: false
DEPLOY_TKG_ON_VSPHERE7: false
ENABLE_TKGS_ON_VSPHERE7: false

#! ---------------------------------------------------------------------
#! Node configuration
#! ---------------------------------------------------------------------

# SIZE:
# CONTROLPLANE_SIZE:
# WORKER_SIZE:
# OS_NAME: ""
# OS_VERSION: ""
# OS_ARCH: ""
# VSPHERE_NUM_CPUS: 2
# VSPHERE_DISK_GIB: 40
# VSPHERE_MEM_MIB: 4096
# VSPHERE_CONTROL_PLANE_NUM_CPUS: 2
# VSPHERE_CONTROL_PLANE_DISK_GIB: 40
# VSPHERE_CONTROL_PLANE_MEM_MIB: 8192
# VSPHERE_WORKER_NUM_CPUS: 2
# VSPHERE_WORKER_DISK_GIB: 40
# VSPHERE_WORKER_MEM_MIB: 4096

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
#! Identity management configuration
#! ---------------------------------------------------------------------

IDENTITY_MANAGEMENT_TYPE: "oidc"

#! Settings for OIDC
# CERT_DURATION: 2160h
# CERT_RENEW_BEFORE: 360h
# OIDC_IDENTITY_PROVIDER_CLIENT_ID:
# OIDC_IDENTITY_PROVIDER_CLIENT_SECRET:
# OIDC_IDENTITY_PROVIDER_GROUPS_CLAIM: groups
# OIDC_IDENTITY_PROVIDER_ISSUER_URL:
# OIDC_IDENTITY_PROVIDER_SCOPES: email
# OIDC_IDENTITY_PROVIDER_USERNAME_CLAIM: email

#! The following two variables are used to configure Pinniped JWTAuthenticator for workload clusters
# SUPERVISOR_ISSUER_URL:
# SUPERVISOR_ISSUER_CA_BUNDLE_DATA:

#! Settings for LDAP
# LDAP_BIND_DN:
# LDAP_BIND_PASSWORD:
# LDAP_HOST:
# LDAP_USER_SEARCH_BASE_DN:
# LDAP_USER_SEARCH_FILTER:
# LDAP_USER_SEARCH_USERNAME: userPrincipalName
# LDAP_USER_SEARCH_ID_ATTRIBUTE: DN
# LDAP_USER_SEARCH_EMAIL_ATTRIBUTE: DN
# LDAP_USER_SEARCH_NAME_ATTRIBUTE:
# LDAP_GROUP_SEARCH_BASE_DN:
# LDAP_GROUP_SEARCH_FILTER:
# LDAP_GROUP_SEARCH_USER_ATTRIBUTE: DN
# LDAP_GROUP_SEARCH_GROUP_ATTRIBUTE:
# LDAP_GROUP_SEARCH_NAME_ATTRIBUTE: cn
# LDAP_ROOT_CA_DATA_B64:

#! ---------------------------------------------------------------------
#! NSX Advanced Load Balancer configuration
#! ---------------------------------------------------------------------

AVI_ENABLE: false
# AVI_NAMESPACE: "tkg-system-networking"
# AVI_DISABLE_INGRESS_CLASS: true
# AVI_AKO_IMAGE_PULL_POLICY: IfNotPresent
# AVI_ADMIN_CREDENTIAL_NAME: avi-controller-credentials
# AVI_CA_NAME: avi-controller-ca
# AVI_CONTROLLER:
# AVI_USERNAME: ""
# AVI_PASSWORD: ""
# AVI_CLOUD_NAME:
# AVI_SERVICE_ENGINE_GROUP:
# AVI_DATA_NETWORK:
# AVI_DATA_NETWORK_CIDR:
# AVI_CA_DATA_B64: ""
# AVI_LABELS: ""
# AVI_INGRESS_DEFAULT_INGRESS_CONTROLLER: false
# AVI_INGRESS_SHARD_VS_SIZE: ""
# AVI_INGRESS_SERVICE_TYPE: ""

#! ---------------------------------------------------------------------
#! Antrea CNI configuration
#! ---------------------------------------------------------------------

# ANTREA_NO_SNAT: false
# ANTREA_TRAFFIC_ENCAP_MODE: "encap"
# ANTREA_PROXY: false
# ANTREA_POLICY: true
# ANTREA_TRACEFLOW: false
```

## <a id="general-vsphere"></a> General vSphere Configuration

Provide information to allow Tanzu Kubernetes Grid to log in to vSphere, and to designate the resources for Tanzu Kubernetes Grid can use.

- Update the `VSPHERE_SERVER`, `VSPHERE_USERNAME`, and `VSPHERE_PASSWORD` settings with the IP address or FQDN of the vCenter Server instance and the credentials to use to log in.
- Provide the full paths to the vSphere datacenter, resource pool, datastores, and folder in which to deploy the management cluster:

   - `VSPHERE_DATACENTER`: `/<MY-DATACENTER>`
   - `VSPHERE_RESOURCE_POOL`: `/<MY-DATACENTER>/host/<CLUSTER>/Resources`
   - `VSPHERE_DATASTORE`: `/<MY-DATACENTER>/datastore/<MY-DATASTORE>`
   - `VSPHERE_FOLDER`: `/<MY-DATACENTER>/vm/<FOLDER>.`
- Set a static virtual IP address for API requests to the Tanzu Kubernetes cluster in the `VSPHERE_CONTROL_PLANE_ENDPOINT` setting. If you mapped a fully qualified domain name (FQDN) to the VIP address, you can specify the FQDN instead of the VIP address.
- Specify a network and a network interface in `VSPHERE_NETWORK` and `VIP_NETWORK_INTERFACE`.
- Optionally uncomment and update `VSPHERE_TEMPLATE` to specify the path to an OVA file if you are using multiple custom OVA images for the same Kubernetes version. Use the format <code>/MY-DC/vm/MY-FOLDER-PATH/MY-IMAGE</code>. For more information, see [Deploy a Cluster with a Custom OVA Image](../tanzu-k8s-clusters/vsphere.md#custom-ova).
- Provide your SSH key in the `VSPHERE_SSH_AUTHORIZED_KEY` option. For information about how to obtain an SSH key, see [Prepare to Deploy Management Clusters to vSphere](vsphere.md).
- Provide the TLS thumbprint in the `VSPHERE_TLS_THUMBPRINT` variable, or set `VSPHERE_INSECURE: true` to skip thumbprint verification.
- Optionally uncomment `VSPHERE_STORAGE_POLICY_ID` and provide the name of a VM storage policy for the management cluster to use.

For example:

```
#! ---------------------------------------------------------------------
#! vSphere configuration
#! ---------------------------------------------------------------------

VSPHERE_SERVER: 10.185.12.154
VSPHERE_USERNAME: tkg-user@vsphere.local
VSPHERE_PASSWORD: <encoded:QWRtaW4hMjM=>
VSPHERE_DATACENTER: /dc0
VSPHERE_RESOURCE_POOL: /dc0/host/cluster0/Resources/tanzu
VSPHERE_DATASTORE: /dc0/datastore/sharedVmfs-1
VSPHERE_FOLDER: /dc0/vm/tanzu
VSPHERE_NETWORK: "VM Network"
VSPHERE_CONTROL_PLANE_ENDPOINT: 10.185.11.134
VIP_NETWORK_INTERFACE: "eth0"
VSPHERE_TEMPLATE: /dc0/vm/tanzu/my-image.ova
VSPHERE_SSH_AUTHORIZED_KEY: ssh-rsa AAAAB3[...]tyaw== user@example.com
VSPHERE_TLS_THUMBPRINT: 47:F5:83:8E:5D:36:[...]:72:5A:89:7D:29:E5:DA
VSPHERE_INSECURE: false
VSPHERE_STORAGE_POLICY_ID: "My storage policy"
```

## <a id="node-sizes"></a> Configure Node Sizes

The Tanzu CLI creates the individual nodes of management clusters and Tanzu Kubernetes clusters according to settings that you provide in the configuration file. On vSphere, you can configure all node VMs to have the same predefined configurations, set different predefined configurations for control plane and worker nodes, or customize the configurations of the nodes. By using these settings, you can create clusters that have nodes with different configurations to the management cluster nodes. You can also create clusters in which the control plane nodes and worker nodes have different configurations.

### Use Predefined Node Configurations

The Tanzu CLI provides the following predefined configurations for cluster nodes:

- `small`: 2 CPUs, 4 GB memory, 20 GB disk
- `medium`: 2 CPUs, 8 GB memory, 40 GB disk
- `large`: 4 CPUs, 16 GB memory, 40 GB disk
- `extra-large`: 8 CPUs, 32 GB memory, 80 GB disk

To create a cluster in which all of the control plane and worker node VMs are the same size, specify the `SIZE` variable. If you set the `SIZE` variable, all nodes will be created with the configuration that you set.

```
SIZE: "large"
```

To create a in which the control plane and worker node VMs are different sizes, specify the `CONTROLPLANE_SIZE` and `WORKER_SIZE` options.

```
CONTROLPLANE_SIZE: "medium"
WORKER_SIZE: "extra-large"
```

You can combine the `CONTROLPLANE_SIZE` and `WORKER_SIZE` options with the `SIZE` option. For example, if you specify `SIZE: "large"` with `WORKER_SIZE: "extra-large"`, the control plane nodes will be set to `large` and worker nodes will be set to `extra-large`.

```
SIZE: "large"
WORKER_SIZE: "extra-large"
```

### Define Custom Node Configurations

You can customize the configuration of the nodes rather than using the predefined configurations.

To use the same custom configuration for all nodes, specify the `VSPHERE_NUM_CPUS`, `VSPHERE_DISK_GIB`, and `VSPHERE_MEM_MIB` options.

```
VSPHERE_NUM_CPUS: 2
VSPHERE_DISK_GIB: 40
VSPHERE_MEM_MIB: 4096
```

To define different custom configurations for control plane nodes and worker nodes, specify the `VSPHERE_CONTROL_PLANE_*` and `VSPHERE_WORKER_*` options.

```
VSPHERE_CONTROL_PLANE_NUM_CPUS: 2
VSPHERE_CONTROL_PLANE_DISK_GIB: 20
VSPHERE_CONTROL_PLANE_MEM_MIB: 8192
VSPHERE_WORKER_NUM_CPUS: 4
VSPHERE_WORKER_DISK_GIB: 40
VSPHERE_WORKER_MEM_MIB: 4096
```

You can override these settings by using the `SIZE`, `CONTROLPLANE_SIZE`, and `WORKER_SIZE` options.

## Configure NSX Advanced Load Balancer

VMware NSX Advanced Load Balancer provides an L4+L7 load balancing solution for vSphere. NSX Advanced Load Balancer includes a Kubernetes operator that integrates with the Kubernetes API to manage the lifecycle of load balancing and ingress resources for workloads. To use NSX Advanced Load Balancer, you must first deploy it in your vSphere environment. For information, see [Install VMware NSX Advanced Load Balancer on a vSphere Distributed Switch](install-nsx-adv-lb.md).

You can configure Tanzu Kubernetes Grid to use NSX Advanced Load Balancer. By default, the management cluster and all workload clusters that it manages will use the load balancer. For information about how to configure the NSX Advanced Load Balancer variables, see [NSX Advanced Load Balancer](../tanzu-config-reference.md#nsx-lb) in the *Tanzu CLI Configuration File Variable Reference*.

```
AVI_ENABLE: true
AVI_NAMESPACE: "tkg-system-networking"
AVI_DISABLE_INGRESS_CLASS: true
AVI_AKO_IMAGE_PULL_POLICY: IfNotPresent
AVI_ADMIN_CREDENTIAL_NAME: avi-controller-credentials
AVI_USERNAME: avi-controller-ca
AVI_CONTROLLER: 10.185.10.217
AVI_USERNAME: "admin"
AVI_PASSWORD: "<password>"
AVI_CLOUD_NAME: "Default-Cloud"
AVI_SERVICE_ENGINE_GROUP: "Default-Group"
AVI_DATA_NETWORK: nsx-alb-dvswitch
AVI_DATA_NETWORK_CIDR: 10.185.0.0/20
AVI_CA_DATA_B64: LS0tLS1CRU[...]UtLS0tLQo=
AVI_LABELS: ""
# AVI_INGRESS_DEFAULT_INGRESS_CONTROLLER: false
# AVI_INGRESS_SHARD_VS_SIZE: ""
# AVI_INGRESS_SERVICE_TYPE: ""
```

## Configure NSX-T Routable Pods

If your vSphere environment uses NSX-T, you can configure it to implement routable, or `NO_NAT`, pods.

**NOTE**: NSX-T Routable Pods is an experimental feature in this release. Information about how to implement NSX-T Routable Pods will be added to this documentation soon.

```
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
```

## <a id="vsphere-7"></a>Management Clusters on vSphere with Tanzu

On vSphere 7, the vSphere with Tanzu option includes a built-in supervisor cluster that works as a management cluster and provides a better experience than a separate management cluster deployed by Tanzu Kubernetes Grid.  Deploying a Tanzu Kubernetes Grid management cluster to vSphere 7 when vSphere with Tanzu is not enabled is supported, but the preferred option is to enable vSphere with Tanzu and use the Supervisor Cluster. VMware Cloud on AWS and Azure VMware Solution do not support a supervisor cluster, so you need to deploy a management cluster. For information, see [Use the Tanzu CLI with a vSphere with Tanzu Supervisor Cluster](../tanzu-k8s-clusters/connect-vsphere7.md).

To reflect the recommendation for using the vSphere with Tanzu supervisor cluster as a management cluster, the Tanzu CLI behaves as follows, controlled by the `DEPLOY_TKG_ON_VSPHERE7` and `ENABLE_TKGS_ON_VSPHERE7` configuration parameters.

- vSphere with Tanzu enabled:
    - `ENABLE_TKGS_ON_VSPHERE7: false` informs you that deploying a management cluster is not possible, and exits.
    - `ENABLE_TKGS_ON_VSPHERE7: true` opens the vSphere Client at the address set by `VSPHERE_SERVER` in your `config.yml` or local environment, so you can configure your supervisor cluster as described in [Enable the Workload Management Platform with the vSphere Networking Stack](https://docs-staging.vmware.com/en/draft/VMware-vSphere/7.0/vmware-vsphere-with-tanzu/GUID-8D7D292B-43E9-4CB8-9E20-E4039B80BF9B.html) in the vSphere documentation.
- vSphere with Tanzu not enabled:
    - `DEPLOY_TKG_ON_VSPHERE7: false` informs you that deploying a Tanzu Kubernetes Grid management cluster is possible but not recommended, and prompts you to either quit the installation or continue to deploy the management cluster.
    - `DEPLOY_TKG_ON_VSPHERE7: true` deploys a TKG management cluster on vSphere 7, against recommendation for vSphere 7, but as required for VMware Cloud on AWS and Azure VMware Solution.

## What to Do Next

After you have finished updating the management cluster configuration file, create the management cluster by following the instructions in [Deploy Management Clusters from a Configuration File](deploy-cli.md#mc-create).
