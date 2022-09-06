// App imports
import { AWSVirtualMachine, AWSAvailabilityZone } from '../../../../../swagger-api';
import { AWSKeyPair } from '../../../../../swagger-api/models/AWSKeyPair';
import { AWS_NODE_PROFILE_NAMES } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';
import { first } from '../../../../../shared/utilities/Array.util';
export class AwsDefaults {
    // The strategy of deciding default os image
    static selectDefaultOsImage = (osImages: AWSVirtualMachine[]) => {
        return first<AWSVirtualMachine>(osImages);
    };

    static selectDefaultEC2KeyPairs = (keyPairs: AWSKeyPair[]) => {
        return first<AWSKeyPair>(keyPairs);
    };

    static getDefaulAvailabilityZones = (azList: AWSAvailabilityZone[], nodeProfile: string) => {
        const defaultAZNameList: AWSAvailabilityZone[] = [];
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
                    console.error(
                        `For profile ${nodeProfile}, we expect to use 3 node profiles, but azList has a length of ${azList.length}`
                    );
                }
                return defaultAZNameList;
            }
        }
    };

    static selectDefaultNodeType(availableNodeTypes: string[], desiredNodeTypes: string[]): string {
        // TODO: implement correct strategy
        const set = new Set(availableNodeTypes);
        for (let i = 0; i < desiredNodeTypes.length; i++) {
            if (set.has(desiredNodeTypes[i])) {
                return desiredNodeTypes[i];
            }
        }
        return '';
    }
}
