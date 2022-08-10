// React imports
import React from 'react';
import { AwsService, AWSVirtualMachine } from '../../../../../swagger-api';

export class AwsDefaults extends React.Component {
    // Find what the default osImage is
    static selectDefalutOsImage = (osImages: AWSVirtualMachine[]) => {
        return osImages[0];
    };

    // Retrieve all osimages
    static retrieveOsImages = (region: string) => {
        return AwsService.getAwsosImages(region);
    };
}
