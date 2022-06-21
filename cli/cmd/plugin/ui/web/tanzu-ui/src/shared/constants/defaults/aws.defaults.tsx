export const AWS_DEFAULT_VALUES = {
    // Cluster Data
    CLUSTER_NAME: 'my-aws-cluster',
    CLUSTER_PLAN: 'dev',

    // VPC New
    // VPC_NAME: 'temp-vpc-name',
    VPC_CIDR: '10.0.0.0/16',

    // other?
    CREATE_CLOUDFORMATION_STACK: true,
    ENABLE_AUDIT_LOGGING: true,
    ENABLE_BASTION_HOST: true,
    ENABLE_CEIP_PARTICIPATION: false,
    ENABLE_MACHINE_HEALTH_CHECK: true,

    // Kubernetes Networking
    CLUSTER_NETWORKING_CNI_PROVIDER: 'antrea',
    CLUSTER_SERVICE_CIDR: '100.64.0.0/13',
    CLUSTER_POD_CIDR: '100.96.0.0/11',

    // HTTP Proxy & Load Balancer
    HTTP_PROXY_ENABLED: false,
    LOAD_BALANCER_SCHEME_INTERNAL: false,
};

const AWS_DEFAULT_INSTANCE_TYPES = new Map<string, string>([
    ['SINGLE_NODE', 't3a.large'],
    ['HIGH_AVAILABILITY', 't3a.large'],
    ['PRODUCTION_READY', 'm6a.xlarge'],
]);

/**
 * @method retrieveAwsInstanceType
 * @param nodeProfile - node profile name set by ManagementClusterSettings.tsx; references key of AWS_DEFAULT_INSTANCE_TYPES
 * defaults map.
 * Returns aws instance type string corresponding to selected node profile.
 */
export function retrieveAwsInstanceType(nodeProfile: string): string {
    if (nodeProfile && AWS_DEFAULT_INSTANCE_TYPES.has(nodeProfile)) {
        return AWS_DEFAULT_INSTANCE_TYPES.get(nodeProfile) as string;
    } else {
        console.warn(`provided node profile value not found in AWS_DEFAULT_INSTANCE_TYPES: ${nodeProfile}`);
        return AWS_DEFAULT_INSTANCE_TYPES.get('SINGLE_NODE') as string;
    }
}
