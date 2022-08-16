import React from 'react';
import { AzureService, AzureVirtualMachine } from '../../../../../swagger-api';
import { AZURE_FIELDS } from '../../azure-mc-basic/AzureManagementClusterBasic.constants';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { AzureDefaults } from '../default-service/AzureDefaults.service';
import { AzureResourceAction, FormAction, StoreDispatch } from '../../../../../shared/types/types';
import { AZURE_ADD_RESOURCES } from '../../../../../state-management/actions/Resources.actions';

interface AzureOrchestratorProps {
    azureState: { [key: string]: any };
    azureDispatch: StoreDispatch;
    errorObject: { [key: string]: any };
    setErrorObject: (newErrorObject: { [key: string]: any }) => void;
}
export class AzureOrchestrator extends React.Component {
    static async initOsImages(props: AzureOrchestratorProps) {
        const { azureState, azureDispatch, setErrorObject, errorObject } = props;
        try {
            const osImages = await AzureService.getAzureOsImages();
            saveCurrentResourceData(azureDispatch, AZURE_FIELDS.OS_IMAGE, osImages);
            setDefaultOsImage(azureDispatch, osImages);
            setErrorObject(removeErrorInfo(errorObject, AZURE_FIELDS.OS_IMAGE));
        } catch (e) {
            clearPreviousResourceData(azureDispatch, AZURE_FIELDS.OS_IMAGE);
            setErrorObject(addErrorInfo(errorObject, e, AZURE_FIELDS.OS_IMAGE));
        }
    }
}

function clearPreviousResourceData(azureDispatch: StoreDispatch, resourceName: AZURE_FIELDS) {
    azureDispatch({
        type: AZURE_ADD_RESOURCES,
        resourceName: resourceName,
        payload: [],
    } as AzureResourceAction);
}

function saveCurrentResourceData(azureDispatch: StoreDispatch, resourceName: AZURE_FIELDS, currentValues: any[]) {
    azureDispatch({
        type: AZURE_ADD_RESOURCES,
        resourceName: resourceName,
        payload: currentValues,
    } as AzureResourceAction);
}

function setDefaultOsImage(azureDispatch: StoreDispatch, osImages: AzureVirtualMachine[]) {
    azureDispatch({
        type: INPUT_CHANGE,
        field: AZURE_FIELDS.OS_IMAGE,
        payload: AzureDefaults.selectDefalutOsImage(osImages),
    } as FormAction);
}

function removeErrorInfo(errorObject: { [key: string]: any }, field: AZURE_FIELDS) {
    const copy = { ...errorObject };
    delete copy[field];
    return copy;
}

function addErrorInfo(errorObject: { [key: string]: any }, error: any, field: AZURE_FIELDS) {
    return {
        ...errorObject,
        [field]: error,
    };
}
