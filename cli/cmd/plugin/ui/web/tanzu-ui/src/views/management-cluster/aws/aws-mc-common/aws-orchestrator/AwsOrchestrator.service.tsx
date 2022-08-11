// React imports
import React from 'react';
import { AwsService } from '../../../../../swagger-api';
import { AWS_FIELDS } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { AwsDefaults } from '../default-service/AwsDefaults.service';
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';
import { AWS_ADD_RESOURCES } from '../../../../../state-management/actions/Resources.actions';
import { AwsResourceAction, FormAction, StoreDispatch } from '../../../../../shared/types/types';
export class AwsOrchestrator extends React.Component {
    static async initOsImage(awsState: { [key: string]: any }, awsDispatch: StoreDispatch) {
        // clear previous os images
        awsDispatch({
            type: AWS_ADD_RESOURCES,
            resourceName: 'osImages',
            payload: [],
        } as AwsResourceAction);
        // retrieve os images
        const osImages = await AwsService.getAwsosImages(awsState[STORE_SECTION_FORM].REGION);
        // save os images in store
        awsDispatch({
            type: AWS_ADD_RESOURCES,
            resourceName: 'osImages',
            payload: osImages,
        } as AwsResourceAction);
        // set default os image
        awsDispatch({
            type: INPUT_CHANGE,
            field: AWS_FIELDS.OS_IMAGE,
            payload: AwsDefaults.selectDefalutOsImage(osImages),
        } as FormAction);
    }
}
