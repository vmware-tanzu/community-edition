// React imports
import React, { useContext } from 'react';

// App imports
import { UmcStore } from '../../../state-management/stores/Store.umc';
import { UMC_ADVANCED_TAB_NAMES } from '../../../shared/constants/NavRoutes.constants';
import Wizard from '../../../shared/components/wizard/Wizard';
import UnmanagedClusterSettingsAdvanced from './unmanaged-cluster-settings-step/UnmanagedClustersSettingsAdvanced';
import UnmanagedClusterNetworkSettings from './unmanaged-cluster-network-settings-step/UnmanagedClustersNetworkSettings';

function UnmanagedClusterWizardAdvanced() {
    const { umcState, umcDispatch } = useContext(UmcStore);

    return (
        <Wizard tabNames={UMC_ADVANCED_TAB_NAMES} state={umcState} dispatch={umcDispatch}>
            <UnmanagedClusterSettingsAdvanced />
            <UnmanagedClusterNetworkSettings />
        </Wizard>
    );
}

export default UnmanagedClusterWizardAdvanced;
