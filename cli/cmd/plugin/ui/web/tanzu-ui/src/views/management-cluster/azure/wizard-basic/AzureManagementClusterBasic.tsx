// React imports
import React, { useContext } from 'react';

// App imports
import { AzureStore } from '../../../../state-management/stores/Azure.store';
import { AWS_MC_BASIC_TAB_NAMES } from '../../../../shared/constants/NavRoutes.constants';
import Wizard from '../../../../shared/components/wizard/Wizard';
import ManagementClusterSettings from '../../../../shared/components/management-cluster-settings-step/ManagementClusterSettings';
import useAzureDeployment from '../../../../shared/services/azureDeployment';
import ManagementCredentials from './management-credential-step/ManagementCredentials';
import { AzureService, AzureVirtualMachine } from '../../../../swagger-api';

function AzureManagementClusterBasic() {
    const { azureState, azureDispatch } = useContext(AzureStore);
    const { deployOnAzure } = useAzureDeployment();
    const getAzureImageMethod = async () => {
        const awsImageList = AzureService.getAzureOsImages();
        return awsImageList;
    };

    return (
        <Wizard tabNames={AWS_MC_BASIC_TAB_NAMES} state={azureState} dispatch={azureDispatch}>
            <ManagementCredentials />
            <ManagementClusterSettings<AzureVirtualMachine>
                deploy={deployOnAzure}
                defaultData={azureState}
                getImageMethod={getAzureImageMethod}
                clusterName={'Azure Machine Image(AMI)'}
            />
        </Wizard>
    );
}

export default AzureManagementClusterBasic;
