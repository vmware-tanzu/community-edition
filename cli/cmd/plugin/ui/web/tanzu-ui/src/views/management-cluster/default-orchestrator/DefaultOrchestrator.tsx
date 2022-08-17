import { ResourceAction, StoreDispatch } from '../../../shared/types/types';
import { RESOURCE } from '../../../state-management/actions/Resources.actions';
import { AWS_FIELDS } from '../aws/aws-mc-basic/AwsManagementClusterBasic.constants';
import { AZURE_FIELDS } from '../azure/azure-mc-basic/AzureManagementClusterBasic.constants';

export function clearPreviousResourceData(azureDispatch: StoreDispatch, actionType: RESOURCE, resourceName: AWS_FIELDS | AZURE_FIELDS) {
    azureDispatch({
        type: actionType,
        resourceName: resourceName,
        payload: [],
    } as ResourceAction);
}

export function saveCurrentResourceData(
    azureDispatch: StoreDispatch,
    actionType: RESOURCE,
    resourceName: AWS_FIELDS | AZURE_FIELDS,
    currentValues: any[]
) {
    azureDispatch({
        type: actionType,
        resourceName: resourceName,
        payload: currentValues,
    } as ResourceAction);
}

export function removeErrorInfo(errorObject: { [key: string]: any }, field: AWS_FIELDS | AZURE_FIELDS) {
    const copy = { ...errorObject };
    delete copy[field];
    return copy;
}

export function addErrorInfo(errorObject: { [key: string]: any }, error: any, field: AWS_FIELDS | AZURE_FIELDS) {
    return {
        ...errorObject,
        [field]: error,
    };
}
