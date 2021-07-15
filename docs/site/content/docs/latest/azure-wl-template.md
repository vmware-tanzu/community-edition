# Azure Guest Cluster Template

When you deploy Tanzu Kubernetes (workload) clusters to Microsoft Azure, you must specify options in the cluster configuration file to connect to your Azure account and identify the resources that the cluster will use.

## <a id="nsg"></a> Create a Network Security Group for Each Cluster

Each workload cluster on Azure requires a Network Security Group (NSG) for its worker nodes named `CLUSTER-NAME-node-nsg`, where `CLUSTER-NAME` is the name of the cluster.

For more information, see [Network Security Groups on Azure](../mgmt-clusters/azure.md#nsgs).

## <a id="private"></a> Azure Private Clusters

By default, Azure management and workload clusters are public.
But you can also configure them to be private, which means their API server uses an Azure internal load balancer (ILB) and is therefore only accessible from within the cluster’s own VNET or peered VNETs.

To make an Azure cluster private, include the following in its configuration file:

* Set `AZURE_ENABLE_PRIVATE_CLUSTER` to `true`.

* (Optional) Set `AZURE_FRONTEND_PRIVATE_IP` to an internal address for the cluster's load balancer.

   - This address must be within the range of its control plane subnet and must not be used by another component.
   - If not set, this address defaults to `10.0.0.100`.

* Set `AZURE_VNET_NAME`, `AZURE_VNET_CIDR`, `AZURE_CONTROL_PLANE_SUBNET_NAME`, `AZURE_CONTROL_PLANE_SUBNET_CIDR`, `AZURE_NODE_SUBNET_NAME`, and `AZURE_NODE_SUBNET_CIDR` to the VNET and subnets that you use for other Azure private clusters.

   - Because Azure private clusters are not accessible outside their VNET, the management cluster and any workload and shared services clusters that it manages must be in the same private VNET.
   - The bootstrap machine, where you run the Tanzu CLI to create and use the private clusters, must also be in the same private VNET.

For more information, see [API Server Endpoint](https://capz.sigs.k8s.io/topics/api-server-endpoint.html) in the Cluster API Provider Azure documentation.

## <a id="template"></a> Tanzu Kubernetes Cluster Template

The template below includes all of the options that are relevant to deploying Tanzu Kubernetes clusters on Azure. You can copy this template and update it to deploy Tanzu Kubernetes clusters to Azure.

Mandatory options are uncommented. Optional settings are commented out. Default values are included where applicable.

The way in which you configure the variables for Tanzu Kubernetes clusters that are specific to Azure is identical for both management clusters and workload clusters. For information about how to configure the variables, see [Create a Management Cluster Configuration File](../mgmt-clusters/create-config-file.md) and [Management Cluster Configuration for Azure](../mgmt-clusters/config-azure.md). Options that are specific to workload clusters that are common to all infrastructure providers are described in [Deploy Tanzu Kubernetes Clusters](deploy.md).

```sh
#! ---------------------------------------------------------------------
#! Cluster creation basic configuration
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
# AZURE_CONTROL_PLANE_MACHINE_TYPE: "Standard_D2s_v3"
# AZURE_NODE_MACHINE_TYPE: "Standard_D2s_v3"
# CONTROL_PLANE_MACHINE_COUNT: 1
# WORKER_MACHINE_COUNT: 1
# WORKER_MACHINE_COUNT_0:
# WORKER_MACHINE_COUNT_1:
# WORKER_MACHINE_COUNT_2:
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
#! Azure Configuration
#! ---------------------------------------------------------------------

AZURE_ENVIRONMENT: "AzurePublicCloud"
AZURE_TENANT_ID:
AZURE_SUBSCRIPTION_ID:
AZURE_CLIENT_ID:
AZURE_CLIENT_SECRET:
AZURE_LOCATION:
AZURE_SSH_PUBLIC_KEY_B64:
# AZURE_CONTROL_PLANE_SUBNET_NAME: ""
# AZURE_CONTROL_PLANE_SUBNET_CIDR: ""
# AZURE_NODE_SUBNET_NAME: ""
# AZURE_NODE_SUBNET_CIDR: ""
# AZURE_RESOURCE_GROUP: ""
# AZURE_VNET_RESOURCE_GROUP: ""
# AZURE_VNET_NAME: ""
# AZURE_VNET_CIDR: ""
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


