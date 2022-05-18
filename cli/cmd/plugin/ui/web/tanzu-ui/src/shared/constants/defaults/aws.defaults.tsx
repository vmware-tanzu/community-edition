export const AWS_DEFAULT_VALUES = {
    // Cluster Data
    CLUSTER_NAME: 'my-aws-cluster',
    CLUSTER_PLAN: 'dev',
    CLUSTER_DEV_NODE_TYPE: 't3a.large',
    CLUSTER_PROD_NODE_TYPE: 'm6a.xlarge',
    CLUSTER_WORKER_NODE_TYPE: 't3a.large',

    // VPC New
    // VPC_NAME: 'temp-vpc-name',
    VPC_CIDR: '10.0.0.0/16',

    // other?
    CREATE_CLOUDFORMATION_STACK: true,
    ENABLE_AUDIT_LOGGING: true,
    ENABLE_BASTION_HOST: true,
    ENABLE_CEIP_PARTICIPATION: true,
    ENABLE_MACHINE_HEALTH_CHECK: true,

    // Kubernetes Networking
    CLUSTER_SERVICE_CIDR: '100.64.0.0/13',
    CLUSTER_POD_CIDR: '100.96.0.0/11',

    // HTTP Proxy & Load Balancer
    HTTP_PROXY_ENABLED: false,
    LOAD_BALANCER_SCHEME_INTERNAL: false,
};
