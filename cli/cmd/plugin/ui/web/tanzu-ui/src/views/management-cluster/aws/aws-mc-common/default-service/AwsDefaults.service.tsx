// React imports
import { AWSVirtualMachine } from '../../../../../swagger-api';

export class AwsDefaults {
    // The strategy of deciding default os image
    static selectDefalutOsImage = (osImages: AWSVirtualMachine[]) => {
        return osImages && osImages.length > 0 ? osImages[0] : undefined;
    };
}
