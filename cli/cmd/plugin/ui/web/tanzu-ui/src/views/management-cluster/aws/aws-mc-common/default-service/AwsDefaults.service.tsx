import { AWSVirtualMachine } from '../../../../../swagger-api';
import { AWSKeyPair } from '../../../../../swagger-api/models/AWSKeyPair';

export class AwsDefaults {
    // The strategy of deciding default os image
    static selectDefalutOsImage = (osImages: AWSVirtualMachine[]) => {
        return osImages && osImages.length > 0 ? osImages[0] : undefined;
    };

    static selectDefalutEC2KeyPairs = (keyPairs: AWSKeyPair[]) => {
        return keyPairs && keyPairs.length > 0 ? keyPairs[0].name : undefined;
    };
}
