// React imports
import React from 'react';

// App imports
import { AzureBasicReviewStep } from './AzureBasicReviewStep';
import { AzureClusterSettingsStep } from '../azure-mc-common/AzureClusterSettingsStep';
import { AZURE_MC_BASIC_TAB_NAMES } from './AzureManagementClusterBasic.constants';
import ManagementCredentials from '../azure-mc-common/management-credential-step/ManagementCredentials';
import useAzureDeployment from '../../../../shared/services/azureDeployment';
import Wizard from '../../../../shared/components/wizard/Wizard';

function AzureManagementClusterBasic() {
    const { deployOnAzure } = useAzureDeployment();

    return (
        <Wizard tabNames={AZURE_MC_BASIC_TAB_NAMES}>
            <ManagementCredentials />
            <AzureClusterSettingsStep />
            <AzureBasicReviewStep deploy={deployOnAzure} />
        </Wizard>
    );
}

export default AzureManagementClusterBasic;
