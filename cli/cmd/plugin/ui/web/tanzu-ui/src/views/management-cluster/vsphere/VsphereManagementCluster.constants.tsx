/* eslint-disable no-unused-vars */
export enum VSPHERE_FIELDS {
    SERVERNAME = 'serverName',
    USERNAME = 'userName',
    PASSWORD = 'password',
    DATACENTER = 'datacenter',
    IPFAMILY = 'ipFamily',
    USETHUMBPRINT = 'useThumbprint',
    CLUSTERNAME = 'clusterName',
    INSTANCETYPE = 'instanceType',
    VMTEMPLATE = 'vmTemplate',
    SSHKEY = 'sshKey',
    CLUSTER_ENDPOINT = 'clusterEndpoint',
}

/* eslint-disable no-unused-vars */
export enum IPFAMILIES {
    IPv4 = 'ipv4',
    IPv6 = 'ipv6',
}

/* eslint-disable no-unused-vars */
export enum ENDPOINT_PROVIDER_IDS {
    KUBE_VIP = 'kube-vip',
    NSX_ADVANCED = 'nsx-advanced',
}

export const ENDPOINT_PROVIDERS = {
    [ENDPOINT_PROVIDER_IDS.KUBE_VIP]: 'Kube-vip',
    [ENDPOINT_PROVIDER_IDS.NSX_ADVANCED]: 'NSX Advanced',
};
