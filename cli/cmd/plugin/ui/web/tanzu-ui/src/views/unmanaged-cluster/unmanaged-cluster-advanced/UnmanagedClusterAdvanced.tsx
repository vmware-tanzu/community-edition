// React imports
import React from 'react';

// App imports
import { UMC_ADVANCED_TAB_NAMES } from './UnmanagedClusterAdvanced.constants';
import Wizard from '../../../shared/components/wizard/Wizard';
import UnmanagedClusterSettingsAdvanced from '../unmanaged-cluster-basic/unmanaged-cluster-settings-step/UnmanagedClustersSettingsAdvanced';
// eslint-disable-next-line max-len
import UnmanagedClusterNetworkSettings from '../unmanaged-cluster-basic/unmanaged-cluster-network-settings-step/UnmanagedClustersNetworkSettings';

function UnmanagedClusterAdvanced() {
    return (
        <Wizard tabNames={UMC_ADVANCED_TAB_NAMES}>
            <UnmanagedClusterSettingsAdvanced />
            <UnmanagedClusterNetworkSettings />
        </Wizard>
    );
}

export default UnmanagedClusterAdvanced;
