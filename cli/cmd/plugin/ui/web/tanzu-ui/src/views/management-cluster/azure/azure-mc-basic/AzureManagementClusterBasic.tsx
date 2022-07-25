// React imports
import React, { useContext } from 'react';

// App imports
import { AzureClusterSettingsStep } from '../azure-mc-common/AzureClusterSettingsStep';
import { AzureStore } from '../../../../state-management/stores/Azure.store';
import { AZURE_MC_BASIC_TAB_NAMES } from './AzureManagementClusterBasic.constants';
import ManagementCredentials from '../azure-mc-common/management-credential-step/ManagementCredentials';
import useAzureDeployment from '../../../../shared/services/azureDeployment';
import Wizard from '../../../../shared/components/wizard/Wizard';

function AzureManagementClusterBasic() {
    const { azureState, azureDispatch } = useContext(AzureStore);
    const { deployOnAzure } = useAzureDeployment();

    return (
        <Wizard tabNames={AZURE_MC_BASIC_TAB_NAMES} state={azureState} dispatch={azureDispatch}>
            <ManagementCredentials />
            <AzureClusterSettingsStep deploy={deployOnAzure} />
        </Wizard>
    );
}

export default AzureManagementClusterBasic;
