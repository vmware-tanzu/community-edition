// App imports
import { AzureService, AzureVirtualMachine, ApiError } from '../../../../../swagger-api';
import { AZURE_FIELDS } from '../../azure-mc-basic/AzureManagementClusterBasic.constants';
import { AzureDefaults } from '../default-service/AzureDefaults.service';
import { DefaultOrchestrator } from '../../../default-orchestrator/DefaultOrchestrator';
import { StoreDispatch } from '../../../../../shared/types/types';

interface AzureOrchestratorProps {
    azureState: { [key: string]: any };
    azureDispatch: StoreDispatch;
    errorObject: { [fieldName: string]: ApiError };
    setErrorObject: (newErrorObject: { [fieldName: string]: ApiError }) => void;
}
export class AzureOrchestrator {
    static async initOsImages(props: AzureOrchestratorProps) {
        const { azureDispatch, setErrorObject, errorObject } = props;
        DefaultOrchestrator.initResources<AzureVirtualMachine>({
            resourceName: AZURE_FIELDS.OS_IMAGE,
            errorObject,
            setErrorObject,
            dispatch: azureDispatch,
            fetcher: AzureService.getAzureOsImages,
            fxnSelectDefault: AzureDefaults.selectDefaultOsImage,
        });
    }
}
