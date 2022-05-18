export interface ClusterClassDefinition {
    name: string;
    requiredVariables?: ClusterClassVariable[];
    basicVariables?: ClusterClassVariable[];
    intermediateVariables?: ClusterClassVariable[];
    advancedVariables?: ClusterClassVariable[];
}

export interface ClusterClassVariable {
    name: string;
    valueType: ClusterClassVariableType;
    description?: string;
    defaultValue?: string;
    possibleValues?: string[];
    required?: boolean;
}

export interface CCDefinition {
    name: string;
    categories?: CCCategory[];
}

export interface CCCategory {
    name: string;
    label: string;
    displayOpen: boolean;
    variables: CCVariable[];
}

export interface CCVariable {
    name: string;
    taxonomy: ClusterClassVariableType; // field classified according to known types
    description?: string;
    default?: any;
    required?: boolean;
    possibleValues?: string[];
    children?: CCVariable[];
}

// NOTE: the string values are literal used in cluster classes; do not change for fun
// for some reason, eslint is reporting these enum values as unused
/* eslint-disable no-unused-vars */
export enum ClusterClassVariableType {
    UNKNOWN = '', // this means the type of this var is unknown, and we generally treat it as a string
    BOOLEAN = 'boolean',
    CIDR = 'cidr',
    INTEGER = 'int',
    IP = 'ip',
    IP_LIST = 'ipList',
    NUMBER = 'number',
    STRING = 'string',
    STRING_K8S_COMPLIANT = 'stringK8sCompliant',
    STRING_PARAGRAPH = 'stringParagraph',

    GROUP = 'group',
    GROUP_OPTIONAL = 'groupOptional',

    PROXY = 'proxy', // httpProxy, httpsProxy, noProxy
    IMAGE_REPOSITORY = 'imageRepo', // host, tlsCertificationValidate
    TAINTS = 'taints', // array of {key, value, effect}
    MAP = 'map', // key/value pairs, cf V nodePoolTaints
    TRUST = 'trust', // array of values with name/data, first is proxy, second imageRepository
    VSPHERE_TRUST = 'vSphereTrust', // vSphere has property additionalTrustedCAs with {name, data}
    PORT = 'port',
    VSPHERE_NETWORK_INTERFACE = 'vSphereNetworkInterface', // eth0 is an example. Is there an enumeration of valid values?
    VSPHERE_NODE_TOPOLOGY = 'vSphereNodeTopology', // count+machine object {diskGiB, memoryMiB, numCPUs} for controlPlane, worker fields
    VSPHERE_SSH_AUTH_KEY = 'vSphereSshAuthKey', //
    VSPHERE_VOLUME = 'vSphereVolume', // array: mountPath, name, capacity {storage: string},
    // used in (VSPHERE) controlPlaneVolumes, nodePoolVolumes
    AWS_IDENTITY = 'awsIdentity',
    ASW_SECURITY_GROUP = 'awsSecurityGroup',
    AWS_REGION = 'awsRegion',
    AWS_SUBNETS = 'awsSubnets', // array of {private:{cidr, id}, public:{cidr, id}, az}
    AWS_VPC = 'awsVpc',
    AWS_NODE_TOPOLOGY = 'awsNodeTopology', // {instanceType, rootVolume{sizeGiB}}
    AZURE_LOCATION = 'azureLocation', // value has to be legal
    AZURE_VNET = 'azureVnet', // {cidr, name, resourceGroup}
    AZURE_IDENTITY = 'azureIdentity', // {name, namespace}
    AZURE_NODE_TOPOLOGY = 'azureNodeTopology', // machineType, osDisk, dataDisks, outboundLB, subnet
}
