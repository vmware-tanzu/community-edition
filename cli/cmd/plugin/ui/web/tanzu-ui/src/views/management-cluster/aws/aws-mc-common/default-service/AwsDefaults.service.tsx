// App imports
import { AwsService, AWSVirtualMachine, AWSAvailabilityZone } from '../../../../../swagger-api';
import { AWSKeyPair } from '../../../../../swagger-api/models/AWSKeyPair';
import { first } from '../../../../../shared/utilities/Array.util';
import { getDefaultNodeTypes } from '../../../../../shared/constants/defaults/aws.defaults';
import { AWS_NODE_PROFILE_NAMES } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';

export interface SelectedAvailabiltyZoneData {
    id: string;
    name: string;
    workerNodeType: string;
    publicSubnetID: string;
    privateSubnetID: string;
}
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
        const defaultAZNameList: { [key: string]: string }[] = [];
        switch (nodeProfile) {
            case AWS_NODE_PROFILE_NAMES.SINGLE_NODE: {
                return [azList[0]];
            }
            default: {
                if (azList.length >= 3) {
                    azList.slice(0, 3).forEach((az) => {
                        defaultAZNameList.push(az);
                    });
                } else {
                    console.error('There is not enough Availability Zone');
                }
                return defaultAZNameList;
            }
        }
    };

    static defaulAvailabilityZoneNodeTypeStrategy = (azList: SelectedAvailabiltyZoneData[], nodeProfile: string, az: string) => {
        const defaultNodeTypes = getDefaultNodeTypes(nodeProfile);
        for (const az of azList) {
            for (let i = 0; i < defaultNodeTypes.length; i++) {
                if (az.workerNodeType === defaultNodeTypes[i]) {
                    return az;
                }
            }
        }
    };

    static async createAZNodeType(defaultAZName: { [key: string]: string }) {
        const azNodeTypes: SelectedAvailabiltyZoneData[] = [];
        const nodeTypes = await AwsService.getAwsNodeTypes(defaultAZName.name);
        nodeTypes.map((nodeType) => {
            azNodeTypes.push({
                id: defaultAZName.id,
                name: defaultAZName.name,
                workerNodeType: nodeType,
                publicSubnetID: '',
                privateSubnetID: '',
            });
        });
        return azNodeTypes;
    }
}
