// React imports
import { useContext, useEffect, useState } from 'react';

// App imports
import { AzureInstanceType, AzureService, AzureVirtualMachine } from '../../../../../swagger-api';
import { AZURE_FIELDS } from '../../azure-mc-basic/AzureManagementClusterBasic.constants';
import { AzureDefaults } from '../default-service/AzureDefaults.service';
import { AzureStore } from '../../store/Azure.store.mc';
import { BATCH_SET } from '../../../../../state-management/actions/Form.actions';
import { DefaultOrchestrator } from '../../../default-orchestrator/DefaultOrchestrator';
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';
import useCompare from '../../../../../shared/hooks/UseCompare';

function usePrerequisite() {
    const [errorObject, setErrorObject] = useState({});
    const { azureState, azureDispatch } = useContext(AzureStore);
    const location = azureState[STORE_SECTION_FORM][AZURE_FIELDS.REGION];
    const locationChanged = useCompare(location);

    return [{ location, locationChanged, azureState, azureDispatch, errorObject, setErrorObject }];
}

function useInitOsImages() {
    const [{ location, locationChanged, azureDispatch, errorObject, setErrorObject }] = usePrerequisite();

    useEffect(() => {
        const fetchOsImages = async () => {
            await DefaultOrchestrator.initResources<AzureVirtualMachine>({
                resourceName: AZURE_FIELDS.OS_IMAGE,
                errorObject,
                setErrorObject,
                dispatch: azureDispatch,
                fetcher: AzureService.getAzureOsImages,
                fxnSelectDefault: AzureDefaults.selectDefaultOsImage,
            });
        };
        if (location && locationChanged) {
            fetchOsImages();
        }
    }, [azureDispatch, errorObject, location, locationChanged, setErrorObject]);

    return [errorObject, setErrorObject];
}

function useInitNodeProfile() {
    const [{ location, locationChanged, azureState, azureDispatch, errorObject, setErrorObject }] = usePrerequisite();

    const selectedNodeProfile = azureState[STORE_SECTION_FORM][AZURE_FIELDS.NODE_PROFILE];

    useEffect(() => {
        const fetchInstanceTypes = async () => {
            const nodeTypes = await DefaultOrchestrator.initResources<AzureInstanceType>({
                resourceName: AZURE_FIELDS.NODE_TYPE,
                errorObject,
                setErrorObject,
                dispatch: azureDispatch,
                fetcher: () => AzureService.getAzureInstanceTypes(location),
            });
            const nodeType = AzureDefaults.getDefaultNodeType(nodeTypes, selectedNodeProfile);

            azureDispatch({
                type: BATCH_SET,
                payload: {
                    [AZURE_FIELDS.CONTROL_PLANE_MACHINE_TYPE]: nodeType,
                    [AZURE_FIELDS.WORKER_MACHINE_TYPE]: nodeType,
                },
            });
        };
        if (location && locationChanged) {
            fetchInstanceTypes();
        }
    }, [azureDispatch, errorObject, location, locationChanged, selectedNodeProfile, setErrorObject]);
    return [errorObject, setErrorObject];
}

const AzureOrchestrator = {
    useInitOsImages,
    useInitNodeProfile,
};
export default AzureOrchestrator;
