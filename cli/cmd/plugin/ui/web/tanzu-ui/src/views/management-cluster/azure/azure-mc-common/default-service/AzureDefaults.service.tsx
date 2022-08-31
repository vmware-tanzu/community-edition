// React imports
import React from 'react';

// App imports
import { getNodeTypesByNodeProfile } from '../../../../../shared/constants/defaults/azure.defaults';
import { findFirstMatchingOption } from '../../../../../shared/utilities/Array.util';
import { AzureInstanceType, AzureVirtualMachine } from '../../../../../swagger-api';

export class AzureDefaults extends React.Component {
    // The strategy of deciding default os image
    static selectDefaultOsImage = (osImages: AzureVirtualMachine[]) => {
        return osImages && osImages.length > 0 ? osImages[0] : undefined;
    };
    static getDefaultNodeType = (nodeTypes: Array<AzureInstanceType>, selectedProfile: string) => {
        if (!nodeTypes) {
            return '';
        }
        const nodeTypeList: Array<string> = nodeTypes.map((instanceType: AzureInstanceType) => instanceType?.name || '');
        return findFirstMatchingOption(nodeTypeList, getNodeTypesByNodeProfile(selectedProfile));
    };
}
