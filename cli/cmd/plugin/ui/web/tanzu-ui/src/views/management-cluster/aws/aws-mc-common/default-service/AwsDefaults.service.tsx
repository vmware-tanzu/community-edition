// App imports
import { AWSVirtualMachine } from '../../../../../swagger-api';
import { AWSKeyPair } from '../../../../../swagger-api/models/AWSKeyPair';
import { first } from '../../../../../shared/utilities/Array.util';
import { validateDefaultNodeType } from '../../../../../shared/constants/defaults/aws.defaults';

export class AwsDefaults {
    // The strategy of deciding default os image
    static selectDefaultOsImage = (osImages: AWSVirtualMachine[]) => {
        return first<AWSVirtualMachine>(osImages);
    };

    static selectDefaultEC2KeyPairs = (keyPairs: AWSKeyPair[]) => {
        return first<AWSKeyPair>(keyPairs);
    };

    static setDefaultNodeType = (nodeTypeList: string[], nodeProfile: string) => {
        if (nodeTypeList.indexOf(validateDefaultNodeType(nodeProfile)) > -1) {
            return validateDefaultNodeType(nodeProfile);
        } else {
            // TODO: refactor to select optimal nodeType when preferred default not found
            return nodeTypeList[Math.round(nodeTypeList.length / 2)];
        }
    };
}
