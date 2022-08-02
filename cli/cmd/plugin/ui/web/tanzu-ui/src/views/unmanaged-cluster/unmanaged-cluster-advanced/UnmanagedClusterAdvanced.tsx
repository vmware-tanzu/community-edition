// React imports
import React, { useContext } from 'react';

// App imports
import { UmcStore } from '../../../state-management/stores/Store.umc';
import { UMC_ADVANCED_TAB_NAMES } from './UnmanagedClusterAdvanced.constants';
import Wizard from '../../../shared/components/wizard/Wizard';
import UnmanagedClusterSettingsAdvanced from '../unmanaged-cluster-basic/unmanaged-cluster-settings-step/UnmanagedClustersSettingsAdvanced';
import UnmanagedClusterNetworkSettings from '../unmanaged-cluster-basic/unmanaged-cluster-network-settings-step/UnmanagedClustersNetworkSettings';

function UnmanagedClusterAdvanced() {
    const { umcState, umcDispatch } = useContext(UmcStore);

    return (
        <Wizard tabNames={UMC_ADVANCED_TAB_NAMES} state={umcState} dispatch={umcDispatch}>
            <UnmanagedClusterSettingsAdvanced />
            <UnmanagedClusterNetworkSettings />
        </Wizard>
    );
}

export default UnmanagedClusterAdvanced;
