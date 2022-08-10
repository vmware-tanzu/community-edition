// React imports
import React, { ReactElement, useState, useContext, useEffect } from 'react';
import { AwsService, AWSVirtualMachine } from '../../../../../swagger-api';
import { AWS_FIELDS } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { AwsStore } from '../../../../../state-management/stores/Store.aws';
import { AwsDefaults } from '../default-service/AwsDefaults.service';
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';
import { AWS_ADD_RESOURCES } from '../../../../../state-management/actions/Resources.actions';
import { AwsResourceAction, FormAction } from '../../../../../shared/types/types';
import { getResource, STORE_SECTION_AWS_RESOURCES } from '../../../../../views/providers/aws/AwsResources.reducer';

export class AwsOrchestrator extends React.Component {
    static async initOsImage(awsState, awsDispatch) {
        // retrieve all os images
        const osImages = await AwsDefaults.retrieveOsImages(awsState[STORE_SECTION_FORM].REGION);
        // save them in store
        if (osImages.length !== 0) {
            awsDispatch({
                type: AWS_ADD_RESOURCES,
                resourceName: 'osImages',
                payload: osImages,
            } as AwsResourceAction);
            // set default
            awsDispatch({
                type: INPUT_CHANGE,
                field: AWS_FIELDS.OS_IMAGE,
                payload: AwsDefaults.selectDefalutOsImage(osImages),
            } as FormAction);
        }
    }
}
