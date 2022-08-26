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

    static setDefaultNodeType = async (nodeProfile: string) => {
        const nodeTypeList: string[] = await AwsService.getAwsNodeTypes();
        return AwsDefaults.defaultNodeTypeStategy(nodeTypeList, nodeProfile);
    };

    static setDefaultAvailabilityZones = async (nodeProfile: string) => {
        const result: { [key: string]: string }[] = [];
        const azList: AWSAvailabilityZone[] = await AwsService.getAwsAvailabilityZones();
        // TODO: Create an error if the azList is empty
        switch (nodeProfile) {
            case AWS_NODE_PROFILE_NAMES.SINGLE_NODE:
                {
                    const defaultAZ: { [key: string]: string } = {};
                    if (azList[0]['name'] !== undefined) {
                        defaultAZ['name'] = azList[0]['name'];
                        const azNodeTypeList = await AwsService.getAwsNodeTypes(azList[0]['name']);
                        defaultAZ['workerNodeType'] = AwsDefaults.defaultNodeTypeStategy(azNodeTypeList, nodeProfile);
                        defaultAZ['publicSubnetID'] = '';
                        defaultAZ['privateSubnetID'] = '';
                        result.push(defaultAZ);
                    } else {
                        console.error(`This Availability Zone ${azList[0]} does not have name`);
                    }
                }
                break;
            default: {
                [...azList].slice(0, 3).forEach(async (az: AWSAvailabilityZone) => {
                    const defaultAZ: { [key: string]: string } = {};
                    if (az['name'] !== undefined) {
                        defaultAZ['name'] = az['name'];
                        const azNodeTypeList = await AwsService.getAwsNodeTypes(az['name']);
                        defaultAZ['workerNodeType'] = AwsDefaults.defaultNodeTypeStategy(azNodeTypeList, nodeProfile);
                        defaultAZ['publicSubnetID'] = '';
                        defaultAZ['privateSubnetID'] = '';
                        result.push(defaultAZ);
                    } else {
                        console.error(`This Availability Zone ${az} does not have name`);
                    }
                });
            }
        }
        return result;
    };

    static defaultNodeTypeStategy = (nodeTypeList: string[], nodeProfile: string) => {
        const defaultNodeTypes = getDefaultNodeTypes(nodeProfile);
        const nodeTypeSet = new Set(nodeTypeList);

        for (let i = 0; i < defaultNodeTypes.length; i++) {
            if (nodeTypeSet.has(defaultNodeTypes[i])) {
                return defaultNodeTypes[i];
            }
        }
        return nodeTypeList[Math.round(nodeTypeList.length / 2)];
    };
}
