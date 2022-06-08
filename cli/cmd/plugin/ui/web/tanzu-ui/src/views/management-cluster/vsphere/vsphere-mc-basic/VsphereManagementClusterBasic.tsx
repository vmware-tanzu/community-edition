// React imports
import React, { useContext } from 'react';
// App imports
import Wizard from '../../../../shared/components/wizard/Wizard';
import { VSPHERE_MC_BASIC_TAB_NAMES } from './VsphereManagementClusterBasic.constants';
import { VsphereCredentialsStep } from '../vsphere-mc-common/VsphereCredentialsStep';
import { VsphereClusterSettingsStep } from '../vsphere-mc-common/VsphereClusterSettingsStep';
import { VsphereClusterResourcesStep } from '../vsphere-mc-common/VsphereClusterResourcesStep';
import { VsphereStore } from '../../../../state-management/stores/Store.vsphere.mc';

function VsphereManagementClusterBasic() {
    const { vsphereState, vsphereDispatch } = useContext(VsphereStore);
    // TODO: create a similar useVsphereDeployment() function
    // const { deployOnAws } = useAwsDeployment();

    return (
        <Wizard tabNames={VSPHERE_MC_BASIC_TAB_NAMES} state={vsphereState} dispatch={vsphereDispatch}>
            <VsphereCredentialsStep />
            <VsphereClusterSettingsStep />
            <VsphereClusterResourcesStep />
        </Wizard>
    );
}

export default VsphereManagementClusterBasic;
