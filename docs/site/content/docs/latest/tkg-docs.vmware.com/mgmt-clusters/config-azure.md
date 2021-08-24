# Management Cluster Configuration for Microsoft Azure

To create a cluster configuration file, you can copy an existing configuration file for a previous deployment to Azure and update it. Alternatively, you can create a file from scratch by using an empty template.

## Management Cluster Configuration Template

The template below includes all of the options that are relevant to deploying management clusters on Azure. You can copy this template and use it to deploy management clusters to Azure.

- For information about how to update the settings that are common to all infrastructure providers, see [Create a Management Cluster Configuration File](create-config-file.md)
- For information about all configuration file variables, see the [Tanzu CLI Configuration File Variable Reference](../tanzu-config-reference.md).
- For examples of how to configure the Azure settings, see the sections below the template.

Mandatory options are uncommented. Optional settings are commented out. Default values are included where applicable.

```
#! ---------------------------------------------------------------------
#! Basic cluster creation configuration
#! ---------------------------------------------------------------------

CLUSTER_NAME:
CLUSTER_PLAN: dev
INFRASTRUCTURE_PROVIDER: azure
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
#! Node configuration
#! ---------------------------------------------------------------------

# SIZE:
# CONTROLPLANE_SIZE:
# WORKER_SIZE:
# AZURE_CONTROL_PLANE_MACHINE_TYPE: "Standard_D2s_v3"
# AZURE_NODE_MACHINE_TYPE: "Standard_D2s_v3"
# OS_NAME: ""
# OS_VERSION: ""
# OS_ARCH: ""
# AZURE_CONTROL_PLANE_DATA_DISK_SIZE_GIB : ""
# AZURE_CONTROL_PLANE_OS_DISK_SIZE_GIB : ""
# AZURE_CONTROL_PLANE_MACHINE_TYPE : ""
# AZURE_CONTROL_PLANE_OS_DISK_STORAGE_ACCOUNT_TYPE : ""
# AZURE_ENABLE_NODE_DATA_DISK : ""
# AZURE_NODE_DATA_DISK_SIZE_GIB : ""
# AZURE_NODE_OS_DISK_SIZE_GIB : ""
# AZURE_NODE_MACHINE_TYPE : ""
# AZURE_NODE_OS_DISK_STORAGE_ACCOUNT_TYPE : ""

#! ---------------------------------------------------------------------
#! Azure configuration
#! ---------------------------------------------------------------------

AZURE_ENVIRONMENT: "AzurePublicCloud"
AZURE_TENANT_ID:
AZURE_SUBSCRIPTION_ID:
AZURE_CLIENT_ID:
AZURE_CLIENT_SECRET:
AZURE_LOCATION:
AZURE_SSH_PUBLIC_KEY_B64:
# AZURE_RESOURCE_GROUP: ""
# AZURE_VNET_RESOURCE_GROUP: ""
# AZURE_VNET_NAME: ""
# AZURE_VNET_CIDR: ""
# AZURE_CONTROL_PLANE_SUBNET_NAME: ""
# AZURE_CONTROL_PLANE_SUBNET_CIDR: ""
# AZURE_NODE_SUBNET_NAME: ""
# AZURE_NODE_SUBNET_CIDR: ""
# AZURE_CUSTOM_TAGS : ""
# AZURE_ENABLE_PRIVATE_CLUSTER : ""
# AZURE_FRONTEND_PRIVATE_IP : ""
# AZURE_ENABLE_ACCELERATED_NETWORKING : ""

#! ---------------------------------------------------------------------
#! Machine Health Check configuration
#! ---------------------------------------------------------------------
ENABLE_MHC: true
MHC_UNKNOWN_STATUS_TIMEOUT: 5m
MHC_FALSE_STATUS_TIMEOUT: 12m

#! ---------------------------------------------------------------------
#! Identity management configuration
#! ---------------------------------------------------------------------

IDENTITY_MANAGEMENT_TYPE: oidc

#! Settings for Pinniped OIDC
# CERT_DURATION: 2160h
# CERT_RENEW_BEFORE: 360h
# OIDC_IDENTITY_PROVIDER_ISSUER_URL:
# OIDC_IDENTITY_PROVIDER_CLIENT_ID:
# OIDC_IDENTITY_PROVIDER_CLIENT_SECRET:
# OIDC_IDENTITY_PROVIDER_SCOPES: "email,profile,groups"
# OIDC_IDENTITY_PROVIDER_USERNAME_CLAIM:
# OIDC_IDENTITY_PROVIDER_GROUPS_CLAIM:

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
#! Antrea CNI configuration
#! ---------------------------------------------------------------------
# ANTREA_NO_SNAT: false
# ANTREA_TRAFFIC_ENCAP_MODE: "encap"
# ANTREA_PROXY: false
# ANTREA_POLICY: true
# ANTREA_TRACEFLOW: false
```

## Azure Connection Settings

Specify information about your Azure account and the region in which you want to deploy the cluster.

For example:

```
AZURE_ENVIRONMENT: "AzurePublicCloud"
AZURE_TENANT_ID: b39138ca-[...]-d9dd62f0
AZURE_SUBSCRIPTION_ID: 3b511ccd-[...]-08a6d1a75d78
AZURE_CLIENT_ID: <encoded:M2ZkYTU4NGM[...]tZmViZjMxOGEyNmU1>
AZURE_CLIENT_SECRET: <encoded:bjVxLUpIUE[...]EN+d0RCd28wfg==>
AZURE_LOCATION: westeurope
AZURE_SSH_PUBLIC_KEY_B64: c3NoLXJzYSBBQUFBQjN[...]XJlLmNvbQ==
```

## <a id="node-sizes"></a> Configure Node Sizes

The Tanzu CLI creates the individual nodes of Tanzu Kubernetes clusters according to settings that you provide in the configuration file. On Amazon EC2, you can configure all node VMs to have the same predefined configurations or set different predefined configurations for control plane and worker nodes. By using these settings, you can create Tanzu Kubernetes clusters that have nodes with different configurations to the management cluster nodes. You can also create clusters in which the control plane nodes and worker nodes have different configurations.

When you created the management cluster, the instance types for the node machines are set in the `AZURE_CONTROL_PLANE_MACHINE_TYPE` and `AZURE_NODE_MACHINE_TYPE` options. By default, these settings are also used for Tanzu Kubernetes clusters. The minimum configuration is 2 CPUs and 8 GB memory. The list of compatible instance types varies in different regions.

```
AZURE_CONTROL_PLANE_MACHINE_TYPE: "Standard_D2s_v3"
AZURE_NODE_MACHINE_TYPE: "Standard_D2s_v3"
```

You can override these settings by using the `SIZE`, `CONTROLPLANE_SIZE` and `WORKER_SIZE` options. To create a Tanzu Kubernetes cluster in which all of the control plane and worker node VMs are the same size, specify the `SIZE` variable. If you set the `SIZE` variable, all nodes will be created with the configuration that you set. Set to `Standard_D2s_v3`, `Standard_D4s_v3`, and so on. For information about node instances for Azure, see [Sizes for virtual machines in Azure](https://docs.microsoft.com/en-us/azure/virtual-machines/sizes).

```
SIZE: Standard_D2s_v3
```

To create a Tanzu Kubernetes cluster in which the control plane and worker node VMs are different sizes, specify the `CONTROLPLANE_SIZE` and `WORKER_SIZE` options.

```
CONTROLPLANE_SIZE: Standard_D2s_v3
WORKER_SIZE: Standard_D4s_v3
```

You can combine the `CONTROLPLANE_SIZE` and `WORKER_SIZE` options with the `SIZE` option. For example, if you specify `SIZE: "Standard_D2s_v3"` with `WORKER_SIZE: "Standard_D4s_v3"`, the control plane nodes will be set to `Standard_D2s_v3` and worker nodes will be set to `Standard_D4s_v3`.

```
SIZE: Standard_D2s_v3
WORKER_SIZE: Standard_D4s_v3
```

## <a id="azure-subnet"></a> Deploy a Cluster with Custom Node Subnets

To specify custom subnets (IP ranges) for the nodes in a cluster, set variables as follows before you create the cluster.
You can define them as environment variables before you run `tanzu cluster create` or include them in the cluster configuration file that you pass in with the `--file` option.

To specify a custom subnet (IP range) for the **control plane node** in a cluster:

- **Subnet already defined in Azure**: Set `AZURE_CONTROL_PLANE_SUBNET_NAME` to the subnet name.
- **Create a new subnet**: Set `AZURE_CONTROL_PLANE_SUBNET_NAME` to name for the new subnet, and optionally set `AZURE_CONTROL_PLANE_SUBNET_CIDR` to a CIDR range within the configured Azure VNET.
    - If you omit `AZURE_CONTROL_PLANE_SUBNET_CIDR` a CIDR is generated automatically.

To specify a custom subnet for the **worker nodes** in a cluster, set environment variables `AZURE_NODE_SUBNET_NAME` and `AZURE_NODE_SUBNET_CIDR` following the same rules as for the control plane node, above.

For example:

```
AZURE_CONTROL_PLANE_SUBNET_CIDR: 10.0.0.0/24
AZURE_CONTROL_PLANE_SUBNET_NAME: my-cp-subnet
AZURE_NODE_SUBNET_CIDR: 10.0.1.0/24
AZURE_NODE_SUBNET_NAME: my-worker-subnet
AZURE_RESOURCE_GROUP: my-rg
AZURE_VNET_CIDR: 10.0.0.0/16
AZURE_VNET_NAME: my-vnet
AZURE_VNET_RESOURCE_GROUP: my-rg
```

## What to Do Next

After you have finished updating the management cluster configuration file, create the management cluster by following the instructions in [Deploy Management Clusters from a Configuration File](deploy-cli.md#mc-create).

**IMPORTANT**: If this is the first time that you are deploying a management cluster to Azure with a new version of Tanzu Kubernetes Grid, for example v1.3.1, make sure that you have accepted the base image license for that version. For information, see [Accept the Base Image License](azure.md#license) in *Prepare to Deploy Management Clusters to Microsoft Azure*.
