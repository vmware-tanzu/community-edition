// React imports
import React, { useContext } from 'react';

// App imports
import { UmcStore } from '../../../state-management/stores/Store.umc';
import { UMC_BASIC_TAB_NAMES } from './UnmanagedClusterBasic.constants';
import Wizard from '../../../shared/components/wizard/Wizard';
import UnmanagedClusterSettingsBasic from './unmanaged-cluster-settings-step/UnmanagedClustersSettingsBasic';
import UnmanagedClusterNetworkSettings from './unmanaged-cluster-network-settings-step/UnmanagedClustersNetworkSettings';

function UnmanagedClusterBasic() {
    const { umcState, umcDispatch } = useContext(UmcStore);

    return (
        <Wizard tabNames={UMC_BASIC_TAB_NAMES} state={umcState} dispatch={umcDispatch}>
            <UnmanagedClusterSettingsBasic />
            <UnmanagedClusterNetworkSettings />
        </Wizard>
    );
}

export default UnmanagedClusterBasic;
