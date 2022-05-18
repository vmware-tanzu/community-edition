export const DOCKER_DEFAULT_VALUES = {
    CLUSTER_NAME: 'my-docker-cluster',
    // Kubernetes Networking
    CNI_TYPE: 'antrea',
    CLUSTER_SERVICE_CIDR: '100.64.0.0/13',
    CLUSTER_POD_CIDR: '100.96.0.0/11',
};