// App imports
import { AWSVirtualMachine } from '../../../../../swagger-api';
import { AWSKeyPair } from '../../../../../swagger-api/models/AWSKeyPair';
import { first } from '../../../../../shared/utilities/Array.util';

export class AwsDefaults {
    // The strategy of deciding default os image
    static selectDefaultOsImage = (osImages: AWSVirtualMachine[]) => {
        return first<AWSVirtualMachine>(osImages);
    };

    static selectDefaultEC2KeyPairs = (keyPairs: AWSKeyPair[]) => {
        return first<AWSKeyPair>(keyPairs);
    };
}
