export const AWS_MC_BASIC_TAB_NAMES = ['AWS Credentials', 'Cluster Settings', 'Review'];

export const enum AWS_FIELDS {
    CREDENTIAL_TYPE = 'credentialType',
    PROFILE = 'PROFILE',
    SECRET_ACCESS_KEY = 'SECRET_ACCESS_KEY',
    SESSION_TOKEN = 'SESSION_TOKEN',
    ACCESS_KEY_ID = 'ACCESS_KEY_ID',
    REGION = 'REGION',
    CLUSTER_NAME = 'CLUSTER_NAME',
    NODE_PROFILE = 'NODE_PROFILE',
    EC2_KEY_PAIR = 'EC2_KEY_PAIR',
    IMAGE_INFO = 'IMAGE_INFO',
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
