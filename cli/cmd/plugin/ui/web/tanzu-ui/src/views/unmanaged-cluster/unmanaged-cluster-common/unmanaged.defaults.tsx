import { UNMANAGED_CLUSTER_FIELDS } from './UnmanagedCluster.constants';
import { K8sProviders } from '../../../shared/constants/K8sProviders.constants';
import { CniProviders } from '../../../shared/constants/CniProviders.constants';
import { ClusterProtocols } from '../../../shared/constants/ClusterProtocols.constants';

export const UNMANAGED_DEFAULT_VALUES = {
    // Cluster Settings advanced
    [UNMANAGED_CLUSTER_FIELDS.CONTROL_PLANE_NODE_COUNT]: 1,
    [UNMANAGED_CLUSTER_FIELDS.WORKER_NODE_COUNT]: 0,
    [UNMANAGED_CLUSTER_FIELDS.CLUSTER_PROVIDER]: K8sProviders.KIND,

    // Cluster Network Settings
    [UNMANAGED_CLUSTER_FIELDS.CNI_PROVIDER]: CniProviders.CALICO,
    [UNMANAGED_CLUSTER_FIELDS.CLUSTER_SERVICE_CIDR]: '100.64.0.0/13',
    [UNMANAGED_CLUSTER_FIELDS.CLUSTER_POD_CIDR]: '100.96.0.0/11',
    [UNMANAGED_CLUSTER_FIELDS.IP_ADDRESS]: '127.0.0.1',
    [UNMANAGED_CLUSTER_FIELDS.NODE_PORT_MAPPING]: '80',
    [UNMANAGED_CLUSTER_FIELDS.HOST_PORT_MAPPING]: '80',
    [UNMANAGED_CLUSTER_FIELDS.CLUSTER_PROTOCOL]: ClusterProtocols.TCP,
};

export const UNMANAGED_PLACEHOLDER_VALUES = {
    // Cluster Settings basic
    [UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME]: 'cluster-name',

    // Cluster Settings advanced
    [UNMANAGED_CLUSTER_FIELDS.CONTROL_PLANE_NODE_COUNT]: 'Control Plane Node Count',
    [UNMANAGED_CLUSTER_FIELDS.WORKER_NODE_COUNT]: 'Worker Node Count',

    // Cluster Network Settings
    [UNMANAGED_CLUSTER_FIELDS.CLUSTER_SERVICE_CIDR]: '100.64.0.0/13',
    [UNMANAGED_CLUSTER_FIELDS.CLUSTER_POD_CIDR]: '100.96.0.0/11',
    [UNMANAGED_CLUSTER_FIELDS.IP_ADDRESS]: '127.0.0.1',
    [UNMANAGED_CLUSTER_FIELDS.NODE_PORT_MAPPING]: '80',
    [UNMANAGED_CLUSTER_FIELDS.HOST_PORT_MAPPING]: '80',
};
