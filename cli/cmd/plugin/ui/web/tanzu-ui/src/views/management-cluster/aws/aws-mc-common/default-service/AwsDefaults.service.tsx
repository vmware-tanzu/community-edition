// React imports
import React, { ReactElement, useState, useContext } from 'react';
import { AwsService, AWSVirtualMachine } from '../../../../../swagger-api';
import { FormAction } from '../../../../../shared/types/types';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { AwsStore } from '../../../../../state-management/stores/Store.aws';

export function AwsDefaults() {
    const { awsDispatch } = useContext(AwsStore);

    const retrieveOsImages = async (region: string) => {
        const result = await AwsService.getAwsosImages(region);
        return result;
    };

    const setDefaultOsImage = (defaultOsImage: AWSVirtualMachine) => {
        if (defaultOsImage) {
            awsDispatch({
                type: INPUT_CHANGE,
                field: 'OS_IMAGE',
                payload: defaultOsImage,
            } as FormAction);
        }
    };

    return { retrieveOsImages: retrieveOsImages, setDefaultOsImage: setDefaultOsImage };
}
