import { KeyOfStringToArray } from '../../types/types';
import { AWS_FIELDS, AWS_NODE_PROFILE_NAMES } from '../../../views/management-cluster/aws/aws-mc-basic/AwsManagementClusterBasic.constants';

export const AWS_DEFAULT_VALUES = {
    // Cluster Name
    [AWS_FIELDS.CLUSTER_NAME]: '',
    // Cluster Data
    [AWS_FIELDS.CLUSTER_PLAN]: 'dev',

    // VPC New
    // VPC_NAME: 'temp-vpc-name',
    [AWS_FIELDS.VPC_CIDR]: '10.0.0.0/16',

    // other?
    [AWS_FIELDS.CREATE_CLOUDFORMATION_STACK]: true,
    [AWS_FIELDS.ENABLE_AUDIT_LOGGING]: true,
    [AWS_FIELDS.ENABLE_BASTION_HOST]: true,
    [AWS_FIELDS.ENABLE_CEIP_PARTICIPATION]: false,
    [AWS_FIELDS.ENABLE_MACHINE_HEALTH_CHECK]: true,

    // Kubernetes Networking
    [AWS_FIELDS.CLUSTER_NETWORKING_CNI_PROVIDER]: 'antrea',
    [AWS_FIELDS.CLUSTER_SERVICE_CIDR]: '100.64.0.0/13',
    [AWS_FIELDS.CLUSTER_POD_CIDR]: '100.96.0.0/11',

    // HTTP Proxy & Load Balancer
    [AWS_FIELDS.HTTP_PROXY_ENABLED]: false,
    [AWS_FIELDS.LOAD_BALANCER_SCHEME_INTERNAL]: false,
};
const SINGLE_NODE_FALL_BACK: Array<string> = [
    't3a.medium',
    't3a.xlarge',
    't4g.large',
    't4g.medium',
    't4g.xlarge',
    't3.large',
    't3.medium',
    't3.xlarge',
    't2.large',
    't2.medium',
    't2.xlarge',
    'm6g.large',
    'm6g.medium',
    'm6g.xlarge',
    'm6i.large',
    'm6i.xlarge',
    'm6a.large',
    'm6a.xlarge',
    'm5.large',
    'm5.xlarge',
    'm5a.large',
    'm5a.xlarge',
    'm5n.large',
    'm5n.xlarge',
    'm5zn.large',
    'm5zn.xlarge',
    'm4.large',
    'm4.xlarge',
    'a1.large',
    'a1.medium',
    'a1.xlarge',
];
const PRODUCTION_READY_FALL_BACK: Array<string> = [
    'm6a.large',
    'm6a.2xlarge',
    't4g.xlarge',
    't4g.large',
    't4g.2xlarge',
    't3.xlarge',
    't3.large',
    't3.2xlarge',
    't3a.xlarge',
    't3a.large',
    't3a.2xlarge',
    't2.xlarge',
    't2.large',
    't2.2xlarge',
    'm6g.xlarge',
    'm6g.large',
    'm6g.2xlarge',
    'm6i.xlarge',
    'm6i.large',
    'm6i.2xlarge',
    'm5.xlarge',
    'm5.large',
    'm5.2xlarge',
    'm5a.xlarge',
    'm5a.large',
    'm5a.2xlarge',
    'm5n.xlarge',
    'm5n.large',
    'm5n.2xlarge',
    'm5zn.xlarge',
    'm5zn.large',
    'm5zn.2xlarge',
    'm4.xlarge',
    'm4.large',
    'm4.2xlarge',
    'a1.xlarge',
    'a1.large',
    'a1.2xlarge',
];
const AWS_DEFAULT_INSTANCE_TYPES: KeyOfStringToArray = {
    [AWS_NODE_PROFILE_NAMES.SINGLE_NODE]: ['t3a.large', ...SINGLE_NODE_FALL_BACK],
    [AWS_NODE_PROFILE_NAMES.HIGH_AVAILABILITY]: ['t3a.large', ...SINGLE_NODE_FALL_BACK],
    [AWS_NODE_PROFILE_NAMES.PRODUCTION_READY]: ['m6a.xlarge', ...PRODUCTION_READY_FALL_BACK],
};

/**
 * @method validateDefaultNodeType
 * @param nodeProfile - node profile name set by ManagementClusterSettings.tsx; references key of AWS_DEFAULT_INSTANCE_TYPES
 * defaults map.
 * Returns default aws instance types according to the selected node profile.
 */
export function getDefaultNodeTypes(nodeProfile: string): Array<string> {
    return AWS_DEFAULT_INSTANCE_TYPES[nodeProfile];
}

// export function createAZList(storedAZObjects: { [key: string]: any }) {
//     const azList: { [key: string]: string }[] = [];
//     const workNodeType = /work-node-type/;
//     const publicSubnetID = /public-subnet-id/;
//     const privateSubnetID = /private-subnet-id/;
//     const keyMap: { [key: string]: string[] } = {};

//     Object.keys(storedAZObjects).map((key: string) => {
//         const nameSplit = key.split('-');
//         if (keyMap[nameSplit[2]] !== undefined) {
//             keyMap[nameSplit[2]].push(key);
//         } else {
//             keyMap[nameSplit[2]] = [key];
//         }
//     });

//     Object.keys(keyMap).map((key) => {
//         const az: { [key: string]: string } = {};
//         keyMap[key].map((key) => {
//             if (typeof storedAZObjects[key] === 'object') {
//                 az['name'] = storedAZObjects[key]['name'];
//             } else if (workNodeType.test(key)) {
//                 az['workNodeType'] = storedAZObjects[key];
//             } else if (publicSubnetID.test(key)) {
//                 az['publicSubnetID'] = storedAZObjects[key];
//             } else if (privateSubnetID.test(key)) {
//                 az['privateSubnetID'] = storedAZObjects[key];
//             }
//         });
//         azList.push(az);
//     });
//     return azList;
// }
