// React imports
import React, { ReactElement, useState, useContext } from 'react';
import { AwsService, AWSVirtualMachine } from '../../../../../swagger-api';
import { AWS_FIELDS } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';
import { FormAction } from '../../../../../shared/types/types';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { AwsStore } from '../../../../../state-management/stores/Store.aws';
import { getResource, STORE_SECTION_AWS_RESOURCES } from '../../../../../views/providers/aws/AwsResources.reducer';
export class AwsDefaults extends React.Component {
    static setDefaultOsImage = (awsState, awsDispatch) => {
        const osImages = (getResource('osImages', awsState) || []) as AWSVirtualMachine[];
        if (osImages[0]) {
            awsDispatch({
                type: INPUT_CHANGE,
                field: AWS_FIELDS.OS_IMAGE,
                payload: osImages[0],
            } as FormAction);
        }
    };
}
