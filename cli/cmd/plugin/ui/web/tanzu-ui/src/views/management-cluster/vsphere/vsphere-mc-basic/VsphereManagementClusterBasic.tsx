// React imports
import React, { useContext } from 'react';

import { VSPHERE_MC_BASIC_TAB_NAMES } from './VsphereManagementClusterBasic.constants';
import { VsphereClusterResourcesStep } from '../vsphere-mc-common/VsphereClusterResourcesStep';
import { VsphereClusterSettingsStep } from '../vsphere-mc-common/VsphereClusterSettingsStep';
import { VsphereCredentialsStep } from '../vsphere-mc-common/VsphereCredentialsStep';
import { VsphereLoadBalancerStep } from '../vsphere-mc-common/VsphereLoadBalancerStep';
import { VsphereStore } from '../Store.vsphere.mc';
// App imports
import Wizard from '../../../../shared/components/wizard/Wizard';
import useVSphereDeployment from '../../../../shared/services/vSphereDeployment';

function VsphereManagementClusterBasic() {
    const { vsphereState, vsphereDispatch } = useContext(VsphereStore);
    const { deployOnVSphere } = useVSphereDeployment();

    return (
        <Wizard tabNames={VSPHERE_MC_BASIC_TAB_NAMES} state={vsphereState} dispatch={vsphereDispatch}>
            <VsphereCredentialsStep />
            <VsphereClusterSettingsStep />
            <VsphereLoadBalancerStep />
            <VsphereClusterResourcesStep deploy={deployOnVSphere} />
        </Wizard>
    );
}

export default VsphereManagementClusterBasic;
