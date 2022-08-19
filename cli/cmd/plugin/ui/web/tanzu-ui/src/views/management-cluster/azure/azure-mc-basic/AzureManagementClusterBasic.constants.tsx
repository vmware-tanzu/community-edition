export const AZURE_MC_BASIC_TAB_NAMES = ['Azure Credentials', 'Cluster Settings', 'Review'];

/* eslint-disable no-unused-vars */
export const enum AZURE_FIELDS {
    CLUSTER_NAME = 'clusterName',
    NODE_PROFILE = 'nodeProfile',
    TENANT_ID = 'tenantId',
    CLIENT_ID = 'clientId',
    CLIENT_SECRET = 'clientSecret',
    SUBSCRIPTION_ID = 'subscriptionId',
    AZURE_ENVIRONMENT = 'azureEnvironment',
    REGION = 'region',
    SSH_PUBLIC_KEY = 'sshPublicKey',
    OS_IMAGE = 'osImage',
    CEIP_OPT_IN = 'ceipOptIn',
    CONTROL_PLANE_FLAVOR = 'controlPlaneFlavor',
    PRIVATE_AZURE_CLUSTER = 'privateAzureCluster',
    ACTIVATE_AUDIT_LOGGING = 'activateAuditLogging',
    MACHINE_HEALTH_CHECK_ENABLED = 'machineHealthCheckEnabled',
    RESOURCE_GROUP = 'resourceGroup',
    VNET_NAME = 'vnetName',
    VNET_CIDR = 'vnetCidr',
    CONTROL_PLANE_SUBNET = 'controlPlaneSubnet',
    CONTROL_PLANE_SUBNET_CIDR = 'controlPlaneSubnetCidr',
    WORKER_NODE_SUBNET = 'workerNodeSubnet',
    WORKER_NODE_SUBNET_CIDR = 'workerNodeSubnetCidr',
    CNI_TYPE = 'cniType',
    CLUSTER_POD_CIDR = 'clusterPodCidr',
    CLUSTER_SERVICE_CIDR = 'clusterServiceCidr',
    ACTIVATE_PROXY_SETTINGS = 'activateProxySettings',
}

export const enum AZURE_NODE_PROFILE_NAMES {
    SINGLE_NODE = 'SINGLE_NODE',
    HIGH_AVAILABILITY = 'HIGH_AVAILABILITY',
    PRODUCTION_READY = 'PRODUCTION_READY',
}