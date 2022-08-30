// App imports
import { AwsService, AWSVirtualMachine, AWSAvailabilityZone } from '../../../../../swagger-api';
import { AWSKeyPair } from '../../../../../swagger-api/models/AWSKeyPair';
import { first } from '../../../../../shared/utilities/Array.util';
import { getDefaultNodeTypes } from '../../../../../shared/constants/defaults/aws.defaults';
import { AWS_NODE_PROFILE_NAMES } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';
export class AwsDefaults {
    // The strategy of deciding default os image
    static selectDefaultOsImage = (osImages: AWSVirtualMachine[]) => {
        return first<AWSVirtualMachine>(osImages);
    };

    static selectDefaultEC2KeyPairs = (keyPairs: AWSKeyPair[]) => {
        return first<AWSKeyPair>(keyPairs);
    };

    static setDefaultNodeType = (nodeTypeList: string[], nodeProfile: string) => {
        const defaultNodeTypes = getDefaultNodeTypes(nodeProfile);
        const nodeTypeSet = new Set(nodeTypeList);

        for (let i = 0; i < defaultNodeTypes.length; i++) {
            if (nodeTypeSet.has(defaultNodeTypes[i])) {
                return defaultNodeTypes[i];
            }
        }
        return nodeTypeList[Math.round(nodeTypeList.length / 2)];
    };

    static defaulAvailabilityZoneNameStrategy = (azList: { [key: string]: string }[], nodeProfile: string) => {
        const defaultAZNameList: string[] = [];
        switch (nodeProfile) {
            case AWS_NODE_PROFILE_NAMES.SINGLE_NODE: {
                return [azList[0].id];
            }
            default: {
                azList.slice(0, 3).forEach((az) => {
                    defaultAZNameList.push(az.id);
                });
                return defaultAZNameList;
            }
        }
    };

    static defaulAvailabilityZoneNodeTypeStrategy = (azList: { [key: string]: string }[], nodeProfile: string) => {
        const defaultNodeTypes = getDefaultNodeTypes(nodeProfile);
        for (const az of azList) {
            for (let i = 0; i < defaultNodeTypes.length; i++) {
                if (az.workerNodeType === defaultNodeTypes[i]) {
                    return az;
                }
            }
        }

        const defaultNodeTypeByAZ: { [key: string]: any } = {
            workerNodeType: defaultNodeTypes[Math.round(defaultNodeTypes.length / 2)],
            publicSubnetID: '',
            privateSubnetID: '',
        };

        return defaultNodeTypeByAZ;
    };

    static async createAZNodeType(defaultAZName: string) {
        const azNodeTypes: { [key: string]: string }[] = [];
        const nodeTypes = await AwsService.getAwsNodeTypes(defaultAZName);
        nodeTypes.forEach((nodeType) => {
            azNodeTypes.push({ name: defaultAZName, workerNodeType: nodeType, publicSubnetID: '', privateSubnetID: '' });
        });
        return azNodeTypes;
    }
}
