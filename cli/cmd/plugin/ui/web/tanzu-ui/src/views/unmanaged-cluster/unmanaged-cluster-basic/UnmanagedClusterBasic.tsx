// React imports
import React from 'react';

// App imports
import { UMC_BASIC_TAB_NAMES } from './UnmanagedClusterBasic.constants';
import Wizard from '../../../shared/components/wizard/Wizard';
import UnmanagedClusterSettingsBasic from './unmanaged-cluster-settings-step/UnmanagedClustersSettingsBasic';
import UnmanagedClusterNetworkSettings from './unmanaged-cluster-network-settings-step/UnmanagedClustersNetworkSettings';

function UnmanagedClusterBasic() {
    return (
        <Wizard tabNames={UMC_BASIC_TAB_NAMES}>
            <UnmanagedClusterSettingsBasic />
            <UnmanagedClusterNetworkSettings />
        </Wizard>
    );
}

export default UnmanagedClusterBasic;
