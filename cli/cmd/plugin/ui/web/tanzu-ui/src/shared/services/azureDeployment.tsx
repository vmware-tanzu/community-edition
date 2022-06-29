// React imports
import { useContext } from 'react';
import { useNavigate } from 'react-router-dom';

// App imports
import { DEPLOYMENT_STATUS_CHANGED } from '../../state-management/actions/Deployment.actions';
import { TOGGLE_APP_STATUS } from '../../state-management/actions/Ui.actions';
import { Store } from '../../state-management/stores/Store';
import { AzureStore } from '../../state-management/stores/Azure.store';
import { AWSManagementClusterParams, AzureService, ConfigFileInfo } from '../../swagger-api';
import { DeploymentStates, DeploymentTypes } from '../constants/Deployment.constants';
import { NavRoutes } from '../constants/NavRoutes.constants';
import { Providers } from '../constants/Providers.constants';
import { STORE_SECTION_FORM } from '../../state-management/reducers/Form.reducer';

const useAzureDeployment = () => {
    const { dispatch } = useContext(Store);
    const { azureState } = useContext(AzureStore);
    const navigate = useNavigate();
    const navigateToProgress = (): void => {
        navigate('/' + NavRoutes.DEPLOY_PROGRESS);
    };

    const getAzureRequestPayload = () => {
        //TODO: Add request payload here.
        return azureState[STORE_SECTION_FORM];
    };

    const deployOnAzure = async () => {
        const awsClusterParams: AWSManagementClusterParams = getAzureRequestPayload();
        try {
            const configFileInfo: ConfigFileInfo = await AzureService.applyTkgConfigForAzure(awsClusterParams);
            await AzureService.createAzureManagementCluster(awsClusterParams);
            dispatch({
                type: DEPLOYMENT_STATUS_CHANGED,
                payload: {
                    type: DeploymentTypes.MANAGEMENT_CLUSTER,
                    status: DeploymentStates.RUNNING,
                    provider: Providers.AZURE,
                    configPath: configFileInfo.path,
                },
            });
        } catch (e) {
            console.log(`Error when calling config or create API: ${e}`);
            return;
        }

        dispatch({
            type: TOGGLE_APP_STATUS,
        });
        navigateToProgress();
    };

    return {
        getAzureRequestPayload,
        deployOnAzure,
    };
};
export default useAzureDeployment;
