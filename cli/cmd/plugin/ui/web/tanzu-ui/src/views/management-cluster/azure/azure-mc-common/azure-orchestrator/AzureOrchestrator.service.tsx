import React from 'react';
import { AzureService } from '../../../../../swagger-api';
import { AZURE_FIELDS } from '../../azure-mc-basic/AzureManagementClusterBasic.constants';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { AzureDefaults } from '../default-service/AzureDefaults.service';
import { AwsResourceAction, FormAction, StoreDispatch } from '../../../../../shared/types/types';
import { AZURE_ADD_RESOURCES } from '../../../../../state-management/actions/Resources.actions';
export class AzureOrchestrator extends React.Component {
    static async initOsImage(azureState: { [key: string]: any }, azureDispatch: StoreDispatch) {
        // clear previous os images
        azureDispatch({
            type: AZURE_ADD_RESOURCES,
            resourceName: 'osImages',
            payload: [],
        } as AwsResourceAction);
        // retrieve os images
        const osImages = await AzureService.getAzureOsImages();
        // save os images in store
        azureDispatch({
            type: AZURE_ADD_RESOURCES,
            resourceName: 'osImages',
            payload: osImages,
        } as AwsResourceAction);
        // set default os image
        azureDispatch({
            type: INPUT_CHANGE,
            field: AZURE_FIELDS.OS_IMAGE,
            payload: AzureDefaults.selectDefalutOsImage(osImages),
        } as FormAction);
    }
}
