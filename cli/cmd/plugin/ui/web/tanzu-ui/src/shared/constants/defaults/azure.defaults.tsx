import {
    AZURE_FIELDS,
    AZURE_NODE_PROFILE_NAMES,
} from '../../../views/management-cluster/azure/azure-mc-basic/AzureManagementClusterBasic.constants';
import { KeyOfStringToString } from '../../types/types';

export const AZURE_DEFAULT_VALUES = {
    [AZURE_FIELDS.NODE_PROFILE]: AZURE_NODE_PROFILE_NAMES.SINGLE_NODE,
    [AZURE_FIELDS.CEIP_OPT_IN]: false,
    [AZURE_FIELDS.CONTROL_PLANE_FLAVOR]: 'dev',
    [AZURE_FIELDS.PRIVATE_AZURE_CLUSTER]: false,
    [AZURE_FIELDS.ACTIVATE_AUDIT_LOGGING]: false,
    [AZURE_FIELDS.MACHINE_HEALTH_CHECK_ENABLED]: true,
    [AZURE_FIELDS.RESOURCE_GROUP]: 'tanzu-resource-group-default',
    // VNET
    [AZURE_FIELDS.VNET_NAME]: 'tanzu-vnet-name-default',
    [AZURE_FIELDS.VNET_CIDR]: '10.0.0.0/16',
    [AZURE_FIELDS.CONTROL_PLANE_SUBNET]: 'tanzu-control-plane-subnet-default',
    [AZURE_FIELDS.CONTROL_PLANE_SUBNET_CIDR]: '10.0.0.0/24',
    [AZURE_FIELDS.WORKER_NODE_SUBNET]: 'tanzu-worker-node-subnet-default',
    [AZURE_FIELDS.WORKER_NODE_SUBNET_CIDR]: '10.0.1.0/24',
    // Network
    [AZURE_FIELDS.CNI_TYPE]: 'antrea', // TODO: refactor to use CniProviders Const
    [AZURE_FIELDS.CLUSTER_POD_CIDR]: '100.96.0.0/11',
    [AZURE_FIELDS.CLUSTER_SERVICE_CIDR]: '100.64.0.0/13',
    [AZURE_FIELDS.ACTIVATE_PROXY_SETTINGS]: false,
};

const AZURE_DEFAULT_INSTANCE_TYPES: KeyOfStringToString = {
    [AZURE_NODE_PROFILE_NAMES.SINGLE_NODE]: 'Standard_D2s_v3',
    [AZURE_NODE_PROFILE_NAMES.HIGH_AVAILABILITY]: 'Standard_D2s_v3',
    [AZURE_NODE_PROFILE_NAMES.PRODUCTION_READY]: 'Standard_D4s_v3',
};

/**
 * @method retrieveAzureInstanceType
 * @param nodeProfile - node profile name set by ManagementClusterSettings.tsx; references key of AZURE_DEFAULT_INSTANCE_TYPES
 * defaults const.
 * Returns aws instance type string corresponding to selected node profile.
 */
export function retrieveAzureInstanceType(nodeProfile: string): string {
    if (nodeProfile && AZURE_DEFAULT_INSTANCE_TYPES[nodeProfile]) {
        return AZURE_DEFAULT_INSTANCE_TYPES[nodeProfile];
    } else {
        console.warn(`provided node profile value not found in AZURE_DEFAULT_INSTANCE_TYPES: ${nodeProfile}`);
        return AZURE_DEFAULT_INSTANCE_TYPES[AZURE_NODE_PROFILE_NAMES.SINGLE_NODE];
    }
}
