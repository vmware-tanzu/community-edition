// React imports
import React, { useContext } from 'react';

// App imports
import { AzureStore } from '../../../../state-management/stores/Azure.store';
import { AWS_MC_BASIC_TAB_NAMES } from '../../../../shared/constants/NavRoutes.constants';
import Wizard from '../../../../shared/components/wizard/Wizard';
import ManagementClusterSettings from '../../../../shared/components/management-cluster-settings-step/ManagementClusterSettings';
import useAzureDeployment from '../../../../shared/services/azureDeployment';
import ManagementCredentials from './management-credential-step/ManagementCredentials';

function AzureManagementClusterBasic() {
    const { azureState, azureDispatch } = useContext(AzureStore);
    const { deployOnAzure } = useAzureDeployment();

    return (
        <Wizard tabNames={AWS_MC_BASIC_TAB_NAMES} state={azureState} dispatch={azureDispatch}>
            <ManagementCredentials />
            <ManagementClusterSettings deploy={deployOnAzure} defaultData={azureState} />
        </Wizard>
    );
}

export default AzureManagementClusterBasic;
