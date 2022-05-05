import { ClusterClassDefinition, ClusterClassVariableType } from '../../shared/models/ClusterClass';
////////////////////////////////////////////////////////////////////////////////////////////////////
// COMMON ccVar objects
//
const ccVarTKG_CUSTOM_IMAGE_REPOSITORY = { name: 'TKG_CUSTOM_IMAGE_REPOSITORY', valueType: ClusterClassVariableType.STRING,
    description: 'custom image repository for TKG images'
}
const ccVarTKG_CUSTOM_IMAGE_REPOSITORY_SKIP_TLS_VERIFY = { name: 'TKG_CUSTOM_IMAGE_REPOSITORY_SKIP_TLS_VERIFY',
    valueType: ClusterClassVariableType.BOOLEAN,
    description: 'skip TLS verfication on custom image repository'
}
const ccVarTKG_CUSTOM_IMAGE_REPOSITORY_CA_CERTIFICATE = { name: 'TKG_CUSTOM_IMAGE_REPOSITORY_CA_CERTIFICATE',
    valueType: ClusterClassVariableType.STRING_PARAGRAPH,
    description: 'CA certificate for custom image repository'
}
const ccVarTKG_HTTP_PROXY = { name: 'HTTP_PROXY', valueType: ClusterClassVariableType.STRING,
    description: 'list of proxy values for HTTP'
}
const ccVarTKG_HTTPS_PROXY = { name: 'TKG_HTTPS_PROXY', valueType: ClusterClassVariableType.STRING,
    description: 'list of proxy values for HTTPS'
}
const ccVarTKG_NO_PROXY = { name: 'TKG_NO_PROXY', valueType: ClusterClassVariableType.STRING,
    description: 'list of IPs that should not use the proxy'
}
const ccVarTKG_PROXY_CA_CERT = { name: 'TKG_PROXY_CA_CERT', valueType: ClusterClassVariableType.STRING,
    description: 'CA cert for proxy server'
}
const ccVarTKG_IP_FAMILY = { name: 'TKG_IP_FAMILY', valueType: ClusterClassVariableType.STRING,
    description: 'IP family',
    defaultValue: 'IPv4',
    possibleValues: ['IPv4', 'IPv6']
}
const ccVarKUBERNETES_VERSION = { name: 'KUBERNETES_VERSION', valueType: ClusterClassVariableType.STRING,
    description: 'kubernetes version for workload cluster'
}
const ccVarCLUSTER_PLAN = { name: 'CLUSTER_PLAN', valueType: ClusterClassVariableType.STRING,
    description: 'plan used for workload cluster: dev or prod',
    defaultValue: 'dev',
    possibleValues: ['dev', 'prod']
}
const ccVarCLUSTER_API_SERVER_PORT = { name: 'CLUSTER_API_SERVER_PORT', valueType: ClusterClassVariableType.STRING,
    description: 'kube-apiserver endpoint (IP) for the workload cluster'
}
const ccVarSIZE = { name: 'SIZE', valueType: ClusterClassVariableType.STRING,
    description: '--um--'
}
const ccVarCONTROLPLANE_SIZE = { name: 'CONTROLPLANE_SIZE', valueType: ClusterClassVariableType.STRING,
    description: '--um--'
}
const ccVarWORKER_SIZE = { name: 'WORKER_SIZE', valueType: ClusterClassVariableType.STRING,
    description: '--um--'
}
const ccVarENABLE_AUDIT_LOGGING = { name: 'ENABLE_AUDIT_LOGGING', valueType: ClusterClassVariableType.BOOLEAN,
    description: 'enable audit logging'
}
const ccVarDOCKER_MACHINE_TEMPLATE_IMAGE = { name: 'DOCKER_MACHINE_TEMPLATE_IMAGE', valueType: ClusterClassVariableType.STRING,
    description: 'image for docker machine'
}

const ccVarCLUSTER_CIDR = { name: 'CLUSTER_CIDR', valueType: ClusterClassVariableType.CIDR,
    description: 'CIDR for IP addresses of nodes in workload cluster'
}
const ccVarSERVICE_CIDR = { name: 'SERVICE_CIDR', valueType: ClusterClassVariableType.CIDR,
    description: 'CIDR for IP addresses of services in workload cluster'
}
const ccVarENABLE_AUTOSCALER = { name: 'ENABLE_AUTOSCALER', valueType: ClusterClassVariableType.BOOLEAN,
    description: 'enable auto scaler'
}
const ccVarAUTOSCALER_MIN_SIZE_0 = { name: 'AUTOSCALER_MIN_SIZE_0', valueType: ClusterClassVariableType.INTEGER_SMALL,
    description: 'auto scaler minimum size for first node'
}
const ccVarAUTOSCALER_MAX_SIZE_0 = { name: 'AUTOSCALER_MAX_SIZE_0', valueType: ClusterClassVariableType.INTEGER_SMALL,
    description: 'auto scaler maximum size for first node'
}
const ccVarAUTOSCALER_MIN_SIZE_1 = { name: 'AUTOSCALER_MIN_SIZE_1', valueType: ClusterClassVariableType.INTEGER_SMALL,
    description: 'auto scaler minimum size for second node'
}
const ccVarAUTOSCALER_MAX_SIZE_1 = { name: 'AUTOSCALER_MAX_SIZE_1', valueType: ClusterClassVariableType.INTEGER_SMALL,
    description: 'auto scaler maximum size for second node'
}
const ccVarAUTOSCALER_MIN_SIZE_2 = { name: 'AUTOSCALER_MIN_SIZE_2', valueType: ClusterClassVariableType.INTEGER_SMALL,
    description: 'auto scaler minimum size for third node'
}
const ccVarAUTOSCALER_MAX_SIZE_2 = { name: 'AUTOSCALER_MAX_SIZE_2', valueType: ClusterClassVariableType.INTEGER_SMALL,
    description: 'auto scaler maximum size for third node'
}
const ccVarCONTROL_PLANE_MACHINE_COUNT = { name: 'CONTROL_PLANE_MACHINE_COUNT', valueType: ClusterClassVariableType.INTEGER_SMALL,
    description: 'number of machines in control plane',
    defaultValue: '1'
}
const ccVarWORKER_MACHINE_COUNT = { name: 'WORKER_MACHINE_COUNT', valueType: ClusterClassVariableType.INTEGER_SMALL,
    description: 'number of machines in worker cluster',
    defaultValue: '1'
}
const ccVarWORKER_MACHINE_COUNT_0 = { name: 'WORKER_MACHINE_COUNT_0', valueType: ClusterClassVariableType.INTEGER_SMALL,
    description: 'number of machines in worker cluster',
    defaultValue: '1'
}
const ccVarWORKER_MACHINE_COUNT_1 = { name: 'WORKER_MACHINE_COUNT_1', valueType: ClusterClassVariableType.INTEGER_SMALL,
    description: 'number of machines in worker cluster',
    defaultValue: '1'
}
const ccVarWORKER_MACHINE_COUNT_2 = { name: 'WORKER_MACHINE_COUNT_2', valueType: ClusterClassVariableType.INTEGER_SMALL,
    description: 'number of machines in worker cluster',
    defaultValue: '1'
}
const ccVarOS_NAME = { name: 'OS_NAME', valueType: ClusterClassVariableType.STRING,
    description: 'name of OS to use on workload cluster'
}
const ccVarOS_VERSION = { name: 'OS_VERSION', valueType: ClusterClassVariableType.STRING,
    description: 'version of OS to use on workload cluster'
}
const ccVarOS_ARCH = { name: 'OS_ARCH', valueType: ClusterClassVariableType.STRING,
    description: 'architecture of OS to use on workload cluster',
    defaultValue: 'ubuntu',
    possibleValues: ['darwin64', 'win32', 'ubuntu']
}
//
// COMMON ccVar objects
////////////////////////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////////////////////////
// VSPHERE-specific ccVar objects
//
const ccVarVSPHERE_NUM_CPUS = { name: 'VSPHERE_NUM_CPUS', valueType: ClusterClassVariableType.INTEGER_SMALL,
    description: 'number of CPUs for workload nodes', defaultValue: '1'
}
const ccVarVSPHERE_DISK_GIB = { name: 'VSPHERE_DISK_GIB', valueType: ClusterClassVariableType.STRING,
    description: 'disk available (in Gb)'
}
const ccVarVSPHERE_MEM_MIB = { name: 'VSPHERE_MEM_MIB', valueType: ClusterClassVariableType.STRING,
    description: 'memory (in Mb)'
}
const ccVarVSPHERE_CONTROL_PLANE_NUM_CPUS = { name: 'VSPHERE_CONTROL_PLANE_NUM_CPUS', valueType: ClusterClassVariableType.INTEGER_SMALL,
    description: 'number of CPUs for control plane nodes', defaultValue: '1'
}
const ccVarVSPHERE_CONTROL_PLANE_DISK_GIB = { name: 'VSPHERE_CONTROL_PLANE_DISK_GIB', valueType: ClusterClassVariableType.INTEGER,
    description: 'disk available (in Gb) for control plane nodes'
}
const ccVarVSPHERE_CONTROL_PLANE_MEM_MIB = { name: 'VSPHERE_CONTROL_PLANE_MEM_MIB', valueType: ClusterClassVariableType.INTEGER,
    description: 'memory (in Mb) for control plane nodes'
}
const ccVarVSPHERE_WORKER_NUM_CPUS = { name: 'VSPHERE_WORKER_NUM_CPUS', valueType: ClusterClassVariableType.INTEGER_SMALL,
    description: 'version of OS to use on workload cluster'
}
const ccVarVSPHERE_WORKER_DISK_GIB = { name: 'VSPHERE_WORKER_DISK_GIB', valueType: ClusterClassVariableType.INTEGER,
    description: 'disk available (in Gb) for workload nodes'
}
const ccVarVSPHERE_WORKER_MEM_MIB = { name: 'VSPHERE_WORKER_MEM_MIB', valueType: ClusterClassVariableType.INTEGER,
    description: 'memory (in Mb) for workload nodes'
}
const ccVarVSPHERE_CLONE_MODE = { name: 'VSPHERE_CLONE_MODE', valueType: ClusterClassVariableType.STRING,
    description: '-- um -- enumeration?'
}
const ccVarVSPHERE_NETWORK = { name: 'VSPHERE_NETWORK', valueType: ClusterClassVariableType.STRING,
    description: '-- um --'
}
const ccVarVSPHERE_TEMPLATE = { name: 'VSPHERE_TEMPLATE', valueType: ClusterClassVariableType.STRING,
    description: 'template to use for worker nodes (non-Windows)'
}
const ccVarVSPHERE_WINDOWS_TEMPLATE = { name: 'VSPHERE_WINDOWS_TEMPLATE', valueType: ClusterClassVariableType.STRING,
    description: 'template to use for worker nodes (Windows)'
}
const ccVarCONTROL_PLANE_NODE_NAMESERVERS = { name: 'CONTROL_PLANE_NODE_NAMESERVERS', valueType: ClusterClassVariableType.STRING,
    description: '--um--'
}
const ccVarWORKER_NODE_NAMESERVERS = { name: 'WORKER_NODE_NAMESERVERS', valueType: ClusterClassVariableType.STRING,
    description: '--um--'
}
const ccVarVSPHERE_CONTROL_PLANE_ENDPOINT = { name: 'VSPHERE_CONTROL_PLANE_ENDPOINT', valueType: ClusterClassVariableType.IP,
    description: 'kube-apiserver endpoint (IP) for the workload cluster'
}
const ccVarIS_WINDOWS_WORKLOAD_CLUSTER = { name: 'IS_WINDOWS_WORKLOAD_CLUSTER', valueType: ClusterClassVariableType.BOOLEAN,
    description: 'Is this a Windows-based workload cluster?'
}
//
// VSPHERE-specific ccVar objects
////////////////////////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////////////////////////
// AWS-specific ccVar objects
//
const ccVarAWS_VPC_ID = { name: 'AWS_VPC_ID', valueType: ClusterClassVariableType.STRING,
    description: 'VPC id',
    defaultValue: '123',
}
const ccVarAWS_REGION = { name: 'AWS_REGION', valueType: ClusterClassVariableType.STRING,
    description: 'Aws region where workload cluster will be deployed',
}
const ccVarAWS_LOAD_BALANCER_SCHEME_INTERNAL = { name: 'AWS_LOAD_BALANCER_SCHEME_INTERNAL', valueType: ClusterClassVariableType.STRING,
    description: '-- um --',
}
const ccVarAWS_NODE_AZ = { name: 'AWS_NODE_AZ', valueType: ClusterClassVariableType.STRING,
    description: 'AZ for first node of workload cluster',
}
const ccVarAWS_NODE_AZ_1 = { name: 'AWS_NODE_AZ_1', valueType: ClusterClassVariableType.STRING,
    description: 'AZ for second node of workload cluster',
}
const ccVarAWS_NODE_AZ_2 = { name: 'AWS_NODE_AZ_2', valueType: ClusterClassVariableType.STRING,
    description: 'AZ for third node of workload cluster',
}
const ccVarAWS_PRIVATE_SUBNET_ID = { name: 'AWS_PRIVATE_SUBNET_ID', valueType: ClusterClassVariableType.IP,
    description: '-- um --',
}
const ccVarAWS_PUBLIC_SUBNET_ID = { name: 'AWS_PUBLIC_SUBNET_ID', valueType: ClusterClassVariableType.IP,
    description: '-- um --',
}
const ccVarAWS_PUBLIC_SUBNET_ID_1 = { name: 'AWS_PUBLIC_SUBNET_ID_1', valueType: ClusterClassVariableType.IP,
    description: '-- um --',
}
const ccVarAWS_PRIVATE_SUBNET_ID_1 = { name: 'PRIVATE_SUBNET_ID_1', valueType: ClusterClassVariableType.IP,
    description: '-- um --',
}
const ccVarAWS_PUBLIC_SUBNET_ID_2 = { name: 'AWS_PUBLIC_SUBNET_ID_2', valueType: ClusterClassVariableType.IP,
    description: '-- um --',
}
const ccVarAWS_PRIVATE_SUBNET_ID_2 = { name: 'AWS_PRIVATE_SUBNET_ID_2', valueType: ClusterClassVariableType.IP,
    description: '-- um --',
}
const ccVarAWS_VPC_CIDR = { name: 'AWS_VPC_CIDR', valueType: ClusterClassVariableType.CIDR,
    description: '-- um --',
}
const ccVarAWS_PRIVATE_NODE_CIDR = { name: 'AWS_PRIVATE_NODE_CIDR', valueType: ClusterClassVariableType.CIDR,
    description: '-- um --',
}
const ccVarAWS_PUBLIC_NODE_CIDR = { name: 'AWS_PUBLIC_NODE_CIDR', valueType: ClusterClassVariableType.CIDR,
    description: '-- um --',
}
const ccVarAWS_PRIVATE_NODE_CIDR_1 = { name: 'AWS_PRIVATE_NODE_CIDR_1', valueType: ClusterClassVariableType.CIDR,
    description: '-- um --',
}
const ccVarAWS_PUBLIC_NODE_CIDR_1 = { name: 'AWS_PUBLIC_NODE_CIDR_1', valueType: ClusterClassVariableType.CIDR,
    description: '-- um --',
}
const ccVarAWS_PRIVATE_NODE_CIDR_2 = { name: 'AWS_PRIVATE_NODE_CIDR_2', valueType: ClusterClassVariableType.CIDR,
    description: '-- um --',
}
const ccVarAWS_PUBLIC_NODE_CIDR_2 = { name: 'AWS_PUBLIC_NODE_CIDR_2', valueType: ClusterClassVariableType.CIDR,
    description: '-- um --',
}
const ccVarAWS_SECURITY_GROUP_APISERVER_LB = { name: 'AWS_SECURITY_GROUP_APISERVER_LB', valueType: ClusterClassVariableType.STRING,
    description: '-- um --',
}
const ccVarAWS_SECURITY_GROUP_BASTION = { name: 'AWS_SECURITY_GROUP_BASTION', valueType: ClusterClassVariableType.STRING,
    description: '-- um --',
}
const ccVarAWS_SECURITY_GROUP_CONTROLPLANE = { name: 'AWS_SECURITY_GROUP_CONTROLPLANE', valueType: ClusterClassVariableType.STRING,
    description: '-- um --',
}
const ccVarAWS_SECURITY_GROUP_LB = { name: 'AWS_SECURITY_GROUP_LB', valueType: ClusterClassVariableType.STRING,
    description: '-- um --',
}
const ccVarAWS_SECURITY_GROUP_NODE = { name: 'AWS_SECURITY_GROUP_NODE', valueType: ClusterClassVariableType.STRING,
    description: '-- um --',
}
const ccVarAWS_IDENTITY_REF_KIND = { name: 'AWS_IDENTITY_REF_KIND', valueType: ClusterClassVariableType.STRING,
    description: '-- um --',
}
const ccVarAWS_IDENTITY_REF_NAME = { name: 'AWS_IDENTITY_REF_NAME', valueType: ClusterClassVariableType.STRING,
    description: '-- um --',
}
const ccVarAWS_CONTROL_PLANE_OS_DISK_SIZE_GIB = { name: 'AWS_CONTROL_PLANE_OS_DISK_SIZE_GIB', valueType: ClusterClassVariableType.INTEGER,
    description: '-- um --',
}
const ccVarAWS_NODE_OS_DISK_SIZE_GIB = { name: 'AWS_NODE_OS_DISK_SIZE_GIB', valueType: ClusterClassVariableType.INTEGER,
    description: '-- um --',
}
const ccVarCONTROL_PLANE_MACHINE_TYPE = { name: 'ccVarCONTROL_PLANE_MACHINE_TYPE', valueType: ClusterClassVariableType.STRING,
    description: '-- um --',
}
const ccVarNODE_MACHINE_TYPE = { name: 'NODE_MACHINE_TYPE', valueType: ClusterClassVariableType.STRING,
    description: '-- um --',
}
const ccVarNODE_MACHINE_TYPE_1 = { name: 'NODE_MACHINE_TYPE_1', valueType: ClusterClassVariableType.STRING,
    description: '-- um --',
}
const ccVarNODE_MACHINE_TYPE_2 = { name: 'NODE_MACHINE_TYPE_2', valueType: ClusterClassVariableType.STRING,
    description: '-- um --',
}
const ccVarAWS_SSH_KEY_NAME = { name: 'AWS_SSH_KEY_NAME', valueType: ClusterClassVariableType.STRING,
    description: '-- um --',
}
const ccVarBASTION_HOST_ENABLED = { name: 'BASTION_HOST_ENABLED', valueType: ClusterClassVariableType.BOOLEAN,
    description: '-- um --',
}
//
// AWS-specific ccVar objects
////////////////////////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////////////////////////
// AWS-specific ccVar objects
//
const ccVarAZURE_RESOURCE_GROUP = { name: 'AZURE_RESOURCE_GROUP', valueType: ClusterClassVariableType.STRING,
    description: 'resource group for creating the workload cluster'
}
const ccVarAZURE_TENANT_ID = { name: 'AZURE_TENANT_ID', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_SUBSCRIPTION_ID = { name: 'AZURE_SUBSCRIPTION_ID', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_CLIENT_ID = { name: 'AZURE_CLIENT_ID', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_CLIENT_SECRET = { name: 'AZURE_CLIENT_SECRET', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_LOCATION = { name: 'AZURE_LOCATION', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_SSH_PUBLIC_KEY_B64 = { name: 'AZURE_SSH_PUBLIC_KEY_B64', valueType: ClusterClassVariableType.STRING_PARAGRAPH,
    description: '-- um-- '
}
const ccVarAZURE_CONTROL_PLANE_MACHINE_TYPE = { name: 'AZURE_CONTROL_PLANE_MACHINE_TYPE', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_NODE_MACHINE_TYPE = { name: 'AZURE_NODE_MACHINE_TYPE', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_ENABLE_ACCELERATED_NETWORKING = { name: 'AZURE_ENABLE_ACCELERATED_NETWORKING', valueType: ClusterClassVariableType.BOOLEAN,
    description: '-- um-- '
}
const ccVarAZURE_VNET_RESOURCE_GROUP = { name: 'AZURE_VNET_RESOURCE_GROUP', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_VNET_NAME = { name: 'AZURE_VNET_NAME', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_VNET_CIDR = { name: 'AZURE_VNET_CIDR', valueType: ClusterClassVariableType.CIDR,
    description: '-- um-- '
}
const ccVarAZURE_CONTROL_PLANE_SUBNET_NAME = { name: 'AZURE_CONTROL_PLANE_SUBNET_NAME', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_CONTROL_PLANE_SUBNET_CIDR = { name: 'AZURE_CONTROL_PLANE_SUBNET_CIDR', valueType: ClusterClassVariableType.CIDR,
    description: '-- um-- '
}
const ccVarAZURE_CONTROL_PLANE_SUBNET_SECURITY_GROUP = { name: 'AZURE_CONTROL_PLANE_SUBNET_SECURITY_GROUP',
    valueType: ClusterClassVariableType.STRING, description: '-- um-- '
}
const ccVarAZURE_NODE_SUBNET_NAME = { name: 'AZURE_NODE_SUBNET_NAME', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_NODE_SUBNET_CIDR = { name: 'AZURE_NODE_SUBNET_CIDR', valueType: ClusterClassVariableType.CIDR,
    description: '-- um-- '
}
const ccVarAZURE_NODE_SUBNET_SECURITY_GROUP = { name: 'AZURE_NODE_SUBNET_SECURITY_GROUP', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_NODE_AZ = { name: 'AZURE_NODE_AZ', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_NODE_AZ_1 = { name: 'AZURE_NODE_AZ_1', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_NODE_AZ_2 = { name: 'AZURE_NODE_AZ_2', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_CUSTOM_TAGS = { name: 'AZURE_CUSTOM_TAGS', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_CONTROL_PLANE_OS_DISK_SIZE_GIB = { name: 'AZURE_CONTROL_PLANE_OS_DISK_SIZE_GIB',
    valueType: ClusterClassVariableType.STRING, description: '-- um-- '
}
const ccVarAZURE_CONTROL_PLANE_OS_DISK_STORAGE_ACCOUNT_TYPE = { name: 'AZURE_CONTROL_PLANE_OS_DISK_STORAGE_ACCOUNT_TYPE',
    valueType: ClusterClassVariableType.STRING, description: '-- um-- '
}
const ccVarAZURE_NODE_OS_DISK_SIZE_GIB = { name: 'AZURE_NODE_OS_DISK_SIZE_GIB', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_NODE_OS_DISK_STORAGE_ACCOUNT_TYPE = { name: 'AZURE_NODE_OS_DISK_STORAGE_ACCOUNT_TYPE',
    valueType: ClusterClassVariableType.STRING, description: '-- um-- '
}
const ccVarAZURE_CONTROL_PLANE_DATA_DISK_SIZE_GIB = { name: 'AZURE_CONTROL_PLANE_DATA_DISK_SIZE_GIB',
    valueType: ClusterClassVariableType.INTEGER, description: '-- um-- '
}
const ccVarAZURE_ENABLE_NODE_DATA_DISK = { name: 'AZURE_ENABLE_NODE_DATA_DISK', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_NODE_DATA_DISK_SIZE_GIB = { name: 'AZURE_NODE_DATA_DISK_SIZE_GIB', valueType: ClusterClassVariableType.INTEGER,
    description: '-- um-- '
}
const ccVarAZURE_ENABLE_PRIVATE_CLUSTER = { name: 'AZURE_ENABLE_PRIVATE_CLUSTER', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_FRONTEND_PRIVATE_IP = { name: 'AZURE_FRONTEND_PRIVATE_IP', valueType: ClusterClassVariableType.IP,
    description: '-- um-- '
}
const ccVarAZURE_ENABLE_CONTROL_PLANE_OUTBOUND_LB = { name: 'AZURE_ENABLE_CONTROL_PLANE_OUTBOUND_LB', 
    valueType: ClusterClassVariableType.STRING, description: '-- um-- '
}
const ccVarAZURE_ENABLE_NODE_OUTBOUND_LB = { name: 'AZURE_ENABLE_NODE_OUTBOUND_LB', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_CONTROL_PLANE_OUTBOUND_LB_FRONTEND_IP_COUNT = { name: 'AZURE_CONTROL_PLANE_OUTBOUND_LB_FRONTEND_IP_COUNT', 
    valueType: ClusterClassVariableType.INTEGER_SMALL, description: '-- um-- '
}
const ccVarAZURE_NODE_OUTBOUND_LB_FRONTEND_IP_COUNT = { name: 'AZURE_NODE_OUTBOUND_LB_FRONTEND_IP_COUNT', 
    valueType: ClusterClassVariableType.INTEGER_SMALL, description: '-- um-- '
}
const ccVarAZURE_NODE_OUTBOUND_LB_IDLE_TIMEOUT_IN_MINUTES = { name: 'AZURE_NODE_OUTBOUND_LB_IDLE_TIMEOUT_IN_MINUTES', 
    valueType: ClusterClassVariableType.INTEGER, description: '-- um-- '
}
const ccVarAZURE_IMAGE_ID = { name: 'AZURE_IMAGE_ID', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_IMAGE_RESOURCE_GROUP = { name: 'AZURE_IMAGE_RESOURCE_GROUP', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_IMAGE_NAME = { name: 'AZURE_IMAGE_NAME', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_IMAGE_SUBSCRIPTION_ID = { name: 'AZURE_IMAGE_SUBSCRIPTION_ID', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_IMAGE_GALLERY = { name: 'AZURE_IMAGE_GALLERY', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_IMAGE_PUBLISHER = { name: 'AZURE_IMAGE_PUBLISHER', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_IMAGE_OFFER = { name: 'AZURE_IMAGE_OFFER', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_IMAGE_SKU = { name: 'AZURE_IMAGE_SKU', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_IMAGE_THIRD_PARTY = { name: 'AZURE_IMAGE_THIRD_PARTY', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_IMAGE_VERSION = { name: 'AZURE_IMAGE_VERSION', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_IDENTITY_NAME = { name: 'AZURE_IDENTITY_NAME', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
const ccVarAZURE_IDENTITY_NAMESPACE = { name: 'AZURE_IDENTITY_NAMESPACE', valueType: ClusterClassVariableType.STRING,
    description: '-- um-- '
}
//
// AWS-specific ccVar objects
////////////////////////////////////////////////////////////////////////////////////////////////////

export const FakeClusterClassVsphere = {
    name: 'tkg-vsphere-default',
    requiredVariables: [
        ccVarVSPHERE_CONTROL_PLANE_ENDPOINT,
        ccVarIS_WINDOWS_WORKLOAD_CLUSTER,
    ],
    optionalVariables: [
        ccVarKUBERNETES_VERSION,
        ccVarCLUSTER_PLAN,
        ccVarCLUSTER_CIDR,
        ccVarSERVICE_CIDR,
        ccVarENABLE_AUDIT_LOGGING,
        ccVarTKG_HTTP_PROXY,
        ccVarTKG_HTTPS_PROXY,
        ccVarTKG_NO_PROXY,
        ccVarTKG_PROXY_CA_CERT,
        ccVarTKG_CUSTOM_IMAGE_REPOSITORY,
        ccVarTKG_CUSTOM_IMAGE_REPOSITORY_SKIP_TLS_VERIFY,
        ccVarTKG_CUSTOM_IMAGE_REPOSITORY_CA_CERTIFICATE,
    ],
    advancedVariables: [
        ccVarTKG_IP_FAMILY,
        ccVarCONTROL_PLANE_NODE_NAMESERVERS,
        ccVarWORKER_NODE_NAMESERVERS,
        ccVarCONTROL_PLANE_MACHINE_COUNT,
        ccVarWORKER_MACHINE_COUNT,
        ccVarWORKER_MACHINE_COUNT_0,
        ccVarWORKER_MACHINE_COUNT_1,
        ccVarWORKER_MACHINE_COUNT_2,
        ccVarOS_NAME,
        ccVarOS_VERSION,
        ccVarOS_ARCH,
        ccVarENABLE_AUTOSCALER,
        ccVarAUTOSCALER_MIN_SIZE_0,
        ccVarAUTOSCALER_MAX_SIZE_0,
        ccVarAUTOSCALER_MIN_SIZE_1,
        ccVarAUTOSCALER_MAX_SIZE_1,
        ccVarAUTOSCALER_MIN_SIZE_2,
        ccVarAUTOSCALER_MAX_SIZE_2,
        ccVarVSPHERE_NUM_CPUS,
        ccVarVSPHERE_DISK_GIB,
        ccVarVSPHERE_MEM_MIB,
        ccVarVSPHERE_CONTROL_PLANE_NUM_CPUS,
        ccVarVSPHERE_CONTROL_PLANE_DISK_GIB,
        ccVarVSPHERE_CONTROL_PLANE_MEM_MIB,
        ccVarVSPHERE_WORKER_NUM_CPUS,
        ccVarVSPHERE_WORKER_DISK_GIB,
        ccVarVSPHERE_WORKER_MEM_MIB,
        ccVarVSPHERE_CLONE_MODE,
        ccVarVSPHERE_NETWORK,
        ccVarVSPHERE_TEMPLATE,
        ccVarVSPHERE_WINDOWS_TEMPLATE,
    ]
} as ClusterClassDefinition

export const FakeClusterClassAws = {
    name: 'tkg-aws-default',
    requiredVariables: [],
    optionalVariables: [
        ccVarKUBERNETES_VERSION,
        ccVarCLUSTER_PLAN,
        ccVarCLUSTER_API_SERVER_PORT,
        ccVarSIZE,
        ccVarCONTROLPLANE_SIZE,
        ccVarWORKER_SIZE,
        ccVarAWS_VPC_ID,
        ccVarCLUSTER_CIDR,
        ccVarSERVICE_CIDR,
        ccVarENABLE_AUDIT_LOGGING,
        ccVarTKG_HTTP_PROXY,
        ccVarTKG_HTTPS_PROXY,
        ccVarTKG_NO_PROXY,
        ccVarTKG_PROXY_CA_CERT,
        ccVarTKG_CUSTOM_IMAGE_REPOSITORY,
        ccVarTKG_CUSTOM_IMAGE_REPOSITORY_SKIP_TLS_VERIFY,
        ccVarTKG_CUSTOM_IMAGE_REPOSITORY_CA_CERTIFICATE,
    ],
    advancedVariables: [
        ccVarTKG_IP_FAMILY,
        ccVarCONTROL_PLANE_MACHINE_COUNT,
        ccVarWORKER_MACHINE_COUNT,
        ccVarWORKER_MACHINE_COUNT_0,
        ccVarWORKER_MACHINE_COUNT_1,
        ccVarWORKER_MACHINE_COUNT_2,
        ccVarOS_NAME,
        ccVarOS_VERSION,
        ccVarOS_ARCH,
        ccVarENABLE_AUTOSCALER,
        ccVarAUTOSCALER_MIN_SIZE_0,
        ccVarAUTOSCALER_MAX_SIZE_0,
        ccVarAUTOSCALER_MIN_SIZE_1,
        ccVarAUTOSCALER_MAX_SIZE_1,
        ccVarAUTOSCALER_MIN_SIZE_2,
        ccVarAUTOSCALER_MAX_SIZE_2,
        ccVarAWS_REGION,
        ccVarAWS_LOAD_BALANCER_SCHEME_INTERNAL,
        ccVarAWS_NODE_AZ,
        ccVarAWS_NODE_AZ_1,
        ccVarAWS_NODE_AZ_2,
        ccVarAWS_VPC_ID,
        ccVarAWS_PRIVATE_SUBNET_ID,
        ccVarAWS_PUBLIC_SUBNET_ID,
        ccVarAWS_PUBLIC_SUBNET_ID_1,
        ccVarAWS_PRIVATE_SUBNET_ID_1,
        ccVarAWS_PUBLIC_SUBNET_ID_2,
        ccVarAWS_PRIVATE_SUBNET_ID_2,
        ccVarAWS_VPC_CIDR,
        ccVarAWS_PRIVATE_NODE_CIDR,
        ccVarAWS_PUBLIC_NODE_CIDR,
        ccVarAWS_PRIVATE_NODE_CIDR_1,
        ccVarAWS_PUBLIC_NODE_CIDR_1,
        ccVarAWS_PRIVATE_NODE_CIDR_2,
        ccVarAWS_PUBLIC_NODE_CIDR_2,
        ccVarAWS_SECURITY_GROUP_APISERVER_LB,
        ccVarAWS_SECURITY_GROUP_BASTION,
        ccVarAWS_SECURITY_GROUP_CONTROLPLANE,
        ccVarAWS_SECURITY_GROUP_LB,
        ccVarAWS_SECURITY_GROUP_NODE,
        ccVarAWS_IDENTITY_REF_KIND,
        ccVarAWS_IDENTITY_REF_NAME,
        ccVarAWS_CONTROL_PLANE_OS_DISK_SIZE_GIB,
        ccVarAWS_NODE_OS_DISK_SIZE_GIB,
        ccVarCONTROL_PLANE_MACHINE_TYPE,
        ccVarNODE_MACHINE_TYPE,
        ccVarNODE_MACHINE_TYPE_1,
        ccVarNODE_MACHINE_TYPE_2,
        ccVarAWS_SSH_KEY_NAME,
        ccVarBASTION_HOST_ENABLED,
    ]
} as ClusterClassDefinition

export const FakeClusterClassAzure = {
    name: 'tkg-azure-default',
    requiredVariables: [
        ccVarAZURE_RESOURCE_GROUP,
    ],
    optionalVariables: [
        ccVarKUBERNETES_VERSION,
        ccVarCLUSTER_PLAN,
        ccVarCLUSTER_API_SERVER_PORT,
        ccVarSIZE,
        ccVarCONTROLPLANE_SIZE,
        ccVarWORKER_SIZE,
        ccVarCLUSTER_CIDR,
        ccVarSERVICE_CIDR,
        ccVarENABLE_AUDIT_LOGGING,
        ccVarTKG_HTTP_PROXY,
        ccVarTKG_HTTPS_PROXY,
        ccVarTKG_NO_PROXY,
        ccVarTKG_PROXY_CA_CERT,
        ccVarTKG_CUSTOM_IMAGE_REPOSITORY,
        ccVarTKG_CUSTOM_IMAGE_REPOSITORY_SKIP_TLS_VERIFY,
        ccVarTKG_CUSTOM_IMAGE_REPOSITORY_CA_CERTIFICATE,
    ],
    advancedVariables: [
        ccVarTKG_IP_FAMILY,
        ccVarCONTROL_PLANE_MACHINE_COUNT,
        ccVarWORKER_MACHINE_COUNT,
        ccVarWORKER_MACHINE_COUNT_0,
        ccVarWORKER_MACHINE_COUNT_1,
        ccVarWORKER_MACHINE_COUNT_2,
        ccVarOS_NAME,
        ccVarOS_VERSION,
        ccVarOS_ARCH,
        ccVarENABLE_AUTOSCALER,
        ccVarAUTOSCALER_MIN_SIZE_0,
        ccVarAUTOSCALER_MAX_SIZE_0,
        ccVarAUTOSCALER_MIN_SIZE_1,
        ccVarAUTOSCALER_MAX_SIZE_1,
        ccVarAUTOSCALER_MIN_SIZE_2,
        ccVarAUTOSCALER_MAX_SIZE_2,
        ccVarAZURE_TENANT_ID,
        ccVarAZURE_SUBSCRIPTION_ID,
        ccVarAZURE_CLIENT_ID,
        ccVarAZURE_CLIENT_SECRET,
        ccVarAZURE_LOCATION,
        ccVarAZURE_SSH_PUBLIC_KEY_B64,
        ccVarAZURE_CONTROL_PLANE_MACHINE_TYPE,
        ccVarAZURE_NODE_MACHINE_TYPE,
        ccVarAZURE_ENABLE_ACCELERATED_NETWORKING,
        ccVarAZURE_RESOURCE_GROUP,
        ccVarAZURE_VNET_RESOURCE_GROUP,
        ccVarAZURE_VNET_NAME,
        ccVarAZURE_VNET_CIDR,
        ccVarAZURE_CONTROL_PLANE_SUBNET_NAME,
        ccVarAZURE_CONTROL_PLANE_SUBNET_CIDR,
        ccVarAZURE_CONTROL_PLANE_SUBNET_SECURITY_GROUP,
        ccVarAZURE_NODE_SUBNET_NAME,
        ccVarAZURE_NODE_SUBNET_CIDR,
        ccVarAZURE_NODE_SUBNET_SECURITY_GROUP,
        ccVarAZURE_NODE_AZ,
        ccVarAZURE_NODE_AZ_1,
        ccVarAZURE_NODE_AZ_2,
        ccVarAZURE_CUSTOM_TAGS,
        ccVarAZURE_CONTROL_PLANE_OS_DISK_SIZE_GIB,
        ccVarAZURE_CONTROL_PLANE_OS_DISK_STORAGE_ACCOUNT_TYPE,
        ccVarAZURE_NODE_OS_DISK_SIZE_GIB,
        ccVarAZURE_NODE_OS_DISK_STORAGE_ACCOUNT_TYPE,
        ccVarAZURE_CONTROL_PLANE_DATA_DISK_SIZE_GIB,
        ccVarAZURE_ENABLE_NODE_DATA_DISK,
        ccVarAZURE_NODE_DATA_DISK_SIZE_GIB,
        ccVarAZURE_ENABLE_PRIVATE_CLUSTER,
        ccVarAZURE_FRONTEND_PRIVATE_IP,
        ccVarAZURE_ENABLE_CONTROL_PLANE_OUTBOUND_LB,
        ccVarAZURE_ENABLE_NODE_OUTBOUND_LB,
        ccVarAZURE_CONTROL_PLANE_OUTBOUND_LB_FRONTEND_IP_COUNT,
        ccVarAZURE_NODE_OUTBOUND_LB_FRONTEND_IP_COUNT,
        ccVarAZURE_NODE_OUTBOUND_LB_IDLE_TIMEOUT_IN_MINUTES,
        ccVarAZURE_IMAGE_ID,
        ccVarAZURE_IMAGE_RESOURCE_GROUP,
        ccVarAZURE_IMAGE_NAME,
        ccVarAZURE_IMAGE_SUBSCRIPTION_ID,
        ccVarAZURE_IMAGE_GALLERY,
        ccVarAZURE_IMAGE_PUBLISHER,
        ccVarAZURE_IMAGE_OFFER,
        ccVarAZURE_IMAGE_SKU,
        ccVarAZURE_IMAGE_THIRD_PARTY,
        ccVarAZURE_IMAGE_VERSION,
        ccVarAZURE_IDENTITY_NAME,
        ccVarAZURE_IDENTITY_NAMESPACE,
    ]
} as ClusterClassDefinition


export const FakeClusterClassDocker = {
    name: 'tkg-docker-default',
    requiredVariables: [],
    optionalVariables: [
        ccVarKUBERNETES_VERSION,
        ccVarCLUSTER_PLAN,
        ccVarCLUSTER_API_SERVER_PORT,
        ccVarSIZE,
        ccVarCONTROLPLANE_SIZE,
        ccVarWORKER_SIZE,
        ccVarENABLE_AUDIT_LOGGING,
        ccVarDOCKER_MACHINE_TEMPLATE_IMAGE,
    ],
    advancedVariables: [
        ccVarCONTROL_PLANE_MACHINE_COUNT,
        ccVarWORKER_MACHINE_COUNT,
        ccVarWORKER_MACHINE_COUNT_0,
        ccVarWORKER_MACHINE_COUNT_1,
        ccVarWORKER_MACHINE_COUNT_2,
        ccVarOS_NAME,
        ccVarOS_VERSION,
        ccVarOS_ARCH,
        ccVarCLUSTER_CIDR,
        ccVarSERVICE_CIDR,
        ccVarENABLE_AUTOSCALER,
        ccVarAUTOSCALER_MIN_SIZE_0,
        ccVarAUTOSCALER_MIN_SIZE_1,
        ccVarAUTOSCALER_MIN_SIZE_2,
    ]
} as ClusterClassDefinition

