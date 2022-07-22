import { AZURE_NODE_PROFILE_NAMES } from '../../../views/management-cluster/azure/azure-mc-basic/AzureManagementClusterBasic.constants';

export const AZURE_DEFAULT_VALUES = {
    CLUSTER_NAME: '',
    NODE_PROFILE: AZURE_NODE_PROFILE_NAMES.SINGLE_NODE,
    CEIP_OPT_IN: false,
    CONTROL_PLANE_FLAVOR: 'dev',
    PRIVATE_AZURE_CLUSTER: false,
    ACTIVATE_AUDIT_LOGGING: false,
    MACHINE_HEALTH_CHECK_ENABLED: true,
    RESOURCE_GROUP: 'tanzu-resource-group-default',
    // VNET
    VNET_NAME: 'tanzu-vnet-name-default',
    VNET_CIDR: '10.0.0.0/16',
    CONTROL_PLANE_SUBNET: 'tanzu-control-plane-subnet-default',
    CONTROL_PLANE_SUBNET_CIDR: '10.0.0.0/24',
    WORKER_NODE_SUBNET: 'tanzu-worker-node-subnet-default',
    WORKER_NODE_SUBNET_CIDR: '10.0.1.0/24',
    // Network
    CNI_TYPE: 'antrea', // TODO: refactor to use CniProviders Const
    CLUSTER_POD_CIDR: '100.96.0.0/11',
    CLUSTER_SERVICE_CIDR: '100.64.0.0/13',
    ACTIVATE_PROXY_SETTINGS: false,
};

interface AzureDefaultInstanceTypes {
    [key: string]: string;
}

const AZURE_DEFAULT_INSTANCE_TYPES: AzureDefaultInstanceTypes = {
    [AZURE_NODE_PROFILE_NAMES.SINGLE_NODE]: 'Standard_D2s_v3',
    [AZURE_NODE_PROFILE_NAMES.HIGH_AVAILABILITY]: 'Standard_D2s_v3',
    [AZURE_NODE_PROFILE_NAMES.PRODUCTION_READY]: 'Standard_D4s_v3',
};

/**
 * @method retrieveAwsInstanceType
 * @param nodeProfile - node profile name set by ManagementClusterSettings.tsx; references key of AZURE_DEFAULT_INSTANCE_TYPES
 * defaults const.
 * Returns aws instance type string corresponding to selected node profile.
 */
export function retrieveAzureInstanceType(nodeProfile: string): string {
    if (nodeProfile && Object.keys(AZURE_DEFAULT_INSTANCE_TYPES).includes(nodeProfile)) {
        return AZURE_DEFAULT_INSTANCE_TYPES[nodeProfile as keyof AzureDefaultInstanceTypes];
    } else {
        console.warn(`provided node profile value not found in AZURE_DEFAULT_INSTANCE_TYPES: ${nodeProfile}`);
        return AZURE_DEFAULT_INSTANCE_TYPES[AZURE_NODE_PROFILE_NAMES.SINGLE_NODE];
    }
}
