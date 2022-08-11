// React imports
import React from 'react';
import { AWSVirtualMachine } from '../../../../../swagger-api';

export class AwsDefaults extends React.Component {
    // The strategy of deciding default os image
    static selectDefalutOsImage = (osImages: AWSVirtualMachine[]) => {
        return osImages && osImages.length > 0 ? osImages[0] : undefined;
    };
}
