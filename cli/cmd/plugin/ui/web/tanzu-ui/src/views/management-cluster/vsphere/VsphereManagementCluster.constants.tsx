/* eslint-disable no-unused-vars */
export enum VSPHERE_FIELDS {
    SERVERNAME = 'serverName',
    USERNAME = 'userName',
    PASSWORD = 'password',
    DATACENTER = 'datacenter',
    IPFAMILY = 'ipFamily',
    THUMBPRINT = 'thumbprint',
    USETHUMBPRINT = 'useThumbprint',
    CLUSTERNAME = 'CLUSTER_NAME',
    NODE_PROFILE_TYPE = 'profileType',
    VMTEMPLATE = 'vmTemplate',
    SSHKEY = 'sshKey',
    CLUSTER_ENDPOINT = 'clusterEndpoint',

    VMFolder = 'vmFolder',
    DataStore = 'datastore',
    Network = 'vSphereNetworkName',
    Pool = 'pool',

    CONTROL_PLANE_FLAVOR = 'controlPlaneFlavor',
    CONTROL_PLANE_INSTANCE_TYPE = 'controlPlaneInstanceType',
    WORKER_INSTANCE_TYPE = 'workerInstanceType',
    NUM_WORKER_NODES = 'numWorkerNodes',
    CNI_TYPE = 'cniType',
    MACHINE_HEALTH_CHECK_ACTIVATED = 'machineHealthCheckActivated',
    CEIP_OPT_IN = 'ceipOptIn',
    ENABLE_AUDIT_LOGGING = 'enableAuditLogging',
    CLUSTER_NODE_CIDR = 'clusterNodeCidr',
    CLUSTER_SERVICE_CIDR = 'clusterServiceCidr',
    CLUSTER_POD_CIDR = 'clusterPodCidr',
}

/* eslint-disable no-unused-vars */
export enum IP_FAMILIES {
    IPv4 = 'ipv4',
    IPv6 = 'ipv6',
}

export const IP_FAMILIES_DISPLAY = {
    [IP_FAMILIES.IPv4]: 'IP v4',
    [IP_FAMILIES.IPv6]: 'IP v6',
};

/* eslint-disable no-unused-vars */
export enum ENDPOINT_PROVIDERS {
    KUBE_VIP = 'kube-vip',
    NSX_ADVANCED = 'nsx-advanced',
}

export const ENDPOINT_PROVIDERS_DISPLAY = {
    [ENDPOINT_PROVIDERS.KUBE_VIP]: 'Kube-vip',
    [ENDPOINT_PROVIDERS.NSX_ADVANCED]: 'NSX Advanced',
};
