import { DOCKER_FIELDS } from '../../../views/management-cluster/docker/docker-mc-basic/DockerManagementClusterBasic.constants';

export const DOCKER_DEFAULT_VALUES = {
    [DOCKER_FIELDS.CLUSTER_NAME]: 'my-docker-cluster',
    // Kubernetes Networking
    [DOCKER_FIELDS.CNI_TYPE]: 'antrea',
    [DOCKER_FIELDS.CLUSTER_SERVICE_CIDR]: '100.64.0.0/13',
    [DOCKER_FIELDS.CLUSTER_POD_CIDR]: '100.96.0.0/11',
};
