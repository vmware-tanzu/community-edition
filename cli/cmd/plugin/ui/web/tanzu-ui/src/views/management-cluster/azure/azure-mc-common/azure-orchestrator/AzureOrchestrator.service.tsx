// App imports
import { AzureService, AzureVirtualMachine } from '../../../../../swagger-api';
import { AZURE_FIELDS } from '../../azure-mc-basic/AzureManagementClusterBasic.constants';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { AzureDefaults } from '../default-service/AzureDefaults.service';
import { FormAction, StoreDispatch } from '../../../../../shared/types/types';
import {
    clearPreviousResourceData,
    saveCurrentResourceData,
    removeErrorInfo,
    addErrorInfo,
} from '../../../default-orchestrator/DefaultOrchestrator';
import { RESOURCE } from '../../../../../state-management/actions/Resources.actions';
interface AzureOrchestratorProps {
    azureState: { [key: string]: any };
    azureDispatch: StoreDispatch;
    errorObject: { [key: string]: any };
    setErrorObject: (newErrorObject: { [key: string]: any }) => void;
}
export class AzureOrchestrator {
    static async initOsImages(props: AzureOrchestratorProps) {
        const { azureDispatch, setErrorObject, errorObject } = props;
        try {
            const osImages = await AzureService.getAzureOsImages();
            saveCurrentResourceData(azureDispatch, RESOURCE.AZURE_ADD_RESOURCES, AZURE_FIELDS.OS_IMAGE, osImages);
            setDefaultOsImage(azureDispatch, osImages);
            setErrorObject(removeErrorInfo(errorObject, AZURE_FIELDS.OS_IMAGE));
        } catch (e) {
            clearPreviousResourceData(azureDispatch, RESOURCE.AZURE_ADD_RESOURCES, AZURE_FIELDS.OS_IMAGE);
            setErrorObject(addErrorInfo(errorObject, e, AZURE_FIELDS.OS_IMAGE));
        }
    }
}

function setDefaultOsImage(azureDispatch: StoreDispatch, osImages: AzureVirtualMachine[]) {
    azureDispatch({
        type: INPUT_CHANGE,
        field: AZURE_FIELDS.OS_IMAGE,
        payload: AzureDefaults.selectDefaultOsImage(osImages),
    } as FormAction);
}
