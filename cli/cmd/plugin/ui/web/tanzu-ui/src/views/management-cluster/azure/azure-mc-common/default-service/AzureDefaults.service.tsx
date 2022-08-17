import React from 'react';
import { AzureVirtualMachine } from '../../../../../swagger-api';

export class AzureDefaults extends React.Component {
    // The strategy of deciding default os image
    static selectDefaultOsImage = (osImages: AzureVirtualMachine[]) => {
        return osImages && osImages.length > 0 ? osImages[0] : undefined;
    };
}
