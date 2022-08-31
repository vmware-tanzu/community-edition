export const AWS_MC_BASIC_TAB_NAMES = ['AWS Credentials', 'Cluster Settings', 'Review'];

export const enum AWS_FIELDS {
    CREDENTIAL_TYPE = 'credentialType',
    PROFILE = 'profile',
    SECRET_ACCESS_KEY = 'secretAccessKey',
    SESSION_TOKEN = 'sessionToken',
    ACCESS_KEY_ID = 'accessKeyId',
    REGION = 'region',
    CLUSTER_NAME = 'clusterName',
    NODE_PROFILE = 'nodeProfile',
    EC2_KEY_PAIR = 'ec2KeyPair',
    OS_IMAGE = 'osImage',
    CLUSTER_PLAN = 'clusterPlan',
    VPC_CIDR = 'vpcCidr',
    CREATE_CLOUDFORMATION_STACK = 'createCloudformationStack',
    ENABLE_AUDIT_LOGGING = 'enableAuditLogging',
    ENABLE_BASTION_HOST = 'enableBastionHost',
    ENABLE_CEIP_PARTICIPATION = 'enableCeipParticipation',
    ENABLE_MACHINE_HEALTH_CHECK = 'enableMachineHealthCheck',
    CLUSTER_NETWORKING_CNI_PROVIDER = 'clusterNetworkingCniProvider',
    CLUSTER_SERVICE_CIDR = 'clusterServiceCidr',
    CLUSTER_POD_CIDR = 'clusterPodCidr',
    HTTP_PROXY_ENABLED = 'httpProxyEnabled',
    LOAD_BALANCER_SCHEME_INTERNAL = 'loadBalancerSchemeInternal',
    NODE_TYPE = 'nodeType',
    AVAILABILITY_ZONES = 'availabilityZones',
    AVAILIABILITY_ZONE_NODE_TYPES = 'availabilityZoneNodeTypes',
    STORED_AZ_OBJECTS = 'storedAZObjects',
}

export const enum AWS_NODE_PROFILE_NAMES {
    SINGLE_NODE = 'SINGLE_NODE',
    HIGH_AVAILABILITY = 'HIGH_AVAILABILITY',
    PRODUCTION_READY = 'PRODUCTION_READY',
}

/* eslint-disable no-unused-vars */
export const enum CREDENTIAL_TYPE {
    PROFILE = 'PROFILE',
    ONE_TIME = 'ONE_TIME',
}
